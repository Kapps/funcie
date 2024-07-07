const { sendMessage } = require("./bastionClient");
const { Message } = require("./models");
const { invokeLambda, info, error } = require("./utils");

// TODO: Proper response codes. This is... gross.
const errNoConsumerActive = 'no consumer is active on this tunnel';
const errApplicationNotFound = 'application not found';

const lambdaProxy = (config, handler) => {
    return async (event, context) => {
        const payload = {
            body: event,
        };
        const forwardRequest = new Message('FORWARD_REQUEST', config.ApplicationId, payload);

        let forwardResponse;
        try {
            forwardResponse = await sendMessage(config.ServerBastionEndpoint, forwardRequest);
        } catch (err) {
            // Failed to send request to bastion. This could be because the bastion is down or other reasons.
            // This shouldn't interrupt standard request flow; in this scenario handle requests directly.
            error(`failed to send request to bastion: ${err}; handling request directly`);
            return invokeLambda(handler, event, context);
        }

        if (forwardResponse.error) {
            // We received a response from the bastion, but an error response.
            // Can be either because we were unable to forward to the client (no client listening),
            // or because we successfully forwarded to the client but the client errored.

            // If there is no consumer active on the bastion, handle the request directly.
            if (forwardResponse.error.message === errNoConsumerActive || forwardResponse.error === errApplicationNotFound) {
                info(`no consumer active on bastion; handling request directly`);
                return invokeLambda(handler, event, context);
            }

            // Otherwise, client erred, so we should throw an error.
            throw new Error(forwardResponse.error.message);
        }

        return forwardResponse.data.body;
    };
};

module.exports = {
    lambdaProxy,
};
