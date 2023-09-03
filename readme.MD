# Funcie

Funcie is a framework for easier local development of AWS Lambda or other serverless functions.
It is inspired by the Live Lambda Development from SST (https://docs.sst.dev/live-lambda-development).

Funcie is designed to seamlessly transition between cloud and local execution of your serverless functions.
For ordinary usage, your cloud functions execute directly, handling any requests with your published code.
Then, when you want to debug or try a new local version, simply launch your application locally and funcie
will automatically recognize that a client endpoint is active, forwarding all incoming requests to that client.
Once you're done your local development / debugging, just close the local application and your cloud function
will resume handling the request directly once more.

Funcie does not care how your function is invoked or what your development environment is.
This means you can trigger your Lambda via S3 triggers, Step Functions, or however else you decide, while debugging locally using VS Code, Goland, or whatever editor/IDE you prefer.

Currently, funcie supports Go and JavaScript/TypeScript, but is designed such that adding new languages requires only minimal effort.
At this time, the only provided wrappers are for AWS Lambda, but there is nothing that inherently ties funcie to Lambda.

## High Level Overview

It works by creating a tunnel where messages can pass through between
a local environment and a Lambda. When a Lambda is invoked, it
checks if there is a local function connected to that tunnel, and if so,
forwards the request to that local function handler. Once the local function returns, the response is returned to the Lambda which sends it back to the cloud provider.

Funcie provides a library that wraps the original Lambda handler or `lambda.Start` call. In addition, a bastion process runs on both the server and on the client, which contains most of the logic for forwarding requests and connection management. On the server side, the funcie wrapper sends each request that comes in to the server bastion, via an HTTP call. The server bastion sends the request to the client bastion through the tunnel, and finally, client bastion sends the request to a small HTTP server the client library starts, which then invokes your Lambda handler.

This means that you can use Step Functions, triggers, or any other method of invoking the Lambda, as the invocation is separated from the implementation. In addition, local changes are immediately applied without needing to republish anything, and you can debug and step through your code just like you would any other local process.

When there is no local function connected, the Lambda will simply execute
as normal, and the local function will not be invoked. Once a local function is connected, the Lambda will forward all requests to the local function until the local function disconnects.

Currently, even when no local function is connected, each Lambda invocation
request has an overhead of a single request to the server Bastion, which makes a Redis call. In general, both to avoid this overhead and to minimize surface area for security reasons, it is recommended to use funcie only in development or staging environments. Running funcie in a production environment is not currently recommended.

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
