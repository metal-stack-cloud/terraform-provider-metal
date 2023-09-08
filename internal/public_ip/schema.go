package ipaddress

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func publicIpResourceAttributes() map[string]resourceschema.Attribute {
	return map[string]resourceschema.Attribute{
		"id": resourceschema.StringAttribute{
			Computed: true,
		},
		"ip": resourceschema.StringAttribute{
			Computed: true,
		},
		"name": resourceschema.StringAttribute{
			Required: true,
		},
		"description": resourceschema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(""),
		},
		"network": resourceschema.StringAttribute{
			Computed: true,
		},
		"project": resourceschema.StringAttribute{
			Computed: true,
		},
		"type": resourceschema.StringAttribute{
			Computed: true,
			Default:  stringdefault.StaticString("ephemeral"),
			Validators: []validator.String{
				stringvalidator.OneOf("ephemeral", "static"),
			},
		},
		"tags": resourceschema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
		},
		"created_at": resourceschema.StringAttribute{
			Computed: true,
		},
		"updated_at": resourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func publicIpDataSourceAttributes() map[string]dataschema.Attribute {
	return map[string]dataschema.Attribute{
		"id": dataschema.StringAttribute{
			Computed: true,
		},
		"ip": dataschema.StringAttribute{
			Computed: true,
		},
		"name": dataschema.StringAttribute{
			Computed: false,
			Required: true,
		},
		"description": dataschema.StringAttribute{
			Computed: true,
		},
		"network": dataschema.StringAttribute{
			Computed: true,
		},
		"project": dataschema.StringAttribute{
			Computed: true,
		},
		"type": dataschema.StringAttribute{
			Computed: true,
		},
		"tags": dataschema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"created_at": dataschema.StringAttribute{
			Computed: true,
		},
		"updated_at": dataschema.StringAttribute{
			Computed: true,
		},
	}
}
