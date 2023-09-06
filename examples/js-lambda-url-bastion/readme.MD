# Lambda URL Sample

This is a very simple JavaScript example that uses Lambda Function URLs to return a JSON response.
This sample uses the bastion deployed from the `deploy` folder, referencing the bastion IP by passing in the LB.

To run this sample, create a `terraform.tfvars` file that contains values for the variables (according to `sample.tfvars`).
Then deploy the Lambda using `terraform apply`.

When running this locally, you may wish to use `node --watch` to watch for changes and rebuild automatically.

For example, to run this locally and automatically restart on any change, you can use:

```bash
FUNCIE_APPLICATION_ID=url FUNCIE_CLIENT_BASTION_ENDPOINT=http://localhost:8081 node --watch index.js
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
