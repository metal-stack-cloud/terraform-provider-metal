data "metal_volume" "volume" {
  name = "pvc-0607cf99-eaf5-412b-ab71-aba468c4219a"
}

output "volume" {
  value = data.metal_volume.volume
}
