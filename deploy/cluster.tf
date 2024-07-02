data "aws_ssm_parameter" "al2023_ami" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
}

resource "aws_ecs_cluster" "funcie_cluster" {
  name = "funcie-cluster"
}

data "aws_ip_ranges" "ec2_instance_connect" {
  services = ["ec2_instance_connect"]
  regions  = [var.region]
}

resource "aws_security_group" "ec2_instance_connect" {
  name        = "ec2-instance-connect"
  description = "Security group for EC2 Instance Connect"
  vpc_id      = local.vpc_id

  dynamic "ingress" {
    for_each = data.aws_ip_ranges.ec2_instance_connect.cidr_blocks
    content {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "nat_instance_sg" {
  name        = "nat-instance-sg"
  description = "Security group for the NAT instance"
  vpc_id      = local.vpc_id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [local.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_launch_template" "bastion_launch_template" {
  name = "bastion-launch-template"

  image_id      = data.aws_ssm_parameter.al2023_ami.value
  instance_type = var.bastion_instance_type

  iam_instance_profile {
    name = aws_iam_instance_profile.instance_profile.name
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups = concat(
      [aws_security_group.server_bastion_sg.id, aws_security_group.ec2_instance_connect.id],
      var.vpc_id == "" ? [aws_security_group.nat_instance_sg.id] : []
    )

    subnet_id = local.public_subnet_ids[0]
  }

  metadata_options {
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
    http_endpoint               = "enabled"
  }

  user_data = base64encode(templatefile("${path.module}/userdata.sh", {
    ECS_CLUSTER    = aws_ecs_cluster.funcie_cluster.name
    REGION         = var.region
    FUNCIE_ENV     = var.funcie_env
    ROUTE_TABLE_ID = var.vpc_id == "" ? aws_route_table.funcie_nat_route_table[0].id : ""
    CREATE_VPC     = var.vpc_id == "" ? "true" : ""
  }))

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "funcie-bastion-ec2-instance"
    }
  }
}

resource "aws_autoscaling_group" "bastion_asg" {
  name = "${aws_launch_template.bastion_launch_template.name}-asg"

  desired_capacity = 1
  max_size         = 1
  min_size         = 1

  launch_template {
    id      = aws_launch_template.bastion_launch_template.id
    version = "$Latest"
  }

  vpc_zone_identifier = local.public_subnet_ids

  tag {
    key                 = "Name"
    value               = "bastion-ec2-instance"
    propagate_at_launch = true
  }
}

resource "aws_iam_role" "instance_role" {
  name = "instance_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "instance_policy" {
  name = "instance_policy"
  role = aws_iam_role.instance_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          // ECS-related permissions to allow the instance to register with the cluster and run tasks / send logs
          "ecs:RegisterContainerInstance",
          "ecs:DeregisterContainerInstance",
          "ecs:DiscoverPollEndpoint",
          "ecs:SubmitContainerStateChange",
          "ecs:SubmitTaskStateChange",
          "ecs:SubmitAttachmentStateChanges",
          "ecs:SubmitInstanceStateChange",
          "ecs:SubmitTaskStateChanges",
          "ecs:Poll",
          "ecs:StartTelemetrySession",
          "ecs:UpdateContainerInstancesState",

          "logs:CreateLogStream",
          "logs:PutLogEvents",
          # EC2 permissions to allow the instance to manage its own EIP
          "ec2:AssociateAddress",
          "ec2:DisassociateAddress",
          "ec2:DescribeAddresses",
          "ec2:DescribeInstances",
          # EC2 permissions to allow the instance to update the route table to register itself as the NAT instance
          "ec2:CreateRoute",
          "ec2:DeleteRoute",
          "ec2:DescribeRouteTables",
          "ec2:ReplaceRoute",
          "ec2:ModifyInstanceAttribute",
        ],
        Resource = "*"
      },
      {
        Effect = "Allow",
        Action = [
          "ssm:PutParameter",
          "ssm:GetParameter",
        ],
        Resource = "arn:aws:ssm:*:*:parameter/funcie/*",
      }
    ],
  })
}

resource "aws_iam_role_policy_attachment" "ssm_role_policy_attachment" {
  role       = aws_iam_role.instance_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "instance_profile" {
  name = "instance_profile"
  role = aws_iam_role.instance_role.name
}

resource "null_resource" "asg_update_trigger" {
  # This is a hack to force the ASG to use the latest launch template version
  triggers = {
    launch_template_version = aws_launch_template.bastion_launch_template.latest_version
    user_data = base64encode(templatefile("${path.module}/userdata.sh", {
      ECS_CLUSTER    = aws_ecs_cluster.funcie_cluster.name
      REGION         = var.region
      FUNCIE_ENV     = var.funcie_env
      ROUTE_TABLE_ID = var.vpc_id == "" ? aws_route_table.funcie_nat_route_table[0].id : ""
      CREATE_VPC     = var.vpc_id == ""
    }))
    instance_policy = aws_iam_role_policy.instance_policy.policy
  }

  provisioner "local-exec" {
    command = <<EOT
      aws autoscaling start-instance-refresh --auto-scaling-group-name ${aws_autoscaling_group.bastion_asg.name} --region ${var.region}
    EOT
  }
}
