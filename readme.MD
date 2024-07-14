# Funcie

Funcie simplifies local development for AWS Lambda and other serverless functions, enabling direct local execution of your code.
Inspired by SST's [Live Lambda Development](https://docs.sst.dev/live-lambda-development), Funcie takes this further and ensures seamless transition between cloud and local execution.

In typical use, your functions execute entirely within your cloud provider. Then, as soon as you launch or debug your application locally, Funcie begins redirecting incoming requests to the local instance. Once you're done, closing the local application immediately resumes execution on the cloud function.

Funcie doesn't care how your function is invoked, what your development environment is, or what frameworks you're using.
This means you can trigger your Lambda via S3 triggers, Step Functions, or however else you decide, while debugging locally using VS Code, GoLand, or anything else.

Currently, Funcie supports AWS Lambdas running Go or JavaScript/TypeScript, but is designed such that adding new languages requires only minimal effort.
It's also designed to be very simple to plug in, just replacing `lambda.Start` in Go with `funcie.Start`, or wrapping your JS handler with our wrapper function.

## Getting Started

To get started, we recommend using the Funcie CLI. The CLI allows you to easily deploy the necessary infrastructure via terraform, and securely connects your local instance to your VPC, facilitating access to VPC resources like RDS without needing a VPN.

The installation process takes approximately 10 minutes when using the CLI.

### Prerequisites

- Go 1.21 or later
- AWS CLI configured
- Terraform installed

### Deploying the Infrastructure

**Warning:** This will deploy a Redis instance and a bastion server to your AWS account, which will incur costs.
Make sure to run `funcie destroy` if you're no longer using Funcie in order to clean up the resources. We estimate the costs at roughly $20-25 USD per month, but it can vary depending on region and other factors.

1. Clone the repository:
    
    ```bash
    git clone git@github.com:Kapps/funcie.git
    ```
   
2. Install the CLI:

    ```bash
    cd funcie/cmd/cli
    go install
    ```
   
3. Run `funcie init` and follow the prompts to deploy the necessary infrastructure to your AWS account.

    ```bash
    funcie init
    ```
    If you encounter a file not found error, ensure your `PATH` environment variable is set for `go install`. Alternatively, run the CLI from the `cli` folder.

4. Now that funcie is deployed, you can run `funcie connect` to allow your local bastion to access the ElastiCache instance.
    
    ```bash
    funcie connect
    ```

5. You can now run your Lambda locally and have it connect to the bastion server. Check out one of the examples to play with it.

### Updating Your Code

To use funcie, you need to update your Lambda handler to use the Funcie library for your language.
The library handles communication between your Lambda and the local environment.
Libraries are provided for AWS Lambda in Go and JavaScript/TypeScript.

Alternatively, use an example from the `examples` directory.

#### Go

1. Install the funcie library:
    ```bash
    go get github.com/Kapps/funcie
    ```

2. Update your Lambda handler to use the funcie library.
    
    Replace:
    ```go
    import "github.com/aws/aws-lambda-go/lambda"

    func main() {
        lambda.Start(handler)
    }
    ```

    with:

    ```go
    import "github.com/Kapps/funcie/clients/go/funcietunnel"

    func main() {
        funcietunnel.Start("my-app", handler)
    }
    ```

3. Run your Lambda locally. You should see the logs in your terminal.

#### JavaScript/TypeScript

1. Install the funcie library:
    ```bash
    npm install @funcie/client
    ```

   
2. Update your Lambda handler to use the funcie library.

    Replace:
    
    ```javascript
    export const handler = async (event, context) => {
        // Your Lambda handler code here
    };
    ```
   
    with:
    ```javascript
    const { lambdaWrapper } = require('@funcie/client');
    export const handler = lambdaWrapper("js-url", async (event, context) => {
        // Your Lambda handler code here
    });
    ```    

3. Run your Lambda locally. You should see the logs in your terminal.

## Accessing VPC Resources

If you have a VPN, often accessing resources like an RDS is trivial. If you don't though, the Funcie CLI can help here via the `funcie connect` command.

For instance, if you have an RDS at `abc.id.us-east-1.rds.amazonaws.com`, which accepts connections from within the VPC, set up a tunnel using:

```
funcie connect -r abc.id.us-east-1.rds.amazonaws.com -p 5432 -l 5432
```

Now, localhost on port 5432 will tunnel to your RDS. This uses SSM port forwarding to securely create the tunnel.

## Security

Funcie is intended for use in development and staging environments, and is not recommended to enable in production.
While the default deployment is protected by being behind a VPC, it's recommended to detect if in a production environment and disable funcie.

For example, you could use an environment variable to decide whether to run `lambda.Start` or `funcie.Start`.

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
