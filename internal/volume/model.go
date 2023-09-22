package volume

import (
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
}

func volumeResponseConvert(volumeP *apiv1.Volume) volumeModel {
	return volumeModel{
		Uuid:         types.StringValue(volumeP.Uuid),
		Name:         types.StringValue(volumeP.Name),
		Project:      types.StringValue(volumeP.Project),
		Partition:    types.StringValue(volumeP.Partition),
		StorageClass: types.StringValue(volumeP.StorageClass),
		ReplicaCount: types.Int64Value(int64(volumeP.ReplicaCount)),
	}
}
