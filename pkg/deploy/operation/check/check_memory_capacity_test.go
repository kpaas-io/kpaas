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

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

const (
	desiredMemoryByteBase float64 = 8
	desiredMemory                 = desiredMemoryByteBase * operation.GiByteUnits
)

// unit test of CheckMemoryCapacity
func TestCheckMemoryCapacity(t *testing.T) {
	testSample := []struct {
		comparedMemory string
		desiredMemory  float64
		want           error
	}{
		{
			comparedMemory: "264116772442",
			desiredMemory:  desiredMemory,
			want:           nil,
		},
		{
			comparedMemory: "16422896121",
			desiredMemory:  desiredMemory,
			want:           nil,
		},
		{
			comparedMemory: "996123214",
			desiredMemory:  desiredMemory,
			want:           fmt.Errorf("amount not enough, desired amount: %.0f, actual amount: 996123214", desiredMemory),
		},
		{
			comparedMemory: "726123214",
			desiredMemory:  desiredMemory,
			want:           fmt.Errorf("amount not enough, desired amount: %.0f, actual amount: 726123214", desiredMemory),
		},
		{
			comparedMemory: "-1241211",
			desiredMemory:  desiredMemory,
			want:           fmt.Errorf("input parameter invalid, input parameter can not be negative, desired amount: %.0f", desiredMemory),
		},
	}

	for _, eachValue := range testSample {
		assert.Equal(t, eachValue.want, CheckMemoryCapacity(eachValue.comparedMemory, eachValue.desiredMemory))
	}

}
