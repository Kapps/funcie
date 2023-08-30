const { sendMessage } = require("./bastionClient");
const { loadConfigFromEnvironment } = require("./config");
const { Message } = require("./models");
const { invokeLambda } = require("./utils");

const lambdaProxy = (handler) => {
    return async (event, context) => {
        const config = loadConfigFromEnvironment();
        const payload = {
            body: event,
        };
        const forwardRequest = new Message('FORWARD_REQUEST', config.ApplicationId, payload);
        let forwardResponse;
        try {
            forwardResponse = await sendMessage(config.ServerBastionEndpoint, forwardRequest);
        } catch (err) {
            console.log(`failed to send request to bastion: ${err}; handling directly`);
            return invokeLambda(handler, event, context);
        }

        return forwardResponse.data.body;
    };
};

module.exports = {
    lambdaProxy,
};
