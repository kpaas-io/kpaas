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
	"sort"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// Processor defines the interface for all task processors
type Processor interface {
	// No need to set the task status in this method, the caller should do that.
	SplitTask(task Task) error
}

// ExtraResult defines the interface to process the task's extra result,
// Task extra result is the specific output of a task, which is not the task's status and err.
// Task processor can implment this interface optionally.
type ExtraResult interface {
	ProcessExtraResult(task Task) error
}

// NewProcessor is a simple factory method to return a task processor based on task type.
func NewProcessor(taskType Type) (Processor, error) {
	var processor Processor
	switch taskType {
	case TaskTypeNodeCheck:
		processor = &nodeCheckProcessor{}
	case TaskTypeNodeInit:
		processor = &nodeInitProcessor{}
	case TaskTypeDeploy:
		processor = &deployProcessor{}
	case TaskTypeFetchKubeConfig:
		processor = &fetchKubeConfigProcessor{}
	case TaskTypeDeployEtcd:
		processor = &deployEtcdProcessor{}
	case TaskTypeDeployWorker:
		processor = &DeployWorkerProcessor{}
	case TaskTypeCheckNetworkRequirements:
		processor = &checkNetworkRequirementsProcessor{}
	default:
		return nil, fmt.Errorf("%s: %s", consts.MsgTaskTypeUnsupported, taskType)
	}

	return processor, nil
}

// StartTask does a basic verification on the task,
// then starts the task's execution and return immediately
func StartTask(t Task) error {
	if err := verifyTask(t); err != nil {
		logrus.Error(err)
		return err
	}

	go ExecuteTask(t)
	return nil
}

func verifyTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	if _, err := NewProcessor(t.GetType()); err != nil {
		return err
	}

	return nil
}

// ExecuteTask starts the task's execution and wait it to finish.
func ExecuteTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	logger.Debug("Start to execute Task")

	logger.Debug("Step 1: Split Task")
	if err := splitTask(t); err != nil {
		logger.Error("Failed in Step 1")
		return err
	}

	logger.Debug("Step 2: Execute Sub Tasks")
	if err := executeSubTasks(t); err != nil {
		logger.Error("Failed in Step 2")
		return err
	}

	logger.Debug("Step 3: Execute Actions")
	if err := executeActions(t); err != nil {
		logger.Error("Failed in Step 3")
		return err
	}

	logger.Debug("Step 4: Stat Task")
	if err := statTask(t); err != nil {
		logger.Error("Failed in Step 4")
		return err
	}

	logger.Debug("Step 5: Process Extra Result")
	if err := processExtraResult(t); err != nil {
		logger.Error("Failed in Step 5")
		return err
	}

	logger.Debug("Finish to execute task")
	return nil
}

func executeTaskWithWG(t Task, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	return ExecuteTask(t)
}

// Create the corresponding processor to split the task.
func splitTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	logger.Debug("Start to split task")
	t.SetStatus(TaskSplitting)

	// Create the task processor
	processor, err := NewProcessor(t.GetType())
	if err != nil {
		t.SetStatus(TaskFailed)
		t.SetErr(&pb.Error{
			Reason: "failed to create task processor",
			Detail: err.Error(),
		})
		logger.Debug("Failed to split task")
		return err
	}

	// Spilt the task
	err = processor.SplitTask(t)
	if err != nil {
		t.SetStatus(TaskFailed)
		t.SetErr(&pb.Error{
			Reason: "failed to split task",
			Detail: err.Error(),
		})
		logger.Debug("Failed to split task")
		return err
	}

	t.SetStatus(TaskSplitted)
	logger.Debug("Finish to split task")
	return nil
}

// Execute the sub tasks of a task
func executeSubTasks(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	if len(t.GetSubTasks()) == 0 {
		logger.Debug("No sub task")
		return nil
	}

	logger.Debug("Start to execute sub tasks")

	// Group the sub tasks by priority firstly.
	priTasks := prioritizeTasks(t.GetSubTasks())
	// Execute the task group sequentially.
	for _, taskGp := range priTasks {
		var wg sync.WaitGroup
		// Execute the tasks in the same group parallelly.
		for _, aSubTask := range taskGp {
			wg.Add(1)
			go executeTaskWithWG(aSubTask, &wg)
		}
		wg.Wait()

		// If any sub task in the current task group was failed, stop to execut other task groups and return.
		for _, aSubTask := range taskGp {
			if err := statTask(aSubTask); err != nil {
				return err
			}
			if aSubTask.GetStatus() != TaskDone {
				return fmt.Errorf("[%s] sub task was failed", aSubTask.GetName())
			}
		}
	}

	logger.Debug("Finish executing sub tasks")
	return nil
}

