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

package docker

import (
	"testing"
)

func TestDockerVersionCheck(t *testing.T) {
	var desiredDockerVersion = "18.06.0"
	testPassedDockerVersionArray := [...]string{"18.07.1-ee-12", "18.09.1", "19.03.05"}
	testFailedDockerVersionOne := "17.03.2-ee-8"
	testFailedDockerVersionTwo := "17.03.1-ee-3"

	for _, i := range testPassedDockerVersionArray {
		err := NewDockerVersionCheck(i, desiredDockerVersion, ".", ">")
		if err != nil {
			t.Errorf("docker version check failed, input version: %v , desired version: %v, errors: %v", i, desiredDockerVersion, err)
		} else {
			t.Logf("docker version check passed, input version: %v, desired version: %v", i, desiredDockerVersion)
		}
	}

	err := NewDockerVersionCheck(testFailedDockerVersionOne, desiredDockerVersion, ".", ">")
	if err != nil {
		t.Errorf("docker version check failed, input version: %v , desired version: %v, errors: %v", testFailedDockerVersionOne, desiredDockerVersion, err)
	} else {
		t.Logf("docker version check passed, input version: %v, desired version: %v", testFailedDockerVersionOne, desiredDockerVersion)
	}


	err = NewDockerVersionCheck(testFailedDockerVersionTwo, desiredDockerVersion, ".", ">")
	if err != nil {
		t.Errorf("docker version check failed, input version: %v , desired version: %v, errors: %v", testFailedDockerVersionTwo, desiredDockerVersion, err)
	} else {
		t.Logf("docker version check passed, input version: %v, desired version: %v", testFailedDockerVersionTwo, desiredDockerVersion)
	}

}
