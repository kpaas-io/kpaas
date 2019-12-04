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

package system

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

func NewMemoryCapacityCheck(comparedMemory string, desiredMemory float64) error {
	var memoryCapacity float64

	memoryCapacityFloat, err := strconv.ParseFloat(comparedMemory, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error_reason": operation.ErrParaInput,
			"actual_amount": comparedMemory,
			"desire_amount": desiredMemory,
		})
		logrus.Error("parameter parse error")
		return fmt.Errorf("%v, desire memory: (%.1f)Gi, actual memory: (%v)Gi", operation.ErrParaInput, desiredMemory, comparedMemory)
	}

	memoryCapacity = memoryCapacityFloat / 1024 / 1024

	if memoryCapacity <= float64(0) {
		logrus.WithFields(logrus.Fields{
			"error_reason": operation.ErrParaInput,
			"actual_amount": comparedMemory,
			"desire_amount": desiredMemory,
		})
		logrus.Error("memory can not be negative")
		return fmt.Errorf("memory can not be negative, actual memory: (%.1f)Gi", memoryCapacity)
	}

	if memoryCapacityFloat >= desiredMemory {
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"error_reason": "memory capacity not enough",
		"actual_amount": comparedMemory,
		"desire_amount": desiredMemory,
	})
	logrus.Error("node memory not satisfied")
	return fmt.Errorf("node memory not enough, desired memory: (%v)Gi, actual memory: (%.1f)Gi", desiredMemory, memoryCapacity)
}
