package ipaddress

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ datasource.DataSource              = &IpDataSource{}
	_ datasource.DataSourceWithConfigure = &IpDataSource{}
)

func NewIpDataSource() datasource.DataSource {
	return &IpDataSource{}
}

// IpDataSource defines the data source implementation.
type IpDataSource struct {
	session *session.Session
}

// Metadata implements datasource.DataSource.
func (d *IpDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_addresses"
}

// Schema implements datasource.DataSource.
func (d *IpDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "IP Address data source",
		Attributes: map[string]schema.Attribute{
			"list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: ipAddressDataSourceAttributes(),
				},
			},
		},
	}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *IpDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ipAddressListDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ipResp, err := d.session.Client.Apiv1().IP().List(ctx, connect.NewRequest(&apiv1.IPServiceListRequest{
		Project: d.session.Project,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to read IP Addresses", err.Error())
		return
	}
	tflog.Trace(ctx, "read ip addresses")

	for _, ip := range ipResp.Msg.Ips {
		data.List = append(data.List, ipAddressFromApi(ip))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
