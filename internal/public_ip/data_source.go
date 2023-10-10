package ipaddress

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ datasource.DataSource              = &PublicIpDataSource{}
	_ datasource.DataSourceWithConfigure = &PublicIpDataSource{}
)

func NewPublicIpDataSource() datasource.DataSource {
	return &PublicIpDataSource{}
}

// PublicIpDataSource defines the data source implementation.
type PublicIpDataSource struct {
	session *session.Session
}

// Metadata implements datasource.DataSource.
func (d *PublicIpDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ips"
}

// Schema implements datasource.DataSource.
func (d *PublicIpDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Each cluster gets an IP automatically provided on the internet gateway for outgoing communication.
Services get an IP automatically on creation.
Services and gateway IPs are dynamic by default.
You can use an IP address in several clusters and locations at the same time.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All public IP addresses",
				NestedObject: schema.NestedAttributeObject{
					Attributes: publicIpDataSourceAttributes(),
				},
			},
		},
	}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *PublicIpDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	session, ok := req.ProviderData.(*session.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *session.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.session = session
}

func (d *PublicIpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicIpListDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ipResp, err := d.session.Client.Apiv1().IP().List(ctx, connect.NewRequest(&apiv1.IPServiceListRequest{
		Project: d.session.Project,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to read public IP Addresses", err.Error())
		return
	}
	tflog.Trace(ctx, "read public ip addresses")

	data.Items = make([]publicIpModel, 0, len(ipResp.Msg.Ips))
	ids := make([]string, 0, len(ipResp.Msg.Ips))
	for _, ip := range ipResp.Msg.Ips {
		data.Items = append(data.Items, publicIpFromApi(ip))
		ids = append(ids, ip.Ip)
	}

	dataId := fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(ids, ""))))
	data.ContentId = types.StringValue(dataId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
