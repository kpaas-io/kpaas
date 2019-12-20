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

package validator_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

type TestParamsStructure struct {
	Key string `json:"key"`
}

func (s *TestParamsStructure) Validate() error {

	return validator.NewWrapper(
		validator.ValidateString(s.Key, "key", validator.ItemNotEmptyLimit, validator.ItemNoLimit),
	).Validate()
}

func TestParams(t *testing.T) {

	tests := []struct {
		Input string
		Want  TestParamsStructure
		Error error
	}{
		{
			Input: `{"key":"a"}`,
			Want:  TestParamsStructure{Key: "a"},
			Error: nil,
		},
		{
			Input: `"sss"`,
			Want:  TestParamsStructure{},
			Error: h.EBindBodyError.WithPayload("json: cannot unmarshal string into Go value of type validator_test.TestParamsStructure"),
		},
		{
			Input: `{"key":""}`,
			Want:  TestParamsStructure{},
			Error: h.EParamsError.WithPayload("\"key\" '' is too short"),
		},
	}

	for _, test := range tests {

		var err error
		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		bodyReader := strings.NewReader(test.Input)
		ctx.Request = httptest.NewRequest("POST", "/test", bodyReader)

		data := &TestParamsStructure{}
		err = validator.Params(ctx, data)

		assert.Equal(t, test.Error, err)
		assert.Equal(t, test.Want, *data)
	}
}
