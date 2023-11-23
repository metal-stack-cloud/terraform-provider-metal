data "metal_cluster" "panda" { # or use your resource
  name = "simon"
}

data "metal_kubeconfig" "panda_kubeconfig" {
  id         = data.metal_cluster.panda.id
  expiration = "1h02m"
}

output "panda_kubeconfig" {
  value = data.metal_kubeconfig.panda_kubeconfig.raw
}
