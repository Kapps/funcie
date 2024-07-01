#!/bin/bash

# Get a session token from IMDSv2
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")

# Get the instance ID using the token
INSTANCE_ID=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/instance-id)
PRIVATE_IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/local-ipv4)

echo ECS_CLUSTER=${ECS_CLUSTER} >> /etc/ecs/ecs.config

aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_host --value $PRIVATE_IP --type String --overwrite --region ${REGION}
aws ssm put-parameter --name /funcie/${FUNCIE_ENV}/bastion_instance_id --value $INSTANCE_ID --type String --overwrite --region ${REGION}

sudo yum install -y ec2-instance-connect

if [ "${CREATE_VPC}" = "true" ]; then
    # Disable source/destination check
    aws ec2 modify-instance-attribute --instance-id $INSTANCE_ID --no-source-dest-check --region ${REGION}

    # Enable IP forwarding
    echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
    sysctl -p /etc/sysctl.conf

    # Configure iptables for NAT
    iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
    iptables-save > /etc/sysconfig/iptables

    # Ensure iptables-persistent is installed and enabled
    yum install -y iptables-services
    systemctl enable iptables
    systemctl start iptables

    # Update the route table to set this instance as the active NAT instance
    OLD_NAT_INSTANCE_ID=$(aws ec2 describe-route-tables --route-table-ids ${ROUTE_TABLE_ID} --region ${REGION} --query "RouteTables[].Routes[?DestinationCidrBlock=='0.0.0.0/0'].InstanceId" --output text)

    if [ ! -z "$OLD_NAT_INSTANCE_ID" ]; then
        aws ec2 delete-route --route-table-id ${ROUTE_TABLE_ID} --destination-cidr-block 0.0.0.0/0 --region ${REGION}
        echo "Deleted route for old NAT instance $OLD_NAT_INSTANCE_ID"
    fi

    aws ec2 create-route --route-table-id ${ROUTE_TABLE_ID} --destination-cidr-block 0.0.0.0/0 --instance-id $INSTANCE_ID --region ${REGION}
    echo "Created route for new NAT instance $INSTANCE_ID"
fi
