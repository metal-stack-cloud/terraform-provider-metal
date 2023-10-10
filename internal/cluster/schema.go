package cluster

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func clusterResourceAttributes() map[string]resourceschema.Attribute {
	return map[string]resourceschema.Attribute{
		"id": resourceschema.StringAttribute{
			Computed: true,
		},
		"name": resourceschema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(2, 11),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"project": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"partition": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"tenant": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"kubernetes": resourceschema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.LengthAtMost(8),
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[0-9]+.[0-9]+.[0-9]+$`), "wrong version pattern",
				),
			},
		},

		"workers": resourceschema.ListNestedAttribute{
			Required:            true,
			MarkdownDescription: "Worker groups settings",
			NestedObject: resourceschema.NestedAttributeObject{
				Attributes: map[string]resourceschema.Attribute{
					"name": resourceschema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The group name of the worker nodes",
						Validators: []validator.String{
							stringvalidator.LengthBetween(2, 128),
						},
					},
					"machine_type": resourceschema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The the type of node for all worker nodes",
					},
					"min_size": resourceschema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "The minimum count of available nodes with type machinetype",
					},
					"max_size": resourceschema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "The maximum count of available nodes with type machinetype for autoscaling",
					},
					// define default
					"max_surge": resourceschema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The maximum count of available nodes which can be updated at once",
					},
					// define default
					"max_unavailable": resourceschema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The maximum count of nodes which can be unavailable during node updates",
					},
				},
			},
		},

		// "maintenance": resourceschema.SingleNestedAttribute{
		// 	Optional:            true,
		// 	MarkdownDescription: "maintenance options",
		// 	PlanModifiers: []planmodifier.Object{
		// 		objectplanmodifier.UseStateForUnknown(),
		// 	},
		// 	Attributes: map[string]resourceschema.Attribute{
		// 		"kubernetes_autoupdate": resourceschema.BoolAttribute{
		// 			Computed:            true,
		// 			MarkdownDescription: "Set kubernetes autoupdate",
		// 		},
		// 		"machineimage_autoupdate": resourceschema.BoolAttribute{
		// 			Computed:            true,
		// 			MarkdownDescription: "Set maschine image autoupdate",
		// 		},
		// 		// "begin": resourceschema.Int64Attribute{
		// 		// 	Optional:            true,
		// 		// 	MarkdownDescription: "Set begin of maintenance window",
		// 		// },
		// 		// "duration": resourceschema.Int64Attribute{
		// 		// 	Optional:            true,
		// 		// 	MarkdownDescription: "Set duration of maintenance window",
		// 		// },
		// 	},
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
			Computed:            true,
			Optional:            true,
			MarkdownDescription: "ID of the cluster",
		},
		"name": datasourceschema.StringAttribute{
			Computed:            true,
			Optional:            true,
			MarkdownDescription: "Name of the cluster",
		},
		"project": datasourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"partition": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"tenant": resourceschema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"kubernetes": resourceschema.StringAttribute{
			Computed: true,
		},
		"workers": resourceschema.SingleNestedAttribute{
			Computed:            true,
			MarkdownDescription: "Worker settings",
			Attributes: map[string]resourceschema.Attribute{
				"machinetype": resourceschema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The the type of node for all worker nodes",
				},
				"minsize": resourceschema.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "The minimum count of available nodes with type machinetype",
				},
				"maxsize": resourceschema.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "The maximum count of available nodes with type machinetype for autoscaling",
				},
				// define default
				"maxsurge": resourceschema.Int64Attribute{
					Computed: true,
				},
				// define default
				"maxunavailable": resourceschema.Int64Attribute{
					Computed: true,
				},
			},
		},
		"maintenance": resourceschema.SingleNestedAttribute{
			Computed:            true,
			MarkdownDescription: "maintenance options",
			Attributes: map[string]resourceschema.Attribute{
				"kubernetesautoupdate": resourceschema.BoolAttribute{
					Computed:            true,
					MarkdownDescription: "Set kubernetes autoupdate",
				},
				"machineimageautoupdate": resourceschema.BoolAttribute{
					Computed:            true,
					MarkdownDescription: "Set maschine image autoupdate",
				},
				// "begin": resourceschema.Int64Attribute{
				// 	MarkdownDescription: "Set begin of maintenance window",
				// },
				// "duration": resourceschema.Int64Attribute{
				// 	MarkdownDescription: "Set duration of maintenance window",
				// },
			},
		},
		"created_at": datasourceschema.StringAttribute{
			Computed: true,
		},
		"updated_at": datasourceschema.StringAttribute{
			Computed: true,
		},
	}
}
