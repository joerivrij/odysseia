variable "cluster_name_odysseia_check" {}

variable "iam_cluster_arn" {
  type    = string
  description = "Arn for the cluster"
}

variable "iam_node_arn" {
  type    = string
  description = "Role for the cluster"
}

variable "subnets" {
  type = list(string)
}

variable "eks_service_policy" {
  description = "Service policy for the cluster"
}

variable "eks_cluster_policy" {

}

variable "security_group_cluster" {}

variable "odysseia_version" {
  default = 1.19
  type    = string
  description = "cluster version for the odysseia check cluster, this has to be set else version 1.16"
}