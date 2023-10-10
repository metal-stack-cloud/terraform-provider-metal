package cluster

import (
	types "github.com/hashicorp/terraform-plugin-framework/types"
)

type clusterModel struct {
	Uuid       types.String         `tfsdk:"id"`
	Name       types.String         `tfsdk:"name"`
	Project    types.String         `tfsdk:"project"`
	Partition  types.String         `tfsdk:"partition"`
	Tenant     types.String         `tfsdk:"tenant"`
	Kubernetes types.String         `tfsdk:"kubernetes"`
	Workers    []clusterWorkerModel `tfsdk:"workers"`
	// Maintenance maintenanceModel     `tfsdk:"maintenance"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

type clusterWorkerModel struct {
	Name           types.String `tfsdk:"name"`
	MachineType    types.String `tfsdk:"machine_type"`
	Minsize        types.Int64  `tfsdk:"min_size"`
	Maxsize        types.Int64  `tfsdk:"max_size"`
	Maxsurge       types.Int64  `tfsdk:"max_surge"`
	Maxunavailable types.Int64  `tfsdk:"max_unavailable"`
}

// type maintenanceModel struct {
// 	KubernetesAutoupdate   types.Bool `tfsdk:"kubernetes_autoupdate"`
// 	MachineimageAutoupdate types.Bool `tfsdk:"machineimage_autoupdate"`
// 	// Begin                  timestamppb.Timestamp `tfsdk:"begin"`
// 	// Duration               durationpb.Duration   `tfsdk:"duration"`
// }
