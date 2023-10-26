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
			Computed:    true,
			Description: "ID of the cluster",
		},
		"name": resourceschema.StringAttribute{
			Required:    true,
			Description: "This is the name of the cluster that will be used to identify it. It can not be changed afterwards.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(2, 11),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"project": resourceschema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "Project ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"partition": resourceschema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "Partition ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"tenant": resourceschema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "Tenant ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"kubernetes": resourceschema.StringAttribute{
			Required: true,
			Description: `Only newer versions can be specified. There is no downgrade possibility.
			Please be aware that it is not possible to skip major and minor updates.
			It is only possible to upgrade in order. For example from 1.23.3 to 1.24.0, not to 1.25.0.`,
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
			Required:    true,
			Description: "Choose the type of server best suited for your cluster.",
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
						MarkdownDescription: "The machine type for this worker group",
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
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "The maximum count of available nodes which can be updated at once",
					},
					// define default
					"max_unavailable": resourceschema.Int64Attribute{
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "The maximum count of nodes which can be unavailable during node updates",
					},
				},
			},
		},

		"maintenance": resourceschema.SingleNestedAttribute{
			Required:            true,
			MarkdownDescription: "maintenance options",
			Attributes: map[string]resourceschema.Attribute{
				"kubernetes_autoupdate": resourceschema.BoolAttribute{
					Computed:            true,
					Optional:            true,
					MarkdownDescription: "Set kubernetes autoupdate",
				},
				"machineimage_autoupdate": resourceschema.BoolAttribute{
					Computed:            true,
					Optional:            true,
					MarkdownDescription: "Set maschine image autoupdate",
				},
				"time_window": resourceschema.SingleNestedAttribute{
					Required:            true,
					MarkdownDescription: "Set time window for maintenance",
					Attributes: map[string]resourceschema.Attribute{
						"begin": resourceschema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Set begin of maintenance window. Use the format 'HH:MM AM/PM' and consider the UTC offset.",
						},
						"duration": resourceschema.Int64Attribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Set duration of maintenance window. The duration must be defined in hours.",
						},
					},
				},
			},
		},

		"created_at": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Creation timestamp of the cluster",
		},
		"updated_at": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Update timestamp of the cluster",
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
		"workers": resourceschema.ListNestedAttribute{
			Computed:            true,
			MarkdownDescription: "Worker settings",
			NestedObject: resourceschema.NestedAttributeObject{
				Attributes: map[string]resourceschema.Attribute{
					"name": resourceschema.StringAttribute{Computed: true,
						MarkdownDescription: "The name of the worker group.",
					},
					"machine_type": resourceschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The the type of node for all worker nodes",
					},
					"min_size": resourceschema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The minimum count of available nodes with type machinetype",
					},
					"max_size": resourceschema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The maximum count of available nodes with type machinetype for autoscaling",
					},
					// TODO: define default
					"max_surge": resourceschema.Int64Attribute{
						Computed: true,
					},
					// TODO: define default
					"max_unavailable": resourceschema.Int64Attribute{
						Computed: true,
					},
				},
			},
		},
		"maintenance": resourceschema.SingleNestedAttribute{
			Required:            true,
			MarkdownDescription: "maintenance options",
			Attributes: map[string]resourceschema.Attribute{
				"kubernetes_autoupdate": resourceschema.BoolAttribute{
					Computed:            true,
					Optional:            true,
					MarkdownDescription: "Set kubernetes autoupdate",
				},
				"machineimage_autoupdate": resourceschema.BoolAttribute{
					Computed:            true,
					Optional:            true,
					MarkdownDescription: "Set maschine image autoupdate",
				},
				"time_window": resourceschema.SingleNestedAttribute{
					Required:            true,
					MarkdownDescription: "Set time window for maintenance",
					Attributes: map[string]resourceschema.Attribute{
						"begin": resourceschema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Set begin of maintenance window",
						},
						"duration": resourceschema.Int64Attribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Set duration of maintenance window",
						},
					},
				},
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
