resource "metal_cluster" "cluster" {
  name       = "terraform01"
  kubernetes = "1.24.14"
  partition  = "eqx-mu4"
  workers = [
    {
      name            = "default"
      machine_type    = "n1-medium-x86"
      min_size        = 1
      max_size        = 2
      max_surge       = 1
      max_unavailable = 1
    }
  ]
  maintenance = {
    # not working yet
    # kubernetes_autoupdate = false
    # machineimage_autoupdate = true
  }
}

output "cluster" {
  value = metal_cluster.cluster
}