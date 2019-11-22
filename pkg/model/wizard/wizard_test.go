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

func TestNewWizardData(t *testing.T) {

	data := NewWizardData()
	assert.IsType(t, &WizardData{}, data)
	assert.Equal(t, ProgressSettingClusterInformation, data.Progress)
	assert.Equal(t, WizardModeNormal, data.WizardMode)
}
