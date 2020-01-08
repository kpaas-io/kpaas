// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/action"
)

func TestActionBelongsToRole(t *testing.T) {
	type inputS struct {
		actionType action.Type
		role       constant.MachineRole
	}
	tests := []struct {
		input inputS
		want  bool
	}{
		{
			input: inputS{action.ActionTypeNodeInit, constant.MachineRoleEtcd},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeNodeInit, constant.MachineRoleMaster},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeNodeInit, constant.MachineRoleWorker},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeNodeInit, constant.MachineRoleIngress},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeDeployEtcd, constant.MachineRoleEtcd},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeInitMaster, constant.MachineRoleMaster},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeJoinMaster, constant.MachineRoleMaster},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeDeployWorker, constant.MachineRoleWorker},
			want:  true,
		},
		{
			input: inputS{action.ActionTypeDeployEtcd, constant.MachineRoleMaster},
			want:  false,
		},
		{
			input: inputS{action.ActionTypeDeployEtcd, constant.MachineRoleWorker},
			want:  false,
		},
		{
			input: inputS{action.ActionTypeDeployEtcd, constant.MachineRoleIngress},
			want:  false,
		},
		{
			input: inputS{action.ActionTypeInitMaster, constant.MachineRoleEtcd},
			want:  false,
		},
		{
			input: inputS{action.ActionTypeJoinMaster, constant.MachineRoleEtcd},
			want:  false,
		},
	}

	for _, tt := range tests {
		result := actionBelongsToRole(tt.input.actionType, tt.input.role)
		assert.Equal(t, tt.want, result)
	}
}
