package cluster

import (
	"context"
	"fmt"

	connect "github.com/bufbuild/connect-go"
	path "github.com/hashicorp/terraform-plugin-framework/path"
	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	session "github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
	pointer "github.com/metal-stack/metal-lib/pkg/pointer"
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
		Attributes: clusterResourceAttributes(),
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
func (clusterP *Cluster) Create(context context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan clusterModel
	diagPlan := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagPlan...)
	if response.Diagnostics.HasError() {
		return
	}

	// map terraform workers arguments to Worker struct
	workersMapping := []*apiv1.Worker{{
		Name:           "test",
		MachineType:    plan.Workers.MachineType.ValueString(),
		Minsize:        uint32(plan.Workers.Minsize.ValueInt64()),
		Maxsize:        uint32(plan.Workers.Maxsize.ValueInt64()),
		Maxsurge:       uint32(plan.Workers.Maxsurge.ValueInt64()),
		Maxunavailable: uint32(plan.Workers.Maxunavailable.ValueInt64()),
	},
	}

	// create requestMessage for client
	requestMessage := apiv1.ClusterServiceCreateRequest{
		Name:       plan.Name.ValueString(),
		Project:    plan.Project.ValueString(),
		Partition:  "eqx-mu4",
		Kubernetes: plan.Kubernetes,
		Workers:    workersMapping,
	}

	// checks
	if requestMessage.Project == "" {
		requestMessage.Project = clusterP.session.Project
	}

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Create(context, connect.NewRequest(&requestMessage))
	if err != nil {
		response.Diagnostics.AddError("Failed to create cluster", err.Error())
		return
	}

	// Save updated data into Terraform state
	data := response.State.Set(context, clusterResponseConvert(clientResponse.Msg.Cluster))
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

	requestMessage := &apiv1.ClusterServiceGetRequest{
		Uuid:    state.Uuid.ValueString(),
		Project: state.Project.ValueString(),
	}

	clientResponse, err := clusterP.session.Client.Apiv1().Cluster().Get(ctx, connect.NewRequest(requestMessage))

	if err != nil {
		response.Diagnostics.AddError("Failed to get cluster", err.Error())
		return
	}

	// Save updated data into Terraform state
	data := response.State.Set(ctx, clusterResponseConvert(clientResponse.Msg.Cluster))
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

	// map terraform workers arguments to WorkerUpdate struct
	workersMapping := []*apiv1.WorkerUpdate{{
		Name:           "test",
		MachineType:    pointer.Pointer(plan.Workers.MachineType.ValueString()),
		Minsize:        pointer.Pointer(uint32(plan.Workers.Minsize.ValueInt64())),
		Maxsize:        pointer.Pointer(uint32(plan.Workers.Maxsize.ValueInt64())),
		Maxsurge:       pointer.Pointer(uint32(plan.Workers.Maxsurge.ValueInt64())),
		Maxunavailable: pointer.Pointer(uint32(plan.Workers.Maxunavailable.ValueInt64())),
	},
	}

	requestMessage := apiv1.ClusterServiceUpdateRequest{
		Uuid:       plan.Uuid.ValueString(),
		Project:    plan.Project.ValueString(),
		Kubernetes: plan.Kubernetes,
		// map plan Workers to WorkerUpdate struct
		Workers:     workersMapping,
		Maintenance: plan.Maintenance,
	}

	// checks
	// if requestMessage.Project == "" {
	// 	requestMessage.Project = clusterP.session.Project
	// }
	// if !plan.Name.IsNull() && plan.Name != state.Name {
	// 	requestMessage.Name = plan.Name.ValueString()
	// }
	// check if kubernetes version is higher than the previous one
	// check Maxsurge and Maxunavailable

	clientResponse, clientError := clusterP.session.Client.Apiv1().Cluster().Update(ctx, connect.NewRequest(&requestMessage))

	if clientError != nil {
		response.Diagnostics.AddError("Failed to update cluster", clientError.Error())
		return
	}

	// Save updated data into Terraform state
	data := response.State.Set(ctx, clusterResponseConvert(clientResponse.Msg.Cluster))
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

	_, clientError := clusterP.session.Client.Apiv1().Cluster().Delete(ctx, connect.NewRequest(&requestMessage))

	if clientError != nil {
		response.Diagnostics.AddError("Failed to delete cluster", clientError.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (*Cluster) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
