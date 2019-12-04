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

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

const (
	DistributionCentos string = "centos"
	DistributionUbuntu string = "ubuntu"
	DistributionRHEL   string = "rhel"
)

// check if system distribution can be supported
func NewSystemDistributionCheck(disName string) error {
	if disName == "" {
		logrus.WithFields(logrus.Fields{
			"error_reason": operation.ErrParaInput,
			"actual_value": disName,
			"desire_value": fmt.Sprintf("supported distribution: '%v' or '%v' or '%v'", DistributionCentos, DistributionUbuntu, DistributionRHEL),
		})
		logrus.Errorf("parameter error")
		return fmt.Errorf("%v, can not be empty", operation.ErrParaInput)
	}

	switch disName {
	case DistributionCentos:
		logrus.WithFields(logrus.Fields{
			"actual_value": disName,
		})
		logrus.Infof("distribution check passed, current version: %v", disName)
		return nil
	case DistributionUbuntu:
		logrus.WithFields(logrus.Fields{
			"actual_value": disName,
		})
		logrus.Infof("distribution check passed, current version: %v", disName)
		return nil
	case DistributionRHEL:
		logrus.WithFields(logrus.Fields{
			"actual_value": disName,
		})
		logrus.Infof("distribution check passed, current version: %v", disName)
		return nil
	default:
		logrus.WithFields(logrus.Fields{
			"error_reason": "distribution error",
			"actual_value": disName,
			"desire_value": fmt.Sprintf("supported distribution: '%v' or '%v' or '%v'", DistributionCentos, DistributionUbuntu, DistributionRHEL),
		})
		logrus.Errorf("unsupported distribution")
		return fmt.Errorf("distribution can not be supported: (%v), desired distribution: (%v, %v, %v)", disName, DistributionCentos, DistributionUbuntu, DistributionRHEL)
	}
}
