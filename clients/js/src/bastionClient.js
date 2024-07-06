const { Response } = require('./models');
const axios = require('axios');
const url = require('node:url');
const { error, info } = require('./utils');

const sendMessage = async (baseUrl, message) => {
    const endpoint = new url.URL('/dispatch', baseUrl);
    info(`Sending message to ${endpoint}`);
    const httpResponse = await axios.post(endpoint, message);
    if (httpResponse.status !== 200) {
        error(`unexpected status code: ${httpResponse.status}`);
        error(`response body: ${httpResponse.data}`);
        throw new Error(`unexpected status code: ${httpResponse.status}`);
    }

    const responseData = httpResponse.data;
    return Response.fromObject(responseData);
};

module.exports = {
    sendMessage,
};
