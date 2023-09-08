resource "metal_cluster" "cluster" {
  name       = "terraform01"
  kubernetes = "1.24.14"
  workers = {
    machinetype    = "n1-medium-x86"
    minsize        = 1
    maxsize        = 2
    maxsurge       = 1
    maxunavailable = 1
  }
}
