// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proxycfgglue

import (
	"context"

	"github.com/shulutkov/yellow-pages/agent/leafcert"
	"github.com/shulutkov/yellow-pages/agent/proxycfg"
)

// LocalLeafCerts satisfies the proxycfg.LeafCertificate interface by sourcing data from
// the given leafcert.Manager.
func LocalLeafCerts(m *leafcert.Manager) proxycfg.LeafCertificate {
	return &localLeafCerts{m}
}

type localLeafCerts struct {
	leafCertManager *leafcert.Manager
}

func (c *localLeafCerts) Notify(ctx context.Context, req *leafcert.ConnectCALeafRequest, correlationID string, ch chan<- proxycfg.UpdateEvent) error {
	return c.leafCertManager.NotifyCallback(ctx, req, correlationID, dispatchCacheUpdate(ch))
}
