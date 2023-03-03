# Lambda URL Sample

This is a very simple example that uses Lambda Function URLs to return a JSON response.

To run this example, you can create a new Lambda function that is configured to support
function URLs. Then, build this example using `go build`, zip it, and upload
the zip file to the Lambda function you created.

You will also need to change the `FUNCIE_REDIS_ADDR` environment variable to point to the
Redis instance you are using, such as `my.remote.host:6379`. This environment variable
will need to be set both on your Lambda, and your local instance.
Remember that both the Lambda
and your local environment need to be able to connect to the Redis instance.

When running this locally, you may wish to use `gow` to watch for changes and rebuild
the binary automatically. You can install `gow` using `go get github.com/mitranim/gow`.

For example, to run this locally and automatically restart on any change, you can use:
```bash
FUNCIE_REDIS_ADDR=10.1.1.1:6379 gow run main.go
```

Then you can simply hit the URL, apply some changes, and see the next invocation automatically
pick up the changes.

If you close the local instance, the Lambda will simply respond using the last code you uploaded
to the Lambda. Then when you start the local instance, it will automatically continue forwarding
requests.