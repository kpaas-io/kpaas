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

func TestDeployHaproxy(t *testing.T) {
	testCases := []struct {
		ipAddresses string
		want        error
	}{
		{
			ipAddresses: "192.168.3.255",
			want:        nil,
		},
		{
			ipAddresses: "255.255.255.255",
			want:        fmt.Errorf(operation.ErrInvalid),
		},
		{
			ipAddresses: "0.0.0.0",
			want:        nil,
		},
		{
			ipAddresses: "-1,-1,-1,-1",
			want:        fmt.Errorf(operation.ErrInvalid),
		},
	}

	for _, cs := range testCases {
		assert.Equal(t, cs.want, CheckHaproxyParameter(cs.ipAddresses))
	}
}
