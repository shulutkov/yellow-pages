// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !consulent
// +build !consulent

package consul

import (
	memdb "github.com/hashicorp/go-memdb"

	"github.com/shulutkov/yellow-pages/acl"
	"github.com/shulutkov/yellow-pages/agent/consul/authmethod"
	"github.com/shulutkov/yellow-pages/agent/consul/state"
	"github.com/shulutkov/yellow-pages/agent/structs"
)

func (a *ACL) tokenUpsertValidateEnterprise(token *structs.ACLToken, existing *structs.ACLToken) error {
	state := a.srv.fsm.State()
	return state.ACLTokenUpsertValidateEnterprise(token, existing)
}

func (a *ACL) policyUpsertValidateEnterprise(policy *structs.ACLPolicy, existing *structs.ACLPolicy) error {
	state := a.srv.fsm.State()
	return state.ACLPolicyUpsertValidateEnterprise(policy, existing)
}

func (a *ACL) roleUpsertValidateEnterprise(role *structs.ACLRole, existing *structs.ACLRole) error {
	state := a.srv.fsm.State()
	return state.ACLRoleUpsertValidateEnterprise(role, existing)
}

func enterpriseAuthMethodValidation(method *structs.ACLAuthMethod, validator authmethod.Validator) error {
	return nil
}

func getTokenNamespaceDefaults(ws memdb.WatchSet, state *state.Store, entMeta *acl.EnterpriseMeta) ([]string, []string, error) {
	return nil, nil, nil
}
