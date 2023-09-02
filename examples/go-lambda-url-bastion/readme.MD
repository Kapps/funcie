# Lambda URL Sample

This is a very simple example that uses Lambda Function URLs to return a JSON response.
This sample uses the bastion deployed from the `deploy` folder, referencing the bastion IP by passing in the LB.

To run this sample, create a `terraform.tfvars` file that contains values for the variables (according to `sample.tfvars`).
Then, build the lambda by typing `make`.
Finally, deploy it via `terraform apply`.

When running this locally, you may wish to use `gow` to watch for changes and rebuild
the binary automatically. See https://github.com/mitranim/gow for details on `gow`.`

For example, to run this locally and automatically restart on any change, you can use:

```bash
FUNCIE_APPLICATION_ID=url FUNCIE_CLIENT_BASTION_ENDPOINT=http://localhost:8081 gow run main.go
```

Then you can simply hit the URL, apply some changes, and see the next invocation automatically
pick up the changes.

If you close the local instance, the Lambda will simply respond using the last code you uploaded
to the Lambda. Then when you start the local instance, it will automatically continue forwarding
requests.
