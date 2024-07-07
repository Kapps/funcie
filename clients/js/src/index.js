import { beginReceiving } from './receiver.js';
import { lambdaProxy } from './proxy.js';
import { info } from './utils.js';

export const lambdaWrapper = (appId, handler, config) => {
    const isRunningInLambda = !!process.env.AWS_LAMBDA_FUNCTION_NAME;

    if (isRunningInLambda) {
        info('Starting Funcie proxy for app:', config.ApplicationId);
        return lambdaProxy(config, handler);
    }

    info('Starting Funcie server for app:', config.ApplicationId);
    return beginReceiving(config, handler);
};
