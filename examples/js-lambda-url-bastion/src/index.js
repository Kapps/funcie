const { lambdaWrapper } = require('funcie-tunnel');

exports.handler = lambdaWrapper(async (event) => {
    return {
        statusCode: 200,
        body: 'Hello, world!',
        headers: {
            'Content-Type': 'text/plain',
        },
    };
});
