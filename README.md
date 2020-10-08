# terraform-aws-prysm-node

[![open-issues](https://img.shields.io/github/issues-raw/insight-infrastructure/terraform-aws-prysm-node?style=for-the-badge)](https://github.com/insight-infrastructure/terraform-aws-prysm-node/issues)
[![open-pr](https://img.shields.io/github/issues-pr-raw/insight-infrastructure/terraform-aws-prysm-node?style=for-the-badge)](https://github.com/insight-infrastructure/terraform-aws-prysm-node/pulls)

> Not ready - WIP

## TODO

- Fix deployment (compose error....)
- port lists 
- ssh into box / manually run user data as root to verify 
- Replace keystore with dummy 
- Build test so that endpoint is tested when successfully deployed 
    - https://github.com/insight-w3f/terraform-polkadot-aws-asg/blob/master/test/terraform_defaults_test.go
   
```bash
cd examples/defaults
terraform init 
terraform apply
terraform destroy  
```

## Features

This module...

## Terraform Versions

For Terraform v0.12.0+

## Usage

```hcl
module "this" {
  source = "github.com/insight-infrastructure/terraform-aws-prysm-node"
}
```
## Examples

- [defaults](https://github.com/insight-infrastructure/terraform-aws-prysm-node/tree/master/examples/defaults)

## Known  Issues
No issue is creating limit on this module.

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Providers

| Name | Version |
|------|---------|
| aws | n/a |
| template | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:-----:|
| additional\_security\_group\_ids | List of security groups | `list(string)` | `[]` | no |
| create | Boolean to create resources or not | `bool` | `true` | no |
| create\_sg | Bool for create security group | `bool` | `true` | no |
| instance\_type | Instance type | `string` | `"t3.small"` | no |
| key\_name | The key pair to import - leave blank to generate new keypair from pub/priv ssh key path | `string` | `""` | no |
| keystore\_password | Password to keystore | `string` | n/a | yes |
| keystore\_path | Path to keystore | `string` | n/a | yes |
| minimum\_volume\_size\_map | Map for networks with min volume size | `map(string)` | <pre>{<br>  "mainnet": 10,<br>  "medalla": 10<br>}</pre> | no |
| name | The name for the label | `string` | `"prep"` | no |
| network\_name | The network name, ie medalla / mainnet | `string` | n/a | yes |
| private\_key\_path | The path to the private ssh key | `string` | n/a | yes |
| private\_port\_cidrs | List of CIDR blocks for private ports | `list(string)` | <pre>[<br>  "172.31.0.0/16"<br>]</pre> | no |
| private\_tcp\_ports | List of publicly tcp open ports | `list(number)` | <pre>[<br>  9100,<br>  9113,<br>  9115,<br>  8080<br>]</pre> | no |
| private\_udp\_ports | List of publicly udp open ports | `list(number)` | `[]` | no |
| public\_key\_path | The path to the public ssh key | `string` | n/a | yes |
| public\_tcp\_ports | List of publicly open ports | `list(number)` | <pre>[<br>  22,<br>  7100,<br>  9000<br>]</pre> | no |
| public\_udp\_ports | List of publicly udp open ports | `list(number)` | <pre>[<br>  7100,<br>  9000<br>]</pre> | no |
| root\_iops | n/a | `string` | n/a | yes |
| root\_volume\_size | Root volume size | `number` | `8` | no |
| root\_volume\_type | n/a | `string` | `"gp2"` | no |
| subnet\_id | The id of the subnet | `string` | `""` | no |
| tags | Map of tags | `map(string)` | `{}` | no |
| vpc\_id | Custom vpc id - leave blank for deault | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| instance\_id | n/a |
| instance\_type | n/a |
| key\_name | n/a |
| network\_name | n/a |
| public\_ip | n/a |
| security\_group\_id | n/a |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

## Testing
This module has been packaged with terratest tests

To run them:

1. Install Go
2. Run `make test-init` from the root of this repo
3. Run `make test` again from root

## Authors

Module managed by [insight-infrastructure](https://github.com/insight-infrastructure)

## Credits

- [Anton Babenko](https://github.com/antonbabenko)

## License

Apache 2 Licensed. See LICENSE for full details.