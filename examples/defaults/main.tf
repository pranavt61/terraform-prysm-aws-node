variable "aws_region" {
  default = "us-east-1"
}

provider "aws" {
  region = var.aws_region
}

variable "private_key_path" {}
variable "public_key_path" {}

module "defaults" {
  source                = "../.."
  network_name          = "medalla"
  private_key_path      = var.private_key_path
  public_key_path       = var.public_key_path
  keystore_password     = "ItWorks!!!!1"
  keystore_path         = "${path.module}/../../test/fixtures/validator/validator_keys/keystore-default.json"
  deposit_path          = "${path.module}/../../test/fixtures/validator/validator_keys/deposit-default.json"
  wallets_dir_path      = "${path.module}/../../test/fixtures/validator/wallets"
  wallet_password_path  = "${path.module}/../../test/fixtures/validator/passwords/wallet-password"
}

output "public_ip" {
  value = module.defaults.public_ip
}
