const { lambdaWrapper } = require('./index');
const { beginReceiving } = require('./receiver');
const { lambdaProxy } = require('./proxy');
const config = require('./config');

jest.mock('./receiver');
jest.mock('./proxy');
jest.mock('./config');

describe('lambdaWrapper', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('should start the proxy if running in AWS Lambda', () => {
        process.env.AWS_LAMBDA_FUNCTION_NAME = 'my-function';
        const handler = jest.fn();
        const expectedProxy = jest.fn();
        lambdaProxy.mockReturnValueOnce(expectedProxy);

        const result = lambdaWrapper(handler);

        expect(result).toBe(expectedProxy);
        expect(lambdaProxy).toHaveBeenCalledWith(handler);
        expect(beginReceiving).not.toHaveBeenCalled();
        expect(config.loadConfigFromEnvironment).toHaveBeenCalled();

        delete process.env.AWS_LAMBDA_FUNCTION_NAME;
    });

    it('should start the server if not running in AWS Lambda', () => {
        const handler = jest.fn();
        const expectedServer = jest.fn();
        beginReceiving.mockReturnValueOnce(expectedServer);

        const result = lambdaWrapper(handler);

        expect(result).toBe(expectedServer);
        expect(beginReceiving).toHaveBeenCalledWith(config.loadConfigFromEnvironment(), handler);
        expect(lambdaProxy).not.toHaveBeenCalled();
        expect(config.loadConfigFromEnvironment).toHaveBeenCalled();
    });
});
