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

package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	desireCPUCore float64 = 8
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
			desiredCore: desireCPUCore,
			want:        nil,
		},
		{
			cpuCore:     "32767",
			desiredCore: desireCPUCore,
			want:        nil,
		},
		{
			cpuCore:     "100000000",
			desiredCore: desireCPUCore,
			want:        nil,
		},
		{
			cpuCore:     "7",
			desiredCore: desireCPUCore,
			want:        nil,
		},
		{
			cpuCore:     "-100",
			desiredCore: desireCPUCore,
			want:        nil,
		},
	}

	for _, eachValue := range testSample {
		assert.Equal(t, eachValue.want, CheckCPUNums(eachValue.cpuCore, eachValue.desiredCore))
	}
}
