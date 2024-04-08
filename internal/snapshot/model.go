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

func snapshotResponseMapping(s *apiv1.Snapshot) snapshotModel {
	return snapshotModel{
		Uuid:             types.StringValue(s.Uuid),
		Name:             types.StringValue(s.Name),
		Project:          types.StringValue(s.Project),
		Partition:        types.StringValue(s.Partition),
		StorageClass:     types.StringValue(s.StorageClass),
		Size:             types.Int64Value(int64(s.Size)),
		Usage:            types.Int64Value(int64(s.Usage)),
		SourceVolumeUuid: types.StringValue(s.SourceVolumeUuid),
	}
}
