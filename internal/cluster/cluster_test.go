package cluster

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	assert "github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_patchKubernetesVersion(t *testing.T) {
	tests := []struct {
		name        string
		respVersion string
		planVersion string
		want        bool
	}{
		{
			name:        "Update Kubernetes version at patch level",
			respVersion: "1.28.11",
			planVersion: "1.28.10",
			want:        true,
		},
		{
			name:        "Prevent Kubernetes version update at minor level",
			respVersion: "1.29.1",
			planVersion: "1.28.10",
			want:        false,
		},
		{
			name:        "Prevent Kubernetes version update at major level",
			respVersion: "2.10.1",
			planVersion: "1.28.10",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := patchKubernetesVersion(tt.respVersion, tt.planVersion)
			if res != tt.want {
				t.Errorf(`patchKubernetesVersion = %v, want %v`, res, tt.want)
			}
		})
	}
}

func Test_clusterCreateRequestMapping(t *testing.T) {
	tests := []struct {
		name         string
		planMock     *clusterModel
		responseMock *resource.CreateResponse
		want         *apiv1.ClusterServiceCreateRequest
	}{
		{
			name: "Request mapping with all cluster fields set",
			planMock: &clusterModel{
				Name:       basetypes.NewStringValue("cluster"),
				Project:    basetypes.NewStringValue("default-project"),
				Partition:  basetypes.NewStringValue("eqx-mu4"),
				Kubernetes: basetypes.NewStringValue("1.28.11"),
				Workers: []clusterWorkerModel{
					{
						Name:           basetypes.NewStringValue("default"),
						MachineType:    basetypes.NewStringValue("n1-medium-x86"),
						Minsize:        basetypes.NewInt64Value(1),
						Maxsize:        basetypes.NewInt64Value(3),
						Maxsurge:       basetypes.NewInt64Value(1),
						Maxunavailable: basetypes.NewInt64Value(0),
					},
				},
				Maintenance: &maintenanceModel{
					TimeWindow: maintenanceTimeWindow{
						Begin: maintenanceTime{
							Hour:     basetypes.NewInt64Value(14),
							Minute:   basetypes.NewInt64Value(30),
							Timezone: basetypes.NewStringValue("UTC"),
						},
						Duration: basetypes.NewInt64Value(1),
					},
				},
			},
			want: &apiv1.ClusterServiceCreateRequest{
				Name:      "cluster",
				Project:   "default-project",
				Partition: "eqx-mu4",
				Kubernetes: &apiv1.KubernetesSpec{
					Version: "1.28.11",
				},
				Workers: []*apiv1.Worker{
					{
						Name:           "default",
						MachineType:    "n1-medium-x86",
						Minsize:        1,
						Maxsize:        3,
						Maxsurge:       1,
						Maxunavailable: 0,
					},
				},
				Maintenance: &apiv1.Maintenance{
					TimeWindow: &apiv1.MaintenanceTimeWindow{
						Begin: &apiv1.Time{
							Hour:     uint32(14),
							Minute:   uint32(30),
							Timezone: "UTC",
						},
						Duration: &durationpb.Duration{
							Seconds: int64(3600),
						},
					},
				},
			},
		},
		{
			name: "Request mapping without cluster maintenance fields set",
			planMock: &clusterModel{
				Name:       basetypes.NewStringValue("cluster"),
				Project:    basetypes.NewStringValue("default-project"),
				Partition:  basetypes.NewStringValue("eqx-mu4"),
				Kubernetes: basetypes.NewStringValue("1.28.11"),
				Workers: []clusterWorkerModel{
					{
						Name:           basetypes.NewStringValue("default"),
						MachineType:    basetypes.NewStringValue("n1-medium-x86"),
						Minsize:        basetypes.NewInt64Value(1),
						Maxsize:        basetypes.NewInt64Value(3),
						Maxsurge:       basetypes.NewInt64Value(1),
						Maxunavailable: basetypes.NewInt64Value(0),
					},
				},
			},
			want: &apiv1.ClusterServiceCreateRequest{
				Name:      "cluster",
				Project:   "default-project",
				Partition: "eqx-mu4",
				Kubernetes: &apiv1.KubernetesSpec{
					Version: "1.28.11",
				},
				Workers: []*apiv1.Worker{
					{
						Name:           "default",
						MachineType:    "n1-medium-x86",
						Minsize:        1,
						Maxsize:        3,
						Maxsurge:       1,
						Maxunavailable: 0,
					},
				},
				Maintenance: &apiv1.Maintenance{
					TimeWindow: &apiv1.MaintenanceTimeWindow{
						Begin: &apiv1.Time{
							Hour:     uint32(1),
							Minute:   uint32(0),
							Timezone: "UTC",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestMapping := clusterCreateRequestMapping(tt.planMock, tt.responseMock)
			assert.Equal(t, &tt.want, &requestMapping)
		})
	}
}

func Test_clusterResponseMapping(t *testing.T) {
	tests := []struct {
		name        string
		clusterMock *apiv1.Cluster
		want        clusterModel
	}{
		{
			name: "Should map correctly",
			clusterMock: &apiv1.Cluster{
				Uuid:      "1",
				Name:      "cluster",
				Project:   "default-project",
				Partition: "eqx-mu4",
				Tenant:    "1",
				Kubernetes: &apiv1.KubernetesSpec{
					Version: "1.28.11",
				},
				Workers: []*apiv1.Worker{
					{
						Name:           "default",
						MachineType:    "n1-medium-x86",
						Minsize:        1,
						Maxsize:        3,
						Maxsurge:       1,
						Maxunavailable: 1,
					},
				},
				Maintenance: &apiv1.Maintenance{
					KubernetesAutoupdate:   basetypes.NewBoolValue(true).ValueBoolPointer(),
					MachineimageAutoupdate: basetypes.NewBoolValue(true).ValueBoolPointer(),
					TimeWindow: &apiv1.MaintenanceTimeWindow{
						Begin: &apiv1.Time{
							Hour:     uint32(1),
							Minute:   uint32(0),
							Timezone: "UTC",
						},
						Duration: &durationpb.Duration{
							Seconds: int64(3600),
						},
					},
				},
				CreatedAt: &timestamppb.Timestamp{
					Seconds: int64(1707382100),
				},
				UpdatedAt: &timestamppb.Timestamp{
					Seconds: int64(1717932877),
				},
			},
			want: clusterModel{
				Uuid:       basetypes.NewStringValue("1"),
				Name:       basetypes.NewStringValue("cluster"),
				Project:    basetypes.NewStringValue("default-project"),
				Partition:  basetypes.NewStringValue("eqx-mu4"),
				Tenant:     basetypes.NewStringValue("1"),
				Kubernetes: basetypes.NewStringValue("1.28.11"),
				Workers: []clusterWorkerModel{
					{
						Name:           basetypes.NewStringValue("default"),
						MachineType:    basetypes.NewStringValue("n1-medium-x86"),
						Minsize:        basetypes.NewInt64Value(1),
						Maxsize:        basetypes.NewInt64Value(3),
						Maxsurge:       basetypes.NewInt64Value(1),
						Maxunavailable: basetypes.NewInt64Value(1),
					},
				},
				Maintenance: &maintenanceModel{
					KubernetesAutoupdate:   basetypes.NewBoolValue(true),
					MachineimageAutoupdate: basetypes.NewBoolValue(true),
					TimeWindow: maintenanceTimeWindow{
						Begin: maintenanceTime{
							Hour:     basetypes.NewInt64Value(1),
							Minute:   basetypes.NewInt64Value(0),
							Timezone: basetypes.NewStringValue("UTC"),
						},
						Duration: basetypes.NewInt64Value(1),
					},
				},
				CreatedAt: basetypes.NewStringValue("2024-02-08 08:48:20 +0000 UTC"),
				UpdatedAt: basetypes.NewStringValue("2024-06-09 11:34:37 +0000 UTC"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseMapping := clusterResponseMapping(tt.clusterMock)
			assert.Equal(t, &tt.want, &responseMapping)
		})
	}
}
