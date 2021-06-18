output "security_group_cluster" {
  value = aws_security_group.odysseia-sc.id
}
//
//output "security_group_node" {
//  value = aws_security_group.odysseia-node-sc.id
//}