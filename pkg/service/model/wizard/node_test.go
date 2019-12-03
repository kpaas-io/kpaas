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

package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
)

func TestNewNode(t *testing.T) {

	node := NewNode()
	assert.IsType(t, &Node{}, node)
	assert.NotNil(t, node.MachineRoles)
	assert.NotNil(t, node.Labels)
	assert.NotNil(t, node.Taints)
	assert.NotNil(t, node.CheckReport)
	assert.NotNil(t, node.DeploymentReports)
	assert.Equal(t, AuthenticationTypePassword, node.AuthenticationType)
}

func TestNewDeploymentReport(t *testing.T) {

	report := NewDeploymentReport()
	assert.NotNil(t, report)
	assert.Equal(t, DeployStatusPending, report.Status)
	assert.Nil(t, report.Error)
}

func TestNewCheckItem(t *testing.T) {

	item := NewCheckItem()
	assert.Equal(t, constant.CheckResultNotRunning, item.CheckResult)
}
