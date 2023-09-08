package cluster

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func clusterResourceAttributes() map[string]resourceschema.Attribute {
	return map[string]resourceschema.Attribute{
		"id": resourceschema.StringAttribute{
			Computed: true,
		},
		// validate string
		"name": resourceschema.StringAttribute{
			Required: true,
		},
		"project": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"partition": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"kubernetes": resourceschema.StringAttribute{
			Optional: true,
		},
		"workers": resourceschema.SingleNestedAttribute{
			Required:            true,
			MarkdownDescription: "Worker settings",
			Attributes: map[string]resourceschema.Attribute{
				"machinetype": resourceschema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The the type of node for all worker nodes",
				},
				"minsize": resourceschema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "The minimum count of available nodes with type machinetype",
				},
				"maxsize": resourceschema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "The maximum count of available nodes with type machinetype for autoscaling",
				},
				// define default
				"maxsurge": resourceschema.Int64Attribute{
					Required: false,
					Optional: true,
				},
				// define default
				"maxunavailable": resourceschema.Int64Attribute{
					Required: false,
					Optional: true,
				},
			},
		},
		"maintenance": resourceschema.MapAttribute{
			Optional:            true,
			MarkdownDescription: "maintenance options",
			ElementType:         types.StringType,
		},
		"created_at": resourceschema.StringAttribute{
			Computed: true,
		},
		"updated_at": resourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func clusterDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Computed: true,
		},
		"name": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "configurable attribute with default value",
		},
		"project": datasourceschema.StringAttribute{
			Computed: true,
		},
		"partition": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"kubernetes": resourceschema.StringAttribute{
			Optional: true,
		},
		"workers": resourceschema.SingleNestedAttribute{
			Required:            true,
			MarkdownDescription: "Worker settings",
			Attributes: map[string]resourceschema.Attribute{
				"machinetype": resourceschema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The the type of node for all worker nodes",
				},
				"minsize": resourceschema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "The minimum count of available nodes with type machinetype",
				},
				"maxsize": resourceschema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "The maximum count of available nodes with type machinetype for autoscaling",
				},
				// define default
				"maxsurge": resourceschema.Int64Attribute{
					Required: false,
					Optional: true,
				},
				// define default
				"maxunavailable": resourceschema.Int64Attribute{
					Required: false,
					Optional: true,
				},
			},
		},
		"maintenance": resourceschema.MapAttribute{
			Optional:            true,
			MarkdownDescription: "maintenance options",
			ElementType:         types.StringType,
		},
		"created_at": datasourceschema.StringAttribute{
			Computed: true,
		},
		"updated_at": datasourceschema.StringAttribute{
			Computed: true,
		},
	}
}
