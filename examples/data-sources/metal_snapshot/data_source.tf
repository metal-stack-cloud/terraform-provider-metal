data "metal_volume" "volume" {
  name = "pvc-0607cf99-eaf5-412b-ab71-aba468c4219a"
}

data "metal_snapshot" "snapshot_name" {
  name = "pvc-0607cf99-eaf5-412b-ab71-aba468c4219a"
}

data "metal_snapshot" "snapshot_volume_name" {
  volume_id = data.metal_volume.volume.id
}

output "snapshot_name" {
  value = data.metal_snapshot.snapshot_name
}

output "snapshot_volume_name" {
  value = data.metal_snapshot.snapshot_volume_name
}
