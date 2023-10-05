package cluster

import (
	"context"
	"fmt"
	"time"

	connect_go "github.com/bufbuild/connect-go"
	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	pointer "github.com/metal-stack/metal-lib/pkg/pointer"
)

func clusterCreateRequestMapping(plan *clusterModel, response *resource.CreateResponse) apiv1.ClusterServiceCreateRequest {
	// map terraform Kubernetes arguments to KubernetesSpec struct
	kubernetesSpecMapping := &apiv1.KubernetesSpec{
		Version: plan.Kubernetes.ValueString(),
	}
	// map maintenance arguments to Maintenance struct
	maintenanceMapping := &apiv1.Maintenance{
		KubernetesAutoupdate:   pointer.Pointer(bool(plan.Maintenance.KubernetesAutoupdate.ValueBool())),
		MachineimageAutoupdate: pointer.Pointer(plan.Maintenance.MachineimageAutoupdate.ValueBool()),
		// TimeWindow:             &apiv1.MaintenanceTimeWindow{
		// 	// todo
		// },
	}

	// map terraform workers list arguments to Worker struct
	var workersSlice []*apiv1.Worker
	for _, v := range plan.Workers {
		workersSlice = append(workersSlice, &apiv1.Worker{
			Name:           v.Name.ValueString(),
			MachineType:    v.MachineType.ValueString(),
			Minsize:        uint32(v.Minsize.ValueInt64()),
			Maxsize:        uint32(v.Maxsize.ValueInt64()),
			Maxsurge:       uint32(v.Maxsurge.ValueInt64()),
			Maxunavailable: uint32(v.Maxunavailable.ValueInt64()),
		})
	}

	// check workersSlice slice length
	if workersSlice == nil {
		response.Diagnostics.AddError("check failed of workersSlice slice", "workersSlice slice length is 0")
		return apiv1.ClusterServiceCreateRequest{}
	}

	// create ClusterServiceCreateRequest for client
	return apiv1.ClusterServiceCreateRequest{
		Name:        plan.Name.ValueString(),
		Project:     plan.Project.ValueString(),
		Partition:   plan.Partition.ValueString(),
		Kubernetes:  kubernetesSpecMapping,
		Workers:     workersSlice,
		Maintenance: maintenanceMapping,
	}
}

func clusterUpdateRequestMapping(state *clusterModel, plan *clusterModel, response *resource.UpdateResponse) apiv1.ClusterServiceUpdateRequest {
	// map terraform Kubernetes arguments to KubernetesSpec struct
	kubernetesSpecMapping := &apiv1.KubernetesSpec{
		Version: plan.Kubernetes.ValueString(),
	}
	// map maintenance arguments to Maintenance struct
	maintenanceMapping := &apiv1.Maintenance{
		KubernetesAutoupdate:   plan.Maintenance.KubernetesAutoupdate.ValueBoolPointer(),
		MachineimageAutoupdate: plan.Maintenance.MachineimageAutoupdate.ValueBoolPointer(),
		// TimeWindow:             &apiv1.MaintenanceTimeWindow{
		// 	// todo
		// },
	}
	// map terraform workers list arguments to WorkerUpdate struct
	var workersSlice []*apiv1.WorkerUpdate
	for _, v := range plan.Workers {
		workersSlice = append(workersSlice, &apiv1.WorkerUpdate{
			Name:           v.Name.ValueString(),
			MachineType:    pointer.Pointer(v.MachineType.ValueString()),
			Minsize:        pointer.Pointer(uint32(v.Minsize.ValueInt64())),
			Maxsize:        pointer.Pointer(uint32(v.Maxsize.ValueInt64())),
			Maxsurge:       pointer.Pointer(uint32(v.Maxsurge.ValueInt64())),
			Maxunavailable: pointer.Pointer(uint32(v.Maxunavailable.ValueInt64())),
		})
	}

	// check workersSlice slice length
	if workersSlice == nil {
		response.Diagnostics.AddError("check failed of workersSlice slice", "workersSlice slice length is 0")
		return apiv1.ClusterServiceUpdateRequest{}
	}

	// update ClusterServiceUpdateRequest for client
	return apiv1.ClusterServiceUpdateRequest{
		Uuid:        state.Uuid.ValueString(),
		Project:     state.Project.ValueString(),
		Kubernetes:  kubernetesSpecMapping,
		Workers:     workersSlice,
		Maintenance: maintenanceMapping,
	}
}

