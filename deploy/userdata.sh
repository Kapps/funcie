#!/bin/bash

# Get a session token from IMDSv2
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")

# Get the instance ID using the token
INSTANCE_ID=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/instance-id)
PRIVATE_IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/local-ipv4)

PRIMARY_INTERFACE=$(ip -o -4 addr show | awk '{print $2 " " $4}' | grep "$PRIVATE_IP" | awk '{print $1}')

echo ECS_CLUSTER=${ECS_CLUSTER} >> /etc/ecs/ecs.config

aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_host --value $PRIVATE_IP --type String --overwrite --region ${REGION}
aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_instance_id --value $INSTANCE_ID --type String --overwrite --region ${REGION}

sudo yum install -y ec2-instance-connect

if [ "${CREATE_VPC}" = "true" ]; then
    # Enable IP forwarding
    echo Enabling IP forwarding
    sysctl -w net.ipv4.ip_forward=1

    # Configure iptables for NAT
    /sbin/iptables -t nat -A POSTROUTING -o $PRIMARY_INTERFACE -j MASQUERADE
    echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf

    # Install iptables-services and save configuration
    yum install -y iptables-services
    iptables-save > /etc/sysconfig/iptables
    systemctl start iptables
    systemctl enable iptables

    echo Ready

    # Disable source/destination check
    aws ec2 modify-instance-attribute --instance-id $INSTANCE_ID --no-source-dest-check --region ${REGION}

    PRIMARY_ENI_ID=$(aws ec2 describe-instances --instance-ids i-03933236642547520 --region ca-central-1 --query "Reservations[0].Instances[0].NetworkInterfaces[?Attachment.DeviceIndex==\`0\`].NetworkInterfaceId" --output text)

    # Update the route table to set this instance as the active NAT instance
    OLD_NAT_ENI_ID=$(aws ec2 describe-route-tables --route-table-ids ${ROUTE_TABLE_ID} --region ${REGION} --query "RouteTables[].Routes[?DestinationCidrBlock=='0.0.0.0/0'].NetworkInterfaceId" --output text)

    if [ ! -z "$OLD_NAT_ENI_ID" ]; then
        aws ec2 replace-route --route-table-id ${ROUTE_TABLE_ID} --destination-cidr-block 0.0.0.0/0 --network-interface-id $PRIMARY_ENI_ID --region ${REGION}
        echo "Deleted route for old NAT instance $OLD_NAT_ENI_ID"
    else
        aws ec2 create-route --route-table-id ${ROUTE_TABLE_ID} --destination-cidr-block 0.0.0.0/0 --network-interface-id $PRIMARY_ENI_ID --region ${REGION}
        echo "Created route for new NAT interface $PRIMARY_ENI_ID"
    fi
fi
