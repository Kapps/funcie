const { Response } = require('./models');
const axios = require('axios');
const url = require('node:url');

const sendMessage = async (baseUrl, message) => {
    const endpoint = new url.URL('/dispatch', baseUrl);
    const httpResponse = await axios.post(endpoint, message);
    if (httpResponse.status !== 200) {
        console.log(`unexpected status code: ${httpResponse.status}`);
        console.log(`response body: ${httpResponse.data}`);
        throw new Error(`unexpected status code: ${httpResponse.status}`);
    }

    const responseData = httpResponse.data;
    const response = Response.fromJson(responseData);

    return response;
};

module.exports = {
    sendMessage,
};
