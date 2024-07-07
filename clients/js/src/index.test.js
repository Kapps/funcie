import { lambdaWrapper } from './index.js';
import { beginReceiving } from './receiver.js';
import { lambdaProxy } from './proxy.js';
import * as config from './config.js';

import { jest } from '@jest/globals';

jest.mock('./receiver.js');
jest.mock('./proxy.js');
jest.mock('./config.js');

describe('lambdaWrapper', () => {
    let conf;

    beforeEach(() => {
        jest.clearAllMocks();

        conf = new config.FuncieConfig(undefined, undefined, undefined, 'app');
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
