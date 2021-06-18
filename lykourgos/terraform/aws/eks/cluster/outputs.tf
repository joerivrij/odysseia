locals {
  config_map_aws_auth = <<CONFIGMAPAWSAUTH
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: kube-system
data:
  mapRoles: |
    - rolearn: ${var.iam_node_arn}
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
CONFIGMAPAWSAUTH
  kubeconfig_odysseia = <<KUBECONFIG
apiVersion: v1
clusters:
- cluster:
    server: ${aws_eks_cluster.odysseia_cluster.endpoint}
    certificate-authority-data: ${aws_eks_cluster.odysseia_cluster.certificate_authority.0.data}
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: aws
  name: aws
current-context: aws
kind: Config
preferences: {}
users:
- name: aws
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: aws-iam-authenticator
      args:
        - "token"
        - "-i"
        - "${var.cluster_name_odysseia_check}"
KUBECONFIG
}

output "config_map_aws_auth" {
  value = local.config_map_aws_auth
}

output "kubeconfig_odysseia" {
  value = local.kubeconfig_odysseia
}

output "eks_certificate_authority_odysseia" {
  value = aws_eks_cluster.odysseia_cluster.certificate_authority.0.data
}

output "eks_endpoint_odysseia" {
  value = aws_eks_cluster.odysseia_cluster.endpoint
}

output "eks_cluster_name_odysseia" {
  value = aws_eks_cluster.odysseia_cluster.name
}

