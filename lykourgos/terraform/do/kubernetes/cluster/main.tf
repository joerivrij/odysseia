terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
    }
  }
}

resource "digitalocean_kubernetes_cluster" "odysseia" {
  name   = var.cluster_name_odysseia
  region = var.region
  version = var.odysseia_version

  node_pool {
    name       = var.node_name
    size       = var.node_type
    node_count = var.node_size
  }
}