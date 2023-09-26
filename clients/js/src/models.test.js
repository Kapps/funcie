const { Message, Response } = require('./models');

describe('Message', () => {
    it('should create a message with the correct properties', () => {
        const kind = 'test';
        const application = 'example';
        const payload = { foo: 'bar' };

        const message = new Message(kind, application, payload);

        expect(message.kind).toBe(kind);
        expect(message.application).toBe(application);
        expect(message.payload).toEqual(payload);
        expect(message.created).toBeInstanceOf(Date);
    });

    it('should parse a message from JSON', () => {
        const json = '{"kind":"test","application":"example","payload":{"foo":"bar"}}';

        const message = Message.fromJson(json);

        expect(message.kind).toBe('test');
        expect(message.application).toBe('example');
        expect(message.payload).toEqual({ foo: 'bar' });
        expect(message.created).toBeInstanceOf(Date);
    });

    it('should convert a message to a string', () => {
        const kind = 'test';
        const application = 'example';
        const payload = { foo: 'bar' };

        const message = new Message(kind, application, payload);

        const expectedString = `Message{ID: ${message.id}, Kind: ${kind}, Application: ${application}, Created: ${message.created.toISOString()}, Payload: ${JSON.stringify(payload)}}`;

        expect(message.toString()).toBe(expectedString);
    });
});

describe('Response', () => {
    it('should create a response with the correct properties', () => {
        const id = '123';
        const data = { baz: 'qux' };
        const received = new Date();

        const response = new Response(id, data, null, received);

        expect(response.id).toBe(id);
        expect(response.data).toEqual(data);
        expect(response.error).toBeNull();
        expect(response.received).toBe(received);
    });

    it('should create a response with an error', () => {
        const id = '123';
        const error = { message: 'something went wrong' };
        const received = new Date();

        const response = new Response(id, null, error, received);

        expect(response.id).toBe(id);
        expect(response.data).toBeNull();
        expect(response.error).toBe(error);
        expect(response.received).toBe(received);
    });

    it('should parse a response from JSON', () => {
        const json = '{"id":"123","data":{"baz":"qux"},"error":null,"received":"2022-01-01T00:00:00.000Z"}';

        const response = Response.fromJson(json);

        expect(response.id).toBe('123');
        expect(response.data).toEqual({ baz: 'qux' });
        expect(response.error).toBeNull();
        expect(response.received).toBeInstanceOf(Date);
    });

    it('should parse a response with an error from JSON', () => {
        const json = '{"id":"123","data":null,"error":{"message":"something went wrong"},"received":"2022-01-01T00:00:00.000Z"}';

        const response = Response.fromJson(json);

        expect(response.id).toBe('123');
        expect(response.data).toBeNull();
        expect(response.error).toEqual({ message: 'something went wrong' });
        expect(response.received).toBeInstanceOf(Date);
    });

    it('should convert a response to a string', () => {
        const id = '123';
        const data = { baz: 'qux' };
        const received = new Date();

        const response = new Response(id, data, null, received);

        const expectedString = `Response{ID: ${id}, Data: ${JSON.stringify(data)}, Error: undefined, Received: ${received.toISOString()}}`;

        expect(response.toString()).toBe(expectedString);
    });

    it('should convert a response with an error to a string', () => {
        const id = '123';
        const error = { message: 'something went wrong' };
        const received = new Date();

        const response = new Response(id, null, error, received);

        const expectedString = `Response{ID: ${id}, Data: null, Error: something went wrong, Received: ${received.toISOString()}}`;

        expect(response.toString()).toBe(expectedString);
    });
});
