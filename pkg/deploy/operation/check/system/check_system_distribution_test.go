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

func TestSystemDistributionCheck(t *testing.T) {
	testPassedSystemDistributionArray := [...]string{"rhel", "centos", "ubuntu"}
	testFailedSystemDistributionOne := "macos"
	testFailedSystemDistributionTwo := "windows"

	for _, i := range testPassedSystemDistributionArray {
		err := NewSystemDistributionCheck(i)
		if err != nil {
			t.Errorf("distribution check failed, input distribution: %v , desired distribution: %v or %v or %v, errors: %v", i, DistributionCentos, DistributionUbuntu, DistributionRHEL, err)
		} else {
			t.Logf("distribution check passed, input distribution: %v", i)
		}
	}

	err := NewSystemDistributionCheck(testFailedSystemDistributionOne)
	if err != nil {
		t.Errorf("distribution check failed, input distribution: %v , desired distribution: %v or %v or %v, errors: %v", testFailedSystemDistributionOne, DistributionCentos, DistributionUbuntu, DistributionRHEL, err)
	} else {
		t.Logf("distribution check passed, input distribution: %v", testFailedSystemDistributionOne)
	}

	err = NewSystemDistributionCheck(testFailedSystemDistributionTwo)
	if err != nil {
		t.Errorf("distribution check failed, input distribution: %v , desired distribution: %v or %v or %v, errors: %v", testFailedSystemDistributionTwo, DistributionCentos, DistributionUbuntu, DistributionRHEL, err)
	} else {
		t.Logf("distribution check passed, input distribution: %v", testFailedSystemDistributionTwo)
	}
}
