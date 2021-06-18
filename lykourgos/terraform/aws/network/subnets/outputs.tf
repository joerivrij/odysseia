output "subnets" {
  value = aws_subnet.odysseia-subnet[*].id
}