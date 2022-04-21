###########################################
# variables related to kubernetes cluster #
##########################################

variable "odysseia_version" {
  default = "1.22.8-do.0"
  type    = string
  description = "cluster version for the odysseia"
}

variable "cluster_name_odysseia" {
  default = "kubernetes-do-odysseia"
  type    = string
  description = "name of the odysseia cluster"
}

variable "region" {
  default = "ams3"
  type    = string
  description = "region to deploy"
}

variable "node_size" {
  default = 4
  type    = string
  description = "size of the node group"
}

variable "node_name" {
  default = "kubernetes-do-odysseia-nodes"
  type    = string
  description = "name of the nodes"
}

variable "node_type" {
  default = "s-2vcpu-2gb"
  type    = string
  description = "type of nodes"
}