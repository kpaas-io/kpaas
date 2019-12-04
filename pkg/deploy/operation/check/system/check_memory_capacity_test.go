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

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

const (
	desiredMemoryBase  float64 = 16
	desiredMemory = desiredMemoryBase * operation.GiByteUnits
)

func TestMemoryCapacityCheck(t *testing.T) {
	testPassedMemoryArray := [...]string{"264116772", "16422896", "16267396"}
	testFailedMemoryOne := "1626123"
	testFailedMemoryTwo := "-1241211"

	for _, i := range testPassedMemoryArray {
		err := NewMemoryCapacityCheck(i, desiredMemory)
		if err != nil {
			t.Errorf("memory check failed, errors: %v", err)
		} else {
			t.Logf("memory check passed, input memory: %v, desired memory: %v", i, desiredMemory)
		}
	}

	err := NewMemoryCapacityCheck(testFailedMemoryOne, desiredMemory)
	if err != nil {
		t.Errorf("memory check failed, errors: %v", err)
	} else {
		t.Logf("memory check passed, input memory: %v, desired memory: %v", testFailedMemoryOne, desiredMemory)
	}

	err = NewMemoryCapacityCheck(testFailedMemoryTwo, desiredMemory)
	if err != nil {
		t.Errorf("memory check failed, errors: %v", err)
	} else {
		t.Logf("memory check passed, input memory: %v, desired memory: %v", testFailedMemoryOne, desiredMemory)
	}
}
