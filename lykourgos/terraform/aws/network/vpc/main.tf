data "aws_availability_zones" "available" {
}

resource "aws_vpc" "odysseia-vpc" {
  cidr_block = var.cidr_block

  tags = {
    "Name"                                      = var.tag-value-name
    "kubernetes.io/cluster/${var.eks_odysseia_cluster_name}" = "shared"
  }
}

resource "aws_internet_gateway" "odysseia-gateway" {
  vpc_id = aws_vpc.odysseia-vpc.id

  tags = {
    Name = var.internet-gateway
  }
}

resource "aws_route_table" "odysseia-route-table" {
  vpc_id = aws_vpc.odysseia-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.odysseia-gateway.id
  }
}
