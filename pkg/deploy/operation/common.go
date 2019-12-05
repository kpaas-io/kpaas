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
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

const (
	SplitSymbol              = "."
	CheckEqual               = "="
	CheckLarge               = ">"
	CheckLess                = "<"
	InfoPassed               = "check passed"
	ErrSplitSym              = "error split symbol found"
	ErrParaInput             = "input parameter invalid"
	ErrTooHigh               = "version too high"
	ErrTooLow                = "version too low"
	ErrNotEqual              = "version not equal"
	ErrNotEnough             = "amount not enough"
	UnclearInputPara         = "input parameter not clear"
	GiByteUnits      float64 = 1000 * 1000
)

// check if version is satisfied with standard version
// checkStandard controls compared method
func CheckVersion(comparedVersion string, standardVersion string, comparedSymbol string) error {
	logger := logrus.WithFields(logrus.Fields{
		"actual_version":  comparedVersion,
		"desired_version": standardVersion,
	})

	if err := checkVersionValid(comparedVersion, standardVersion); err != nil {
		return err
	}

	comparedVerStr := strings.Split(strings.TrimSpace(comparedVersion), "-")[0]
	standardVerStr := strings.Split(strings.TrimSpace(standardVersion), "-")[0]

	switch comparedSymbol {
	case CheckEqual:

		if comparedVersion == standardVersion {
			return nil
		}

		logger.Errorf("%v", ErrNotEqual)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrNotEqual, standardVersion, comparedVersion)

	case CheckLarge:

		result := versionLarger(comparedVerStr, standardVerStr)
		if result >= 0 {
			logger.Infof("check version passed")
			return nil
		}

		logger.Errorf("%v", ErrTooLow)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrTooLow, standardVersion, comparedVersion)

	case CheckLess:

		result := versionLarger(comparedVerStr, standardVerStr)
		if result <= 0 {
			logger.Infof("check version passed")
			return nil
		}

		logger.Errorf("%v", ErrTooHigh)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrTooHigh, standardVersion, comparedVersion)

	default:
		logger.Error("%v", UnclearInputPara)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", UnclearInputPara, standardVersion, comparedVersion)
	}
}

// check if first version larger than second version
func versionLarger(firstVer string, secondVer string) int {
	firstArray := strings.Split(firstVer, ".")
	secondArray := strings.Split(secondVer, ".")

	for i := 0; i < findMaxLength(firstArray, secondArray); i++ {
		var firstInt int
		var secondInt int
		var compare = 0

		if i < len(firstArray) {
			firstInt, _ = strconv.Atoi(firstArray[i])
		}
		if i < len(secondArray) {
			secondInt, _ = strconv.Atoi(secondArray[i])
		}
		if firstInt > secondInt {
			compare = 1
		} else if firstInt < secondInt {
			compare = -1
		}
		if compare != 0 {
			return compare
		}
	}
	return 0
}

// check if entity resource satisfied minimal requirements
func CheckEntity(comparedEntity string, desiredEntity float64) error {
	logger := logrus.WithFields(logrus.Fields{
		"actual_amount":  comparedEntity,
		"desired_amount": desiredEntity,
	})

	comparedEntityFloat64, err := strconv.ParseFloat(comparedEntity, 64)
	if err != nil {
		logger.Errorf("%v", ErrParaInput)
		return fmt.Errorf("%v, desired amount: %v, actual amount: %v", ErrParaInput, desiredEntity, comparedEntity)
	}

	if comparedEntityFloat64 < float64(0) {
		logger.Errorf("%v", ErrParaInput)
		return fmt.Errorf("%v, input parameter can not be negative, desired amount: %.1f", ErrParaInput, desiredEntity)
	}

	if comparedEntityFloat64 >= desiredEntity {
		logger.Infof("compared satisfied")
		return nil
	}

	logger.Errorf("%v", ErrNotEnough)
	return fmt.Errorf("%v, desired amount: %v, actual amount: %v", ErrNotEnough, desiredEntity, comparedEntity)
}

// check if raw input contains non-digit character
func checkContainsNonDigit(rawInput string) bool {
	bareRawInput := strings.ReplaceAll(rawInput, ".", "")
	for _, eachChar := range bareRawInput {
		if !unicode.IsDigit(eachChar) {
			return false
		}
	}
	return true
}

// check if input is invalid
func checkVersionValid(comparedVersion string, standardVersion string) error {
	logger := logrus.WithFields(logrus.Fields{
		"actual_version":  comparedVersion,
		"desired_version": standardVersion,
	})

	// check if input is empty
	if comparedVersion == "" || standardVersion == "" {
		logger.Errorf("%v", ErrParaInput)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrParaInput, comparedVersion, standardVersion)
	}

	// check if input not contains split symbol
	if !strings.Contains(comparedVersion, SplitSymbol) || !strings.Contains(standardVersion, SplitSymbol) {
		logger.Error("%v,", ErrSplitSym)
		return fmt.Errorf("%v: split symbol: %v", ErrSplitSym, SplitSymbol)
	}

	comparedVersion = strings.Split(strings.TrimSpace(comparedVersion), "-")[0]
	standardVersion = strings.Split(strings.TrimSpace(standardVersion), "-")[0]

	// check if input contains non-digit char
	if ok := checkContainsNonDigit(comparedVersion); !ok {
		logger.Error("%v", ErrParaInput)
		return fmt.Errorf("%v, contains non-digit char, desired version: %v, actual version: %v", ErrParaInput, comparedVersion, standardVersion)
	}
	if ok := checkContainsNonDigit(standardVersion); !ok {
		logger.Error("%v", ErrParaInput)
		return fmt.Errorf("%v, contains non-digit char, desired version: %v, actual version: %v", ErrParaInput, comparedVersion, standardVersion)
	}

	return nil
}

// find max length of two arrays
func findMaxLength(firstArr []string, secondArr []string) int {
	if len(firstArr) >= len(secondArr) {
		return len(firstArr)
	}
	return len(secondArr)
}
