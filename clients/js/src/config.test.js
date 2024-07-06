const { FuncieConfig, getConfigPurpose, loadConfig} = require('./config');
const { URL } = require('url');

jest.mock('@aws-sdk/client-ssm', () => {
    const originalModule = jest.requireActual('@aws-sdk/client-ssm');
    return {
        ...originalModule,
        SSMClient: jest.fn().mockImplementation(() => ({
            send: jest.fn().mockResolvedValue({
                Parameter: {
                    Value: "example.org"
                }
            })
        }))
    };
});

describe('FuncieConfig Module Tests', () => {
    let originalEnv;

    beforeEach(() => {
        originalEnv = { ...process.env };

        process.env = {
            AWS_LAMBDA_FUNCTION_NAME: undefined,
            FUNCIE_CLIENT_BASTION_ENDPOINT: "http://localhost/client",
            FUNCIE_LISTEN_ADDRESS: "http://127.0.0.1:3000",
            FUNCIE_APPLICATION_ID: "test_app_id",
        };
    });

    afterEach(() => {
        process.env = originalEnv;
        jest.clearAllMocks();
        jest.resetModules();
    });

    describe('getConfigPurpose', () => {
        it('should return "server" if AWS_LAMBDA_FUNCTION_NAME is set', () => {
            process.env.AWS_LAMBDA_FUNCTION_NAME = "test_function";
            expect(getConfigPurpose()).toBe("server");
        });

        it('should return "client" if AWS_LAMBDA_FUNCTION_NAME is not set', () => {
            process.env.AWS_LAMBDA_FUNCTION_NAME = undefined;
            expect(getConfigPurpose()).toBe("client");
        });
    });

    describe('FuncieConfig', () => {
        it('should create a config object correctly with the inputs provided', () => {
            const config = new FuncieConfig(
                new URL("http://localhost/client"),
                new URL("http://localhost/server"),
                "127.0.0.1:3000",
                "test_app_id",
            );

            expect(config.ClientBastionEndpoint).toEqual(new URL("http://localhost/client"));
            expect(config.ServerBastionEndpoint).toEqual(new URL("http://localhost/server"));
            expect(config.ListenAddress).toBe("127.0.0.1:3000");
            expect(config.ApplicationId).toBe("test_app_id");
        });
    });

    describe('loadConfig', () => {
        afterEach(() => {
            delete process.env.FUNCIE_SERVER_BASTION_ENDPOINT
        });

        it('should load a config from environment variables and SSM', async () => {
            const config = await loadConfig("test_app_id", "test_env");

            expect(config.ClientBastionEndpoint).toEqual(new URL("http://localhost/client"));
            expect(config.ServerBastionEndpoint).toEqual(new URL("http://example.org:8082/dispatch"));
            expect(config.ListenAddress.toString()).toBe("http://127.0.0.1:3000/");
            expect(config.ApplicationId).toBe("test_app_id");
        });

        it('should load a config from only environment variables if all are set', async () => {
            process.env.FUNCIE_SERVER_BASTION_ENDPOINT = "http://localhost/serverenv";

            const config = await loadConfig("test_app_id", "test_env");

            expect(config.ClientBastionEndpoint).toEqual(new URL("http://localhost/client"));
            expect(config.ServerBastionEndpoint).toEqual(new URL("http://localhost/serverenv"));
            expect(config.ListenAddress.toString()).toBe("http://127.0.0.1:3000/");
            expect(config.ApplicationId).toBe("test_app_id");
        });
    });
});
