terraform {
  required_providers {
    metal = {
      source = "metalstack.cloud/terraform/metal"
    }
  }
}

provider "metal" {
    organization = "x-cellent@github"
    project = "f8e67080-ba68-41d2-ad44-59dc65a09a33"
}

data "metal_ip_addresses" "example" {
}

# output "default_ips" {
#   value = data.metal_ip_addresses.example.list
# }

# resource "metal_ip_address" "terrify" {
#   name = "terrify-lemming"
#   description = "Hallo Welt!"
#   type = "static" // just eph -> static
#   tags = ["test"]
# }

# output "example_ip" {
#   value = metal_ip_address.terrify
# }

resource "metal_ip_address" "imported" {
  name = "tadpole"
  type = "static"
  description = "registry"
  tags = []
  # tags = ["cluster.metal-stack.io/id/namespace/service=6fd03373-5e9a-45a7-8607-7492a25fa850/ingress-nginx/nginx-ingress-ingress-nginx-controller"]
}
