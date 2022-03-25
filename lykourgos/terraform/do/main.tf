module "odysseia-cluster" {
  source           = "./kubernetes/cluster"
  odysseia_version = var.odysseia_version
  cluster_name_odysseia = var.cluster_name_odysseia
  region = var.region
  node_size = var.node_size
  node_name = var.node_name
  node_type = var.node_type
}
