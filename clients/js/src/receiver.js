const http = require('http');
const { promisify } = require('util');
const { Message, Response } = require('./models');
const { sendMessage } = require('./bastionClient');
const { invokeLambda } = require('./utils');

const beginReceiving = async (config, handler) => {
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
            console.log(`Sending response: ${response}`);
            res.writeHead(200, {
                'Content-Type': 'application/json',
            });
    
            res.write(JSON.stringify(response));
            res.end();
        });
    });

    server.on('error', (err) => {
        console.error('Server error: ');
    });

    server.on('close', () => {
        console.log('Server closed');
    });

    await promisify(server.listen).bind(server)(config.ListenAddress.port, config.ListenAddress.hostname);

    console.log('Server started: ', server.address());

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

    console.log(`Registered with registration ID ${resp.data.RegistrationId}`);
};

module.exports = {
    beginReceiving,
}
