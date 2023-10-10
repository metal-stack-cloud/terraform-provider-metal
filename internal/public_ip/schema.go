package ipaddress

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func publicIpDataSourceAttributes() map[string]dataschema.Attribute {
	return map[string]dataschema.Attribute{
		"id": dataschema.StringAttribute{
			Computed:    true,
			Description: "The ID that represents this public IP address.",
		},
		"ip": dataschema.StringAttribute{
			Computed:    true,
			Description: "The publicly accessible IP address.",
		},
		"name": dataschema.StringAttribute{
			Computed:    true,
			Description: "You can give your IP address a freely chosen name to identify it in the future.",
		},
		"description": dataschema.StringAttribute{
			Computed:    true,
			Description: "Here you can give your IP an optional description for your own use.",
		},
		"network": dataschema.StringAttribute{
			Computed:    true,
			Description: "The network this address is bound to.",
		},
		"project": dataschema.StringAttribute{
			Computed:    true,
			Description: "The project this address is part of. Cannot be moved.",
		},
		"type": dataschema.StringAttribute{
			Computed: true,
			Description: `Determines the type of the public ip address. 
	If you want the IP to outlive the cluster lifecycle, mark it as static. Otherwise it will be deleted along with the cluster. 
	Another use case would be if you want to have a stable egress address on the internet gateway for your cluster.
			`,
		},
		"tags": dataschema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "The tags used to organize this address.",
		},
		"created_at": dataschema.StringAttribute{
			Computed:    true,
			Description: "Indicates when this IP address has initially been claimed.",
		},
		"updated_at": dataschema.StringAttribute{
			Computed:    true,
			Description: "Indicates when this IP address has been updated.",
		},
	}
}

func publicIpResourceAttributes() map[string]resourceschema.Attribute {
	return map[string]resourceschema.Attribute{
		"id": resourceschema.StringAttribute{
			Computed:    true,
			Description: "The ID that represents this public IP address.",
		},
		"ip": resourceschema.StringAttribute{
			Computed:    true,
			Description: "The publicly accessible IP address.",
		},
		"name": resourceschema.StringAttribute{
			Required:    true,
			Description: "You can give your IP address a freely chosen name to identify it in the future.",
			Validators: []validator.String{
				stringvalidator.LengthAtMost(32),
			},
		},
		"description": resourceschema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			Description: "Here you can give your IP an optional description for your own use.",
		},
		"network": resourceschema.StringAttribute{
			Computed:    true,
			Description: "The network this address is bound to.",
		},
		"project": resourceschema.StringAttribute{
			Computed:    true,
			Description: "The project this address is part of. Cannot be moved.",
		},
		"type": resourceschema.StringAttribute{
			Computed: true,
			Default:  stringdefault.StaticString("ephemeral"),
			Description: `Determines the type of the public ip address. 
	If you want the IP to outlive the cluster lifecycle, mark it as static. Otherwise it will be deleted along with the cluster. 
	Another use case would be if you want to have a stable egress address on the internet gateway for your cluster.
			`,
			Validators: []validator.String{
				stringvalidator.OneOf("ephemeral", "static"),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, sr planmodifier.StringRequest, rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse) {
					rrifr.RequiresReplace = sr.StateValue.ValueString() == "static" &&
						sr.PlanValue.ValueString() == "ephemeral"
				}, "desc", "mddesc"),
			},
		},
		"tags": resourceschema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			Description: "The tags used to organize this address.",
		},
		"created_at": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Indicates when this IP address has initially been claimed.",
		},
		"updated_at": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Indicates when this IP address has been updated.",
		},
	}
}
