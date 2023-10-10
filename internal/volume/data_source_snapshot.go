package volume

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
	_ datasource.DataSource              = &SnapshotDataSource{}
	_ datasource.DataSourceWithConfigure = &SnapshotDataSource{}
)

func NewSnapshotDataSource() datasource.DataSource {
	return &SnapshotDataSource{}
}

type SnapshotDataSource struct {
	session *session.Session
}

// Metadata implements datasource.datasource.
func (*SnapshotDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_snapshot"
}

// Schema implements datasource.datasource.
func (*SnapshotDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: SnapshotDataSourceAttributes(),
	}
}

// Configure implements datasource.ResourceWithConfigure.
func (snapshotP *SnapshotDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

	snapshotP.session = client
}

// Read implements datasource.datasource.
func (snapshotP *SnapshotDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data snapshotModel
	diagState := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diagState...)
	if response.Diagnostics.HasError() {
		return
	}

	// set project
	var project string
	if data.Project.ValueString() == "" {
		project = snapshotP.session.Project
	}

	// get all snapshots and select snapshot by name if uuid is not set
	var uuidString string
	if data.Uuid.ValueString() == "" {
		listRequestMessage := &apiv1.SnapshotServiceListRequest{
			Project: project,
		}
		// get snapshotList type snapshots []*snapshot
		snapshotList, err := snapshotP.session.Client.Apiv1().Snapshot().List(ctx, connect.NewRequest(listRequestMessage))
		if err != nil {
			response.Diagnostics.AddError("Failed to get snapshot list", err.Error())
			return
		}
		// find uuid and set uuidString
		list := snapshotList.Msg.Snapshots
		returnString, err := findSnapshotUuid(list, data.Name.ValueString())
		if err != nil {
			response.Diagnostics.AddError(fmt.Sprintf("Failed to find snapshot with name %v", data.Name.ValueString()), err.Error())
			return
		} else {
			uuidString = returnString
		}
	} else {
		uuidString = data.Uuid.ValueString()
	}

	// get snapshot by uuid
	requestMessage := &apiv1.SnapshotServiceGetRequest{
		Uuid:    uuidString,
		Project: project,
	}
	clientResponse, err := snapshotP.session.Client.Apiv1().Snapshot().Get(ctx, connect.NewRequest(requestMessage))
	if err != nil {
		response.Diagnostics.AddError("Failed to get snapshot", err.Error())
		return
	}

	// save updated data into terraform state
	state := response.State.Set(ctx, snapshotResponseMapping(clientResponse.Msg.Snapshot))
	response.Diagnostics.Append(state...)
}

func findSnapshotUuid(list []*apiv1.Snapshot, name string) (string, error) {
	for _, e := range list {
		if e.Name == name {
			return e.Uuid, nil
		}
	}
	return "", fmt.Errorf("snapshot name not found in list")
}

// TODO: actually use.
// func findSnapshotSourceVolumeUuid(list []*apiv1.Snapshot, name string) (string, error) {
// 	for _, e := range list {
// 		if e.Name == name {
// 			return e.Uuid, nil
// 		}
// 	}
// 	return "", fmt.Errorf("snapshot name not found in list")
// }
