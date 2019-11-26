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

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

// Processor defines the interface for all task processors
type Processor interface {
	StartTask(task Task) error
}

// NewProcessor is a simple factory method to return a task processor based on task type.
func NewProcessor(taskType Type) (Processor, error) {
	var processor Processor
	switch taskType {
	case TaskTypeNodeCheck:
		processor = &nodeCheckProcessor{}
	default:
		return nil, fmt.Errorf("%s: %s", consts.MsgTaskTypeUnsupported, taskType)
	}

	return processor, nil
}
