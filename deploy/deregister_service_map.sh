#!/bin/bash

[ $# -ne 1 ] && echo "Usage: $0 <service-id>" && exit 1

SERVICE_ID="--service-id=$1"

echo "Draining servicediscovery instances from $1 ..."
INSTANCE_IDS="$(aws servicediscovery list-instances $SERVICE_ID --query 'Instances[].Id' --output text | tr '\t' ' ')"

FOUND=
for INSTANCE_ID in $INSTANCE_IDS; do
  if [ -n "$INSTANCE_ID" ]; then
    echo "Deregistering $1 / $INSTANCE_ID ..."
    aws servicediscovery deregister-instance $SERVICE_ID --instance-id "$INSTANCE_ID" --region "$REGION"
    FOUND=1
  fi
done

[ -n "$FOUND" ] && sleep 5 || true
