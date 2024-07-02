#!/bin/bash

set -euxo pipefail

# Get a session token from IMDSv2
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")

# Get the instance ID using the token
INSTANCE_ID=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/instance-id)
PRIVATE_IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/local-ipv4)

yum install -y https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_amd64/amazon-ssm-agent.rpm
systemctl start amazon-ssm-agent

echo Setting up ECS agent with cluster ${ECS_CLUSTER}

mkdir -p /etc/ecs/
echo "ECS_CLUSTER=${ECS_CLUSTER}" >> /etc/ecs/ecs.config

yum install -y ecs-init

aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_host --value $PRIVATE_IP --type String --overwrite --region ${REGION}
aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_instance_id --value $INSTANCE_ID --type String --overwrite --region ${REGION}
echo "Updated SSM parameters for bastion host $PRIVATE_IP and instance ID $INSTANCE_ID"

yum install -y ec2-instance-connect

if [ "${CREATE_VPC}" = "true" ]; then
    # Configure iptables for NAT
    PRIMARY_INTERFACE=$(ip -o -4 addr show | awk '{print $2 " " $4}' | grep "$PRIVATE_IP" | awk '{print $1}')

    echo Enabling IP forwarding through $PRIMARY_INTERFACE

    sysctl -w net.ipv4.ip_forward=1
    echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf

    # Install iptables-services and save configuration
    yum install -y iptables-services

    /sbin/iptables -t nat -A POSTROUTING -o $PRIMARY_INTERFACE -j MASQUERADE
    iptables-save > /etc/sysconfig/iptables

    systemctl start iptables
    systemctl enable iptables

    # Disable source/destination check
    aws ec2 modify-instance-attribute --instance-id $INSTANCE_ID --no-source-dest-check --region ${REGION}

    PRIMARY_ENI_ID=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --region ${REGION} --query "Reservations[0].Instances[0].NetworkInterfaces[?Attachment.DeviceIndex==\`0\`].NetworkInterfaceId" --output text)

    # Update the route table to set this instance as the active NAT instance
    OLD_NAT_ENI_ID=$(aws ec2 describe-route-tables --route-table-ids ${ROUTE_TABLE_ID} --region ${REGION} --query "RouteTables[].Routes[?DestinationCidrBlock=='0.0.0.0/0'].NetworkInterfaceId" --output text)

    if [ ! -z "$OLD_NAT_ENI_ID" ]; then
        aws ec2 replace-route --route-table-id ${ROUTE_TABLE_ID} --destination-cidr-block 0.0.0.0/0 --network-interface-id $PRIMARY_ENI_ID --region ${REGION}
        echo "Replaced route for old NAT instance $OLD_NAT_ENI_ID to $PRIMARY_ENI_ID"
    else
        aws ec2 create-route --route-table-id ${ROUTE_TABLE_ID} --destination-cidr-block 0.0.0.0/0 --network-interface-id $PRIMARY_ENI_ID --region ${REGION}
        echo "Created route for new NAT interface $PRIMARY_ENI_ID"
    fi
fi

# Has to happen after NAT is configured to not mess with the iptables rules.
systemctl enable --now --no-block ecs.service

echo Completed userdata script
