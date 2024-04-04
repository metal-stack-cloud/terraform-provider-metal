package asset

import (
	types "github.com/hashicorp/terraform-plugin-framework/types"
)

type AssetListDataSourceModel struct {
	Items []assetModel `tfsdk:"items"`
}

type assetModel struct {
	Region       *region              `tfsdk:"region"`
	MachineTypes []*machineType       `tfsdk:"machine_types"`
	Kubernetes   []*kubernetesVersion `tfsdk:"kubernetes"`
}

type region struct {
	Id          types.String  `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Address     types.String  `tfsdk:"address"`
	Active      types.Bool    `tfsdk:"active"`
	Partitions  []*partition  `tfsdk:"partitions"`
	Defaults    *assetDefault `tfsdk:"defaults"`
	Description types.String  `tfsdk:"description"`
}

type partition struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	Active      types.Bool   `tfsdk:"active"`
	Description types.String `tfsdk:"description"`
}

type assetDefault struct {
	MachineType       types.String `tfsdk:"machine_type"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
	WorkerMin         types.Int64  `tfsdk:"worker_min"`
	WorkerMax         types.Int64  `tfsdk:"worker_max"`
	Partition         types.String `tfsdk:"partition"`
}

type machineType struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Cpus               types.Int64  `tfsdk:"cpus"`
	Memory             types.Int64  `tfsdk:"memory"`
	Storage            types.Int64  `tfsdk:"storage"`
	CpuDescription     types.String `tfsdk:"cpu_description"`
	StorageDescription types.String `tfsdk:"storage_description"`
}

type kubernetesVersion struct {
	Version types.String `tfsdk:"version"`
}
