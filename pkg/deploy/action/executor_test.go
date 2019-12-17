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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const ActionTypeTestExecutorMockup Type = "ActionTypeMockupForExecutorTest"

type actionMockupForExecutorTest struct {
	Base

	goodResult bool
}

type executorMockupForExecutorTest struct{}

func (e *executorMockupForExecutorTest) Execute(act Action) *pb.Error {
	mockupAct, ok := act.(*actionMockupForExecutorTest)
	if !ok {
		return new(pb.Error)
	}
	if !mockupAct.goodResult {
		return new(pb.Error)
	}
	return nil
}

func TestRegisterExecutor(t *testing.T) {
	err := RegisterExecutor(ActionTypeTestExecutorMockup, new(executorMockupForExecutorTest))
	assert.NoError(t, err)

	exec, err := NewExecutor(ActionTypeTestExecutorMockup)
	assert.NoError(t, err)
	assert.NotNil(t, exec)
	assert.IsType(t, new(executorMockupForExecutorTest), exec)

	err = RegisterExecutor(ActionTypeTestExecutorMockup, new(executorMockupForExecutorTest))
	assert.Error(t, err)

	// cleanup
	_executorRegistry = nil
}

func TestExecuteAction(t *testing.T) {
	err := RegisterExecutor(ActionTypeTestExecutorMockup, new(executorMockupForExecutorTest))
	assert.NoError(t, err)

	input := []struct {
		action     Action
		wantStatus Status
		wantErr    bool
	}{
		{
			action: &actionMockupForExecutorTest{
				Base: Base{
					Name:       "action1",
					ActionType: ActionTypeTestExecutorMockup,
				},
				goodResult: true,
			},
			wantStatus: ActionDone,
			wantErr:    false,
		},
		{
			action: &actionMockupForExecutorTest{
				Base: Base{
					Name:       "action2",
					ActionType: ActionTypeTestExecutorMockup,
				},
				goodResult: false,
			},
			wantStatus: ActionFailed,
			wantErr:    true,
		},
	}

	for _, tt := range input {
		var wg sync.WaitGroup
		wg.Add(1)
		ExecuteAction(tt.action, &wg)
		wg.Wait()

		assert.Equal(t, tt.wantStatus, tt.action.GetStatus())
		assert.Equal(t, tt.wantErr, tt.action.GetErr() != nil)
	}

	// cleanup
	_executorRegistry = nil
}
