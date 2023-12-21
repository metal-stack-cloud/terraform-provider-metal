package cluster

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"connectrpc.com/connect"
	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	pointer "github.com/metal-stack/metal-lib/pkg/pointer"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	clusterStatusOperationTypeCreate    = "Create"
	clusterStatusOperationTypeReconcile = "Reconcile"
	clusterStatusOperationTypeDelete    = "Delete"
	clusterStatusStateProcessing        = "Processing"
	clusterStatusStateSucceeded         = "Succeeded"
	clusterStatusStateError             = "Error"
	clusterStatusStateFailed            = "Failed"
	clusterStatusStatePending           = "Pending"
	clusterStatusStateAborted           = "Aborted"
)

func clusterCreateRequestMapping(plan *clusterModel, response *resource.CreateResponse) apiv1.ClusterServiceCreateRequest {
	// map terraform Kubernetes arguments to KubernetesSpec struct
	kubernetesSpecMapping := &apiv1.KubernetesSpec{
		Version: plan.Kubernetes.ValueString(),
	}
	// FIXME: format will change, and default might not be required later
	// https://github.com/metal-stack-cloud/terraform-provider-metal/issues/51
	if plan.maintenance == nil {
		plan.maintenance = &maintenanceModel{
			TimeWindow: maintenanceTimeWindow{
				Begin:    types.StringValue("2:00 AM"),
				Duration: types.Int64Value(2),
			},
		}
	}
	maintenanceMapping := &apiv1.Maintenance{
		KubernetesAutoupdate:   plan.maintenance.KubernetesAutoupdate.ValueBoolPointer(),   //TODO: default to true and delete from schema?
		MachineimageAutoupdate: plan.maintenance.MachineimageAutoupdate.ValueBoolPointer(), //TODO: default to true and delete from schema?
		TimeWindow: &apiv1.MaintenanceTimeWindow{
			Begin: &timestamppb.Timestamp{
				Seconds: computeBegin(plan.maintenance.TimeWindow.Begin.ValueString()),
			},
			Duration: computeDuration(plan.maintenance.TimeWindow.Duration.ValueInt64()),
		},
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

	// FIXME: format will change, and default might not be required later
	// https://github.com/metal-stack-cloud/terraform-provider-metal/issues/51
	if plan.maintenance == nil {
		plan.maintenance = &maintenanceModel{
			TimeWindow: maintenanceTimeWindow{
				Begin:    types.StringValue("2:00 AM"),
				Duration: types.Int64Value(2),
			},
		}
	}
	// map maintenance arguments to Maintenance struct
	maintenanceMapping := &apiv1.Maintenance{
		KubernetesAutoupdate:   plan.maintenance.KubernetesAutoupdate.ValueBoolPointer(),
		MachineimageAutoupdate: plan.maintenance.MachineimageAutoupdate.ValueBoolPointer(),
		TimeWindow: &apiv1.MaintenanceTimeWindow{
			Begin: &timestamppb.Timestamp{
				Seconds: computeBegin(plan.maintenance.TimeWindow.Begin.ValueString()),
			},
			Duration: computeDuration(plan.maintenance.TimeWindow.Duration.ValueInt64()),
		},
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
		worker := clusterWorkerModel{
			Name:        types.StringValue(v.Name),
			MachineType: types.StringValue(v.MachineType),
			Minsize:     types.Int64Value(int64(v.Minsize)),
			Maxsize:     types.Int64Value(int64(v.Maxsize)),
		}
		if v.Maxsurge != 0 {
			worker.Maxsurge = types.Int64Value(int64(v.Maxsurge))
		}
		if v.Maxunavailable != 0 {
			worker.Maxunavailable = types.Int64Value(int64(v.Maxunavailable))
		}
		workersSlice = append(workersSlice, worker)
	}

	// map terraform Kubernetes arguments to maintenance struct
	maintenanceMapping := maintenanceModel{
		KubernetesAutoupdate:   types.BoolValue(*clusterP.Maintenance.KubernetesAutoupdate),
		MachineimageAutoupdate: types.BoolValue(*clusterP.Maintenance.MachineimageAutoupdate),
		TimeWindow: maintenanceTimeWindow{
			Begin:    types.StringValue(convertTimestamp(clusterP.Maintenance.TimeWindow.Begin)),
			Duration: types.Int64Value(convertDuration(clusterP.Maintenance.TimeWindow.Duration.Seconds)),
		},
	}

	return clusterModel{
		Uuid:        types.StringValue(clusterP.Uuid),
		Name:        types.StringValue(clusterP.Name),
		Project:     types.StringValue(clusterP.Project),
		Partition:   types.StringValue(clusterP.Partition),
		Tenant:      types.StringValue(clusterP.Tenant),
		Kubernetes:  types.StringValue(kubernetesVersion),
		Workers:     workersSlice,
		maintenance: &maintenanceMapping,
		CreatedAt:   types.StringValue(clusterP.CreatedAt.AsTime().String()),
		UpdatedAt:   types.StringValue(clusterP.UpdatedAt.AsTime().String()),
	}
}

func clusterOperationWaitStatus(ctx context.Context, clusterP *ClusterResource, statusRequest *apiv1.ClusterServiceWatchStatusRequest, operationWhitelist []string) error {
	// add timeout to context
	watchCtx, watchCancel := context.WithTimeout(ctx, 20*time.Minute)
	defer watchCancel()

	// It might take a while until expected cluster operations are reflected
	var hadValidOperationType bool
	for {
		// cluster status wait functions
		clusterStatusStream, err := clusterP.session.Client.Apiv1().Cluster().WatchStatus(watchCtx, connect.NewRequest(statusRequest))
		if err != nil {
			return fmt.Errorf("cluster watch status response failed %w", err)
		}

		var statusMsg *apiv1.ClusterStatus
		for clusterStatusStream.Receive() {
			statusMsg = clusterStatusStream.Msg().Status
			if !hadValidOperationType && !slices.Contains(operationWhitelist, statusMsg.Type) {
				continue
			}
			if !hadValidOperationType {
				hadValidOperationType = true
			}

			tflog.Debug(ctx, "waiting for cluster status change", map[string]any{
				"progress": statusMsg.Progress,
				"type":     statusMsg.Type,
				"state":    statusMsg.State,
			})

			// check operation type of cluster
			if !slices.Contains(operationWhitelist, statusMsg.Type) && statusMsg.Progress > 0 {
				tflog.Debug(ctx, fmt.Sprintf("statusMsg check of type not %q", operationWhitelist), map[string]any{
					"progress": statusMsg.Progress,
					"type":     statusMsg.Type,
					"state":    statusMsg.State,
				})

				return fmt.Errorf("expected operation type of %q, got %q", operationWhitelist, statusMsg.Type)
			}

			if statusMsg.State == clusterStatusStateSucceeded {
				tflog.Debug(ctx, fmt.Sprintf("statusMsg check of state %v successful", clusterStatusStateSucceeded), map[string]any{
					"progress": statusMsg.Progress,
					"type":     statusMsg.Type,
					"state":    statusMsg.State,
				})
				return nil
			}
		}

		err = clusterStatusStream.Err()
		if errors.Is(err, io.ErrUnexpectedEOF) {
			// reconnect if EOF error
			continue
		}
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("unknown stream connection error encountered with cluster status %v", clusterStatusStateSucceeded), map[string]any{
				"error": err.Error(),
			})
			return err
		}
		if statusMsg.State != clusterStatusStateSucceeded {
			tflog.Debug(ctx, fmt.Sprintf("statusMsg check of state %v failed", clusterStatusStateSucceeded), map[string]any{
				"progress": statusMsg.Progress,
				"type":     statusMsg.Type,
				"state":    statusMsg.State,
			})
			return err
		}
	}
}

func computeBegin(s string) int64 {
	layout := "03:04 PM"
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}
	utcTime := parsedTime.UTC()
	hours := utcTime.Hour()
	minutes := utcTime.Minute()
	seconds := utcTime.Second()
	totalSeconds := hours*3600 + minutes*60 + seconds

	return int64(totalSeconds)
}

func computeDuration(hours int64) *durationpb.Duration {
	return durationpb.New(time.Duration(hours) * time.Hour)
}

func convertTimestamp(t *timestamppb.Timestamp) string {
	timeObj := t.AsTime()
	return timeObj.Format("03:04 PM")
}

func convertDuration(d int64) int64 {
	return d / 3600
}
