variable "iam-cluster-name" {
  default = "odysseia-iam-cluster-eks"
  type = string
  description = "name for the Iam account"
}

variable "iam-node-name" {
  default = "odysseia-iam-node-eks"
  type = string
  description = "name for the Iam account"
}

variable "eks_odysseia_cluster_name" {}