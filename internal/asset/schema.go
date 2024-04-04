package asset

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func assetDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"region": datasourceschema.SingleNestedAttribute{
			Computed:    true,
			Description: "The location of a datacenter.",
			Attributes: map[string]datasourceschema.Attribute{
				"id": datasourceschema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The id of the region.",
				},
				"name": datasourceschema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The name of the region.",
				},
				"address": datasourceschema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The address of the region.",
				},
				"active": datasourceschema.BoolAttribute{
					Computed:            true,
					MarkdownDescription: "Indicates if the region is usable.",
				},
				"partitions": datasourceschema.ListNestedAttribute{
					Computed:            true,
					MarkdownDescription: "Available partitions in this region",
					NestedObject: datasourceschema.NestedAttributeObject{
						Attributes: map[string]datasourceschema.Attribute{
							"id": datasourceschema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The id of the partition.",
							},
							"name": datasourceschema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The name of the partition.",
							},
							"address": datasourceschema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The address of the partition.",
							},
							"active": datasourceschema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates if the partition is usable.",
							},
							"description": datasourceschema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Description of the partition.",
							},
						},
					},
				},
				"defaults": datasourceschema.SingleNestedAttribute{
					Computed:            true,
					MarkdownDescription: "The defaults for assets, if not specified otherwise.",
					Attributes: map[string]datasourceschema.Attribute{
						"machine_type": datasourceschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The default machine type used.",
						},
						"kubernetes_version": datasourceschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The default Kubernetes version used.",
						},
						"worker_min": datasourceschema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The minimum servers specified.",
						},
						"worker_max": datasourceschema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The maximum servers specified.",
						},
						"partition": datasourceschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The partition where the cluster is created by default.",
						},
					},
				},
				"description": datasourceschema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Description of the region.",
				},
			},
		},
		"machine_types": datasourceschema.ListNestedAttribute{
			Computed:    true,
			Description: "The machine types available in a region.",
			NestedObject: datasourceschema.NestedAttributeObject{
				Attributes: map[string]datasourceschema.Attribute{
					"id": datasourceschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The id of the machine type.",
					},
					"name": datasourceschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the machine type.",
					},
					"cpus": datasourceschema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "CPUs in this machine.",
					},
					"memory": datasourceschema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Memory in this machine.",
					},
					"storage": datasourceschema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Storage in this machine.",
					},
					"cpu_description": datasourceschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The description of the CPUs in this machine.",
					},
					"storage_description": datasourceschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The description of the disks in this machine.",
					},
				},
			},
		},
		"kubernetes": datasourceschema.ListNestedAttribute{
			Computed:    true,
			Description: "The list of supported Kubernetes versions.",
			NestedObject: datasourceschema.NestedAttributeObject{
				Attributes: map[string]datasourceschema.Attribute{
					"version": datasourceschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version of Kubernetes.",
					},
					// "expiration": necessary?
				},
			},
		},
	}
}
