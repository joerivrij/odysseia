data "aws_availability_zones" "available" {}

resource "aws_subnet" "odysseia-subnet" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = var.vpc_id
  map_public_ip_on_launch = true

  tags = {
    "Name"                                      = var.tag-value-name
    "kubernetes.io/cluster/${var.eks_odysseia_cluster_name}" = "shared"
  }
}