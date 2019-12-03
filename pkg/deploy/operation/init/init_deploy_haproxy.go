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

package init

import (
	"fmt"

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	"github.com/sirupsen/logrus"
)

func NewDeployHaproxy(IPaddress ...string) error {
	if len(IPaddress) == 0 {
		logrus.WithFields(logrus.Fields{
			"error_reason": "parameter error",
		})
		logrus.Error("IP address can not be empty")
		return fmt.Errorf("input IP address can not be empty, input: (%v)", IPaddress)
	}
	for _, v := range IPaddress {
		if ok := operation.IPValidationCheck(v); ok {
			return nil
		}

		logrus.WithFields(logrus.Fields{
			"error_reason": "input IP address error",
		})
		logrus.Error("IP address is not valid: %v", v)
		return fmt.Errorf("input IP address is not validate: (%v)", IPaddress)
	}

	logrus.WithFields(logrus.Fields{
		"error_reason": "input IP address error",
	})
	logrus.Error("IP addresses is not valid: %v", IPaddress)
	return fmt.Errorf("input IP address (%v) is not validate", IPaddress)
}
