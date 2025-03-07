// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proxycfgglue

import (
	"context"

	"github.com/hashicorp/go-memdb"

	"github.com/shulutkov/yellow-pages/acl"
	"github.com/shulutkov/yellow-pages/agent/cache"
	cachetype "github.com/shulutkov/yellow-pages/agent/cache-types"
	"github.com/shulutkov/yellow-pages/agent/consul/watch"
	"github.com/shulutkov/yellow-pages/agent/proxycfg"
	"github.com/shulutkov/yellow-pages/agent/structs"
	"github.com/shulutkov/yellow-pages/agent/structs/aclfilter"
)

// CacheGatewayServices satisfies the proxycfg.GatewayServices interface by
// sourcing data from the agent cache.
func CacheGatewayServices(c *cache.Cache) proxycfg.GatewayServices {
	return &cacheProxyDataSource[*structs.ServiceSpecificRequest]{c, cachetype.GatewayServicesName}
}

// ServerGatewayServices satisfies the proxycfg.GatewayServices interface by
// sourcing data from a blocking query against the server's state store.
func ServerGatewayServices(deps ServerDataSourceDeps) proxycfg.GatewayServices {
	return &serverGatewayServices{deps}
}

type serverGatewayServices struct {
	deps ServerDataSourceDeps
}

func (s *serverGatewayServices) Notify(ctx context.Context, req *structs.ServiceSpecificRequest, correlationID string, ch chan<- proxycfg.UpdateEvent) error {
	return watch.ServerLocalNotify(ctx, correlationID, s.deps.GetStore,
		func(ws memdb.WatchSet, store Store) (uint64, *structs.IndexedGatewayServices, error) {
			var authzContext acl.AuthorizerContext
			authz, err := s.deps.ACLResolver.ResolveTokenAndDefaultMeta(req.Token, &req.EnterpriseMeta, &authzContext)
			if err != nil {
				return 0, nil, err
			}
			if err := authz.ToAllowAuthorizer().ServiceReadAllowed(req.ServiceName, &authzContext); err != nil {
				return 0, nil, err
			}

			index, services, err := store.GatewayServices(ws, req.ServiceName, &req.EnterpriseMeta)
			if err != nil {
				return 0, nil, err
			}

			response := &structs.IndexedGatewayServices{
				Services: services,
				QueryMeta: structs.QueryMeta{
					Backend: structs.QueryBackendBlocking,
					Index:   index,
				},
			}
			aclfilter.New(authz, s.deps.Logger).Filter(response)

			return index, response, nil
		},
		dispatchBlockingQueryUpdate[*structs.IndexedGatewayServices](ch),
	)
}
