package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

type clusterModel struct {
	Uuid    types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Project types.String `tfsdk:"project"`
	// Partition  types.String          `tfsdk:"partition"`
	Kubernetes  *apiv1.KubernetesSpec `tfsdk:"kubernetes"`
	Workers     clusterWorkersModel   `tfsdk:"workers"`
	Maintenance *apiv1.Maintenance    `tfsdk:"maintenance"`
	CreatedAt   types.String          `tfsdk:"created_at"`
	UpdatedAt   types.String          `tfsdk:"updated_at"`
	// Tags        []types.String `tfsdk:"tags"`
}

type clusterWorkersModel struct {
	MachineType    types.String `tfsdk:"MachineType"`
	Minsize        types.Int64  `tfsdk:"Minsize"`
	Maxsize        types.Int64  `tfsdk:"Maxsize"`
	Maxsurge       types.Int64  `tfsdk:"Maxsurge"`
	Maxunavailable types.Int64  `tfsdk:"Maxunavailable"`
}

func clusterResponseConvert(clusterPointer *apiv1.Cluster) clusterModel {
	// tags := make([]types.String, len(clusterPointer.Tags))
	// for i, tag := range clusterPointer.Tags {
	// 	tags[i] = types.StringValue(tag)
	// }

	return clusterModel{
		Uuid: types.StringValue(clusterPointer.Uuid),
		// Ip:          types.StringValue(clusterPointer.Ip),
		Name: types.StringValue(clusterPointer.Name),
		// Network:     types.StringValue(clusterPointer.Network),
		Project: types.StringValue(clusterPointer.Project),
		// Tags:      tags,
		CreatedAt: types.StringValue(clusterPointer.CreatedAt.AsTime().String()),
		UpdatedAt: types.StringValue(clusterPointer.UpdatedAt.AsTime().String()),
	}
}
