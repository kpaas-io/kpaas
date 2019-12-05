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

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

const (
	desiredRootDiskVolume float64 = 200
	diskStandard          float64 = desiredRootDiskVolume * operation.GiByteUnits
)

func TestCheckRootDiskVolume(t *testing.T) {
	testSample := []struct {
		rootDiskVolume    string
		desiredDiskVolume float64
		want              error
	}{
		{
			rootDiskVolume:    "287700360",
			desiredDiskVolume: diskStandard,
			want:              nil,
		},
		{
			rootDiskVolume:    "228492713",
			desiredDiskVolume: diskStandard,
			want:              nil,
		},
		{
			rootDiskVolume:    "51473368",
			desiredDiskVolume: diskStandard,
			want:              nil,
		},
		{
			rootDiskVolume:    "-103079200",
			desiredDiskVolume: diskStandard,
			want:              nil,
		},
	}

	for _, eachValue := range testSample {
		assert.Equal(t, eachValue.want, CheckRootDiskVolume(eachValue.rootDiskVolume, eachValue.desiredDiskVolume))
	}
}
