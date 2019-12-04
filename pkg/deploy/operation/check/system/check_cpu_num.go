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

// new CPU numbers check task, compare with desired CPU core
func NewCPUNumsCheck(cpuCore string, desireCPUCore int) error {
	coreNums, err := strconv.Atoi(cpuCore)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error_reason": operation.ErrParaInput,
			"actual_amount": cpuCore,
			"desired_amount": desireCPUCore,
		})
		logrus.Error("parameter error")
		return fmt.Errorf("%v, desire core: %v, input core: %v cores", operation.ErrParaInput, desireCPUCore, cpuCore)
	}

	if coreNums < 0 {
		logrus.WithFields(logrus.Fields{
			"error_reason": operation.ErrParaInput,
			"actual_amount": cpuCore,
			"desired_amount": desireCPUCore,
		})
		logrus.Error("cpu core can not be negative")
		return fmt.Errorf("%v, cpu core can not be negative, input: %v cores", operation.ErrParaInput, cpuCore)
	}

	if coreNums >= desireCPUCore {
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"error_reason": "cores not enough",
		"actual_amount": cpuCore,
		"desired_amount": desireCPUCore,
	})
	logrus.Error("node CPU Core numbers not satisfied")
	return fmt.Errorf("node CPU numbers not enough, desired CPU nums: %v, actual CPU nums: %v cores", desireCPUCore, coreNums)
}
