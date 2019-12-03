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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {

	failureDetail := NewFailureDetail()
	assert.NotNil(t, failureDetail)
}

func TestFailureDetail_Clone(t *testing.T) {

	failureDetail1 := NewFailureDetail()
	failureDetail1.LogId = 1
	failureDetail1.Reason = "reason"
	failureDetail1.Detail = "detail"
	failureDetail1.FixMethods = "fixMethods"
	failureDetail2 := failureDetail1.Clone()
	assert.EqualValues(t, failureDetail1, failureDetail2)
	assert.False(t, failureDetail1 == failureDetail2)
}
