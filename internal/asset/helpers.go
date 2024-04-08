package asset

import (
	"sort"

	types "github.com/hashicorp/terraform-plugin-framework/types"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func assetResponseMapping(a *apiv1.Asset) assetModel {
	var partitions []*partition
	for _, p := range a.Region.Partitions {
		partition := partition{
			Id:          types.StringValue(p.Id),
			Name:        types.StringValue(p.Name),
			Address:     types.StringValue(p.Address),
			Active:      types.BoolValue(p.Active),
			Description: types.StringValue(p.Description),
		}
		partitions = append(partitions, &partition)
	}

	defaults := assetDefault{
		MachineType:       types.StringValue(a.Region.Defaults.MachineType),
		KubernetesVersion: types.StringValue(a.Region.Defaults.KubernetesVersion),
		WorkerMin:         types.Int64Value(int64(a.Region.Defaults.WorkerMin)),
		WorkerMax:         types.Int64Value(int64(a.Region.Defaults.WorkerMax)),
		Partition:         types.StringValue(a.Region.Defaults.Partition),
	}

	region := region{
		Id:          types.StringValue(a.Region.Id),
		Name:        types.StringValue(a.Region.Name),
		Address:     types.StringValue(a.Region.Address),
		Active:      types.BoolValue(a.Region.Active),
		Partitions:  partitions,
		Defaults:    &defaults,
		Description: types.StringValue(a.Region.Description),
	}

	var machineTypes []*machineType
	for _, m := range a.MachineTypes {
		machineType := machineType{
			Id:                 types.StringValue(m.Id),
			Name:               types.StringValue(m.Name),
			Cpus:               types.Int64Value(int64(m.Cpus)),
			Memory:             types.Int64Value(int64(m.Memory)),
			Storage:            types.Int64Value(int64(m.Storage)),
			CpuDescription:     types.StringValue(m.CpuDescription),
			StorageDescription: types.StringValue(m.StorageDesription), //TODO: fix typo when moving to latest api version
		}
		machineTypes = append(machineTypes, &machineType)
	}
	sort.Slice(machineTypes, func(i, j int) bool {
		return int(machineTypes[i].Memory.ValueInt64()) < int(machineTypes[j].Memory.ValueInt64())
	})

	var kubernetesVersions []*kubernetesVersion
	for _, k := range a.Kubernetes {
		kv := kubernetesVersion{
			Version: types.StringValue(k.Version),
		}
		kubernetesVersions = append(kubernetesVersions, &kv)
	}

	return assetModel{
		Region:       &region,
		MachineTypes: machineTypes,
		Kubernetes:   kubernetesVersions,
	}
}
