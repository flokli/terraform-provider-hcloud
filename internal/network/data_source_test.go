package network_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/terraform-providers/terraform-provider-hcloud/internal/network"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-hcloud/internal/loadbalancer"
	"github.com/terraform-providers/terraform-provider-hcloud/internal/testsupport"
	"github.com/terraform-providers/terraform-provider-hcloud/internal/testtemplate"
)

func TestAccHcloudDataSourceNetworkTest(t *testing.T) {
	tmplMan := testtemplate.Manager{}

	res := &network.RData{
		Name:    "network-ds-test",
		IPRange: "10.0.0.0/16",
		Labels: map[string]string{
			"key": strconv.Itoa(acctest.RandInt()),
		},
	}
	res.SetRName("network-ds-test")
	networkByName := &network.DData{
		NetworkName: res.TFID() + ".name",
	}
	networkByName.SetRName("network_by_name")
	networkByID := &network.DData{
		NetworkID: res.TFID() + ".id",
	}
	networkByID.SetRName("network_by_id")
	networkBySel := &network.DData{
		LabelSelector: fmt.Sprintf("key=${%s.labels[\"key\"]}", res.TFID()),
	}
	networkBySel.SetRName("network_by_sel")

	resource.Test(t, resource.TestCase{
		PreCheck:     testsupport.AccTestPreCheck(t),
		Providers:    testsupport.AccTestProviders(),
		CheckDestroy: testsupport.CheckResourcesDestroyed(loadbalancer.ResourceType, loadbalancer.ByID(t, nil)),
		Steps: []resource.TestStep{
			{
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_network", res,
					"testdata/d/hcloud_network", networkByName,
					"testdata/d/hcloud_network", networkByID,
					"testdata/d/hcloud_network", networkBySel,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkByName.TFID(),
						"name", fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt)),
					resource.TestCheckResourceAttr(networkByName.TFID(), "ip_range", res.IPRange),

					resource.TestCheckResourceAttr(networkByID.TFID(),
						"name", fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt)),
					resource.TestCheckResourceAttr(networkByID.TFID(), "ip_range", res.IPRange),

					resource.TestCheckResourceAttr(networkBySel.TFID(),
						"name", fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt)),
					resource.TestCheckResourceAttr(networkBySel.TFID(), "ip_range", res.IPRange),
				),
			},
		},
	})
}