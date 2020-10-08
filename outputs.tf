

output "instance_type" {
  value = var.instance_type
}

output "network_name" {
  value = var.network_name
}

output "instance_id" {
  value = join("", aws_instance.this.*.id)
}

output "key_name" {
  value = join("", aws_key_pair.this.*.key_name)
}

output "public_ip" {
  value = join("", aws_instance.this.*.public_ip)
}

output "security_group_id" {
  value = join("", aws_security_group.this.*.id)
}