variable "key_name" {
  description = "The key pair to import - leave blank to generate new keypair from pub/priv ssh key path"
  type        = string
  default     = ""
}

variable "root_volume_size" {
  description = "Root volume size"
  type        = number
  default     = 8
}

variable "root_volume_type" {
  description = ""
  type        = string
  default     = "gp2"
}

variable "root_iops" {
  description = ""
  type        = string
  default     = null
}

variable "instance_type" {
  description = "Instance type"
  type        = string
  default     = "t3.small"
}

variable "public_key_path" {
  description = "The path to the public ssh key"
  type        = string
}

variable "private_key_path" {
  description = "The path to the private ssh key"
  type        = string
}


variable "subnet_id" {
  description = "The id of the subnet"
  type        = string
  default     = ""
}

// TODO: Fill out volume size requirements for each network
variable "minimum_volume_size_map" {
  description = "Map for networks with min volume size "
  type        = map(string)
  default = {
    mainnet = 10,
    medalla = 10
  }
}

variable "keystore_password" {
  description = "Password to keystore"
  type        = string
}

variable "keystore_path" {
  description = "Path to keystore"
  type        = string
}

module "ami" {
  source = "github.com/insight-infrastructure/terraform-aws-ami.git?ref=v0.1.0"
}

resource "aws_key_pair" "this" {
  count      = var.public_key_path != "" && var.create ? 1 : 0
  public_key = file(pathexpand(var.public_key_path))
}

locals {
  root_volume_size = var.root_volume_size == "8" ? var.root_volume_size : lookup(var.minimum_volume_size_map, var.network_name)
  tags             = merge(var.tags, { Name = var.name })
}

data "template_file" "user_data" {
  template = file("${path.module}/data/user-data.sh")
  vars = {
    keystore_password = var.keystore_password
  }
}

resource "aws_instance" "this" {
  count         = var.create ? 1 : 0
  ami           = module.ami.ubuntu_1804_ami_id
  instance_type = var.instance_type

  root_block_device {
    volume_size = local.root_volume_size
    volume_type = var.root_volume_type
    iops        = var.root_iops
  }

  subnet_id              = var.subnet_id
  vpc_security_group_ids = compact(concat(aws_security_group.this.*.id, var.additional_security_group_ids))
  key_name               = var.public_key_path == "" ? var.key_name : aws_key_pair.this.*.key_name[0]
  tags                   = merge({ name = var.name }, local.tags)

  user_data = data.template_file.user_data.rendered

  provisioner "remote-exec" {
    script = "${path.module}/data/wait-for-apt-on-startup.sh"
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = "ubuntu"
      private_key = file(var.private_key_path)
    }
  }

  provisioner "remote-exec" {
    inline = [
      "docker-compose -f /home/ubuntu/prysm-docker-compose/create-account.yaml run validator-import-launchpad"
    ]
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = "ubuntu"
      private_key = file(var.private_key_path)
    }
  }

  provisioner "file" {
    source      = var.keystore_path
    destination = "/home/ubuntu/prysm-docker-compose/launchpad/eth2.0-deposit-cli/validator_keys/keystore"
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = "ubuntu"
      private_key = var.private_key_path
    }
  }

  provisioner "remote-exec" {
    inline = [
      "docker-compose -f /home/ubuntu/prysm-docker-compose/create-account.yaml run validator-import-launchpad"
    ]
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = "ubuntu"
      private_key = file(var.private_key_path)
    }
  }

  provisioner "remote-exec" {
    inline = [
      "docker-compose -f /home/ubuntu/prysm-docker-compose/docker-compose.yaml"
    ]
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = "ubuntu"
      private_key = file(var.private_key_path)
    }
  }
}
