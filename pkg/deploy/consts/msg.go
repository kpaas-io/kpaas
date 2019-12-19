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

package consts

const (
	MsgRequestFailed string = "the request was failed"

	MsgUnknownFixMethod string = "unknown"

	// Task related messages
	MsgTaskTypeUnsupported         string = "unsupported task type"
	MsgTaskSplitFailed             string = "failed to split task"
	MsgTaskTypeMismatched          string = "task type mismatched"
	MsgEmptyTask                   string = "empty task"
	MsgTaskProcessorCreationFailed string = "failed to create task processor"
	MsgTaskGenSummaryFailed        string = "failed to generate task summary"

	// Action related messages
	MsgActionTypeUnsupported         string = "unsupported action type"
	MsgActionExecutorCreationFailed  string = "failed to create action executor"
	MsgActionExecutionFailed         string = "failed to execute aciton"
	MsgActionTypeMismatched          string = "action type mismatched"
	MsgActionTypeMismatchedDetail    string = "the action type is not match: should be %T, but is %T"
	MsgActionInvalidConfig           string = "the action config is invalid"
	MsgActionInvalidConfigNodeNotSet string = "the action's target node is not set"
	MsgEmptyAction                   string = "empty action"
)
