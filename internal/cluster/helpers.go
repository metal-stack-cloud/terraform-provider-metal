package cluster

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"connectrpc.com/connect"
	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	pointer "github.com/metal-stack/metal-lib/pkg/pointer"
	"google.golang.org/protobuf/types/known/durationpb"
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

func clusterCreateRequestMapping(plan *clusterModel, response *resource.CreateResponse) *apiv1.ClusterServiceCreateRequest {
	// map terraform Kubernetes arguments to KubernetesSpec struct
	kubernetesSpecMapping := &apiv1.KubernetesSpec{
		Version: plan.Kubernetes.ValueString(),
	}
	var (
		maintenanceMapping apiv1.Maintenance
	)
	// defaults for the maintenance time window, according to the console
	if plan.Maintenance == nil {
		maintenanceMapping = apiv1.Maintenance{
			TimeWindow: &apiv1.MaintenanceTimeWindow{
				Begin: &apiv1.Time{
					Hour:     uint32(*types.Int64Value(1).ValueInt64Pointer()),
					Minute:   uint32(*types.Int64Value(0).ValueInt64Pointer()),
					Timezone: *types.StringValue("UTC").ValueStringPointer(),
				},
			},
		}
	} else {
		maintenanceMapping = apiv1.Maintenance{
			TimeWindow: &apiv1.MaintenanceTimeWindow{
				Begin: &apiv1.Time{
					Hour:     uint32(*plan.Maintenance.TimeWindow.Begin.Hour.ValueInt64Pointer()),
					Minute:   uint32(*plan.Maintenance.TimeWindow.Begin.Minute.ValueInt64Pointer()),
					Timezone: plan.Maintenance.TimeWindow.Begin.Timezone.ValueString(),
				},
				Duration: computeDuration(plan.Maintenance.TimeWindow.Duration.ValueInt64()),
			},
		}
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
		return &apiv1.ClusterServiceCreateRequest{}
	}

	// create ClusterServiceCreateRequest for client
	return &apiv1.ClusterServiceCreateRequest{
		Name:        plan.Name.ValueString(),
		Project:     plan.Project.ValueString(),
		Partition:   plan.Partition.ValueString(),
		Kubernetes:  kubernetesSpecMapping,
		Workers:     workersSlice,
		Maintenance: &maintenanceMapping,
	}
}

func clusterUpdateRequestMapping(state *clusterModel, plan *clusterModel, response *resource.UpdateResponse) apiv1.ClusterServiceUpdateRequest {
	// map terraform Kubernetes arguments to KubernetesSpec struct
	kubernetesSpecMapping := &apiv1.KubernetesSpec{
		Version: plan.Kubernetes.ValueString(),
	}

	// map maintenance arguments to Maintenance struct
	maintenanceMapping := &apiv1.Maintenance{
		TimeWindow: &apiv1.MaintenanceTimeWindow{
			Begin: &apiv1.Time{
				Hour:     uint32(*plan.Maintenance.TimeWindow.Begin.Hour.ValueInt64Pointer()),
				Minute:   uint32(*plan.Maintenance.TimeWindow.Begin.Minute.ValueInt64Pointer()),
				Timezone: plan.Maintenance.TimeWindow.Begin.Timezone.ValueString(),
			},
			Duration: computeDuration(plan.Maintenance.TimeWindow.Duration.ValueInt64()),
		},
	}
	// map terraform workers list arguments to WorkerUpdate struct
	var workersSlice []*apiv1.WorkerUpdate
	for _, v := range plan.Workers {
		workerUpdate := &apiv1.WorkerUpdate{
			Name:        v.Name.ValueString(),
			MachineType: pointer.Pointer(v.MachineType.ValueString()),
			Minsize:     pointer.Pointer(uint32(v.Minsize.ValueInt64())),
			Maxsize:     pointer.Pointer(uint32(v.Maxsize.ValueInt64())),
		}
		if !v.Maxsurge.IsNull() {
			workerUpdate.Maxsurge = pointer.Pointer(uint32(v.Maxsurge.ValueInt64()))
		}
		if !v.Maxunavailable.IsNull() {
			workerUpdate.Maxunavailable = pointer.Pointer(uint32(v.Maxunavailable.ValueInt64()))
		}
		workersSlice = append(workersSlice, workerUpdate)

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
			Begin: maintenanceTime{
				Hour:     types.Int64Value(int64(clusterP.Maintenance.TimeWindow.Begin.Hour)),
				Minute:   types.Int64Value(int64(clusterP.Maintenance.TimeWindow.Begin.Minute)),
				Timezone: types.StringValue(clusterP.Maintenance.TimeWindow.Begin.Timezone),
			},
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
		Maintenance: &maintenanceMapping,
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

func computeDuration(hours int64) *durationpb.Duration {
	return durationpb.New(time.Duration(hours) * time.Hour)
}

func convertDuration(d int64) int64 {
	return d / 3600
}

func patchKubernetesVersion(respV string, planV string) bool {
	respVSplit := strings.Split(respV, ".")
	planVSplit := strings.Split(planV, ".")
	if respVSplit[0] == planVSplit[0] && respVSplit[1] == planVSplit[1] && respVSplit[2] != planVSplit[2] {
		return true
	}
	return false
}
