/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mqttprot

import (
	"bytes"
	"io"

	"github.com/megaease/easegress/pkg/protocols"
)

type (
	// Response contains MQTT response.
	Response struct {
		drop       bool
		disconnect bool
		payload    []byte
	}
)

var _ protocols.Response = (*Response)(nil)

// NewResponse returns a new MQTT response.
func NewResponse() *Response {
	return &Response{}
}

// SetDrop means the packet in context will be drop.
func (r *Response) SetDrop() {
	r.drop = true
}

// Drop return true if the packet in context will be drop.
// For example, if SetDrop and the packet in Request is subscribe packet, MQTTProxy will
// not subscribe the topics in the packet.
func (r *Response) Drop() bool {
	return r.drop
}

// SetDisconnect means the MQTT client will be disconnect.
func (r *Response) SetDisconnect() {
	r.disconnect = true
}

// Disconnect return true if the MQTT client will be disconnect.
func (r *Response) Disconnect() bool {
	return r.disconnect
}

// Header return MQTT response header
func (r *Response) Header() protocols.Header {
	// TODO: what header to return?
	return nil
}

// SetPayload set the payload of the response to payload.
func (r *Response) SetPayload(payload []byte) {
	r.payload = payload
}

// GetPayload returns a new payload reader.
func (r *Response) GetPayload() io.Reader {
	return bytes.NewReader(r.payload)
}

// RawPayload returns the payload in []byte, the caller should
// not modify its content.
func (r *Response) RawPayload() []byte {
	return r.payload
}

// PayloadLength returns the length of the payload.
func (r *Response) PayloadLength() int {
	return len(r.payload)
}

// Close closes the response.
func (r *Response) Close() {
}
