import { sendMessage } from './bastionClient.js';
import { Response } from './models.js';
import { URL } from 'url';

import { jest } from '@jest/globals';

jest.unstable_mockModule('axios', () => ({
    post: jest.fn(),
}));

import axios from 'axios';

describe('bastionClient', () => {
    const baseUrl = 'http://localhost:3000';

    afterEach(() => {
        jest.resetAllMocks();
    });

    describe('sendMessage', () => {
        it('should send a message and return a response', async () => {
            const message = { type: 'test', payload: { foo: 'bar' } };
            const expectedResponse = new Response('success', { baz: 'qux' });

            axios.post.mockResolvedValueOnce({ status: 200, data: expectedResponse });

            const response = await sendMessage(baseUrl, message);

            expect(axios.post).toHaveBeenCalledTimes(1);
            expect(axios.post).toHaveBeenCalledWith(new URL('/dispatch', baseUrl), message);
            expect(response).toEqual(expectedResponse);
        });

        it('should throw an error if the response status code is not 200', async () => {
            const message = { type: 'test', payload: { foo: 'bar' } };
            const expectedResponse = new Response('error', { message: 'something went wrong' });

            axios.post.mockResolvedValueOnce({ status: 500, data: expectedResponse });

            await expect(sendMessage(baseUrl, message)).rejects.toThrow('unexpected status code: 500');
        });
    });
});
