package cluster

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/hashicorp/go-uuid"
	path "github.com/hashicorp/terraform-plugin-framework/path"
	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	session "github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ resource.Resource                = &ClusterResource{}
	_ resource.ResourceWithConfigure   = &ClusterResource{}
	_ resource.ResourceWithImportState = &ClusterResource{}
)

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

type ClusterResource struct {
	session *session.Session
}

// Metadata implements resource.Resource.
func (*ClusterResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cluster"
}

// Schema implements resource.Resource.
func (*ClusterResource) Schema(ctx context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes:          clusterResourceAttributes(),
		MarkdownDescription: "Managing Clusters of worker nodes. Required permissions: `Cluster *`. Can be imported by ID or name.",
	}
}

// Configure implements resource.ResourceWithConfigure.
func (clusterP *ClusterResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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
func (clusterP *ClusterResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan clusterModel
	diagPlan := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diagPlan...)
	if response.Diagnostics.HasError() {
		return
	}

	// create requestMessage for client
	requestMessage := clusterCreateRequestMapping(&plan, response)

	// check if project is set
	if requestMessage.Project == "" {
		requestMessage.Project = clusterP.session.Project
	}
	if requestMessage.Partition == "" {
		requestMessage.Partition = "eqx-mu4" // TODO: Partition
	}

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Create(ctx, connect.NewRequest(&requestMessage))
	if err != nil {
		response.Diagnostics.AddError("failed to create cluster", err.Error())
		return
	}

	clusterStatus := apiv1.ClusterServiceWatchStatusRequest{
		Uuid:    &clientResponse.Msg.Cluster.Uuid,
		Project: clientResponse.Msg.Cluster.Project,
	}
	err = clusterOperationWaitStatus(ctx, clusterP, &clusterStatus, []string{clusterStatusOperationTypeCreate, clusterStatusOperationTypeReconcile})
	if err != nil {
		response.Diagnostics.AddError("cluster created inconsistently", err.Error())
	}

	// Save updated data into Terraform state
	data := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(data...)
}

// Read implements resource.Resource.
func (clusterP *ClusterResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
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

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Get(ctx, connect.NewRequest(&requestMessage))

	if err != nil {
		response.Diagnostics.AddError("failed to get cluster", err.Error())
		return
	}

	// Save updated data into Terraform state
	data := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(data...)
}

// Update implements resource.Resource.
func (clusterP *ClusterResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
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

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Update(ctx, connect.NewRequest(&requestMessage))
	if err != nil {
		response.Diagnostics.AddError("failed to update cluster", err.Error())
		return
	}

	clusterStatus := apiv1.ClusterServiceWatchStatusRequest{
		Uuid:    &clientResponse.Msg.Cluster.Uuid,
		Project: clientResponse.Msg.Cluster.Project,
	}
	err = clusterOperationWaitStatus(ctx, clusterP, &clusterStatus, []string{clusterStatusOperationTypeCreate, clusterStatusOperationTypeReconcile})
	if err != nil {
		response.Diagnostics.AddError("cluster update status inconsistent", err.Error())
	}

	// Save updated data into Terraform state
	data := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
	response.Diagnostics.Append(data...)
}

// Delete implements resource.Resource.
func (clusterP *ClusterResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
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

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Delete(ctx, connect.NewRequest(&requestMessage))
	if err != nil {
		response.Diagnostics.AddError("failed to delete cluster", err.Error())
		return
	}

	clusterStatus := apiv1.ClusterServiceWatchStatusRequest{
		Uuid:    &clientResponse.Msg.Cluster.Uuid,
		Project: clientResponse.Msg.Cluster.Project,
	}
	err = clusterOperationWaitStatus(ctx, clusterP, &clusterStatus, []string{clusterStatusOperationTypeDelete})
	if err != nil && !strings.Contains(err.Error(), fmt.Sprintf("no entity with uuid:%q found", state.Uuid.ValueString())) {
		response.Diagnostics.AddError("cluster delete status inconsistent", err.Error())
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := uuid.ParseUUID(req.ID); err == nil {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
		return
	}

	name := req.ID
	listRequestMessage := &apiv1.ClusterServiceListRequest{
		Project: r.session.Project,
	}
	clusterList, err := r.session.Client.Apiv1().Cluster().List(ctx, connect.NewRequest(listRequestMessage))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get cluster list", err.Error())
		return
	}
	// find uuid and set uuidString
	list := clusterList.Msg.Clusters
	uuidStr, err := findUuidByName(list, name)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to find cluster with name %v", req.ID), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), uuidStr)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
}
