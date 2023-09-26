# Funcie Client

This package provides a wrapper for Lambda functions to utilize [Funcie](https://github.com/Kapps/funcie).

Funcie is a framework for easier local development of Lambda functions. It allows you to run your Lambda functions locally, and automatically forward requests to the Lambda function when it is deployed.

## Usage

First, you will need to set up funcie. See the instructions in the readme at the root of this repository.

Next, you will need to install the funcie client package:

```bash
npm install --save @funcie/client
```

Then, simply import it and wrap your existing Lambda handler:

```javascript
const { lambdaWrapper } = require("@funcie/client");

exports.handler = lambdaWrapper(async (event) => {
  // Your existing Lambda handler code here
});
```

Now if you run your Lambda locally, the deployed Lambda will automatically invoke your local code when hit.
