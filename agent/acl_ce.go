// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !consulent
// +build !consulent

package agent

import (
	"github.com/hashicorp/serf/serf"

	"github.com/shulutkov/yellow-pages/acl"
	"github.com/shulutkov/yellow-pages/api"
)

func serfMemberFillAuthzContext(m *serf.Member, ctx *acl.AuthorizerContext) {
	// no-op
}

func agentServiceFillAuthzContext(s *api.AgentService, ctx *acl.AuthorizerContext) {
	// no-op
}
