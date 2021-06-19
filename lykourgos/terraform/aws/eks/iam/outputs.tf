output "iam_cluster_arn" {
  value = aws_iam_role.odysseia-iam-cluster.arn
}

//output "iam_instance_profile" {
//  value = aws_iam_instance_profile
//}

output "iam_node_arn" {
  value = aws_iam_role.odysseia-iam-node.arn
}

output "eks_cluster_policy" {
  value = aws_iam_role_policy_attachment.odysseia-AmazonEKSClusterPolicy
}

output "eks_service_policy" {
  value = aws_iam_role_policy_attachment.odysseia-AmazonEKSServicePolicy
}

output "eks_node_policy" {
  value = aws_iam_role_policy_attachment.odysseia-node-AmazonEKSWorkerNodePolicy
}

output "eks_cni_policy" {
  value = aws_iam_role_policy_attachment.odysseia-node-AmazonEKS_CNI_Policy
}

output "eks_registry_policy" {
  value = aws_iam_role_policy_attachment.odysseia-node-AmazonEC2ContainerRegistryReadOnly
}
