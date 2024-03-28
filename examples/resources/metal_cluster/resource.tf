resource "metal_cluster" "cluster" {
  name       = "cluster"
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
  maintenance = {
    time_window = {
      begin = {
        hour   = 18
        minute = 30
      }
      duration = 2
    }
  }
}

output "cluster" {
  value = metal_cluster.cluster
}
