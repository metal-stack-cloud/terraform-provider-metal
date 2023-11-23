data "metal_snapshot" "snapshot_by_name" {
  name = "snapshot-21e68bcc-daa0-489a-aeb4-224fe150d5f7"
}

output "snapshot_by_name" {
  value = data.metal_snapshot.snapshot_by_name
}
