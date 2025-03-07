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
	"github.com/shulutkov/yellow-pages/proto/private/pbpeering"
)

// CachePeeringList satisfies the proxycfg.PeeringList interface by sourcing
// data from the agent cache.
func CachePeeringList(c *cache.Cache) proxycfg.PeeringList {
	return &cacheProxyDataSource[*cachetype.PeeringListRequest]{c, cachetype.PeeringListName}
}

// ServerPeeringList satisfies the proxycfg.PeeringList interface by sourcing
// data from a blocking query against the server's state store.
func ServerPeeringList(deps ServerDataSourceDeps) proxycfg.PeeringList {
	return &serverPeeringList{deps}
}

type serverPeeringList struct {
	deps ServerDataSourceDeps
}

func (s *serverPeeringList) Notify(ctx context.Context, req *cachetype.PeeringListRequest, correlationID string, ch chan<- proxycfg.UpdateEvent) error {
	entMeta := structs.DefaultEnterpriseMetaInPartition(req.Request.Partition)

	return watch.ServerLocalNotify(ctx, correlationID, s.deps.GetStore,
		func(ws memdb.WatchSet, store Store) (uint64, *pbpeering.PeeringListResponse, error) {
			var authzCtx acl.AuthorizerContext
			authz, err := s.deps.ACLResolver.ResolveTokenAndDefaultMeta(req.Token, entMeta, &authzCtx)
			if err != nil {
				return 0, nil, err
			}
			if err := authz.ToAllowAuthorizer().PeeringReadAllowed(&authzCtx); err != nil {
				return 0, nil, err
			}

			index, peerings, err := store.PeeringList(ws, *entMeta)
			if err != nil {
				return 0, nil, err
			}
			return index, &pbpeering.PeeringListResponse{
				OBSOLETE_Index: index,
				Peerings:       peerings,
			}, nil
		},
		dispatchBlockingQueryUpdate[*pbpeering.PeeringListResponse](ch),
	)
}
