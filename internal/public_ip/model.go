package ipaddress

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

// PublicIpListDataSourceModel describes the data source data model.
type PublicIpListDataSourceModel struct {
	ContentId types.String    `tfsdk:"id"`
	Items     []publicIpModel `tfsdk:"items"`
}

type publicIpModel struct {
	Uuid        types.String   `tfsdk:"id"`
	Ip          types.String   `tfsdk:"ip"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Network     types.String   `tfsdk:"network"`
	Project     types.String   `tfsdk:"project"`
	Type        types.String   `tfsdk:"type"` // TODO: make enum; unspecified, ephemeral, static
	Tags        []types.String `tfsdk:"tags"`
	CreatedAt   types.String   `tfsdk:"created_at"`
	UpdatedAt   types.String   `tfsdk:"updated_at"`
}

func publicIpFromApi(ip *apiv1.IP) publicIpModel {
	ipType := "unspecified"
	switch ip.Type {
	case apiv1.IPType_IP_TYPE_STATIC:
		ipType = "static"
	case apiv1.IPType_IP_TYPE_EPHEMERAL:
		ipType = "ephemeral"
	}
	tags := make([]types.String, len(ip.Tags))
	for i, tag := range ip.Tags {
		tags[i] = types.StringValue(tag)
	}
	return publicIpModel{
		Uuid:        types.StringValue(ip.Uuid),
		Ip:          types.StringValue(ip.Ip),
		Name:        types.StringValue(ip.Name),
		Description: types.StringValue(ip.Description),
		Network:     types.StringValue(ip.Network),
		Project:     types.StringValue(ip.Project),
		Type:        types.StringValue(ipType),
		Tags:        tags,
		CreatedAt:   types.StringValue(ip.CreatedAt.AsTime().String()),
		UpdatedAt:   types.StringValue(ip.UpdatedAt.AsTime().String()),
	}
}
