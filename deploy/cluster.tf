data "aws_ami" "ecs_optimized_amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-*-arm64-ebs"]
  }
}

resource "aws_ecs_cluster" "funcie_cluster" {
  name = "funcie-cluster"
}

resource "aws_eip" "bastion_eip" {
  # Don't explicitly attach here, because user_data will do it
  domain = "vpc"
}

resource "aws_launch_template" "bastion_launch_template" {
  name = "bastion-launch-template"

  image_id      = data.aws_ami.ecs_optimized_amazon_linux.id
  instance_type = "t4g.micro"
  key_name      = var.bastion_public_key_path

  iam_instance_profile {
    name = aws_iam_instance_profile.instance_profile.name
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.server_bastion_sg.id]

    subnet_id = var.public_subnet_ids[0]
  }

  user_data = base64encode(templatefile("${path.module}/userdata.sh", {
    ECS_CLUSTER       = aws_ecs_cluster.funcie_cluster.name
    EIP_ALLOCATION_ID = aws_eip.bastion_eip.id
    REGION            = var.region
  }))

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "funcie-bastion-ec2-instance"
    }
  }
}

resource "aws_autoscaling_group" "bastion_asg" {
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

  depends_on = [aws_eip.bastion_eip]
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
          "ec2:AssociateAddress",
          "ec2:DescribeAddresses",
          "ec2:DescribeInstances"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_instance_profile" "instance_profile" {
  name = "instance_profile"
  role = aws_iam_role.instance_role.name
}
