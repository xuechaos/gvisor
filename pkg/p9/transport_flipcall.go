// Copyright 2019 The gVisor Authors.
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

package p9

import (
	"runtime"
	"syscall"

	"gvisor.dev/gvisor/pkg/fd"
	"gvisor.dev/gvisor/pkg/fdchannel"
	"gvisor.dev/gvisor/pkg/flipcall"
	"gvisor.dev/gvisor/pkg/log"
)

// channelsPerClient is the number of channels to create per client.
//
// While the client and server will generally agree on this number, in reality
// it's completely up to the server. We simply define a minimum of 2, and a
// maximum of 4, and select the number of available processes as a tie-breaker.
// Note that we don't want the number of channels to be too large, because each
// will account for channelSize memory used, which can be large.
var channelsPerClient = func() int {
	n := runtime.NumCPU()
	if n < 2 {
		return 2
	}
	if n > 4 {
		return 4
	}
	return n
}()

// channelsTotalMaximum is the total number of channels allowed.
const channelsTotalMaximum = 32

// channelsTotal is the current number of channels created.
var channelsTotal int32

// channelHeader is the extra space reserved in a channel.
const channelHeader = 32

// channelSize is the channel size to create.
//
// We simply ensure that this is larger than the largest possible message size.
// We don't really bother calculating the precise header, because at the time
// of writing, the flipcall channel itself required an additional 16 bytes of
// space, and so we provide a generous buffer.
const channelSize = int(maximumLength + channelHeader)

// channel is a fast IPC channel.
//
// The same object is used by both the server and client implementations. In
// general, the client will use only the send and recv methods.
type channel struct {
	desc flipcall.PacketWindowDescriptor
	data *flipcall.Endpoint
	fds  *fdchannel.Endpoint

	// -- server only --
	client *fd.FD
	done   chan struct{}
}

// service services the channel.
func (ch *channel) service(cs *connState) error {
	defer close(ch.done)
	rsz, err := ch.data.RecvFirst()
	if err != nil {
		return err
	}
	for rsz > 0 {
		m, err := ch.recv(nil, rsz)
		if err != nil {
			return err
		}
		r := cs.handle(m)
		msgRegistry.put(m)
		rsz, err = ch.send(r)
		if err != nil {
			return err
		}
	}
	return nil // Done.
}

// Close closes the channel.
//
// This must only be called once, and cannot return an error.
func (ch *channel) Close() error {
	// Stop inflight messages.
	ch.data.Shutdown()
	if ch.done != nil {
		// Wait for the service routine, if required, which may be
		// accessing the underlying mappings. If there is no service
		// routine (i.e. this is the client calling), then they
		// shouldn't be calling sendRecv and Close simultaneously.
		<-ch.done
	}

	// Close all backing transports.
	ch.fds.Destroy()
	ch.data.Destroy()
	if ch.client != nil {
		ch.client.Close()
	}
	return nil
}

// send sends the given message.
//
// The return value is the size of the received response. Not that in the
// server case, this is the size of the next request.
func (ch *channel) send(m message) (uint32, error) {
	if log.IsLogging(log.Debug) {
		log.Debugf("send %s", m.String())
	}

	// Send any file payload.
	sentFD := false
	if filer, ok := m.(filer); ok {
		if f := filer.FilePayload(); f != nil {
			if err := ch.fds.SendFD(f.FD()); err != nil {
				return 0, err
			}
			f.Close()     // Per sendRecvLegacy.
			sentFD = true // To mark below.
		}
	}

	// Encode the message.
	//
	// Note that IPC itself encodes the length of messages, so we don't
	// need to encode a standard 9P header. We write only the message type.
	data := ch.data.Data()
	buf := buffer{data: data[:0]}
	buf.WriteMsgType(m.Type())
	if sentFD {
		buf.Write8(1) // Incoming FD.
	} else {
		buf.Write8(0) // No incoming FD.
	}
	m.Encode(&buf)
	ssz := uint32(len(buf.data)) // Updated below.

	// Is there a payload?
	if payloader, ok := m.(payloader); ok {
		p := payloader.Payload()
		copy(data[ssz:], p)
		ssz += uint32(len(p))
	}

	// Perform the one-shot communication.
	return ch.data.SendRecv(ssz)
}

// recv decodes a message that exists on the channel.
//
// If the passed r is non-nil, then the type must match or an error will be
// generated. If the passed r is nil, then a new message will be created and
// returned.
func (ch *channel) recv(r message, rsz uint32) (message, error) {
	// Decode the response from the inline buffer.
	data := ch.data.Data()
	buf := buffer{data: data[:rsz]}
	t := buf.ReadMsgType()
	hasFD := buf.Read8() != 0
	if t == MsgRlerror {
		// Change the message type. We check for this special case
		// after decoding below, and transform into an error.
		r = &Rlerror{}
	} else if r == nil {
		nr, err := msgRegistry.get(0, t)
		if err != nil {
			return nil, err
		}
		r = nr // New message.
	} else if t != r.Type() {
		// Not an error and not the expected response; propagate.
		return nil, &ErrBadResponse{Got: t, Want: r.Type()}
	}

	// Is there a payload? Set to the latter portion.
	if payloader, ok := r.(payloader); ok {
		fs := payloader.FixedSize()
		payloader.SetPayload(buf.data[fs:])
		buf.data = buf.data[:fs]
	}

	r.Decode(&buf)
	if buf.isOverrun() {
		// Nothing valid was available.
		log.Debugf("recv [got %d bytes, needed more]", rsz)
		return nil, ErrNoValidMessage
	}

	// Read any FD result.
	if hasFD {
		if rfd, err := ch.fds.RecvFDNonblock(); err == nil {
			f := fd.New(rfd)
			if filer, ok := r.(filer); ok {
				// Set the payload.
				filer.SetFilePayload(f)
			} else {
				// Don't want the FD.
				f.Close()
			}
		}
	}

	// Log a message.
	if log.IsLogging(log.Debug) {
		log.Debugf("recv %s", r.String())
	}

	// Convert errors appropriately; see above.
	if rlerr, ok := r.(*Rlerror); ok {
		return nil, syscall.Errno(rlerr.Error)
	}

	return r, nil
}

// sendRecv sends the given message over the channel.
//
// This is used by the client.
func (ch *channel) sendRecv(c *Client, m, r message) error {
	rsz, err := ch.send(m)
	if err != nil {
		return err
	}
	_, err = ch.recv(r, rsz)
	return err
}
