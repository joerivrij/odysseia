#
# EKS Worker Nodes Resources
#  * IAM role allowing Kubernetes actions to access other AWS services
#  * EKS Node Group to launch manager nodes
#

resource "aws_eks_node_group" "odysseia_check_nodes" {
  cluster_name    = var.eks_cluster_name_odysseia
  node_group_name = var.eks_node_group_name_odysseia
  node_role_arn   = var.iam_node_arn
  subnet_ids      = var.subnets
  version         = var.odysseia_version
  instance_types  = ["t3.medium"]
  disk_size = var.disk_size_odysseia

  scaling_config {
    desired_size = var.desired_size_odysseia
    max_size     = var.max_size_odysseia
    min_size     = var.min_size_odysseia
  }

  tags = {
    key                 = "kubernetes.io/cluster/${var.eks_cluster_name_odysseia}"
    value               = "owned"
    "k8s.io/cluster-autoscaler/enabled" = true
    "k8s.io/cluster-autoscaler/${var.eks_cluster_name_odysseia}" = true
  }

  depends_on = [
    var.eks_node_policy,
    var.eks_cni_policy,
    var.eks_registry_policy,
  ]
}
