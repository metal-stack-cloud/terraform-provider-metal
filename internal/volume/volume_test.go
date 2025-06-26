package volume

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/stretchr/testify/assert"
)

func Test_volumeFromApi(t *testing.T) {
	vol := &apiv1.Volume{
		Uuid:         "ea5ba51d-11ca-4fae-8494-aa78523cbe26",
		Name:         "test-volume",
		Size:         100,
		Project:      "project-a",
		Partition:    "partition-a",
		StorageClass: "default",
		ReplicaCount: 2,
		ClusterName:  "my-cluster",
		Labels: []*apiv1.VolumeLabel{
			{
				Key:   "hello",
				Value: "world",
			},
			{
				Key:   "foo",
				Value: "bar",
			},
		},
	}

	want := volumeModel{
		Uuid:         types.StringValue("ea5ba51d-11ca-4fae-8494-aa78523cbe26"),
		Name:         types.StringValue("test-volume"),
		Project:      types.StringValue("project-a"),
		Partition:    types.StringValue("partition-a"),
		StorageClass: types.StringValue("default"),
		ReplicaCount: types.Int64Value(2),
		ClusterName:  types.StringValue("my-cluster"),
		Labels: types.MapValueMust(basetypes.StringType{},
			map[string]attr.Value{
				"hello": types.StringValue("world"),
				"foo":   types.StringValue("bar"),
			},
		),
	}

	got := volumeResponseMapping(vol)
	assert.Equal(t, want, got)
}
