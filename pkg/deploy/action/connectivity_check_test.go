// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	machine.IsTesting = true
}

func TestConnectivityCheck(t *testing.T) {
	executor := new(connectivityCheckExecutor)

	normalAction, err := NewConnectivityCheckAction(&ConnectivityCheckActionConfig{
		SourceNode: &pb.Node{
			Name: "normal-1",
			Ip:   "10.10.10.10",
		},
		DestinationNode: &pb.Node{
			Name: "normal-2",
			Ip:   "10.10.10.11",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, normalAction)

	pbErr := executor.Execute(normalAction)
	assert.Nil(t, pbErr)

	errorAction, err := NewConnectivityCheckAction(&ConnectivityCheckActionConfig{
		SourceNode: &pb.Node{
			Name: "error",
			Ip:   "10.10.10.10",
		},
		DestinationNode: &pb.Node{
			Name: "normal-2",
			Ip:   "10.10.10.11",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, errorAction)

	pbErr = executor.Execute(errorAction)
	assert.NoError(t, err)
	assert.NotNil(t, pbErr)
}
