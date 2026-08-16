// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nclient4

import (
	"context"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// TestLeaseNilFieldGuards checks that Renew and Release reject a Lease with a nil
// field with an error instead of dereferencing it and panicking. Renew reads ACK
// (request) and Offer (server identifier); Release reads ACK.
func TestLeaseNilFieldGuards(t *testing.T) {
	valid := &dhcpv4.DHCPv4{}
	c := &Client{}

	t.Run("Renew", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			lease *Lease
		}{
			{"nil lease", nil},
			{"nil Offer", &Lease{Offer: nil, ACK: valid}},
			{"nil ACK", &Lease{Offer: valid, ACK: nil}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if _, err := c.Renew(context.Background(), tc.lease); err == nil {
					t.Fatal("Renew returned a nil error, want a nil-field error")
				}
			})
		}
	})

	t.Run("Release", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			lease *Lease
		}{
			{"nil lease", nil},
			{"nil ACK", &Lease{Offer: valid, ACK: nil}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if err := c.Release(tc.lease); err == nil {
					t.Fatal("Release returned a nil error, want a nil-field error")
				}
			})
		}
	})
}
