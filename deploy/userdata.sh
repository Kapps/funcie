#!/bin/bash

# Get a session token from IMDSv2
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")

# Get the instance ID using the token
INSTANCE_ID=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/instance-id)

echo ECS_CLUSTER=${ECS_CLUSTER} >> /etc/ecs/ecs.config

# SSH tunneling to Redis

# Associate the EIP with this instance
EIP_ALLOCATION_ID="${EIP_ALLOCATION_ID}"
REGION="${REGION}"

# Retry logic for EIP association
MAX_RETRIES=5
RETRY_COUNT=0
SUCCESS=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if aws ec2 associate-address --instance-id $${INSTANCE_ID} --allocation-id $${EIP_ALLOCATION_ID} --region $${REGION}; then
    SUCCESS=1
    break
  else
    RETRY_COUNT=$((RETRY_COUNT + 1))
    sleep 10
  fi
done

if [ $SUCCESS -ne 1 ]; then
  echo "Failed to associate EIP after $MAX_RETRIES attempts" >&2
fi

echo "EIP associated with instance"
