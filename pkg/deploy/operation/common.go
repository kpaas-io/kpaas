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
	"net"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	SplitSymbol              = "."
	CheckEqual               = "="
	CheckLarge               = ">"
	CheckLess                = "<"
	ErrParaEmpty             = "parameter empty"
	ErrPara                  = "parameter error"
	ErrInvalid               = "parameter invalid"
	ErrSplitSym              = "error split symbol found"
	ErrParaInput             = "input parameter invalid"
	ErrTooHigh               = "version too high"
	ErrTooLow                = "version too low"
	ErrNotEqual              = "version not equal"
	ErrNotEnough             = "amount not enough"
	UnclearInputPara         = "input parameter not clear"
	GiByteUnits      float64 = 1000 * 1000 * 1000
)

type NodeInitAction struct {
	NodeInitConfig *pb.NodeDeployConfig
	NodesConfig    []*pb.NodeDeployConfig
	ClusterConfig  *pb.ClusterConfig
}

// check if version is satisfied with standard version
// checkStandard controls compared method
func CheckVersion(comparedVersion string, standardVersion string, comparedSymbol string) error {
	logger := logrus.WithFields(logrus.Fields{
		"actual_version":  comparedVersion,
		"desired_version": standardVersion,
	})

	if err := checkVersionValid(comparedVersion); err != nil {
		return err
	}
	if err := checkVersionValid(standardVersion); err != nil {
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

		result := versionLargerAndEqual(comparedVerStr, standardVerStr)
		if result >= 0 {
			logger.Infof("check version passed")
			return nil
		}

		logger.Errorf("%v", ErrTooLow)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrTooLow, standardVersion, comparedVersion)

	case CheckLess:

		result := versionLargerAndEqual(comparedVerStr, standardVerStr)
		if result <= 0 {
			logger.Infof("check version passed")
			return nil
		}

		logger.Errorf("%v", ErrTooHigh)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", ErrTooHigh, standardVersion, comparedVersion)

	default:
		logger.Errorf("%v", UnclearInputPara)
		return fmt.Errorf("%v, desired version: %v, actual version: %v", UnclearInputPara, standardVersion, comparedVersion)
	}
}

// check if first version larger than second version
func versionLargerAndEqual(firstVersion string, secondVersion string) int {
	firstArray := strings.Split(firstVersion, ".")
	secondArray := strings.Split(secondVersion, ".")

	for i := 0; i < findMaxLength(firstArray, secondArray); i++ {
		var firstInt int
		var secondInt int

		if i < len(firstArray) {
			firstInt, _ = strconv.Atoi(firstArray[i])
		}
		if i < len(secondArray) {
			secondInt, _ = strconv.Atoi(secondArray[i])
		}
		if firstInt > secondInt {
			return 1
		} else if firstInt < secondInt {
			return -1
		}
	}
	return 0
}

// check if entity resource satisfied minimal requirements
func CheckEntity(comparedEntity string, desiredEntity float64) error {
	logger := logrus.WithFields(logrus.Fields{
		"actual_amount":  comparedEntity,
		"desired_amount": fmt.Sprintf("%.0f", desiredEntity),
	})

	comparedEntityFloat64, err := strconv.ParseFloat(comparedEntity, 64)
	if comparedEntity == "" {
		logger.Errorf("%v: input entity is empty", ErrParaEmpty)
		return fmt.Errorf("%v, desired amount: %.0f, actual amount: %v", ErrParaInput, desiredEntity, comparedEntity)
	}

	if err != nil {
		logger.Errorf("%v: parse float failed", ErrParaInput)
		return fmt.Errorf("%v, desired amount: %.0f, actual amount: %v", ErrParaInput, desiredEntity, comparedEntity)
	}

	if comparedEntityFloat64 < float64(0) {
		logger.Errorf("%v: value can not be negative", ErrParaInput)
		return fmt.Errorf("%v, input parameter can not be negative, desired amount: %.0f", ErrParaInput, desiredEntity)
	}

	if comparedEntityFloat64 >= desiredEntity {
		logger.Infof("compared satisfied")
		return nil
	}

	logger.Errorf("%v", ErrNotEnough)
	return fmt.Errorf("%v, desired amount: %.0f, actual amount: %v", ErrNotEnough, desiredEntity, comparedEntity)
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
func checkVersionValid(rawVersion string) error {
	logger := logrus.WithFields(logrus.Fields{
		"input_version": rawVersion,
	})

	// check if version is empty
	if rawVersion == "" {
		logger.Errorf("%v: input version empty", ErrParaInput)
		return fmt.Errorf("%v, input version: %v", ErrParaInput, rawVersion)
	}

	// check if not contains split symbol
	if !strings.Contains(rawVersion, SplitSymbol) {
		logger.Errorf("%v,", ErrSplitSym)
		return fmt.Errorf("%v: split symbol: %v", ErrSplitSym, SplitSymbol)
	}

	splitedVersion := strings.Split(strings.TrimSpace(rawVersion), "-")[0]

	// check if input contains non-digit char
	if ok := checkContainsNonDigit(splitedVersion); !ok {
		logger.Errorf("%v: inputs can not contain non-digit character", ErrParaInput)
		return fmt.Errorf("%v, contains non-digit char, input version: %v", ErrParaInput, rawVersion)
	}
	if ok := checkContainsNonDigit(splitedVersion); !ok {
		logger.Errorf("%v", ErrParaInput)
		return fmt.Errorf("%v, contains non-digit char, input version: %v", ErrParaInput, rawVersion)
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

// check if ip valid as 0.0.0.0 or defined in RFC1122, RFC4632, RFC4291
func CheckIPValid(rawIP string) bool {
	if rawIP == "0.0.0.0" {
		return true
	}

	parsedRawIP := net.ParseIP(rawIP)
	if ok := parsedRawIP.IsGlobalUnicast(); ok {
		return true
	}
	return false
}
