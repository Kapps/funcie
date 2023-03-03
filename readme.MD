# Funcie

Funcie is a framework for easier local development of AWS Lambda or other serverless functions. 
It is heavily inspired by the Live Lambda Development from SST (https://docs.sst.dev/live-lambda-development).

It works by creating a tunnel where messages can pass through between
a local environment and a Lambda. Then, when a Lambda is invoked, it
checks if there is a local function connected to that tunnel, and if so,
forwards the request to that local function handler. Once the local function
returns, the response is returned to the Lambda which sends it back
to the cloud provider.

This means that you can use Step Functions, triggers, and any other method of 
invoking the Lambda, as the invocation is separated from the implementation.
In addition, local changes are immediately applied without needing to
republish anything, and you can debug and step through your code as normal.

When there is no local function connected, the Lambda will simply execute
as normal, and the local function will not be invoked. Once a local function
is connected, the Lambda will forward all requests to the local function
until the local function disconnects.

Currently only Go is supported.

## Backends

Backends are the mechanism for sending messages through a tunnel.
Various backends will exist to allow for flexibility in how to send
messages back and forth.

### Redis
Currently the only supported backend is Redis.

The Redis backend uses a Redis instance as a queue to send messages
between the Lambda and the local environment using the PUBSUB mechanism.

Responses are stored as Redis keys, and the Lambda will wait for a
response to be available before returning.

## Consumer State

Because the system is meant to have standard execution within the Lambda
itself until a local function is active, each backend should know whether
a message was able to reach any consumer.

If a consumer disconnects or crashes, the Lambda should fail any outstanding
requests and handle all future requests itself.

This means that the default state of the Lambda should be to handle all
requests itself, and only forward requests to a local function if it is
available.

In addition, a client reconnecting should only process requests that
come in _after_ the client reconnects, not retroactively process requests
from previous attempts.
