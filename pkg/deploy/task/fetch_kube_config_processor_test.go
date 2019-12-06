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

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func TestVerifyTask(t *testing.T) {
	tests := []struct {
		task Task
		want bool
	}{
		{
			task: nil,
			want: false,
		},
		{
			task: new(nodeCheckTask),
			want: false,
		},
		{
			task: &FetchKubeConfigTask{
				node: new(pb.Node),
			},
			want: true,
		},
	}
	processor := new(fetchKubeConfigProcessor)
	for _, tt := range tests {
		err := processor.verifyTask(tt.task)
		assert.Equal(t, tt.want, err == nil)
	}
}

func TestSplitTask(t *testing.T) {
	taskCfg := FetchKubeConfigTaskConfig{
		Node: &pb.Node{
			Name: "testnodename",
			Ip:   "0.0.0.0",
		},
	}

	kubeConfigTask, err := NewFetchKubeConfigTask("test-task", &taskCfg)
	assert.NoError(t, err)

	processor := new(fetchKubeConfigProcessor)
	err = processor.SplitTask(kubeConfigTask)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(kubeConfigTask.GetActions()))
}

func TestProcessExtraResult(t *testing.T) {
	node := &pb.Node{
		Name: "testnodename",
		Ip:   "0.0.0.0",
	}
	taskCfg := FetchKubeConfigTaskConfig{
		Node: node,
	}

	kubeConfigTask, err := NewFetchKubeConfigTask("test-task", &taskCfg)
	assert.NoError(t, err)

	actCfg := &action.FetchKubeConfigActionConfig{
		Node: node,
	}
	kubeconfigAct, err := action.NewFetchKubeConfigAction(actCfg)
	assert.NoError(t, err)
	kubeConfigContent := "test content"
	kubeconfigAct.(*action.FetchKubeConfigAction).KubeConfig = kubeConfigContent
	kubeConfigTask.(*FetchKubeConfigTask).actions = []action.Action{kubeconfigAct}

	processor := new(fetchKubeConfigProcessor)
	err = processor.ProcessExtraResult(kubeConfigTask)
	assert.NoError(t, err)
	assert.Equal(t, kubeConfigContent, kubeConfigTask.(*FetchKubeConfigTask).KubeConfig)
}
