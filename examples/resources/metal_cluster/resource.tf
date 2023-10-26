resource "metal_cluster" "cluster" {
  name       = "terraform01"
  kubernetes = "1.26.9"
  partition  = "eqx-mu4"
  workers = [
    {
      name            = "default"
      machine_type    = "n1-medium-x86"
      min_size        = 1
      max_size        = 1
      max_surge       = 1
      max_unavailable = 1
    }
  ]
  maintenance = {
    kubernetes_autoupdate   = true
    machineimage_autoupdate = false
    time_window = {
      begin    = "05:00 AM"
      duration = 5
    }
  }
}

output "cluster" {
  value = metal_cluster.cluster
}
