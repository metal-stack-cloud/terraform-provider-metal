package volume

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func VolumeDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Computed: false,
			Optional: true,
		},
		"name": datasourceschema.StringAttribute{
			Computed:            false,
			Optional:            true,
			MarkdownDescription: "Name of the cluster",
		},
		"project": datasourceschema.StringAttribute{
			Computed: true,
		},
		"partition": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"storageclass": datasourceschema.StringAttribute{
			Computed: true,
		},
		"replicacount": datasourceschema.Int64Attribute{
			Computed: true,
		},
	}
}

func SnapshotDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Computed: false,
			Optional: true,
		},
		"name": datasourceschema.StringAttribute{
			Computed:            false,
			Optional:            true,
			MarkdownDescription: "Name of the cluster",
		},
		"project": datasourceschema.StringAttribute{
			Computed: true,
		},
		"partition": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"size": datasourceschema.StringAttribute{
			Computed: true,
		},
		"usage": datasourceschema.Int64Attribute{
			Computed: true,
		},
		"volume_id": datasourceschema.Int64Attribute{
			Computed: true,
		},
	}
}
