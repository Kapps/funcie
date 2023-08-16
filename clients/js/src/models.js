const uuid = require('uuid');

class Message {
    constructor(kind, application, payload) {
        this.id = uuid.v4();
        this.kind = kind;
        this.application = application;
        this.payload = payload;
        this.created = new Date();
    }

    static fromJson(json) {
        const data = JSON.parse(json);
        return new Message(data.kind, data.application, data.payload);
    }

    toString() {
        const marshaledPayload = JSON.stringify(this.payload);
        return `Message{ID: ${this.id}, Kind: ${this.kind}, Application: ${this.application}, Created: ${this.created.toISOString()}, Payload: ${marshaledPayload}}`;
    }
}

class Response {
    constructor(id, data, error, received) {
        this.id = id;
        this.data = data;
        this.error = error;
        this.received = received;
    }

    static fromJson(json) {
        const data = JSON.parse(json);
        return new Response(data.id, data.data, data.error, new Date(data.received));
    }

    toString() {
        const marshaledData = JSON.stringify(this.data);
        return `Response{ID: ${this.id}, Data: ${marshaledData}, Error: ${this.error}, Received: ${this.received.toISOString()}}`;
    }
}

module.exports = {
    Message,
    Response,
};
