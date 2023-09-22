package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

type clusterModel struct {
	Uuid        types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	Project     types.String         `tfsdk:"project"`
	Partition   types.String         `tfsdk:"partition"`
	Kubernetes  types.String         `tfsdk:"kubernetes"`
	Workers     *clusterWorkersModel `tfsdk:"workers"`
	Maintenance *apiv1.Maintenance   `tfsdk:"maintenance"`
	CreatedAt   types.String         `tfsdk:"created_at"`
	UpdatedAt   types.String         `tfsdk:"updated_at"`
}

// type clusterKubernetesModel struct {
// 	Version types.String `tfsdk:"version"`
// }

type clusterWorkersModel struct {
	MachineType    types.String `tfsdk:"machinetype"`
	Minsize        types.Int64  `tfsdk:"minsize"`
	Maxsize        types.Int64  `tfsdk:"maxsize"`
	Maxsurge       types.Int64  `tfsdk:"maxsurge"`
	Maxunavailable types.Int64  `tfsdk:"maxunavailable"`
}

func clusterResponseConvert(clusterP *apiv1.Cluster) clusterModel {
	kubernetesVersion := clusterP.Kubernetes.Version
	// check if workersSlice slice is length 1
	// check for null values
	workersSlice := clusterP.Workers
	workersMapper := clusterWorkersModel{
		MachineType:    types.StringValue(workersSlice[0].MachineType),
		Minsize:        types.Int64Value(int64(workersSlice[0].Minsize)),
		Maxsize:        types.Int64Value(int64(workersSlice[0].Maxsize)),
		Maxsurge:       types.Int64Value(int64(workersSlice[0].Maxsurge)),
		Maxunavailable: types.Int64Value(int64(workersSlice[0].Maxunavailable)),
	}

	return clusterModel{
		Uuid:       types.StringValue(clusterP.Uuid),
		Name:       types.StringValue(clusterP.Name),
		Project:    types.StringValue(clusterP.Project),
		Partition:  types.StringValue(clusterP.Partition),
		Kubernetes: types.StringValue(kubernetesVersion),
		Workers:    &workersMapper,
		CreatedAt:  types.StringValue(clusterP.CreatedAt.AsTime().String()),
		UpdatedAt:  types.StringValue(clusterP.UpdatedAt.AsTime().String()),
	}
}
