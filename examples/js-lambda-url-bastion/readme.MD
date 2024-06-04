# Lambda URL Sample

This is a very simple JavaScript example that uses Lambda Function URLs to return a JSON response.
This sample requires funcie to be deployed via the deploy folder in the root of the repository.

## PREREQUISITES

- To run this sample, you will need to have deployed the terraform in the `deploy` folder.

To run this sample, create a `terraform.tfvars` file that contains values for the variables (according to `sample.tfvars`).
Then deploy the Lambda using `terraform apply`.

When running this locally, you may wish to use `node --watch` to watch for changes and rebuild automatically.

For example, to run this locally and automatically restart on any change, you can use:

```bash
FUNCIE_APPLICATION_ID=url node --watch index.js
```

Then you can simply hit the URL, apply some changes, and see the next invocation automatically
pick up the changes.

If you close the local instance, the Lambda will simply respond using the last code you uploaded
to the Lambda. Then when you start the local instance, it will automatically continue forwarding
requests.

## Deployment Variables

When deploying this Lambda, you will need to pass the following Terraform variables in your terraform.tfvars file:

- subnet_ids: List of (private) subnet IDs for the Lambda function. Example: ["subnet-12345678", "subnet-87654321"]
- redis_host: Address of the Redis host, including port. Example: "abc.cache.amazonaws.com:6379"
- region: AWS region to deploy to. Example: "us-east-1"
- vpc_id: VPC ID for the Lambda function. Example: "vpc-12345678"
- bastion_host: Host name for the bastion ECS. Usually the bastion_host output from the deployment file. Example: "server-bastion.funcie.local"
