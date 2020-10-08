variable "create" {
  description = "Boolean to create resources or not"
  type        = bool
  default     = true
}

variable "name" {
  description = "The name for the label"
  type        = string
  default     = "prep"
}

variable "tags" {
  description = "Map of tags"
  type        = map(string)
  default     = {}
}

variable "network_name" {
  description = "The network name, ie medalla / mainnet"
  type        = string
}

