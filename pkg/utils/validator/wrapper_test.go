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

package validator

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWrapper(t *testing.T) {

	noneErrorFunc := func() error {
		return nil
	}
	tests := []struct {
		Input []ValidateFunc
		Want  *ValidateWrapper
	}{
		{
			Input: []ValidateFunc{
				noneErrorFunc,
			},
			Want: &ValidateWrapper{
				items: []ValidateFunc{
					noneErrorFunc,
				},
			},
		},
	}

	for _, test := range tests {

		assert.Equal(t, fmt.Sprint(test.Want), fmt.Sprint(NewWrapper(test.Input...)))
	}
}

func TestValidateWrapper_AddValidateFunc(t *testing.T) {

	errFunc1 := func() error {
		return nil
	}
	errFunc2 := func() error {
		return nil
	}
	tests := []struct {
		Original ValidateWrapper
		Input    []ValidateFunc
		Want     ValidateWrapper
	}{
		{
			Original: ValidateWrapper{},
			Input: []ValidateFunc{
				errFunc1,
			},
			Want: ValidateWrapper{
				items: []ValidateFunc{
					errFunc1,
				},
			},
		},
		{
			Original: ValidateWrapper{
				items: []ValidateFunc{
					errFunc1,
				},
			},
			Input: []ValidateFunc{
				errFunc2,
			},
			Want: ValidateWrapper{
				items: []ValidateFunc{
					errFunc1,
					errFunc2,
				},
			},
		},
	}

	for _, test := range tests {

		wrapper := test.Original
		wrapper.AddValidateFunc(test.Input...)
		assert.Equal(t, fmt.Sprint(test.Want), fmt.Sprint(wrapper))
	}
}

func TestValidateWrapper_Validate(t *testing.T) {

	errFunc1 := func() error {
		return nil
	}
	errFunc2 := func() error {
		return errors.New("error message")
	}
	tests := []struct {
		Input ValidateWrapper
		Want  error
	}{
		{
			Input: ValidateWrapper{},
			Want:  nil,
		},
		{
			Input: ValidateWrapper{
				items: []ValidateFunc{
					errFunc1,
				},
			},
			Want: nil,
		},
		{
			Input: ValidateWrapper{
				items: []ValidateFunc{
					errFunc2,
				},
			},
			Want: errors.New("error message"),
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.Want, test.Input.Validate())
	}
}

