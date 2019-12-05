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

func TestCheckSystemDistribution(t *testing.T) {
	testSample := []struct {
		disName string
		want    error
	}{
		{
			disName: "rhel",
			want:    nil,
		},
		{
			disName: "centos",
			want:    nil,
		},
		{
			disName: "guess what",
			want:    nil,
		},
		{
			disName: "macos",
			want:    nil,
		},
		{
			disName: "",
			want:    nil,
		},
	}

	for _, eachValue := range testSample {
		assert.Equal(t, eachValue.want, CheckSystemDistribution(eachValue.disName))
	}
}
