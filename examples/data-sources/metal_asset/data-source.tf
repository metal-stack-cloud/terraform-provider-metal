data "metal_asset" "asset" {}

output "asset" {
  value = data.metal_asset.asset
}
