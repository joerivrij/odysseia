####################################
# variables related to eks cluster #
####################################

variable "cluster-name-odysseia" {
  default = "eks-odysseia-cluster"
  type    = string
  description = "Master cluster name for the odysseia cluster"
}

variable "node-name-odysseia" {
  default = "eks-odysseia-node-groups"
  type    = string
  description = "Name for the node groups"
}

variable "desired_size_odysseia" {
  default = 3
  type = number
  description = "Number of worker nodes for the odysseia cluster"
}

variable "max_size_odysseia" {
  default = 7
  type = number
  description = "Maximum number of worker nodes for the odysseia cluster"
}

variable "min_size_odysseia" {
  default = 3
  type = number
  description = "Maximum number of worker nodes for the odysseia cluster"
}


####################################
# variables related to network     #
####################################

variable "internet-gateway" {
  default = "odysseia-gateway-eks"
  type = string
}

variable "cidr_block" {
  default = "10.0.0.0/16"
  type = string
  description = "CIDR for the vpn"
}

variable "security-group-name" {
  default = "odysseia-security-eks"
  type = string
  description = "name for the Security group"
}

variable "security-group-node-name" {
  default = "odysseia-security-eks-nodes"
  type = string
  description = "name for the Security group"
}

####################################
# tags                             #
####################################

variable "tag-value-name" {
  default = "odysseia-eks"
  type = string
  description = "Tag to add to the instances"
}
