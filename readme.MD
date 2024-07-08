# Funcie

Funcie is a framework for easier local development of AWS Lambda or other serverless functions.
It is inspired by the Live Lambda Development from SST (https://docs.sst.dev/live-lambda-development).

Funcie is designed to seamlessly transition between cloud and local execution of your serverless functions.
Under ordinary usage, your cloud functions execute directly in your cloud provider, handling any requests with your published code.
Then, when you want to debug or try a new local version, simply launch your application locally and funcie
will instantly recognize that a client endpoint is active, forwarding all incoming requests to that client.
Once you're done your local development / debugging, just close the local application and your cloud function
will immediately resume handling the request directly once more.

Funcie doesn't care how your function is invoked, what your development environment is, or what framework you're using.
This means you can trigger your Lambda via S3 triggers, Step Functions, or however else you decide, while debugging locally using VS Code, Goland, or whatever editor/IDE you prefer.

Currently, funcie supports Go and JavaScript/TypeScript, but is designed such that adding new languages requires only minimal effort.
At this time, the only provided wrappers are for AWS Lambda, but there is nothing that inherently ties funcie to Lambda.

## Getting Started

To get started, we recommend using the CLI. The CLI allows you to easily deploy the necessary infrastructure via terraform,
and also provides a simple and secure way for your local instance to connect to your AWS VPC.

The installation process is split up into two parts: deploying the infrastructure and updating your code to use the funcie library.

### Deploying the Infrastructure

**Warning:** This will deploy a Redis instance and a bastion server to your AWS account, which will incur costs.
Make sure to run `funcie destroy` if you're no longer using funcie to clean up the resources.

The primary costs come from a t3.micro EC2 instance, and a t4g.micro Elasticache instance.

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
    
    You can also deploy the infrastructure manually using the terraform files at https://github.com/Kapps/terraform-aws-funcie

    If you receive a file not found error, you may need to add the go bin directory to your PATH. You can do this by adding the following line to your shell profile:
   
    ```bash
    export PATH=$PATH:$(go env GOPATH)/bin
    ```

4. Now that funcie is deployed, you can run `funcie connect` to allow your local bastion to access the Elasticache instance.
    
    ```bash
    funcie connect
    ```

5. You can now run your Lambda locally and have it connect to the bastion server.

### Updating Your Code

To use funcie, you need to update your Lambda handler to use the funcie library.
The library will handle the communication between your Lambda and the local environment.
Libraries are provided for AWS Lambda in Go and JavaScript/TypeScript.

Alternatively you can use one of the examples in the `examples` directory as a starting point.

1. Go to the directory where your Lambda handler is located.

2. Install the funcie library for your language:

    Go:
    
    ```bash
    go get github.com/Kapps/funcie
    ```
    
    JavaScript/TypeScript:
    
    ```bash
    npm install @funcie/client
    ```
   
3. Update your Lambda handler to use the funcie library.

    For Go, replace:
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
   
    For JavaScript/TypeScript, replace:
    
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

4. Run your Lambda locally.

## Security

Funcie is intended for use in development and staging environments, and is not recommended to enable in production.
While the default deployment is protected by being behind a VPC, it's recommended to detect if in a production environment and disable funcie.

For example, you could use an environment variable to decide to run `lambda.Start` or `funcietunnel.Start`.

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
