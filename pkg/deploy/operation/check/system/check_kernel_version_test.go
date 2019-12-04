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

import "testing"

func TestKernelVersionCheck(t *testing.T) {
	const desiredKernelVersion = "4.19.46"
	testPassedKernelVersionArray := [...]string{"4.19.46", "4.20"}
	testFailedKernelVersionOne := "3.10.0-957.21.3.el7.x86_64"
	testFailedKernelVersionTwo := "4.18.5-041805-generic"

	for _, i := range testPassedKernelVersionArray {
		err := NewKernelVersionCheck(i, desiredKernelVersion, ".", ">")
		if err != nil {
			t.Errorf("kernel version check failed, errors: %v", err)
		} else {
			t.Logf("kernel version check passed, input version: %v, desired version: %v", i, desiredKernelVersion)
		}
	}

	err := NewKernelVersionCheck(testFailedKernelVersionOne, desiredKernelVersion, ".", ">")
	if err != nil {
		t.Errorf("kernel version check failed, errors: %v", err)
	} else {
		t.Logf("kernel version check passed, input version: %v, desired version: %v", testFailedKernelVersionOne, desiredKernelVersion)
	}

	err = NewKernelVersionCheck(testFailedKernelVersionTwo, desiredKernelVersion, ".", ">")
	if err != nil {
		t.Errorf("kernel version check failed, errors: %v", err)
	} else {
		t.Logf("kernel version check passed, input version: %v, desired version: %v", testFailedKernelVersionTwo, desiredKernelVersion)
	}

}
