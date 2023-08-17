const { beginReceiving } = require('./receiver');
const { lambdaProxy } = require('./proxy');


const lambdaWrapper = (handler) => {
    const config = require('./config').loadConfigFromEnvironment();
    const isRunningInLambda = !!process.env.AWS_LAMBDA_FUNCTION_NAME;

    if (isRunningInLambda) {
        console.log('Starting proxy with config: ', JSON.stringify(config));
        return lambdaProxy(handler);
    }

    console.log('Starting server with config: ', JSON.stringify(config));
    return beginReceiving(config, handler);
};

module.exports = {
    lambdaWrapper,
};