// Execute the actions of a task
func executeActions(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	if len(t.GetActions()) == 0 {
		logger.Debug("No action")
		return nil
	}

	logger.Debug("Start to execute actions")

	var wg sync.WaitGroup
	// execute the actions parallelly
	for _, act := range t.GetActions() {
		wg.Add(1)
		go action.ExecuteAction(act, &wg)
	}
	wg.Wait()

	logger.Debug("Finish to execute actions")
	return nil
}

type taskGroup []Task

// Reorder the tasks by priority, the tasks with the same priority will be in a taskGroup. Return
// a slice of taskGroup, the taskGroup in the slice is ordered: higher priority taskGroup will come
// first
func prioritizeTasks(tasks []Task) []taskGroup {
	if len(tasks) == 0 {
		return nil
	}

	mapTaskGroup := make(map[int]taskGroup)
	// First, group tasks by priority into a map
	for _, t := range tasks {
		priority := t.GetPriority()
		mapTaskGroup[priority] = append(mapTaskGroup[priority], t)
	}

	// Collect all priority values into an slice
	keys := make([]int, 0, len(mapTaskGroup))
	for k := range mapTaskGroup {
		keys = append(keys, k)
	}
	// Sort the slice
	sort.Ints(keys)

	// Iterate the sorted slice and add the corresponding taskGoup to a new slice, the taskGroup
	// in the new slice will be in order too.
	prioritizedTasks := make([]taskGroup, 0, len(keys))
	for _, k := range keys {
		prioritizedTasks = append(prioritizedTasks, mapTaskGroup[k])
	}
	return prioritizedTasks
}

// Analyze the task status according to its sub tasks and actions.
func statTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	logger.Debug("Start to gen task summary")

	// If the task has sub taks, analyze the status of all of its sub tasks one by one.
	for _, subTask := range t.GetSubTasks() {
		if err := statTask(subTask); err != nil {
			t.SetStatus(TaskFailed)
			t.SetErr(&pb.Error{
				Reason: "failed to generate the task summary",
				Detail: err.Error(),
			})
			logger.Error("Finish to gen task summary")
			return err
		}
	}

	done := 0
	failed := 0
	// combined error message in sub tasks and actions
	var errMsgs []string

	for _, subTask := range t.GetSubTasks() {
		switch subTask.GetStatus() {
		case TaskFailed:
			failed++
			errMsgs = append(errMsgs, fmt.Sprintf("%v", subTask.GetErr()))
		case TaskDone:
			done++
		}
	}

	for _, act := range t.GetActions() {
		switch act.GetStatus() {
		case action.ActionFailed:
			failed++
			errMsgs = append(errMsgs, fmt.Sprintf("%v", act.GetErr()))
		case action.ActionDone:
			done++
		}
	}

	// if any subtask or action is failed, the task is failed
	if failed > 0 {
		t.SetStatus(TaskFailed)
		t.SetErr(&pb.Error{
			Reason:     "one or more checks failed",
			Detail:     fmt.Sprintf("%v", errMsgs),
			FixMethods: "check the detail mssage",
		})
	} else if done == len(t.GetSubTasks())+len(t.GetActions()) {
		// if all subtasks/actions are done, the task is done
		t.SetStatus(TaskDone)
	}

	logger.Debug("Finish to gen task summary")
	return nil
}

func processExtraResult(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	// Create the task processor
	processor, err := NewProcessor(t.GetType())
	if err != nil {
		t.SetStatus(TaskFailed)
		t.SetErr(&pb.Error{
			Reason: "failed to create task processor",
			Detail: err.Error(),
		})
		logger.Debug("Failed to split task")
		return err
	}

	// Check if the processor implemented the ExraResult interface
	extraResult, ok := processor.(ExtraResult)
	if !ok {
		return nil
	}
	return extraResult.ProcessExtraResult(t)
}
