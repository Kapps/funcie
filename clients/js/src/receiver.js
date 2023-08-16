const http = require('http');
const url = require('node:url');
const { promisify } = require('util');
const { Message, Response } = require('./models');
const axios = require('axios');

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
            // console.log(`Received message: ${message}`);
            const responseData = await invokeLambda(handler, message.payload);
            const response = new Response(message.id, {
                body: responseData,
            }, undefined, new Date());
            // console.log(`Sending response: ${response}`);
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
    // return;

    const resp = await axios.post(new url.URL('/dispatch', config.ClientBastionEndpoint), req);
    //const resp = axios.post(`${config.ClientBastionEndpoint}/dispatch`, req);
    /*const resp = await axios.post({
        url: new url.URL('/dispatch', config.ClientBastionEndpoint),
        data: req,
        headers: {
            'Content-Type': 'application/json',
        },
    });*/
    
    if (resp.error) {
        throw new Error(resp.error);
    }

    console.log(`Registered with registration ID ${resp.data.data.RegistrationId}`);

    //console.log('Registered with bastion: ', resp.data);
};

const invokeLambda = (handler, event, context) => {
    // Callback-based handlers are supported, but we wrap them in a promise
    if (handler.length >= 3) {
        return new Promise((resolve, reject) => {
            handler(event, context, (err, data) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(data);
                }
            });
        });
    }

    // Async handlers are just invoked directly.
    return handler(event, context);
}

module.exports = {
    beginReceiving,
}
