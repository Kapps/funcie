let done = false;
const { exit } = require('process');
const { beginReceiving } = require('./receiver');
const config = require('./config').loadConfigFromEnvironment();

console.log('Starting server with config: ', JSON.stringify(config));

const handler = (event) => {
    console.log('Receiving Event: ', JSON.stringify(event));

    return {
        statusCode: 200,
        body: JSON.stringify({
            message: `Hello ${event.body?.queryStringParameters?.name} from Funcie!`,
        }),
    };
};


beginReceiving(config, handler).then((server) => {
    server.on('close', () => {
        done = true;
    });
}).catch((err) => {
    console.log(err.message);
    exit(1);
});
