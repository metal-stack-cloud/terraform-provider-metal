package ipaddress

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

var (
	_ resource.Resource                = &PublicIpResource{}
	_ resource.ResourceWithConfigure   = &PublicIpResource{}
	_ resource.ResourceWithImportState = &PublicIpResource{}
)

func NewPublicIpResource() resource.Resource {
	return &PublicIpResource{}
}

type PublicIpResource struct {
	session *session.Session
}

// Metadata implements resource.Resource.
func (*PublicIpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

// Schema implements resource.Resource.
func (*PublicIpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: publicIpResourceAttributes(),
		MarkdownDescription: "Each cluster gets an IP automatically provided on the internet gateway for outgoing communication. \n" +
			"Services get an IP automatically on creation. \n" +
			"Services and gateway IPs are dynamic by default. \n" +
			"You can use an IP address in several clusters and locations at the same time. \n" +
			"Required permissions: `IP *`. Can be imported by ID, name or ip address.",
	}
}

// Configure implements resource.ResourceWithConfigure.
func (ip *PublicIpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	session, ok := req.ProviderData.(*session.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *session.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	ip.session = session
}

// Create implements resource.Resource.
func (ip *PublicIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan publicIpModel
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
		ipReq.Project = ip.session.Project
	}
	switch plan.Type.ValueString() {
	case "ephemeral":
		ipReq.Static = false
	case "static":
		ipReq.Static = true
	case "unspecified", "": // let the api decide
	default:
		resp.Diagnostics.AddError("Invalid ip type", fmt.Sprintf("ip type %q is invalid", plan.Type.ValueString()))
		return
	}
	for _, tag := range plan.Tags {
		ipReq.Tags = append(ipReq.Tags, tag.ValueString())
	}
	createdIp, err := ip.session.Client.Apiv1().IP().Allocate(ctx, connect.NewRequest(ipReq))
	if err != nil {
		resp.Diagnostics.AddError("Failed to allocate IP address", err.Error())
		return
	}
	diags = resp.State.Set(ctx, publicIpFromApi(createdIp.Msg.Ip))
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (ip *PublicIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state publicIpModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipResp, err := ip.session.Client.Apiv1().IP().Get(ctx, connect.NewRequest(&apiv1.IPServiceGetRequest{
		Uuid:    state.Uuid.ValueString(),
		Project: ip.session.Project,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get IP address", err.Error())
		return
	}

	diags = resp.State.Set(ctx, publicIpFromApi(ipResp.Msg.Ip))
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (ip *PublicIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state publicIpModel
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
		ipUpdate.Project = ip.session.Project
	}
	switch state.Type.ValueString() {
	case "ephemeral":
		ipUpdate.Type = apiv1.IPType_IP_TYPE_EPHEMERAL
	case "static":
		ipUpdate.Type = apiv1.IPType_IP_TYPE_STATIC
	case "unspecified", "":
		ipUpdate.Type = apiv1.IPType_IP_TYPE_UNSPECIFIED
	default:
		resp.Diagnostics.AddError("Invalid ip type", fmt.Sprintf("ip type %q is invalid", state.Type.ValueString()))
		return
	}
	for _, tag := range state.Tags {
		ipUpdate.Tags = append(ipUpdate.Tags, tag.ValueString())
	}

	var plan publicIpModel
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
			if ipUpdate.Type == apiv1.IPType_IP_TYPE_STATIC {
				resp.Diagnostics.AddError("Cannot update static IPs to ephemeral", "Static IP addresses cannot be declared ephemeral.")
				return
			}
			ipUpdate.Type = apiv1.IPType_IP_TYPE_EPHEMERAL
		case "static":
			ipUpdate.Type = apiv1.IPType_IP_TYPE_STATIC
		case "unspecified", "":
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

	updatedIp, err := ip.session.Client.Apiv1().IP().Update(ctx, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
		Project: ip.session.Project,
		Ip:      ipUpdate,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to update IP address", err.Error())
		return
	}
	diags = resp.State.Set(ctx, publicIpFromApi(updatedIp.Msg.Ip))
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (ip *PublicIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state publicIpModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := ip.session.Client.Apiv1().IP().Delete(ctx, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
		Uuid:    state.Uuid.ValueString(),
		Project: state.Project.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete IP address", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (ip *PublicIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := uuid.ParseUUID(req.ID); err == nil {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
		return
	}

	listRequestMessage := &apiv1.IPServiceListRequest{
		Project: ip.session.Project,
	}
	clusterList, err := ip.session.Client.Apiv1().IP().List(ctx, connect.NewRequest(listRequestMessage))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get all public ips", err.Error())
		return
	}
	// find uuid and set uuidString
	list := clusterList.Msg.Ips
	uuidStr, err := findUuidByName(list, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to find IP with address or name %v", req.ID), err.Error())
		return
	}
	req.ID = uuidStr
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func findUuidByName(list []*apiv1.IP, nameOrIP string) (string, error) {
	for _, e := range list {
		if e.Name == nameOrIP || e.Ip == nameOrIP {
			return e.Uuid, nil
		}
	}
	return "", fmt.Errorf("ip address or name not found")
}
