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

package task

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func TestNewFetchKubeConfigTask(t *testing.T) {

	// test invalid paramters
	tests := []struct {
		taskName   string
		taskConfig *FetchKubeConfigTaskConfig
	}{
		{
			taskName:   "",
			taskConfig: &FetchKubeConfigTaskConfig{},
		},
		{
			taskName:   "test",
			taskConfig: nil,
		},
		{
			taskName:   "",
			taskConfig: &FetchKubeConfigTaskConfig{},
		},
	}
	for _, test := range tests {
		_, err := NewFetchKubeConfigTask(test.taskName, test.taskConfig)
		assert.Error(t, err)
	}

	cfg := &FetchKubeConfigTaskConfig{
		Node: &pb.Node{},
	}
	aTask, err := NewFetchKubeConfigTask("test", cfg)
	assert.NoError(t, err)
	assert.NotNil(t, aTask)
	assert.IsType(t, &FetchKubeConfigTask{}, aTask)
	assert.Equal(t, TaskTypeFetchKubeConfig, aTask.GetType())
	assert.Equal(t, TaskPending, aTask.GetStatus())
	assert.Equal(t, cfg.Node, aTask.(*FetchKubeConfigTask).node)
}
