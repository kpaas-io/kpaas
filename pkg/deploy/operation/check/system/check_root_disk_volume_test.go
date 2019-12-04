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
	desiredRootDiskVolume float64 = 200
	diskStandard          float64 = desiredRootDiskVolume * operation.GiByteUnits
)

func TestRootDiskVolumeCheck(t *testing.T) {
	testPassedRootVolumeArray := [...]string{"287700360", "228492713"}
	testFailedRootVolumeOne := "51473368"
	testFailedRootVolumeTwo := "-103079200"

	for _, i := range testPassedRootVolumeArray {
		err := NewRootDiskVolumeCheck(i, diskStandard)
		if err != nil {
			t.Errorf("root disk volume check failed, errors: %v", err)
		} else {
			t.Logf("root disk volume check passed, input disk volume: %v, desired disk volume: (%.1f)", i, diskStandard)
		}
	}

	err := NewRootDiskVolumeCheck(testFailedRootVolumeOne, diskStandard)
	if err != nil {
		t.Errorf("root disk volume check failed, errors: %v", err)
	} else {
		t.Logf("root disk volume check passed, input disk volume: %v, desired disk volume: (%.1f)", testFailedRootVolumeOne, diskStandard)
	}

	err = NewRootDiskVolumeCheck(testFailedRootVolumeTwo, diskStandard)
	if err != nil {
		t.Errorf("root disk volume check failed, errors: %v", err)
	} else {
		t.Logf("root disk volume check passed, input disk volume: %v, desired disk volume: (%.1f)", testFailedRootVolumeTwo, diskStandard)
	}

}
