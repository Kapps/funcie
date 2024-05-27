#!/bin/bash

# Get a session token from IMDSv2
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")

# Get the instance ID using the token
INSTANCE_ID=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/instance-id)
PRIVATE_IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/local-ipv4)

echo ECS_CLUSTER=${ECS_CLUSTER} >> /etc/ecs/ecs.config

aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_host --value $PRIVATE_IP --type String --overwrite --region ${REGION}
aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_instance_id --value $INSTANCE_ID --type String --overwrite --region ${REGION}

sudo yum install ec2-instance-connect
