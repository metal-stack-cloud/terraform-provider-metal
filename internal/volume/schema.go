package volume

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func VolumeDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Computed:    false,
			Optional:    true,
			Description: "The id of the volume.",
		},
		"name": datasourceschema.StringAttribute{
			Computed:            false,
			Optional:            true,
			MarkdownDescription: "Name of the volume.",
		},
		"project": datasourceschema.StringAttribute{
			Computed:    true,
			Description: "The project id of the volume.",
		},
		"partition": resourceschema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The partition of the volume.",
		},
		"storageclass": datasourceschema.StringAttribute{
			Computed:    true,
			Description: "The used storage class of the volume.",
		},
		"replicacount": datasourceschema.Int64Attribute{
			Computed:    true,
			Description: "The amount of replicas used for the volume.",
		},
		"clustername": datasourceschema.StringAttribute{
			Computed:    true,
			Description: "The cluster name a volume is attached to.",
		},
		"labels": datasourceschema.MapAttribute{
			Computed:            true,
			MarkdownDescription: "The labels of a volume.",
			ElementType:         types.StringType,
		},
	}
}
