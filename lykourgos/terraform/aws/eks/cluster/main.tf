#
# EKS Cluster Resources
#  * EKS Cluster
#

resource "aws_eks_cluster" "odysseia_cluster" {
  name     = var.cluster_name_odysseia_check
  role_arn = var.iam_cluster_arn
  version = var.odysseia_version

  vpc_config {
    security_group_ids = var.security_group_cluster[*]
    subnet_ids         = var.subnets[*]
  }

  depends_on = [
    var.eks_service_policy,
    var.eks_cluster_policy,
  ]
}