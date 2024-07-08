const http = require('http');
const { promisify } = require('util');
const { Message, Response } = require('./models');
const { sendMessage } = require('./bastionClient');
const { invokeLambda, info, error } = require('./utils');
const { loadConfig } = require("./config");

/**
 * Helper function to load the default configuration and pass it to `beginReceivingWithConfig`.
 * @param appId - arbitrary unique application identifier.
 * @param handler - lambda handler to be wrapped.
 * @returns {Promise<Server<IncomingMessage, ServerResponse>>}
 */
const beginReceiving = async (appId, handler) => {
    const conf = await loadConfig(appId);
    return beginReceivingWithConfig(conf, handler);
};

/**
 * Starts a server that listens for incoming requests and forwards them to the provided handler.
 * @param config - configuration object, often loaded from `loadConfig`
 * @param handler - lambda handler to be wrapped.
 * @returns {Promise<Server<typeof IncomingMessage, typeof ServerResponse>>}
 */
const beginReceivingWithConfig = async (config, handler) => {
    if (config.ListenAddress.protocol !== 'http:') {
        throw new Error('Only HTTP is supported');
    }

    const server = http.createServer(async (req, res) => {
        const body = [];
        req.on('data', (chunk) => {
            body.push(chunk);
        }).on('end', async () => {
            const data = Buffer.concat(body).toString();
            const message = Message.fromJson(data);
            let response;
            try {
                const responseData = await invokeLambda(handler, message.payload.body);
                response = new Response(message.id, {
                    body: responseData,
                }, undefined, new Date());
            } catch (err) {
                response = new Response(message.id, undefined, { message: err.message }, new Date());
            }
            res.writeHead(200, {
                'Content-Type': 'application/json',
            });

            res.write(JSON.stringify(response));
            res.end();
        });
    });

    server.on('error', (err) => {
        error(`Server error: ${err.message}`);
    });

    await promisify(server.listen).bind(server)(config.ListenAddress.port, config.ListenAddress.hostname);

    info('Funcie Server Started: ', server.address());

    await subscribe(config, server.address());

    return server;
};

const subscribe = async (config, address) => {
    const app = {
        name: config.ApplicationId,
        endpoint: {
            protocol: 'http',
            host: address.address,
            port: address.port,
        },
    };

    const req = new Message('REGISTER', config.Application, app);
    const resp = await sendMessage(config.ClientBastionEndpoint, req);

    info(`Funcie registered with registration ID ${resp.data.RegistrationId}`);
};

module.exports = {
    beginReceiving,
    beginReceivingWithConfig,
}
