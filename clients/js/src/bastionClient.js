const { Response } = require('./models');
const axios = require('axios');
const url = require('node:url');

const sendMessage = async (endpoint, message) => {
    const httpResponse = await axios.post(new url.URL('/dispatch', endpoint), message);
    if (httpResponse.status !== 200) {
        console.log(`unexpected status code: ${httpResponse.status}`);
        console.log(`response body: ${httpResponse.data}`);
        throw new Error(`unexpected status code: ${httpResponse.status}`);
    }

    const response = httpResponse.data;

    return response;
};

module.exports = {
    sendMessage,
};
