package snapshot

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func SnapshotDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Computed:    false,
			Optional:    true,
			Description: "The id of the snapshot. Can be used to query the snapshot.",
		},
		"name": datasourceschema.StringAttribute{
			Computed:    false,
			Optional:    true,
			Description: "The name of the snapshot. Can be used to query the snapshot. Typically starts with `pvc`.",
		},
		"volume_id": datasourceschema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The original volume for this snapshot.",
		},
		"project": datasourceschema.StringAttribute{
			Computed:    true,
			Description: "The project the snapshot is in.",
		},
		"partition": resourceschema.StringAttribute{
			Computed:    true,
			Description: "The partition of the snapshot.",
		},
		"storage_class": datasourceschema.StringAttribute{
			Computed:    true,
			Description: "The storage class of the snapshot.",
		},
		"size": datasourceschema.Int64Attribute{
			Computed:    true,
			Description: "The size of the snapshot.",
		},
		"usage": datasourceschema.Int64Attribute{
			Computed:    true,
			Description: "The usage of the snapshot",
		},
	}
}
