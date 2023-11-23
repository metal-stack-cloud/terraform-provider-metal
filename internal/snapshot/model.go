package snapshot

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

type snapshotModel struct {
	Uuid             types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Project          types.String `tfsdk:"project"`
	Partition        types.String `tfsdk:"partition"`
	StorageClass     types.String `tfsdk:"storage_class"`
	Size             types.Int64  `tfsdk:"size"`
	Usage            types.Int64  `tfsdk:"usage"`
	SourceVolumeUuid types.String `tfsdk:"volume_id"`
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
