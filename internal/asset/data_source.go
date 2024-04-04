package asset

import (
	"context"
	"fmt"

	connect "connectrpc.com/connect"
	datasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	session "github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ datasource.DataSource = &AssetDataSource{}
)

func NewAssetDataSource() datasource.DataSource {
	return &AssetDataSource{}
}

type AssetDataSource struct {
	session *session.Session
}

func (*AssetDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_asset"
}

func (*AssetDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Show the available assets.",
		MarkdownDescription: "Show the available assets.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of assets.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: assetDataSourceAttributes(),
				},
			},
		},
	}
}

func (assetP *AssetDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

	assetP.session = client
}

func (assetP *AssetDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data AssetListDataSourceModel
	diagState := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diagState...)
	if response.Diagnostics.HasError() {
		return
	}

	listRequestMessage := &apiv1.AssetServiceListRequest{}

	assetResp, err := assetP.session.Client.Apiv1().Asset().List(ctx, connect.NewRequest(listRequestMessage))

	if err != nil {
		response.Diagnostics.AddError("Failed to get asset list", err.Error())
	}

	data.Items = make([]assetModel, 0, len(assetResp.Msg.Assets))
	for _, asset := range assetResp.Msg.Assets {
		data.Items = append(data.Items, assetResponseMapping(asset))
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
