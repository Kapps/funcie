const http = require("http");
const { beginReceiving, beginReceivingWithConfig} = require("./receiver");
const { sendMessage } = require("./bastionClient");
const { loadConfig } = require("./config");

jest.mock("http");
jest.mock("./bastionClient");
jest.mock("./utils");
jest.mock("./config");

describe("beginReceiving[WithConfig]", () => {
    const mockConfig = {
        ListenAddress: {
            protocol: "http:",
            port: 8080,
            hostname: "localhost",
        },
        ApplicationId: "app-id",
        ClientBastionEndpoint: "some-endpoint",
    };
    const mockHandler = jest.fn();
    let mockServer;
    let mockReq;
    let mockRes;

    beforeEach(() => {
        mockServer = {
            listen: jest.fn(),
            on: jest.fn(),
            address: jest.fn(),
        };
        mockReq = {
            on: jest.fn(),
        };
        mockRes = {
            writeHead: jest.fn(),
            write: jest.fn(),
            end: jest.fn(),
        };
        http.createServer.mockReturnValue(mockServer);
    });

    it("should throw error if protocol is not http", async () => {
        await expect(
            beginReceivingWithConfig(
                { ...mockConfig, ListenAddress: { protocol: "https:" } },
                mockHandler
            )
        ).rejects.toThrow("Only HTTP is supported");
    });

    it("should set up server and register application", async () => {
        const mockAddress = { address: "localhost", port: 8080 };
        const mockRegisterResp = { data: { RegistrationId: "reg-id" } };

        mockServer.address.mockReturnValue(mockAddress);
        sendMessage.mockResolvedValue(mockRegisterResp);
        mockServer.listen.mockImplementation((port, hostname, cb) => {
            cb();
        });

        await beginReceivingWithConfig(mockConfig, mockHandler);

        expect(mockServer.listen).toHaveBeenCalledWith(
            mockConfig.ListenAddress.port,
            mockConfig.ListenAddress.hostname,
            expect.any(Function)
        );
        expect(mockServer.on).toHaveBeenCalledTimes(1);
    });

    it('should forward from beginReceiving to beginReceivingWithConfig', async () => {
        loadConfig.mockReturnValue(mockConfig);

        const mockAddress = { address: "localhost", port: 8080 };
        const mockRegisterResp = { data: { RegistrationId: "reg-id" } };

        mockServer.address.mockReturnValue(mockAddress);
        sendMessage.mockResolvedValue(mockRegisterResp);
        mockServer.listen.mockImplementation((port, hostname, cb) => {
            cb();
        });

        await beginReceiving('app-id', mockHandler);
        expect(http.createServer).toHaveBeenCalled();
        expect(loadConfig).toHaveBeenCalledWith('app-id');
    });
});
