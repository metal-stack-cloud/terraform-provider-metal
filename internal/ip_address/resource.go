package ipaddress

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ resource.Resource                = &IpAddressResource{}
	_ resource.ResourceWithConfigure   = &IpAddressResource{}
	_ resource.ResourceWithImportState = &IpAddressResource{}
)

func NewIpResource() resource.Resource {
	return &IpAddressResource{}
}

type IpAddressResource struct {
	session *session.Session
}

// Metadata implements resource.Resource.
func (*IpAddressResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_address"
}

// Schema implements resource.Resource.
func (*IpAddressResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: ipAddressResourceAttributes(),
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *IpAddressResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.session = session
}

// Create implements resource.Resource.
func (r *IpAddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ipAddressModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipReq := &apiv1.IPServiceAllocateRequest{
		Project:     plan.Project.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}
	if ipReq.Project == "" {
		ipReq.Project = r.session.Project
	}
	switch plan.Type.ValueString() {
	case "ephemeral":
		ipReq.Static = false
	case "static":
		ipReq.Static = true
	case "unspecified": // let the api decide
	default:
		resp.Diagnostics.AddError("Invalid ip type", fmt.Sprintf("ip type %q is invalid", plan.Type.ValueString()))
		return
	}
	for _, tag := range plan.Tags {
		ipReq.Tags = append(ipReq.Tags, tag.ValueString())
	}
	createdIp, err := r.session.Client.Apiv1().IP().Allocate(ctx, connect.NewRequest(ipReq))
	if err != nil {
		resp.Diagnostics.AddError("Failed to allocate IP address", err.Error())
		return
	}
	diags = resp.State.Set(ctx, ipAddressFromApi(createdIp.Msg.Ip))
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *IpAddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ipAddressModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipResp, err := r.session.Client.Apiv1().IP().Get(ctx, connect.NewRequest(&apiv1.IPServiceGetRequest{
		Uuid:    state.Uuid.ValueString(),
		Project: r.session.Project,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get IP address", err.Error())
		return
	}

	diags = resp.State.Set(ctx, ipAddressFromApi(ipResp.Msg.Ip))
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *IpAddressResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state ipAddressModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipUpdate := &apiv1.IP{
		Uuid:        state.Uuid.ValueString(),
		Ip:          state.Ip.ValueString(),
		Name:        state.Name.ValueString(),
		Description: state.Description.ValueString(),
		Network:     state.Network.ValueString(),
		Project:     state.Project.ValueString(),
	}
	if ipUpdate.Project == "" {
		ipUpdate.Project = r.session.Project
	}
	switch state.Type.ValueString() {
	case "ephemeral":
		ipUpdate.Type = apiv1.IPType_IP_TYPE_EPHEMERAL
	case "static":
		ipUpdate.Type = apiv1.IPType_IP_TYPE_STATIC
	case "unspecified":
		ipUpdate.Type = apiv1.IPType_IP_TYPE_UNSPECIFIED
	default:
		resp.Diagnostics.AddError("Invalid ip type", fmt.Sprintf("ip type %q is invalid", state.Type.ValueString()))
		return
	}
	for _, tag := range state.Tags {
		ipUpdate.Tags = append(ipUpdate.Tags, tag.ValueString())
	}

	var plan ipAddressModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.IsNull() && plan.Name != state.Name {
		ipUpdate.Name = plan.Name.ValueString()
	}
	if !plan.Description.IsNull() && plan.Description != state.Description {
		ipUpdate.Description = plan.Description.ValueString()
	}

	if !plan.Type.IsNull() && plan.Type != state.Type {
		switch plan.Type.ValueString() {
		case "ephemeral":
			ipUpdate.Type = apiv1.IPType_IP_TYPE_EPHEMERAL
		case "static":
			ipUpdate.Type = apiv1.IPType_IP_TYPE_STATIC
		case "unspecified":
			ipUpdate.Type = apiv1.IPType_IP_TYPE_UNSPECIFIED
		default:
			resp.Diagnostics.AddError("Invalid ip type", fmt.Sprintf("ip type %q is invalid", plan.Type.ValueString()))
			return
		}
	}
	if plan.Tags != nil {
		ipUpdate.Tags = make([]string, len(plan.Tags))
		for i, tag := range plan.Tags {
			ipUpdate.Tags[i] = tag.ValueString()
		}
	}

	updatedIp, err := r.session.Client.Apiv1().IP().Update(ctx, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
		Project: r.session.Project,
		Ip:      ipUpdate,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to update IP address", err.Error())
		return
	}
	diags = resp.State.Set(ctx, ipAddressFromApi(updatedIp.Msg.Ip))
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *IpAddressResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ipAddressModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.session.Client.Apiv1().IP().Delete(ctx, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
		Uuid:    state.Uuid.ValueString(),
		Project: state.Project.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete IP address", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (*IpAddressResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
