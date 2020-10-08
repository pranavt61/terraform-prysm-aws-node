variable "vpc_id" {
  description = "Custom vpc id - leave blank for deault"
  type        = string
  default     = ""
}

variable "create_sg" {
  type        = bool
  description = "Bool for create security group"
  default     = true
}

variable "public_tcp_ports" {
  description = "List of publicly open ports"
  type        = list(number)
  default = [
    22,
    7100,
    9000,
  ]
}

variable "public_udp_ports" {
  description = "List of publicly udp open ports"
  type        = list(number)
  default = [
    7100,
    9000,
  ]
}


variable "private_tcp_ports" {
  description = "List of publicly tcp open ports"
  type        = list(number)
  default = [
    9100,
    9113,
    9115,
    8080,
  ]
}

variable "private_udp_ports" {
  description = "List of publicly udp open ports"
  type        = list(number)
  default     = []
}


variable "private_port_cidrs" {
  description = "List of CIDR blocks for private ports"
  type        = list(string)
  default     = ["172.31.0.0/16"]
}

variable "additional_security_group_ids" {
  description = "List of security groups"
  type        = list(string)
  default     = []
}

resource "aws_security_group" "this" {
  count  = var.create_sg && var.create ? 1 : 0
  vpc_id = var.vpc_id == "" ? null : var.vpc_id

}

resource "aws_security_group_rule" "public_tcp_ports" {
  count = var.create_sg && var.create ? length(var.public_tcp_ports) : 0

  type              = "ingress"
  security_group_id = join("", aws_security_group.this.*.id)
  protocol          = "tcp"
  from_port         = var.public_tcp_ports[count.index]
  to_port           = var.public_tcp_ports[count.index]
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "public_udp_ports" {
  count = var.create_sg && var.create ? length(var.public_udp_ports) : 0

  type              = "ingress"
  security_group_id = join("", aws_security_group.this.*.id)
  protocol          = "udp"
  from_port         = var.public_udp_ports[count.index]
  to_port           = var.public_udp_ports[count.index]
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "private_tcp_ports" {
  count = var.create_sg && var.create ? length(var.private_tcp_ports) : 0

  type              = "ingress"
  security_group_id = join("", aws_security_group.this.*.id)
  protocol          = "tcp"
  from_port         = var.private_tcp_ports[count.index]
  to_port           = var.private_tcp_ports[count.index]
  cidr_blocks       = var.private_port_cidrs
}

resource "aws_security_group_rule" "private_udp_ports" {
  count = var.create_sg && var.create ? length(var.private_udp_ports) : 0

  type              = "ingress"
  security_group_id = join("", aws_security_group.this.*.id)
  protocol          = "udp"
  from_port         = var.private_udp_ports[count.index]
  to_port           = var.private_udp_ports[count.index]
  cidr_blocks       = var.private_port_cidrs
}

resource "aws_security_group_rule" "egress" {
  count             = var.create_sg && var.create ? 1 : 0
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = join("", aws_security_group.this.*.id)
  type              = "egress"
}