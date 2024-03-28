package cluster_test

import (
	"strings"
	"testing"

	"math/rand/v2"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/provider"
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"metal": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

func TestAccClusterResourceAndDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleClusterSeed,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("metal_cluster.acctest", "name", "tf-c-"+runId),
				),
			},
			{
				Config: testAccExampleClusterSeed + testAccExampleDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_cluster.acctest-data", "kubernetes", "1.27.11"),
				),
			},
			{
				Config: testAccExampleClusterSeedWithAllFields,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("metal_cluster.acctest", "name", "tf-c-"+runId),
				),
			},
		},
	})
}

var (
	runId = func(n int) string {
		const letters = "abcdefghijklmnopqrstuvwxyz1234567890"
		var str strings.Builder
		for range n {
			str.WriteByte(letters[rand.N(len(letters))])
		}
		return str.String()
	}(5)

	testAccExampleClusterSeed = `
resource "metal_cluster" "acctest" {
	name = "tf-c-` + runId + `"
	kubernetes = "1.27.11"
	workers = [
		{
			name = "group-0"
			machine_type = "c1-medium-x86"
			max_size = 2
			min_size = 1
		}
	]
	maintenance = {
		time_window = {
			begin = {
			  hour   = 18
			  minute = 30
			}
			duration = 2
		  }
	}
}
`

	testAccExampleDataSource = `
data "metal_cluster" "acctest-data" {
	name = "tf-c-` + runId + `"
}
`

	testAccExampleClusterSeedWithAllFields = `
resource "metal_cluster" "acctest" {
	name = "tf-c-` + runId + `"
	kubernetes = "1.27.11"
	workers = [
		{
			name = "group-0"
			machine_type = "c1-medium-x86"
			max_size = 5
			min_size = 1
			max_surge = 3
			max_unavailable = 2
		}
	]
	maintenance = {
		time_window = {
			begin = {
			  hour   = 18
			  minute = 30
			}
			duration = 2
		  }
	}
}
`
)
