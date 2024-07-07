import { Response } from './models.js';
import { error, info } from './utils.js';
import axios from 'axios';
import { URL } from 'url';

export const sendMessage = async (baseUrl, message) => {
    const endpoint = new URL('/dispatch', baseUrl);
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
