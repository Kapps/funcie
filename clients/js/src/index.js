const { beginReceiving } = require('./receiver');
const { lambdaProxy } = require('./proxy');
const { info } = require('./utils');


const lambdaWrapper = (appId, handler) => {
    const isRunningInLambda = !!process.env.AWS_LAMBDA_FUNCTION_NAME;

    if (isRunningInLambda) {
        info('Starting Funcie proxy for app:', appId);
        return lambdaProxy(appId, handler);
    }

    info('Starting Funcie server for app:', appId);
    return beginReceiving(appId, handler);
};

module.exports = {
    lambdaWrapper,
};
