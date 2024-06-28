data "metal_cluster" "cluster" {
  name = "cluster"
}

output "cluster" {
  value = data.metal_cluster.cluster
}
