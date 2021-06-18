## Network
# Create VPC
module "vpc" {
  source           = "./network/vpc"
  eks_odysseia_cluster_name = var.cluster-name-odysseia
  cidr_block       = var.cidr_block
  tag-value-name   = var.tag-value-name
  internet-gateway = var.internet-gateway
}

# Create Subnets
module "subnets" {
  source           = "./network/subnets"
  eks_odysseia_cluster_name = var.cluster-name-odysseia
  vpc_id           = module.vpc.vpc_id
  vpc_cidr_block   = module.vpc.vpc_cidr_block
  tag-value-name   = var.tag-value-name
}

# Configure Routes
module "route" {
  source = "./network/routes"
  main_route_table_id = module.vpc.main_route_table_id
  gateway_id = module.vpc.gw_id
  subnets = module.subnets.subnets
}
module "eks_iam_roles" {
  source = "./eks/iam"
  eks_odysseia_cluster_name = var.cluster-name-odysseia
}

module "security_group_eks" {
  source           = "./eks/security_groups"
  eks_odysseia_cluster_name = var.cluster-name-odysseia
  vpc_id           = module.vpc.vpc_id
  vpc_cidr_block   = module.vpc.vpc_cidr_block
  tag-value-name   = var.tag-value-name
  security-group-eks-name = var.security-group-name
  security-group-eks-node-name = var.security-group-node-name
}

module "eks_cluster" {
  source           = "./eks/cluster"
  cluster_name_odysseia_check = var.cluster-name-odysseia
  iam_cluster_arn  = module.eks_iam_roles.iam_cluster_arn
  iam_node_arn     = module.eks_iam_roles.iam_node_arn
  subnets = module.subnets.subnets
  security_group_cluster = module.security_group_eks.security_group_cluster
  eks_service_policy =  module.eks_iam_roles.eks_service_policy
  eks_cluster_policy = module.eks_iam_roles.eks_cluster_policy
}

module "eks_nodes" {
  source = "./eks/nodes"
  eks_cluster_name_odysseia = module.eks_cluster.eks_cluster_name_odysseia
  eks_certificate_authority_odysseia = module.eks_cluster.eks_certificate_authority_odysseia
  eks_endpoint_odysseia = module.eks_cluster.eks_endpoint_odysseia
  eks_node_group_name_odysseia = var.node-name-odysseia
  iam_node_arn = module.eks_iam_roles.iam_node_arn
  desired_size_odysseia = var.desired_size_odysseia
  max_size_odysseia = var.max_size_odysseia
  min_size_odysseia = var.min_size_odysseia
  eks_node_policy = module.eks_iam_roles.eks_node_policy
  eks_cni_policy = module.eks_iam_roles.eks_cni_policy
  eks_registry_policy = module.eks_iam_roles.eks_registry_policy
  subnets = module.subnets.subnets
  security_group_eks = module.security_group_eks.security_group_cluster
}
