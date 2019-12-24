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
	desiredRootDiskByteVolume float64 = 200
	diskStandard              float64 = desiredRootDiskByteVolume * operation.GiByteUnits
)

// unit test of CheckRootDiskVolume
func TestCheckRootDiskVolume(t *testing.T) {
	testSample := []struct {
		rootDiskVolume    string
		desiredDiskVolume float64
		want              error
	}{
		{
			rootDiskVolume:    "287700360772",
			desiredDiskVolume: diskStandard,
			want:              nil,
		},
		{
			rootDiskVolume:    "228492713842",
			desiredDiskVolume: diskStandard,
			want:              nil,
		},
		{
			rootDiskVolume:    "51473368953",
			desiredDiskVolume: diskStandard,
			want:              fmt.Errorf("amount not enough, desired amount: %.0f, actual amount: 51473368953", diskStandard),
		},
		{
			rootDiskVolume:    "-1",
			desiredDiskVolume: diskStandard,
			want:              fmt.Errorf("input parameter invalid, input parameter can not be negative, desired amount: %.0f", diskStandard),
		},
	}

	for _, eachValue := range testSample {
		assert.Equal(t, eachValue.want, CheckRootDiskVolume(eachValue.rootDiskVolume, eachValue.desiredDiskVolume))
	}
}
