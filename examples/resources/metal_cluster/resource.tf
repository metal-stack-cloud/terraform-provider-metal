resource "metal_cluster" "cluster" {
  name = "terraform-01"
  workers = {
    # MachineType = "n1-medium-x86"
    Minsize     = 1
    Maxsize     = 2
  }
}
