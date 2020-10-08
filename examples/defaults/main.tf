variable "aws_region" {
  default = "us-east-1"
}

provider "aws" {
  region = var.aws_region
}

variable "private_key_path" {}
variable "public_key_path" {}

//variable "keystore_path" {
//  default = ""
//}
//locals {
//  keystore_path = var.keystore_path == "" ? "${path.module}/../../test/fixtures/keystores/keystore-default" : var.keystore_path
//}

module "defaults" {
  source            = "../.."
  network_name      = "medalla"
  private_key_path  = var.private_key_path
  public_key_path   = var.public_key_path
  keystore_password = "ItWorks!!!!1"
  keystore_path     = "${path.module}/../../test/fixtures/keystores/keystore-default"
}

output "public_ip" {
  value = module.defaults.public_ip
}