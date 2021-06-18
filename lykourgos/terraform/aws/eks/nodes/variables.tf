
variable "eks_cluster_name_odysseia" {
  description = "cluster name"
}

variable "eks_endpoint_odysseia" {
  description = "eks cluster endpoint"
}

variable "eks_node_group_name_odysseia" {}

variable "eks_certificate_authority_odysseia" {
  description = "eks certificate authority"
}

variable "iam_node_arn" {}

variable "subnets" {}

variable "max_size_odysseia" {}

variable "min_size_odysseia" {}

variable "desired_size_odysseia" {}

variable "eks_node_policy" {}

variable "eks_cni_policy" {}

variable "eks_registry_policy" {}

variable "odysseia_version" {
  default = 1.19
  type    = string
  description = "cluster version for the odysseia check cluster, this has to be set else version 1.16"
}

variable "security_group_eks" {}

variable "disk_size_odysseia" {
  default = 100
}