func TestValidLength(t *testing.T) {

	tests := []struct {
		InputString  string
		InputKeyName string
		InputMinimum int
		InputMaximum int
		Want         error
	}{
		{
			InputString:  "",
			InputKeyName: "empty",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
		{
			InputString:  "",
			InputKeyName: "empty",
			InputMinimum: ItemNotEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         errors.New(`"empty" '' is too short`),
		},
		{
			InputString:  "to",
			InputKeyName: "keyName",
			InputMinimum: 3,
			InputMaximum: ItemNoLimit,
			Want:         errors.New(`"keyName" 'to' is too short`),
		},
		{
			InputString:  "one",
			InputKeyName: "keyName",
			InputMinimum: 3,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
		{
			InputString:  "one",
			InputKeyName: "keyName",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: 3,
			Want:         nil,
		},
		{
			InputString:  "four",
			InputKeyName: "keyName",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: 3,
			Want:         errors.New(`"keyName" 'four' is too long`),
		},
		{
			InputString:  "one",
			InputKeyName: "keyName",
			InputMinimum: ItemNotEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.Want, ValidLength(test.InputString, test.InputKeyName, test.InputMinimum, test.InputMaximum))
	}
}

func TestValidateString(t *testing.T) {

	tests := []struct {
		InputString  string
		InputKeyName string
		InputMinimum int
		InputMaximum int
		Want         error
	}{
		{
			InputString:  "",
			InputKeyName: "empty",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
		{
			InputString:  "",
			InputKeyName: "empty",
			InputMinimum: ItemNotEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         errors.New(`"empty" '' is too short`),
		},
		{
			InputString:  "to",
			InputKeyName: "keyName",
			InputMinimum: 3,
			InputMaximum: ItemNoLimit,
			Want:         errors.New(`"keyName" 'to' is too short`),
		},
		{
			InputString:  "one",
			InputKeyName: "keyName",
			InputMinimum: 3,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
		{
			InputString:  "one",
			InputKeyName: "keyName",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: 3,
			Want:         nil,
		},
		{
			InputString:  "four",
			InputKeyName: "keyName",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: 3,
			Want:         errors.New(`"keyName" 'four' is too long`),
		},
		{
			InputString:  "one",
			InputKeyName: "keyName",
			InputMinimum: ItemNotEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
	}

	for _, test := range tests {

		validateFunction := ValidateString(test.InputString, test.InputKeyName, test.InputMinimum, test.InputMaximum)
		assert.Equal(t, test.Want, validateFunction())
	}
}

func TestValidateStringPointer(t *testing.T) {

	toStringPointer := func(str string) *string {
		return &str
	}

	tests := []struct {
		InputString  *string
		InputKeyName string
		InputMinimum int
		InputMaximum int
		Want         error
	}{
		{
			InputString:  toStringPointer(""),
			InputKeyName: "empty",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
		{
			InputString:  nil,
			InputKeyName: "nil",
			InputMinimum: ItemEmptyLimit,
			InputMaximum: ItemNoLimit,
			Want:         nil,
		},
	}

	for _, test := range tests {

		validateFunction := ValidateStringPointer(test.InputString, test.InputKeyName, test.InputMinimum, test.InputMaximum)
		assert.Equal(t, test.Want, validateFunction())
	}
}

func TestValidateSameString(t *testing.T) {

	tests := []struct {
		InputString1  string
		InputString2  string
		InputKeyName1 string
		InputKeyName2 string
		Want          error
	}{
		{
			InputString1:  "a",
			InputString2:  "a",
			InputKeyName1: "keyName1",
			InputKeyName2: "keyName2",
			Want:          nil,
		},
		{
			InputString1:  "1",
			InputString2:  "2",
			InputKeyName1: "keyName1",
			InputKeyName2: "keyName2",
			Want:          errors.New(`"keyName1" & "keyName2" must equal`),
		},
	}

	for _, test := range tests {

		validateFunction := ValidateSameString(test.InputString1, test.InputKeyName1, test.InputString2, test.InputKeyName2)
		assert.Equal(t, test.Want, validateFunction())
	}
}

func TestValidateRegexp(t *testing.T) {

	tests := []struct {
		InputRegularExpression *regexp.Regexp
		InputString            string
		InputKeyName           string
		Want                   error
	}{
		{
			InputRegularExpression: regexp.MustCompile(`\w+`),
			InputString:            "word",
			InputKeyName:           "keyName",
			Want:                   nil,
		},
		{
			InputRegularExpression: regexp.MustCompile(`[a-zA-Z]+`),
			InputString:            "0123456",
			InputKeyName:           "keyName",
			Want:                   errors.New(`"keyName" illegal`),
		},
	}

	for _, test := range tests {

		validateFunction := ValidateRegexp(test.InputRegularExpression, test.InputString, test.InputKeyName)
		assert.Equal(t, test.Want, validateFunction())
	}
}

func TestValidateEmail(t *testing.T) {

	tests := []struct {
		Input string
		Want  bool
	}{
		{
			Input: "lucky@kpaas.io",
			Want:  true,
		},
		{
			Input: "kpaas.io",
			Want:  false,
		},
		{
			Input: "dev-support@kpaas.io",
			Want:  true,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.Want, ValidateEmail(test.Input))
	}
}

func TestValidateMobile(t *testing.T) {

	tests := []struct {
		Input string
		Want  bool
	}{
		{
			Input: "13800138000",
			Want:  true,
		},
		{
			Input: "90086",
			Want:  false,
		},
		{
			Input: "12345678901",
			Want:  false,
		},
		{
			Input: "18910010000",
			Want:  true,
		},
		{
			Input: "a word",
			Want:  false,
		},
		{
			Input: "a word include 13800138000",
			Want:  false,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.Want, ValidateMobile(test.Input))
	}
}

func TestValidateStringOptions(t *testing.T) {

	tests := []struct {
		InputOptions []string
		InputString  string
		InputKeyName string
		Want         error
	}{
		{
			InputOptions: []string{"option1", "option2"},
			InputString:  "option1",
			InputKeyName: "keyName",
			Want:         nil,
		},
		{
			InputOptions: []string{"option1", "option2"},
			InputString:  "0123456",
			InputKeyName: "keyName",
			Want:         errors.New(`keyName not in specify options`),
		},
	}

	for _, test := range tests {

		validateFunction := ValidateStringOptions(test.InputString, test.InputKeyName, test.InputOptions)
		assert.Equal(t, test.Want, validateFunction())
	}
}

func TestValidateStringArrayOptions(t *testing.T) {

	tests := []struct {
		InputOptions []string
		InputStrings []string
		InputKeyName string
		Want         error
	}{
		{
			InputOptions: []string{"option1", "option2"},
			InputStrings: []string{"option1"},
			InputKeyName: "keyName",
			Want:         nil,
		},
		{
			InputOptions: []string{"option1", "option2"},
			InputStrings: []string{"0123456"},
			InputKeyName: "keyName",
			Want:         errors.New(`keyName not in specify options`),
		},
		{
			InputOptions: []string{"option1", "option2"},
			InputStrings: []string{},
			InputKeyName: "keyName",
			Want:         errors.New(`keyName is empty`),
		},
	}

	for _, test := range tests {

		validateFunction := ValidateStringArrayOptions(test.InputStrings, test.InputKeyName, test.InputOptions)
		assert.Equal(t, test.Want, validateFunction())
	}
}

func TestValidateIntRange(t *testing.T) {

	tests := []struct {
		InputInt     int
		InputKeyName string
		InputMinimum int
		InputMaximum int
		Want         error
	}{
		{
			InputInt:     0,
			InputKeyName: "keyName",
			InputMinimum: 0,
			InputMaximum: 100,
			Want:         nil,
		},
		{
			InputInt:     8,
			InputKeyName: "keyName",
			InputMinimum: 0,
			InputMaximum: 100,
			Want:         nil,
		},
		{
			InputInt:     100,
			InputKeyName: "keyName",
			InputMinimum: 0,
			InputMaximum: 100,
			Want:         nil,
		},
		{
			InputInt:     101,
			InputKeyName: "keyName",
			InputMinimum: 0,
			InputMaximum: 100,
			Want:         errors.New("keyName out of range: minimum: 0, maximum: 100"),
		},
		{
			InputInt:     -1,
			InputKeyName: "keyName",
			InputMinimum: 0,
			InputMaximum: 100,
			Want:         errors.New("keyName out of range: minimum: 0, maximum: 100"),
		},
	}

	for _, test := range tests {

		validateFunction := ValidateIntRange(test.InputInt, test.InputKeyName, test.InputMinimum, test.InputMaximum)
		assert.Equal(t, test.Want, validateFunction())
	}
}

func TestValidateIP(t *testing.T) {

	tests := []struct {
		InputString  string
		InputKeyName string
		Want         error
	}{
		{
			InputString:  "192.168.1.1",
			InputKeyName: "keyName",
			Want:         nil,
		},
		{
			InputString:  "192.168.0.0",
			InputKeyName: "keyName",
			Want:         nil,
		},
		{
			InputString:  "string",
			InputKeyName: "keyName",
			Want:         errors.New("keyName is invalid ip"),
		},
	}

	for _, test := range tests {

		validateFunction := ValidateIP(test.InputString, test.InputKeyName)
		assert.Equal(t, test.Want, validateFunction())
	}
}
