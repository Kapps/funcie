const { lambdaWrapper } = require('./index');
const { beginReceiving } = require('./receiver');
const { lambdaProxy } = require('./proxy');

jest.mock('./receiver');
jest.mock('./proxy');

describe('lambdaWrapper', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('should start the proxy if running in AWS Lambda', async () => {
        process.env.AWS_LAMBDA_FUNCTION_NAME = 'my-function';
        const handler = jest.fn();
        const expectedProxy = jest.fn();
        lambdaProxy.mockReturnValueOnce(expectedProxy);

        const result = await lambdaWrapper("app", handler);

        expect(result).toBe(expectedProxy);
        expect(lambdaProxy).toHaveBeenCalledWith('app', handler);
        expect(beginReceiving).not.toHaveBeenCalled();

        delete process.env.AWS_LAMBDA_FUNCTION_NAME;
    });

    it('should start the server if not running in AWS Lambda', async () => {
        const handler = jest.fn();
        const expectedServer = jest.fn();
        beginReceiving.mockReturnValueOnce(expectedServer);

        const result = await lambdaWrapper('app', handler);

        expect(result).toBe(expectedServer);
        expect(beginReceiving).toHaveBeenCalledWith('app', handler);
        expect(lambdaProxy).not.toHaveBeenCalledWith(handler);
    });
});
