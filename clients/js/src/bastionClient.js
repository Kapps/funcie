const { Response } = require('./models');
const axios = require('axios');
const url = require('node:url');

const sendMessage = async (config, message) => {
    const httpResponse = await axios.post(new url.URL('/dispatch', config.ClientBastionEndpoint), message);
    
    const response = httpResponse.data;
    if (response.error) {
        throw new Error(response.error);
    }

    return response;
};

module.exports = {
    sendMessage,
};