func clusterResponseMapping(clusterP *apiv1.Cluster) clusterModel {
	kubernetesVersion := clusterP.Kubernetes.Version
	// check if workersSlice slice is length > 1
	// check for null values
	var workersSlice []clusterWorkerModel
	for _, v := range clusterP.Workers {
		workersSlice = append(workersSlice, clusterWorkerModel{
			Name:           types.StringValue(v.Name),
			MachineType:    types.StringValue(v.MachineType),
			Minsize:        types.Int64Value(int64(v.Minsize)),
			Maxsize:        types.Int64Value(int64(v.Maxsize)),
			Maxsurge:       types.Int64Value(int64(v.Maxsurge)),
			Maxunavailable: types.Int64Value(int64(v.Maxunavailable)),
		})
	}

	// map terraform Kubernetes arguments to maintenance struct
	maintenanceMapping := maintenanceModel{
		KubernetesAutoupdate:   types.BoolValue(*clusterP.Maintenance.KubernetesAutoupdate),
		MachineimageAutoupdate: types.BoolValue(*clusterP.Maintenance.MachineimageAutoupdate),
		// Begin:                  *clusterP.Maintenance.TimeWindow.Begin,
		// Duration:               *clusterP.Maintenance.TimeWindow.Duration,
	}

	return clusterModel{
		Uuid:        types.StringValue(clusterP.Uuid),
		Name:        types.StringValue(clusterP.Name),
		Project:     types.StringValue(clusterP.Project),
		Partition:   types.StringValue(clusterP.Partition),
		Tenant:      types.StringValue(clusterP.Tenant),
		Kubernetes:  types.StringValue(kubernetesVersion),
		Workers:     workersSlice,
		Maintenance: maintenanceMapping,
		CreatedAt:   types.StringValue(clusterP.CreatedAt.AsTime().String()),
		UpdatedAt:   types.StringValue(clusterP.UpdatedAt.AsTime().String()),
	}
}

func clusterCreateWaitStatus(ctx context.Context, clusterP *Cluster, clientResponse *connect_go.Response[apiv1.ClusterServiceCreateResponse], response *resource.CreateResponse) {
	clusterStatus := apiv1.ClusterServiceWatchStatusRequest{
		Uuid:    &clientResponse.Msg.Cluster.Uuid,
		Project: clientResponse.Msg.Cluster.Project,
	}

	// add timeout to context
	watchCtx, watchCancel := context.WithTimeout(ctx, 30*time.Minute)

	defer func() {
		// context cancel function
		watchCancel()
		// Save updated data into Terraform state
		// todo knabel, update status
		data := response.State.Set(ctx, clusterResponseMapping(clientResponse.Msg.Cluster))
		response.Diagnostics.Append(data...)
	}()

	// cluster status wait functions
	clusterStatusStream, err := clusterP.session.Client.Apiv1().Cluster().WatchStatus(watchCtx, connect_go.NewRequest(&clusterStatus))
	if err != nil {
		response.Diagnostics.AddError("cluster watch status reponse failed", err.Error())
		return
	}

	var statusMsg *apiv1.ClusterStatus
	for clusterStatusStream.Receive() {
		statusMsg = clusterStatusStream.Msg().Status

		tflog.Debug(ctx, "waiting for cluster to become ready", map[string]any{
			"progress": statusMsg.Progress,
			"type":     statusMsg.Type,
			"state":    statusMsg.State,
		})

		// check operation type of cluster
		// todo - condition check on empty type not working
		if statusMsg.Type != clusterStatusOperationTypeCreate && statusMsg.Type != clusterStatusOperationTypeReconcile && statusMsg.Type != "" && statusMsg.Progress > 0 {
			tflog.Debug(ctx, fmt.Sprintf("statusMsg check of type not %v and %v", clusterStatusOperationTypeCreate, clusterStatusOperationTypeReconcile), map[string]any{
				"progress":          statusMsg.Progress,
				"type":              statusMsg.Type,
				"state":             statusMsg.State,
				"ApiServerReady":    statusMsg.ApiServerReady,
				"ControlPlaneReady": statusMsg.ControlPlaneReady,
			})
			response.Diagnostics.AddError(
				"created cluster is in unexpected operation type", fmt.Sprintf("expected create or reconcile operation type, got %q", statusMsg.Type),
			)

			err := clusterStatusStream.Close()
			if err != nil {
				response.Diagnostics.AddError("could not close cluster status stream gracefully", err.Error())
				tflog.Debug(ctx, "could not close cluster status stream gracefully", map[string]any{
					"error":             err,
					"progress":          statusMsg.Progress,
					"type":              statusMsg.Type,
					"state":             statusMsg.State,
					"ApiServerReady":    statusMsg.ApiServerReady,
					"ControlPlaneReady": statusMsg.ControlPlaneReady,
				})
			}

			return
		}

		if statusMsg.State == clusterStatusStateSucceeded {
			_ = clusterStatusStream.Close()
			tflog.Debug(ctx, fmt.Sprintf("statusMsg check of state %v successful", clusterStatusStateSucceeded), map[string]any{
				"progress":       statusMsg.Progress,
				"type":           statusMsg.Type,
				"state":          statusMsg.State,
				"ApiServerReady": statusMsg.ApiServerReady,
			})
			break
		}
	}

	err = clusterStatusStream.Err()
	if err != nil {
		// response.Diagnostics.AddError("could not determine cluster status", err.Error())
		tflog.Debug(ctx, fmt.Sprintf("stream error encountered with cluster status %v", clusterStatusStateSucceeded), map[string]any{
			"error": err.Error(),
		})
		// return
	}

	if statusMsg.State != clusterStatusStateSucceeded {
		response.Diagnostics.AddError("created cluster is in unexpected state", err.Error())
		_ = clusterStatusStream.Close()
		tflog.Debug(ctx, fmt.Sprintf("statusMsg check of state %v failed", clusterStatusStateSucceeded), map[string]any{
			"progress":       statusMsg.Progress,
			"type":           statusMsg.Type,
			"state":          statusMsg.State,
			"ApiServerReady": statusMsg.ApiServerReady,
		})
		return
	}
}
