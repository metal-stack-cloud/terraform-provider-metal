package ipaddress

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	assert "github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_publicIpFromApi(t *testing.T) {
	ip := &apiv1.IP{
		Uuid:        "1",
		Ip:          "1.2.3",
		Name:        "ip",
		Description: "Test ip",
		Network:     "internet",
		Project:     "default-project",
		Type:        apiv1.IPType(1),
		Tags: []string{
			"tag-1",
			"tag-2",
		},
		CreatedAt: &timestamppb.Timestamp{
			Seconds: int64(1707382100),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: int64(1717932877),
		},
	}
	want := publicIpModel{
		Uuid:        basetypes.NewStringValue("1"),
		Ip:          basetypes.NewStringValue("1.2.3"),
		Name:        basetypes.NewStringValue("ip"),
		Description: basetypes.NewStringValue("Test ip"),
		Network:     basetypes.NewStringValue("internet"),
		Project:     basetypes.NewStringValue("default-project"),
		Type:        basetypes.NewStringValue("ephemeral"),
		Tags: []basetypes.StringValue{
			basetypes.NewStringValue("tag-1"),
			basetypes.NewStringValue("tag-2"),
		},
		CreatedAt: basetypes.NewStringValue("2024-02-08 08:48:20 +0000 UTC"),
		UpdatedAt: basetypes.NewStringValue("2024-06-09 11:34:37 +0000 UTC"),
	}

	ipModel := publicIpFromApi(ip)
	assert.Equal(t, want, ipModel)
}
