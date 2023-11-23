package kubeconfig

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	_ datasource.DataSource              = &KubeconfigDataSource{}
	_ datasource.DataSourceWithConfigure = &KubeconfigDataSource{}
)

func NewKubeconfigDataSource() datasource.DataSource {
	return &KubeconfigDataSource{}
}

type KubeconfigDataSource struct {
	session *session.Session
}

// Metadata implements datasource.DataSource.
func (d *KubeconfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubeconfig"
}

// Schema implements datasource.DataSource.
func (d *KubeconfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: ``,
		Attributes:  kubeconfigDataSourceAttributes(),
	}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *KubeconfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KubeconfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data kubeconfigDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	project := data.Project.ValueString()
	if data.Project.IsNull() {
		project = d.session.Project
	}

	expiration, err := time.ParseDuration(data.Expiration.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("expiration"), "invalid duration", "must be of form `1h30m`")
		return
	}

	kcResp, err := d.session.Client.Apiv1().Cluster().GetCredentials(ctx, connect.NewRequest(&apiv1.ClusterServiceGetCredentialsRequest{
		Project:    project,
		Uuid:       data.Uuid.ValueString(),
		Expiration: durationpb.New(expiration),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to generate kubeconfig", err.Error())
		return
	}
	tflog.Trace(ctx, "generated kubeconfig")

	data.Raw = types.StringValue(kcResp.Msg.GetKubeconfig())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
