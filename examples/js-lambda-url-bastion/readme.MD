# Lambda URL Sample

This is a very simple example that uses Lambda Function URLs to return a JSON response. This is a good starting point for building more complex Lambda applications.

## Prerequisites

To run this sample, you will need the following:

- Active Funcie deployment (see the [Funcie README](https://github.com/Kapps/funcie/blob/main/readme.MD) for instructions).
- Client bastion running locally (automatically launched by the CLI if you enabled it).
- Access to the Redis host (run `funcie connect` to open a tunnel if you used the CLI).

## Deployment

1. Create a `terraform.tfvars` file with values for the variables (refer to `sample.tfvars`).

2. Build and deploy the Lambda:

    ```bash
    make deploy
    ```

3. Note the `lambda_url` output from `terraform apply`. This is the URL to access the Lambda function.

## Running Locally

To run this locally, simply run node on the index file (preferably with --watch to reload on changes).

```bash
node --watch src/index.js
```

Open the URL to see requests in your terminal.
If you close the local instance, the Lambda will respond using the last uploaded code.
Restarting the local instance resumes request forwarding.

When you're done, don't forget to clean up by running:

```bash
make destroy
```

## Deployment Variables

When deploying this Lambda, you will need to pass the following Terraform variables in your terraform.tfvars file. These are all also available from the outputs of the funcie deployment terraform module.

- subnet_ids: List of (private) subnet IDs for the Lambda function
- redis_host: Address of the Redis host, including port.
- region: AWS region to deploy to.
