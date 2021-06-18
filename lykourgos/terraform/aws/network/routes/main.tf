resource "aws_route_table_association" "odysseia-table-association" {
  count = 2

  subnet_id      = var.subnets[count.index]
  route_table_id = var.main_route_table_id
}

//# Internet access
//resource "aws_route" "route_internet" {
//  route_table_id         = var.main_route_table_id
//  destination_cidr_block = "0.0.0.0/0"
//  gateway_id             = var.gateway_id
//  depends_on             = ["aws_route_table_association.odysseia-table-association"]
//}