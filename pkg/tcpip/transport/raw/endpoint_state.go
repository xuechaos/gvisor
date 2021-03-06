// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package raw

import (
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/buffer"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

// saveData saves packet.data field.
func (p *packet) saveData() buffer.VectorisedView {
	// We cannot save p.data directly as p.data.views may alias to p.views,
	// which is not allowed by state framework (in-struct pointer).
	return p.data.Clone(nil)
}

// loadData loads packet.data field.
func (p *packet) loadData(data buffer.VectorisedView) {
	// NOTE: We cannot do the p.data = data.Clone(p.views[:]) optimization
	// here because data.views is not guaranteed to be loaded by now. Plus,
	// data.views will be allocated anyway so there really is little point
	// of utilizing p.views for data.views.
	p.data = data
}

// beforeSave is invoked by stateify.
func (ep *endpoint) beforeSave() {
	// Stop incoming packets from being handled (and mutate endpoint state).
	// The lock will be released after saveRcvBufSizeMax(), which would have
	// saved ep.rcvBufSizeMax and set it to 0 to continue blocking incoming
	// packets.
	ep.rcvMu.Lock()
}

// saveRcvBufSizeMax is invoked by stateify.
func (ep *endpoint) saveRcvBufSizeMax() int {
	max := ep.rcvBufSizeMax
	// Make sure no new packets will be handled regardless of the lock.
	ep.rcvBufSizeMax = 0
	// Release the lock acquired in beforeSave() so regular endpoint closing
	// logic can proceed after save.
	ep.rcvMu.Unlock()
	return max
}

// loadRcvBufSizeMax is invoked by stateify.
func (ep *endpoint) loadRcvBufSizeMax(max int) {
	ep.rcvBufSizeMax = max
}

// afterLoad is invoked by stateify.
func (ep *endpoint) afterLoad() {
	// StackFromEnv is a stack used specifically for save/restore.
	ep.stack = stack.StackFromEnv

	// If the endpoint is connected, re-connect via the save/restore stack.
	if ep.connected {
		var err *tcpip.Error
		ep.route, err = ep.stack.FindRoute(ep.registeredNIC, ep.boundAddr, ep.route.RemoteAddress, ep.netProto, false)
		if err != nil {
			panic(*err)
		}
	}

	// If the endpoint is bound, re-bind via the save/restore stack.
	if ep.bound {
		if ep.stack.CheckLocalAddress(ep.registeredNIC, ep.netProto, ep.boundAddr) == 0 {
			panic(tcpip.ErrBadLocalAddress)
		}
	}

	if err := ep.stack.RegisterRawTransportEndpoint(ep.registeredNIC, ep.netProto, ep.transProto, ep); err != nil {
		panic(*err)
	}
}
