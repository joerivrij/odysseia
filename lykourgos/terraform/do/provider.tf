# Configure the DigitalOcean Provider
variable "do_token" {}

provider "digitalocean" {
  token = var.do_token
}