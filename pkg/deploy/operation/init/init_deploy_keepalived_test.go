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

package init

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

func TestDeployKeepalived(t *testing.T) {
	testCases := []struct {
		ipAddress string
		ethernet  string
		want      error
	}{
		{
			ipAddress: "192.100.35.64",
			ethernet:  "",
			want:      fmt.Errorf(operation.ErrParaEmpty),
		},
		{
			ipAddress: "100.100.22.13",
			ethernet:  "bond0",
			want:      nil,
		},
		{
			ipAddress: "",
			ethernet:  "espn2",
			want:      fmt.Errorf(operation.ErrParaEmpty),
		},
		{
			ipAddress: "",
			ethernet:  "",
			want:      fmt.Errorf(operation.ErrParaEmpty),
		},
		{
			ipAddress: "0.0.0.0",
			ethernet:  "bond4",
			want:      nil,
		},
	}

	for _, cs := range testCases {
		assert.Equal(t, cs.want, CheckKeepalivedParameter(cs.ipAddress, cs.ethernet))
	}
}
