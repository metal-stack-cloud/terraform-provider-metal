package volume

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
	_ datasource.DataSource              = &VolumeDataSource{}
	_ datasource.DataSourceWithConfigure = &VolumeDataSource{}
)

func NewVolumeDataSource() datasource.DataSource {
	return &VolumeDataSource{}
}

type VolumeDataSource struct {
	session *session.Session
}

// Metadata implements datasource.datasource.
func (*VolumeDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_volume"
}

// Schema implements datasource.datasource.
func (*VolumeDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: VolumeDataSourceAttributes(),
	}
}

// Configure implements datasource.ResourceWithConfigure.
func (volumeP *VolumeDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

	volumeP.session = client
}

// Read implements datasource.datasource.
func (volumeP *VolumeDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data volumeModel
	diagState := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diagState...)
	if response.Diagnostics.HasError() {
		return
	}

	// set project
	var project string
	if data.Project.ValueString() == "" {
		project = volumeP.session.Project
	}

	// get all volumes and select volume by name if uuid is not set
	var uuidString string
	if data.Uuid.ValueString() == "" {
		listRequestMessage := &apiv1.VolumeServiceListRequest{
			Project: project,
		}
		// get volumeList type volumes []*volume
		volumeList, err := volumeP.session.Client.Apiv1().Volume().List(ctx, connect_go.NewRequest(listRequestMessage))
		if err != nil {
			response.Diagnostics.AddError("Failed to get volume list", err.Error())
			return
		}
		// find uuid and set uuidString
		list := volumeList.Msg.Volumes
		returnString, err := findUuid(list, data.Name.ValueString())
		if err != nil {
			response.Diagnostics.AddError(fmt.Sprintf("Failed to find volume with name %v", data.Name.ValueString()), err.Error())
			return
		} else {
			uuidString = returnString
		}
	} else {
		uuidString = data.Uuid.ValueString()
	}

	// get volume by uuid
	getRequestMessage := &apiv1.VolumeServiceGetRequest{
		Uuid:    uuidString,
		Project: project,
	}
	clientResponse, err := volumeP.session.Client.Apiv1().Volume().Get(ctx, connect_go.NewRequest(getRequestMessage))
	if err != nil {
		response.Diagnostics.AddError("Failed to get volume", err.Error())
		return
	}

	// save updated data into terraform state
	state := response.State.Set(ctx, volumeResponseConvert(clientResponse.Msg.Volume))
	response.Diagnostics.Append(state...)
}

func findUuid(list []*apiv1.Volume, name string) (string, error) {
	for _, e := range list {
		if e.Name == name {
			return e.Uuid, nil
		}
	}
	return "", fmt.Errorf("volume name not found in list")
}
