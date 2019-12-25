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

package check

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	desiredCPUCore float64 = 4
)

// unit test of CheckCPUNums
func TestCheckCPUNums(t *testing.T) {
	testSample := []struct {
		cpuCore     string
		desiredCore float64
		want        error
	}{
		{
			cpuCore:     "8",
			desiredCore: desiredCPUCore,
			want:        nil,
		},
		{
			cpuCore:     "32767",
			desiredCore: desiredCPUCore,
			want:        nil,
		},
		{
			cpuCore:     "100000000",
			desiredCore: desiredCPUCore,
			want:        nil,
		},
		{
			cpuCore:     "3",
			desiredCore: desiredCPUCore,
			want:        fmt.Errorf("amount not enough, desired amount: %.0f, actual amount: 3", desiredCPUCore),
		},
		{
			cpuCore:     "-100",
			desiredCore: desiredCPUCore,
			want:        fmt.Errorf("input parameter invalid, input parameter can not be negative, desired amount: %.0f", desiredCPUCore),
		},
	}

	for _, eachValue := range testSample {
		assert.Equal(t, eachValue.want, CheckCPUNums(eachValue.cpuCore, eachValue.desiredCore))
	}
}
