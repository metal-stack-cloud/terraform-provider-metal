data "metal_assets" "assets" {}

output "assets" {
  value = data.metal_assets.assets
}
