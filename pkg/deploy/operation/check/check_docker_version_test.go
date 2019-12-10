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
	desiredDockerVersion = "18.09.0"
	versionTooLow        = "version too low"
)

// test docker if satisfied with minimal version
func TestCheckDockerVersion(t *testing.T) {
	testSample := []struct {
		comparedVersion string
		desiredVersion  string
		want            error
	}{
		{
			comparedVersion: "18.07.1-ee-12",
			desiredVersion:  desiredDockerVersion,
			want:            fmt.Errorf("%v, desired version: %v, actual version: 18.07.1-ee-12", versionTooLow, desiredDockerVersion),
		},
		{
			comparedVersion: "18.09.1",
			desiredVersion:  desiredDockerVersion,
			want:            nil,
		},
		{
			comparedVersion: "19.03.05",
			desiredVersion:  desiredDockerVersion,
			want:            nil,
		},
		{
			comparedVersion: "17.03.2-ee-8",
			desiredVersion:  desiredDockerVersion,
			want:            fmt.Errorf("%v, desired version: %v, actual version: 17.03.2-ee-8", versionTooLow, desiredDockerVersion),
		},
		{
			comparedVersion: "17.03.1-ee-7",
			desiredVersion:  desiredDockerVersion,
			want:            fmt.Errorf("%v, desired version: %v, actual version: 17.03.1-ee-7", versionTooLow, desiredDockerVersion),
		},
	}

	for _, eachValue := range testSample {
		assert.Equal(t, eachValue.want, CheckDockerVersion(eachValue.comparedVersion, eachValue.desiredVersion, ">"))
	}
}
