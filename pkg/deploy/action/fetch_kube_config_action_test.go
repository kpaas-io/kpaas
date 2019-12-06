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

package action

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func TestNewFetchKubeConfigAction(t *testing.T) {
	// test invalid paramters
	tests := []*FetchKubeConfigActionConfig{
		nil,
		&FetchKubeConfigActionConfig{},
	}
	for _, test := range tests {
		_, err := NewFetchKubeConfigAction(test)
		assert.Error(t, err)
	}

	cfg := &FetchKubeConfigActionConfig{
		Node: &pb.Node{},
	}
	act, err := NewFetchKubeConfigAction(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, act)
	assert.IsType(t, &FetchKubeConfigAction{}, act)
	assert.Equal(t, ActionTypeFetchKubeConfig, act.GetType())
	assert.Equal(t, ActionPending, act.GetStatus())
	assert.Equal(t, cfg.Node, act.(*FetchKubeConfigAction).node)
}
