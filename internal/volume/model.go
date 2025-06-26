package volume

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

type volumeModel struct {
	Uuid         types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Project      types.String `tfsdk:"project"`
	Partition    types.String `tfsdk:"partition"`
	StorageClass types.String `tfsdk:"storageclass"`
	ReplicaCount types.Int64  `tfsdk:"replicacount"`
	ClusterName  types.String `tfsdk:"clustername"`
	Labels       types.Map    `tfsdk:"labels"`
}

func volumeResponseMapping(v *apiv1.Volume) volumeModel {
	m := volumeModel{
		Uuid:         types.StringValue(v.Uuid),
		Name:         types.StringValue(v.Name),
		Project:      types.StringValue(v.Project),
		Partition:    types.StringValue(v.Partition),
		StorageClass: types.StringValue(v.StorageClass),
		ReplicaCount: types.Int64Value(int64(v.ReplicaCount)),
		ClusterName:  types.StringValue(v.ClusterName),
	}

	labels := make(map[string]attr.Value)
	for _, l := range v.Labels {
		labels[l.Key] = types.StringValue(l.Value)
	}
	m.Labels = types.MapValueMust(types.StringType, labels)

	return m
}
