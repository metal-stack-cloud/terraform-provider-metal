package cluster

import (
	"context"
	"fmt"

	connect_go "github.com/bufbuild/connect-go"
	datasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	session "github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ datasource.DataSource              = &ClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &ClusterDataSource{}
)

func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

type ClusterDataSource struct {
	session *session.Session
}

// Metadata implements datasource.datasource.
func (*ClusterDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cluster"
}

// Schema implements datasource.datasource.
func (*ClusterDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: clusterDataSourceAttributes(),
	}
}

// Configure implements datasource.ResourceWithConfigure.
func (clusterP *ClusterDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*session.Session)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *session.Session, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	clusterP.session = client
}

// Read implements datasource.datasource.
func (clusterP *ClusterDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data clusterModel
	diagState := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diagState...)

	if response.Diagnostics.HasError() {
		return
	}

	requestMessage := &apiv1.ClusterServiceGetRequest{
		Uuid:    data.Uuid.ValueString(),
		Project: data.Project.ValueString(),
	}

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Get(ctx, connect_go.NewRequest(requestMessage))

	if err != nil {
		response.Diagnostics.AddError("Failed to get cluster", err.Error())
		return
	}

	// Save updated data into Terraform state
	state := response.State.Set(ctx, clusterResponseConvert(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(state...)
}
