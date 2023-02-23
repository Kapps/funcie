# Funcie

Funcie is a framework for local development of AWS Lambda or other serverless functions. 
It is heavily inspired by the Live Lambda Development from SST (https://docs.sst.dev/live-lambda-development).

It works by creating a tunnel where messages can pass through between
a local environment and a Lambda. Then, when a Lambda is invoked, it
checks if there is a local function connected to that tunnel, and if so,
it forwards the request to that local function. Once the local function
returns, the response is returned to the Lambda which forwards it back
to the cloud provider.

This means that you can use Step Functions, triggers, and any other approach
for invoking the Lambda, as the invocation is separated from the implementation.
In addition, local changes are immediately applied without needing to
republish anything.

## Backends

Backends are the mechanism for sending messages through a tunnel.
Various backends will exist to allow for flexibility in how to send
messages back and forth.

### Redis
The first backend that will be supported is a Redis instance.
The sender will use Redis as a queue, while the receiver will either
use long polling or eventing to receive messages.

## Consumer State

POTENTIALLY. This may be unnecessary as the backend should have a way
of understanding if a message is delivered.

Because the system is meant to have standard execution within the Lambda
itself until a local function is active, we need to store the consumer
state of each active consumer. If a client disconnects or crashes,
the Lambda should fail any outstanding requests and handle all future
requests itself.

In addition, a client reconnecting should only process requests that
come in _after_ the client reconnects, not retroactively process requests
from previous attempts.
