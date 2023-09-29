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

type snapshotModel struct {
	Uuid             types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Project          types.String `tfsdk:"project"`
	Partition        types.String `tfsdk:"partition"`
	StorageClass     types.String `tfsdk:"storageclass"`
	Size             types.Int64  `tfsdk:"size"`
	Usage            types.Int64  `tfsdk:"usage"`
	SourceVolumeUuid types.String `tfsdk:"volume_id"`
}

func volumeResponseMapping(volumeP *apiv1.Volume) volumeModel {
	return volumeModel{
		Uuid:         types.StringValue(volumeP.Uuid),
		Name:         types.StringValue(volumeP.Name),
		Project:      types.StringValue(volumeP.Project),
		Partition:    types.StringValue(volumeP.Partition),
		StorageClass: types.StringValue(volumeP.StorageClass),
		ReplicaCount: types.Int64Value(int64(volumeP.ReplicaCount)),
	}
}

func snapshotResponseMapping(snapshotP *apiv1.Snapshot) snapshotModel {
	return snapshotModel{
		Uuid:             types.StringValue(snapshotP.Uuid),
		Name:             types.StringValue(snapshotP.Name),
		Project:          types.StringValue(snapshotP.Project),
		Partition:        types.StringValue(snapshotP.Partition),
		StorageClass:     types.StringValue(snapshotP.StorageClass),
		Size:             types.Int64Value(int64(snapshotP.Size)),
		Usage:            types.Int64Value(int64(snapshotP.Usage)),
		SourceVolumeUuid: types.StringValue(snapshotP.SourceVolumeUuid),
	}
}
