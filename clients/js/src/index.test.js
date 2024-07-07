const { lambdaWrapper } = require('./index');
const { beginReceiving } = require('./receiver');
const { lambdaProxy } = require('./proxy');
const config = require('./config');
const {FuncieConfig} = require("./config");

jest.mock('./receiver');
jest.mock('./proxy');
jest.mock('./config');

describe('lambdaWrapper', () => {
    let conf;

    beforeEach(() => {
        jest.clearAllMocks();

        conf = new FuncieConfig(undefined, undefined, undefined, 'app');
        config.loadConfig.mockReturnValue(conf);
    });

    it('should start the proxy if running in AWS Lambda', async () => {
        process.env.AWS_LAMBDA_FUNCTION_NAME = 'my-function';
        const handler = jest.fn();
        const expectedProxy = jest.fn();
        lambdaProxy.mockReturnValueOnce(expectedProxy);

        const result = await lambdaWrapper("app", handler);

        expect(result).toBe(expectedProxy);
        expect(lambdaProxy).toHaveBeenCalledWith(conf, handler);
        expect(beginReceiving).not.toHaveBeenCalled();
        expect(config.loadConfig).toHaveBeenCalledWith("app");

        delete process.env.AWS_LAMBDA_FUNCTION_NAME;
    });

    it('should start the server if not running in AWS Lambda', async () => {
        const handler = jest.fn();
        const expectedServer = jest.fn();
        beginReceiving.mockReturnValueOnce(expectedServer);

        const result = await lambdaWrapper('app', handler);

        expect(result).toBe(expectedServer);
        expect(beginReceiving).toHaveBeenCalledWith(conf, handler);
        expect(lambdaProxy).not.toHaveBeenCalledWith(handler);
        expect(config.loadConfig).toHaveBeenCalledWith('app');
    });
});
