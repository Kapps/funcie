# Lambda URL Sample

This is a very simple example that uses Lambda Function URLs to return a JSON response.
This sample uses the bastion deployed from the `deploy` folder, referencing the bastion IP by passing in the LB.

## Prerequisites

- To run this sample, you will need to have deployed the terraform in the `deploy` folder.
- You will also need to have the bastion docker image running locally.
  - You can run the bastion docker image with the following command:
    ```bash
    docker run -e FUNCIE_REDIS_ADDRESS=<elasticachehost>:6379 -e FUNCIE_LISTEN_ADDRESS=0.0.0.0:24193 -p 24193:24193 public.ecr.aws/w1h1o7p8/funcie-client-bastion:0.3.3
    ```
- You will need to be able to access the Redis host from your local machine.
  - One approach is using an OpenVPN Access Server connection, or following the AWS guide for external connections.

## Deployment

- Create a `terraform.tfvars` file that contains values for the variables (according to `sample.tfvars`)
- Build the lambda by running `make`
- Deploy the Lambda function via `terraform apply`
- Note the lambda_url output from the terraform apply -- this is the URL you will use to access the Lambda function.

## Running Locally

When running this locally, you may wish to use `gow` to watch for changes and rebuild
the binary automatically. See https://github.com/mitranim/gow for details on `gow`.`

For example, to run this locally and automatically restart on any change, you can use:

```bash
FUNCIE_APPLICATION_ID=url gow run main.go
```

Then you can simply hit the URL, apply some changes, and see the next invocation automatically
pick up the changes.

If you close the local instance, the Lambda will simply respond using the last code you uploaded
to the Lambda. Then when you start the local instance, it will automatically continue forwarding
requests.

## Deployment Variables

When deploying this Lambda, you will need to pass the following Terraform variables in your terraform.tfvars file:

- subnet_ids: List of (private) subnet IDs for the Lambda function
- security_group_ids: List of security group IDs for the Lambda function
- redis_host: Address of the Redis host, including port.
- bastion_lb_arn: The ARN for the load balancer created to be used for the bastion (can be viewed in the outputs).
- region: AWS region to deploy to.
