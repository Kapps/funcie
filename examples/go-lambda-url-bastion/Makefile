compile:
	GOOS=linux GOARCH=amd64 go build -o bin/bootstrap main.go

deploy:
	# Compile the Go code
	make compile
	# Deploy the Lambda function via terraform
	terraform init
	terraform apply -auto-approve

