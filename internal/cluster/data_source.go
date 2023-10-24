package cluster

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
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
		Attributes:          clusterDataSourceAttributes(),
		Description:         "Allows querying a specific cluster that already exists and is not yet managed.",
		MarkdownDescription: "Allows querying a specific cluster that already exists and is not yet managed. Either `id` or `project` and `name` are required.",
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

	// set project
	var project string
	if data.Project.ValueString() == "" {
		project = clusterP.session.Project
	}

	// get all clusters and select cluster by name if uuid is not set
	var uuidString string
	if data.Uuid.ValueString() == "" {
		listRequestMessage := &apiv1.ClusterServiceListRequest{
			Project: project,
		}
		// get clusterList type Clusters []*Cluster
		clusterList, err := clusterP.session.Client.Apiv1().Cluster().List(ctx, connect.NewRequest(listRequestMessage))
		if err != nil {
			response.Diagnostics.AddError("Failed to get cluster list", err.Error())
			return
		}
		// find uuid and set uuidString
		list := clusterList.Msg.Clusters
		uuidStr, err := findUuidByName(list, data.Name.ValueString())
		if err != nil {
			response.Diagnostics.AddError(fmt.Sprintf("Failed to find cluster with name %v", data.Name.ValueString()), err.Error())
			return
		} else {
			uuidString = uuidStr
		}
	} else {
		uuidString = data.Uuid.ValueString()
	}

	// get Cluster by uuid
	getRequestMessage := &apiv1.ClusterServiceGetRequest{
		Uuid:    uuidString,
		Project: project,
	}
	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Get(ctx, connect.NewRequest(getRequestMessage))
	if err != nil {
		response.Diagnostics.AddError("Failed to get cluster", err.Error())
		return
	}

	// save updated data into terraform state
	state := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(state...)
}

func findUuidByName(list []*apiv1.Cluster, name string) (string, error) {
	for _, e := range list {
		if e.Name == name {
			return e.Uuid, nil
		}
	}
	return "", fmt.Errorf("cluster name not found in list")
}
