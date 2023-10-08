const { beginReceiving } = require('./receiver');
const { lambdaProxy } = require('./proxy');
const { info } = require('./utils');


const lambdaWrapper = (handler) => {
    const config = require('./config').loadConfigFromEnvironment();
    const isRunningInLambda = !!process.env.AWS_LAMBDA_FUNCTION_NAME;

    if (isRunningInLambda) {
        info('Starting Funcie proxy with config: ', JSON.stringify(config));
        return lambdaProxy(handler);
    }

    info('Starting Funcie server with config: ', JSON.stringify(config));
    return beginReceiving(config, handler);
};

module.exports = {
    lambdaWrapper,
};
