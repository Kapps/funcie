data "aws_ssm_parameter" "ecs_al2023_arm64_ami" {
  name = "/aws/service/ecs/optimized-ami/amazon-linux-2023/recommended/image_id"
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
  vpc_id      = var.vpc_id

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

resource "aws_launch_template" "bastion_launch_template" {
  name = "bastion-launch-template"

  image_id      = data.aws_ssm_parameter.ecs_al2023_arm64_ami.value
  instance_type = "t3.micro"

  iam_instance_profile {
    name = aws_iam_instance_profile.instance_profile.name
  }

  network_interfaces {
    associate_public_ip_address = true # Even though we have an EIP, we still need this to be able to use the CLI to associate it
    security_groups             = [aws_security_group.server_bastion_sg.id, aws_security_group.ec2_instance_connect.id]

    subnet_id = var.public_subnet_ids[0]
  }

  metadata_options {
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
    http_endpoint               = "enabled"
  }

  user_data = base64encode(templatefile("${path.module}/userdata.sh", {
    ECS_CLUSTER = aws_ecs_cluster.funcie_cluster.name
    REGION      = var.region
    FUNCIE_ENV  = var.funcie_env
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

  vpc_zone_identifier = var.public_subnet_ids

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
          "ec2:AssociateAddress",
          "ec2:DisassociateAddress",
          "ec2:DescribeAddresses",
          "ec2:DescribeInstances",
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
      ECS_CLUSTER = aws_ecs_cluster.funcie_cluster.name
      REGION      = var.region
      FUNCIE_ENV  = var.funcie_env
    }))
    instance_policy = aws_iam_role_policy.instance_policy.policy
  }

  provisioner "local-exec" {
    command = <<EOT
      aws autoscaling start-instance-refresh --auto-scaling-group-name ${aws_autoscaling_group.bastion_asg.name} --region ${var.region}
    EOT
  }
}
