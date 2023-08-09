resource "metal_public_ip" "my_ip" {
  name        = "my_ip"
  description = "Some description"
  type        = "ephemeral" # either ephemeral or static
  tags        = ["test"]
}
