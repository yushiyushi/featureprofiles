/*
 Copyright 2022 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package bgp_base_test

import (
	"testing"
	"time"

	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/telemetry"
	"github.com/openconfig/ygot/ygot"
)

const (
	dutDesc    = "To ATE"
	dutIPv4    = "10.244.0.70"
	dutIPv4Len = 30

	ateDesc    = "To DUT"
	ateIPv4    = "10.244.0.71"
	ateIPv4Len = 30

	dutAS = 64500
	ateAS = 64501
)

func bgpWithNbr(as uint32, routerID string, nbr *telemetry.NetworkInstance_Protocol_Bgp_Neighbor) *telemetry.NetworkInstance_Protocol_Bgp {
	bgp := &telemetry.NetworkInstance_Protocol_Bgp{}
	bgp.GetOrCreateGlobal().As = ygot.Uint32(as)
	if routerID != "" {
		bgp.Global.RouterId = ygot.String(routerID)
	}
	bgp.AppendNeighbor(nbr)
	return bgp
}

func TestEstablish(t *testing.T) {
	dut := ondatra.DUT(t, "fakedut")
	ate := ondatra.DUT(t, "fakedut2")

	dutConfPath := dut.Config().NetworkInstance("default").Protocol(telemetry.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	ateConfPath := ate.Config().NetworkInstance("default").Protocol(telemetry.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	// Remove any existing BGP config
	dutConfPath.Delete(t)
	ateConfPath.Delete(t)

	statePath := dut.Telemetry().NetworkInstance("default").Protocol(telemetry.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	nbrPath := statePath.Neighbor(ateIPv4)
	// Start a new session
	dutConf := bgpWithNbr(dutAS, dutIPv4, &telemetry.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(ateIPv4),
	})
	ateConf := bgpWithNbr(dutAS, ateIPv4, &telemetry.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(dutIPv4),
	})
	dutConfPath.Replace(t, dutConf)
	ateConfPath.Replace(t, ateConf)
	nbrPath.SessionState().Await(t, time.Second*15, telemetry.Bgp_Neighbor_SessionState_ESTABLISHED)
}