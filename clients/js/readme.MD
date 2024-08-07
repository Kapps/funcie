# Funcie Client

This package provides a wrapper for Lambda functions to utilize [Funcie](https://github.com/Kapps/funcie).

Funcie is a framework for easier local development of Lambda functions. It allows you to run your Lambda functions locally, and automatically forward requests to the Lambda function when it is deployed.

## Usage

First, you will need to set up funcie. See the instructions in the readme at the root of this repository.

Next, you will need to install the funcie client package:

```bash
npm install --save @funcie/client
```

Then, simply import it and wrap your existing Lambda handler along-side a user-defined application name:

```javascript
const { lambdaWrapper } = require("@funcie/client");

exports.handler = lambdaWrapper('app-name', async (event) => {
  // Your existing Lambda handler code here
});
```

Now if you run your Lambda locally, the deployed Lambda will automatically invoke your local code when hit.

You may find it useful to use the --watch flag when running your Lambda locally, so that it will automatically reload when you make changes:

```bash
node --watch src/index.js
```

## Configuration

The funcie client can be configured via environment variables:

- `FUNCIE_CLIENT_BASTION_ENDPOINT` (optional): The endpoint for the client bastion -- defaults to `http://localhost:27193`.
- `FUNCIE_SERVER_BASTION_ENDPOINT`: The endpoint for the server bastion (only used when not local).
- `FUNCIE_LISTEN_ADDRESS` (optional): The address to listen on -- defaults to `0.0.0.0` on a random port).

## Troubleshooting

By default, most errors will be logged to the console, with other logging disabled.
Additional debug logging can be enabled via the `FUNCIE_DEBUG` environment variable:

```bash
FUNCIE_DEBUG=true node index.js
```
