const { beginReceiving } = require('./receiver');
const { lambdaProxy } = require('./proxy');
const { info } = require('./utils');


const lambdaWrapper = async (appId, handler) => {
    const config = await require('./config').loadConfig(appId);
    const isRunningInLambda = !!process.env.AWS_LAMBDA_FUNCTION_NAME;

    if (isRunningInLambda) {
        info('Starting Funcie proxy for app:', config.ApplicationId);
        return lambdaProxy(config, handler);
    }

    info('Starting Funcie server for app:', config.ApplicationId);
    return beginReceiving(config, handler);
};

module.exports = {
    lambdaWrapper,
};
