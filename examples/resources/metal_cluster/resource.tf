resource "metal_cluster" "my-cluster" {
  name       = "my-cluster"
  kubernetes = "1.27.8"
  partition  = "eqx-mu4"
  workers = [
    {
      name         = "default"
      machine_type = "n1-medium-x86"
      min_size     = 1
      max_size     = 3
    }
  ]
}

output "cluster" {
  value = metal_cluster.cluster
}
