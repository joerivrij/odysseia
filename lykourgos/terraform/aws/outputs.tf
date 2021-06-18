output "vpc_id" {
  value = module.vpc.vpc_id
}

output "config_map_aws_auth" {
  value = module.eks_cluster.config_map_aws_auth
}

output "kubeconfig_app" {
  value = module.eks_cluster.kubeconfig_app
}

output "kubeconfig_odysseia" {
  value = module.eks_cluster.kubeconfig_odysseia
}

output "db_instance_address" {
  value = module.rds.db_instance_address
}

output "db_instance_username" {
  value = module.rds.db_instance_username
}

output "db_instance_password" {
  value = module.rds.db_instance_password
}

output "db_instance_port" {
  value = module.rds.db_instance_port
}