package cluster

import (
	"context"
	"fmt"

	connect_go "github.com/bufbuild/connect-go"
	path "github.com/hashicorp/terraform-plugin-framework/path"
	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	session "github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ resource.Resource                = &Cluster{}
	_ resource.ResourceWithConfigure   = &Cluster{}
	_ resource.ResourceWithImportState = &Cluster{}
)

func NewClusterResource() resource.Resource {
	return &Cluster{}
}

type Cluster struct {
	session *session.Session
}

// Metadata implements resource.Resource.
func (*Cluster) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cluster"
}

// Schema implements resource.Resource.
func (*Cluster) Schema(ctx context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes:          clusterResourceAttributes(),
		MarkdownDescription: "Hello World",
	}
}

// Configure implements resource.ResourceWithConfigure.
func (clusterP *Cluster) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

// Create implements resource.Resource.
func (clusterP *Cluster) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan clusterModel
	diagPlan := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diagPlan...)
	if response.Diagnostics.HasError() {
		return
	}

	// create requestMessage for client
	requestMessage := clusterCreateRequestMapping(&plan, response)

	// todo - checks: check partition name, check Kubernetes version and apply default if not set, check Maxsurge and Maxunavailable
	// check if project is set
	if requestMessage.Project == "" {
		requestMessage.Project = clusterP.session.Project
	}

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Create(ctx, connect_go.NewRequest(&requestMessage))
	if err != nil {
		response.Diagnostics.AddError("failed to create cluster", err.Error())
		return
	}

	clusterStatus := apiv1.ClusterServiceWatchStatusRequest{
		Uuid:    &clientResponse.Msg.Cluster.Uuid,
		Project: clientResponse.Msg.Cluster.Project,
	}
	err = clusterCreateWaitStatus(ctx, clusterP, &clusterStatus, []string{clusterStatusOperationTypeCreate, clusterStatusOperationTypeReconcile})
	if err != nil {
		response.Diagnostics.AddError("cluster created inconsistently", err.Error())
	}

	// Save updated data into Terraform state
	// todo knabel, update status
	data := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(data...)
}

// Read implements resource.Resource.
func (clusterP *Cluster) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state clusterModel
	diagState := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diagState...)
	if response.Diagnostics.HasError() {
		return
	}

	requestMessage := apiv1.ClusterServiceGetRequest{
		Uuid:    state.Uuid.ValueString(),
		Project: state.Project.ValueString(),
	}

	// check if project is set
	if requestMessage.Project == "" {
		requestMessage.Project = clusterP.session.Project
	}

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Get(ctx, connect_go.NewRequest(&requestMessage))

	if err != nil {
		response.Diagnostics.AddError("failed to get cluster", err.Error())
		return
	}

	// Save updated data into Terraform state
	data := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(data...)
}

// Update implements resource.Resource.
func (clusterP *Cluster) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// Read Terraform prior state data into the model
	var state clusterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan data into the model
	var plan clusterModel
	diagPlan := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diagPlan...)
	if response.Diagnostics.HasError() {
		return
	}

	// create requestMessage for client
	requestMessage := clusterUpdateRequestMapping(&state, &plan, response)

	// checks
	// check if kubernetes version is higher than the previous one

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Update(ctx, connect_go.NewRequest(&requestMessage))
	if err != nil {
		response.Diagnostics.AddError("failed to update cluster", err.Error())
		return
	}

	clusterStatus := apiv1.ClusterServiceWatchStatusRequest{
		Uuid:    &clientResponse.Msg.Cluster.Uuid,
		Project: clientResponse.Msg.Cluster.Project,
	}
	err = clusterCreateWaitStatus(ctx, clusterP, &clusterStatus, []string{clusterStatusOperationTypeReconcile})
	if err != nil {
		response.Diagnostics.AddError("cluster update status inconsistent", err.Error())
	}

	// Save updated data into Terraform state
	data := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(data...)
}

// Delete implements resource.Resource.
func (clusterP *Cluster) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state clusterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	requestMessage := apiv1.ClusterServiceDeleteRequest{
		Uuid:    state.Uuid.ValueString(),
		Project: state.Project.ValueString(),
	}

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Delete(ctx, connect_go.NewRequest(&requestMessage))
	if err != nil {
		response.Diagnostics.AddError("failed to delete cluster", err.Error())
		return
	}

	clusterStatus := apiv1.ClusterServiceWatchStatusRequest{
		Uuid:    &clientResponse.Msg.Cluster.Uuid,
		Project: clientResponse.Msg.Cluster.Project,
	}
	err = clusterCreateWaitStatus(ctx, clusterP, &clusterStatus, []string{clusterStatusOperationTypeDelete})
	if err != nil {
		response.Diagnostics.AddError("cluster delete status inconsistent", err.Error())
	}
}

// ImportState implements resource.ResourceWithImportState.
func (*Cluster) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
