package cluster

import (
	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	pointer "github.com/metal-stack/metal-lib/pkg/pointer"
)

type clusterModel struct {
	Uuid        types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	Project     types.String         `tfsdk:"project"`
	Partition   types.String         `tfsdk:"partition"`
	Tenant      types.String         `tfsdk:"tenant"`
	Kubernetes  types.String         `tfsdk:"kubernetes"`
	Workers     []clusterWorkerModel `tfsdk:"workers"`
	Maintenance maintenanceModel     `tfsdk:"maintenance"`
	CreatedAt   types.String         `tfsdk:"created_at"`
	UpdatedAt   types.String         `tfsdk:"updated_at"`
}

type clusterWorkerModel struct {
	Name           types.String `tfsdk:"name"`
	MachineType    types.String `tfsdk:"machine_type"`
	Minsize        types.Int64  `tfsdk:"min_size"`
	Maxsize        types.Int64  `tfsdk:"max_size"`
	Maxsurge       types.Int64  `tfsdk:"max_surge"`
	Maxunavailable types.Int64  `tfsdk:"max_unavailable"`
}

type maintenanceModel struct {
	KubernetesAutoupdate   types.Bool `tfsdk:"kubernetes_autoupdate"`
	MachineimageAutoupdate types.Bool `tfsdk:"machineimage_autoupdate"`
	// Begin                  timestamppb.Timestamp `tfsdk:"begin"`
	// Duration               durationpb.Duration   `tfsdk:"duration"`
}

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
		Kubernetes:  types.StringValue(kubernetesVersion),
		Workers:     workersSlice,
		Maintenance: maintenanceMapping,
		Tenant:      types.StringValue(clusterP.Tenant),
		CreatedAt:   types.StringValue(clusterP.CreatedAt.AsTime().String()),
		UpdatedAt:   types.StringValue(clusterP.UpdatedAt.AsTime().String()),
	}
}