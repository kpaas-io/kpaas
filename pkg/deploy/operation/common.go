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

package operation

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	CheckEqual           = "="
	CheckLarge           = ">"
	CheckLess            = "<"
	ErrSplitSym          = "error split symbol found"
	ErrParaInput         = "input parameter invalid"
	ErrVersionHigh       = "version too high"
	ErrVersionLow        = "version too low"
	UnclearInputPara     = "input parameter not clear"
	GiByteUnits  float64 = 1000 * 1000
)

// check if version is satisfied with standard version
// checkStandard controls compared method
func VersionSatisfiedStandard(comparedVersion string, standardVersion string, splitSymbol string, comparedSymbol string) error {
	var loop int
	if comparedVersion == "" || standardVersion == "" {
		logrus.WithFields(logrus.Fields{
			"error_reason": ErrParaInput,
			"actual_version": comparedVersion,
			"desired_version": standardVersion,
		})
		logrus.Error(ErrParaInput)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrParaInput, comparedVersion, standardVersion)
	}

	if !strings.Contains(comparedVersion, splitSymbol) || !strings.Contains(standardVersion, splitSymbol) {
		logrus.WithFields(logrus.Fields{
			"error_reason": ErrSplitSym,
			"actual_version": comparedVersion,
			"desired_version": standardVersion,
		})
		logrus.Error(ErrSplitSym)
		return fmt.Errorf("%v: split symbol: %v", ErrSplitSym, splitSymbol)
	}

	comparedVersion = strings.TrimSpace(comparedVersion)
	standardVersion = strings.TrimSpace(standardVersion)

	comparedVerArray := strings.Split(strings.Split(comparedVersion, "-")[0], splitSymbol)
	standardVerArray := strings.Split(strings.Split(standardVersion, "-")[0], splitSymbol)

	loop = len(standardVerArray)

	switch comparedSymbol {
	case CheckEqual:
		if comparedVersion == standardVersion {
			return nil
		}

	case CheckLarge:

		for i := 0; i < loop; i++ {
			comparedInt, _ := strconv.Atoi(comparedVerArray[i])
			standardInt, _ := strconv.Atoi(standardVerArray[i])
			if comparedInt >= standardInt {
				return nil
			}

			logrus.WithFields(logrus.Fields{
				"error_reason": ErrVersionLow,
				"actual_version": comparedVersion,
				"desired_version": standardVersion,
			})
			logrus.Error("version too low")
			return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrVersionLow, standardVersion, comparedVersion)
		}

	case CheckLess:

		for i := 0; i < loop; i++ {
			comparedInt, _ := strconv.Atoi(comparedVerArray[i])
			standardInt, _ := strconv.Atoi(standardVerArray[i])
			if comparedInt <= standardInt {
				return nil
			}

			logrus.WithFields(logrus.Fields{
				"error_reason": ErrVersionHigh,
				"actual_version": comparedVersion,
				"desired_version": standardVersion,
			})
			logrus.Error("version too high")
			return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrVersionHigh, standardVersion, comparedVersion)
		}

	default:
		logrus.WithFields(logrus.Fields{
			"error_reason": UnclearInputPara,
			"actual_version": comparedVersion,
			"desired_version": standardVersion,
		})
		logrus.Error("version not clear")
		return fmt.Errorf("%v, desired version: %v, actual version: %v", UnclearInputPara, standardVersion, comparedVersion)
	}

	logrus.WithFields(logrus.Fields{
		"error_reason": UnclearInputPara,
		"actual_version": comparedVersion,
		"desired_version": standardVersion,
	})
	logrus.Error("version not clear")
	return fmt.Errorf("%v, desired version: %v, actual version: %v", UnclearInputPara, standardVersion, comparedVersion)
}
