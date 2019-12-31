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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// Mockup an action and excutor
const ActionTypeTestProcessorMockup action.Type = "ActionTypeMockupForProcessorTest"

type actionMockupForProcessorTest struct {
	action.Base
}
type executorMockupForProcessorTest struct{}

func (e *executorMockupForProcessorTest) Execute(act action.Action) *pb.Error {
	_, ok := act.(*actionMockupForProcessorTest)
	if !ok {
		return new(pb.Error)
	}
	// just for testing: do nothing.
	return nil
}

// Mockup three task types and processors
const TaskTypeTestProcessorMockup1 Type = "TaskTypeMockupForProcessorTest1"
const TaskTypeTestProcessorMockup2 Type = "TaskTypeMockupForProcessorTest2"
const TaskTypeTestProcessorMockup3 Type = "TaskTypeMockupForProcessorTest3"

type taskMockupForProcessorTest1 struct {
	Base
}
type taskMockupForProcessorTest2 struct {
	Base
}
type taskMockupForProcessorTest3 struct {
	Base
}

type processorMockupForProcessorTest1 struct{}
type processorMockupForProcessorTest2 struct{}
type processorMockupForProcessorTest3 struct{}

func (e *processorMockupForProcessorTest1) SplitTask(t Task) error {
	task1, ok := t.(*taskMockupForProcessorTest1)
	if !ok {
		return fmt.Errorf("mismatched task type")
	}
	// Split it into two mockup sub tasks
	task2 := &taskMockupForProcessorTest2{
		Base: Base{
			Name:     "task2",
			TaskType: TaskTypeTestProcessorMockup2,
			Status:   TaskPending,
			Parent:   t.GetName(),
		},
	}
	task3 := &taskMockupForProcessorTest3{
		Base: Base{
			Name:     "task3",
			TaskType: TaskTypeTestProcessorMockup3,
			Status:   TaskPending,
			Parent:   t.GetName(),
		},
	}

	task1.SubTasks = []Task{task2, task3}

	return nil
}

func (e *processorMockupForProcessorTest2) SplitTask(t Task) error {
	tsk, ok := t.(*taskMockupForProcessorTest2)
	if !ok {
		return fmt.Errorf("mismatched task type")
	}

	// Split it into one mockup action
	act1 := &actionMockupForProcessorTest{
		Base: action.Base{
			Name:       "action1",
			ActionType: ActionTypeTestProcessorMockup,
			Status:     action.ActionPending,
		},
	}
	tsk.Actions = []action.Action{act1}
	return nil
}

func (e *processorMockupForProcessorTest3) SplitTask(t Task) error {
	tsk, ok := t.(*taskMockupForProcessorTest3)
	if !ok {
		return fmt.Errorf("mismatched task type")
	}

	// Split it into one mockup action
	act2 := &actionMockupForProcessorTest{
		Base: action.Base{
			Name:       "action2",
			ActionType: ActionTypeTestProcessorMockup,
			Status:     action.ActionPending,
		},
	}
	tsk.Actions = []action.Action{act2}
	return nil
}

func TestRegisterProcessor(t *testing.T) {
	err := RegisterProcessor(TaskTypeTestProcessorMockup1, new(processorMockupForProcessorTest1))
	assert.NoError(t, err)

	proc, err := NewProcessor(TaskTypeTestProcessorMockup1)
	assert.NoError(t, err)
	assert.NotNil(t, proc)
	assert.IsType(t, new(processorMockupForProcessorTest1), proc)

	err = RegisterProcessor(TaskTypeTestProcessorMockup1, new(processorMockupForProcessorTest1))
	assert.Error(t, err)

	// cleanup
	_processRegistry = nil
}

func TestExecuteTask(t *testing.T) {
	err := action.RegisterExecutor(ActionTypeTestProcessorMockup, new(executorMockupForProcessorTest))
	assert.NoError(t, err)

	err = RegisterProcessor(TaskTypeTestProcessorMockup1, new(processorMockupForProcessorTest1))
	assert.NoError(t, err)
	err = RegisterProcessor(TaskTypeTestProcessorMockup2, new(processorMockupForProcessorTest2))
	assert.NoError(t, err)
	err = RegisterProcessor(TaskTypeTestProcessorMockup3, new(processorMockupForProcessorTest3))
	assert.NoError(t, err)

	task1 := &taskMockupForProcessorTest1{
		Base: Base{
			Name:     "task1",
			TaskType: TaskTypeTestProcessorMockup1,
			Status:   TaskPending,
		},
	}

	err = ExecuteTask(task1)
	assert.NoError(t, err)
	assert.Equal(t, TaskSuccessful, task1.GetStatus())
	assert.Nil(t, task1.GetErr())
	assert.Equal(t, 2, len(task1.GetSubTasks()))
	assert.Equal(t, 0, len(task1.GetActions()))
	for _, subTask := range task1.GetSubTasks() {
		assert.Equal(t, TaskSuccessful, subTask.GetStatus())
		assert.Nil(t, subTask.GetErr())
		assert.Equal(t, 0, len(subTask.GetSubTasks()))
		assert.Equal(t, 1, len(subTask.GetActions()))
		for _, act := range subTask.GetActions() {
			assert.Equal(t, action.ActionDone, act.GetStatus())
			assert.Nil(t, act.GetErr())
		}
	}

	// cleanup
	_processRegistry = nil
}
