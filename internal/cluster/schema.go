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
		"name": resourceschema.StringAttribute{
			Required: true,
		},
		"project": resourceschema.StringAttribute{
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
				"MachineType": resourceschema.StringAttribute{
					Required: true,
				},
				"Minsize": resourceschema.Int64Attribute{
					Required: true,
				},
				"Maxsize": resourceschema.Int64Attribute{
					Required: true,
				},
				"Maxsurge": resourceschema.Int64Attribute{
					Required: false,
					Optional: true,
				},
				"Maxunavailable": resourceschema.Int64Attribute{
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
		// "type": resourceschema.StringAttribute{
		// 	Computed: true,
		// 	Default:  stringdefault.StaticString("ephemeral"),
		// 	Validators: []validator.String{
		// 		stringvalidator.OneOf("ephemeral", "static"),
		// 	},
		// },
		// "tags": resourceschema.ListAttribute{
		// 	Computed:    true,
		// 	ElementType: types.StringType,
		// 	Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
		// },
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
		// "tags": datasourceschema.ListAttribute{
		// 	Computed:    true,
		// 	ElementType: types.StringType,
		// },
		"created_at": datasourceschema.StringAttribute{
			Computed: true,
		},
		"updated_at": datasourceschema.StringAttribute{
			Computed: true,
		},
	}
}
