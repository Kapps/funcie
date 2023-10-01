const { lambdaWrapper } = require('@funcie/client');

exports.handler = lambdaWrapper(async (event) => {
    if (event.queryStringParameters && event.queryStringParameters.name) {
        if (event.queryStringParameters.name === 'error') {
            throw new Error('error being forwarded');
        }
        if (event.queryStringParameters.name === 'sleep') {
            await new Promise((resolve) => setTimeout(resolve, 10000));
        }
        return {
            statusCode: 200,
            body: `Hello, ${event.queryStringParameters.name}!`,
            headers: {
                'Content-Type': 'text/plain',
            },
        };
    }
    return {
        statusCode: 200,
        body: 'Hello, world!',
        headers: {
            'Content-Type': 'text/plain',
        },
    };
});
