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

type (
	WizardData struct {
		Progress   Progress
		WizardMode WizardMode
	}

	Progress   string
	WizardMode string
)

const (
	ProgressSettingClusterInformation Progress = "settingClusterInformation"
	ProgressSettingNodesInformation   Progress = "settingNodesInformation"
	ProgressCheckingNodes             Progress = "checkingNodes"
	ProgressDeploying                 Progress = "deploying"
	ProgressDeployCompleted           Progress = "deployCompleted"

	WizardModeNormal   WizardMode = "normal"
	WizardModeAdvanced WizardMode = "advanced"
)

func NewWizardData() *WizardData {

	data := new(WizardData)
	data.WizardMode = WizardModeNormal
	data.Progress = ProgressSettingClusterInformation
	return data
}
