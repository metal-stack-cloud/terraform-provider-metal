data "metal_cluster" "cluster" {
  name = "cb-infra"
}

output "cluster" {
  value = data.metal_cluster.cluster
}
