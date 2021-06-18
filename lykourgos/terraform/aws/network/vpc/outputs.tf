output "vpc_id" {
  value = aws_vpc.odysseia-vpc.id
}

output "vpc_cidr_block" {
  value = aws_vpc.odysseia-vpc.cidr_block
}

output "gw_id" {
  value = aws_internet_gateway.odysseia-gateway.id
}

output "main_route_table_id" {
  value = aws_route_table.odysseia-route-table.id
}
