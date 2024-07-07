import http from 'http';
import { promisify } from 'util';
import { Message, Response } from './models.js';
import { sendMessage } from './bastionClient.js';
import { invokeLambda, info, error } from './utils.js';

export const beginReceiving = async (config, handler) => {
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
