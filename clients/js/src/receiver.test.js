const http = require("http");
const { promisify } = require("util");
const { beginReceiving } = require("./receiver");
const { sendMessage } = require("./bastionClient");
const { invokeLambda } = require("./utils");
const { Message, Response } = require("./models");

jest.mock("http");
jest.mock("./bastionClient");
jest.mock("./utils");

describe("beginReceiving", () => {
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
            beginReceiving(
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

        await beginReceiving(mockConfig, mockHandler);

        expect(mockServer.listen).toHaveBeenCalledWith(
            mockConfig.ListenAddress.port,
            mockConfig.ListenAddress.hostname,
            expect.any(Function)
        );
        expect(mockServer.on).toHaveBeenCalledTimes(2);
    });

    // Add more test cases as required, for example, to test the request handling logic, error paths, etc.
});
