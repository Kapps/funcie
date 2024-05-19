data "aws_ssm_parameter" "ecs_al2023_arm64_ami" {
  name = "/aws/service/ecs/optimized-ami/amazon-linux-2023/arm64/recommended/image_id"
}

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

  image_id      = data.aws_ssm_parameter.ecs_al2023_arm64_ami.value
  instance_type = "t4g.micro"
  key_name      = aws_key_pair.bastion_key.key_name

  iam_instance_profile {
    name = aws_iam_instance_profile.instance_profile.name
  }

  network_interfaces {
    associate_public_ip_address = true # Even though we have an EIP, we still need this to be able to use the CLI to associate it
    security_groups             = [aws_security_group.server_bastion_sg.id]

    subnet_id = var.public_subnet_ids[0]
  }

  metadata_options {
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
    http_endpoint               = "enabled"
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

resource "tls_private_key" "bastion_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "bastion_key" {
  key_name   = "bastion-key"
  public_key = tls_private_key.bastion_key.public_key_openssh

}

resource "null_resource" "asg_update_trigger" {
  # This is a hack to force the ASG to use the latest launch template version
  triggers = {
    launch_template_version = aws_launch_template.bastion_launch_template.latest_version
    user_data = base64encode(templatefile("${path.module}/userdata.sh", {
      ECS_CLUSTER       = aws_ecs_cluster.funcie_cluster.name
      EIP_ALLOCATION_ID = aws_eip.bastion_eip.id
      REGION            = var.region
    }))
  }

  provisioner "local-exec" {
    command = <<EOT
      aws autoscaling start-instance-refresh --auto-scaling-group-name ${aws_autoscaling_group.bastion_asg.name} --region ${var.region}
    EOT
  }
}
