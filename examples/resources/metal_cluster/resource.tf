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
  # Manually setting the maintenance and the time window is currently not support.
  # This is going to change in the future.
  # https://github.com/metal-stack-cloud/terraform-provider-metal/issues/51
  # maintenance = {
  #   kubernetes_autoupdate   = true
  #   machineimage_autoupdate = false
  #   time_window = {
  #     begin    = "05:00 AM"
  #     duration = 1
  #   }
  # }
}

output "cluster" {
  value = metal_cluster.cluster
}
