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
)

func TestNewCluster(t *testing.T) {

	assert.IsType(t, &Cluster{}, NewCluster())
}

func TestNewClusterInfo(t *testing.T) {

	assert.IsType(t, &ClusterInfo{}, NewClusterInfo())
}

func TestGetCurrentWizard(t *testing.T) {

	cluster := GetCurrentWizard()
	assert.Equal(t, DeployClusterStatusNotRunning, cluster.DeployClusterStatus)
	assert.NotNil(t, cluster.Wizard)
	assert.NotNil(t, cluster.Nodes)
	assert.NotNil(t, cluster.Info)
	assert.Equal(t, ProgressSettingClusterInformation, cluster.Wizard.Progress)
	assert.Equal(t, WizardModeNormal, cluster.Wizard.WizardMode)
}

func TestClearCurrentWizardData(t *testing.T) {

	var cluster *Cluster
	cluster = GetCurrentWizard()
	cluster.Wizard.Progress = ProgressSettingNodesInformation
	cluster = GetCurrentWizard()
	assert.Equal(t, ProgressSettingNodesInformation, cluster.Wizard.Progress)

	ClearCurrentWizardData()
	cluster = GetCurrentWizard()
	assert.Equal(t, ProgressSettingClusterInformation, cluster.Wizard.Progress)
}

func TestNewKubeAPIServerConnectionData(t *testing.T) {

	data := NewKubeAPIServerConnectionData()
	assert.IsType(t, &KubeAPIServerConnectionData{}, data)
	assert.Equal(t, KubeAPIServerConnectTypeFirstMasterIP, data.KubeAPIServerConnectType)
}
