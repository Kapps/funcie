const { FuncieConfig, loadConfigFromEnvironment, getConfigPurpose } = require('./config');
const { URL } = require('url');

describe('FuncieConfig Module Tests', () => {
    let originalEnv;

    beforeEach(() => {
        originalEnv = { ...process.env };

        process.env = {
            AWS_LAMBDA_FUNCTION_NAME: undefined,
            FUNCIE_CLIENT_BASTION_ENDPOINT: "http://localhost/client",
            FUNCIE_SERVER_BASTION_ENDPOINT: "http://localhost/server",
            FUNCIE_LISTEN_ADDRESS: "127.0.0.1:3000",
            FUNCIE_APPLICATION_ID: "test_app_id"
        };
    });

    afterEach(() => {
        process.env = originalEnv;
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
        it('should create a config object correctly', () => {
            const config = new FuncieConfig(
                new URL("http://localhost/client"),
                new URL("http://localhost/server"),
                "127.0.0.1:3000",
                "test_app_id"
            );

            expect(config.ClientBastionEndpoint).toEqual(new URL("http://localhost/client"));
            expect(config.ServerBastionEndpoint).toEqual(new URL("http://localhost/server"));
            expect(config.ListenAddress).toBe("127.0.0.1:3000");
            expect(config.ApplicationId).toBe("test_app_id");
        });
    });

    describe('loadConfigFromEnvironment', () => {
        it('should load a config from environment variables', () => {
            const config = loadConfigFromEnvironment();

            expect(config.ClientBastionEndpoint).toEqual(new URL("http://localhost/client"));
            expect(config.ServerBastionEndpoint).toEqual(new URL("http://localhost/server"));
            expect(config.ListenAddress).toBe("127.0.0.1:3000");
            expect(config.ApplicationId).toBe("test_app_id");
        });
    });
});
