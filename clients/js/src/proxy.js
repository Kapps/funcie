const { sendMessage } = require("./bastionClient");
const { Message } = require("./models");
const { invokeLambda, info, error, debug} = require("./utils");
const { loadConfig } = require('./config');

// TODO: Proper response codes. This is... gross.
const errNoConsumerActive = 'no consumer is active on this tunnel';
const errApplicationNotFound = 'application not found';

/**
 * Returns a lambda handler much like `lambdaProxyWithConfig`, but lazily loads the configuration before first execution.
 * @param appId - arbitrary unique application identifier.
 * @param handler - lambda handler to be wrapped.
 * @returns {function(*, *): Promise<*>}
 */
const lambdaProxy = (appId, handler) => {
    let conf;
    return async (event, context) => {
        if (!conf) {
            conf = await loadConfig(appId);
        }
        return await lambdaProxyWithConfig(conf, handler)(event, context);
    }
};

/**
 * Returns a lambda handler that forwards requests to the bastion server or handles them directly when required.
 * @param config - configuration object, often loaded from `loadConfig`
 * @param handler - lambda handler to be wrapped.
 * @returns {(function(*, *): Promise<unknown>)|*}
 */
const lambdaProxyWithConfig = (config, handler) => {
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
                debug(`no consumer active on bastion; handling request directly`);
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
    lambdaProxyWithConfig,
};
