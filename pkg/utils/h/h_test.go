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

package h

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWithMore(t *testing.T) {
	err := ENotFound.WithPayload(map[string]interface{}{"text": "user not found", "id": 1})
	assert.NotNil(t, err)

	s1 := fmt.Sprintf("%s", err)
	s2 := err.Error()
	assert.Equal(t, s1, s2, "should be equal")
}

func TestWrapErr(t *testing.T) {
	tests := []struct {
		Input error
		msg   []string
		Want  int
	}{
		{EExists, []string{}, 409},
		{fmt.Errorf("some error"), []string{"some error"}, 500},
		{fmt.Errorf("some error"), nil, 500},
		{EParamsError, nil, 400},
		{EParamsError, []string{"something error"}, 400},
	}

	for _, item := range tests {
		res := WrapErr(item.Input, item.msg...)
		assert.Equal(t, item.Want, res.Status)
	}

}

func TestNewAppErr(t *testing.T) {

	tests := [] struct {
		Input struct {
			Code    int
			Msg     string
			Payload interface{}
		}
		Want *AppErr
	}{
		{
			Input: struct {
				Code    int
				Msg     string
				Payload interface{}
			}{
				Code:    http.StatusBadRequest,
				Msg:     "BadRequest",
				Payload: nil,
			},
			Want: &AppErr{
				Status:  http.StatusBadRequest,
				Msg:     "BadRequest",
				Payload: nil,
			},
		},
	}

	for _, item := range tests {
		res := NewAppErr(item.Input.Code, item.Input.Msg, item.Input.Payload)
		assert.Equal(t, item.Want, res)
	}
}

func TestAppErr_Error(t *testing.T) {

	tests := [] struct {
		Input *AppErr
		Want  string
	}{
		{
			Input: NewAppErr(http.StatusBadRequest, "BadRequest", nil),
			Want:  `{"msg":"BadRequest","payload":null}`,
		},
		{
			Input: NewAppErr(http.StatusInternalServerError, "InternalServerError", make(chan int)),
			Want:  ``,
		},
	}

	for _, item := range tests {
		assert.Equal(t, item.Want, item.Input.Error())
	}
}

func TestE(t *testing.T) {

	tests := [] struct {
		Input error
		Want  struct {
			Code           int
			ResponseString string
		}
	}{
		{
			Input: ENotFound,
			Want: struct {
				Code           int
				ResponseString string
			}{Code: ENotFound.Status, ResponseString: `{"msg":"NotFound","payload":null}`},
		},
		{
			Input: EParamsError,
			Want: struct {
				Code           int
				ResponseString string
			}{Code: EParamsError.Status, ResponseString: `{"msg":"ParamsError","payload":null}`},
		},
		{
			Input: errors.New("BadRequest"),
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusInternalServerError, ResponseString: `{"msg":"BadRequest","payload":null}`},
		},
	}

	for _, item := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		E(ctx, item.Input)
		resp.Flush()

		assert.Equal(t, item.Want.Code, resp.Code)
		assert.Equal(t, item.Want.ResponseString, resp.Body.String())
	}
}

func TestR(t *testing.T) {

	tests := [] struct {
		Input struct {
			Body       interface{}
			HTTPMethod string
		}
		Want struct {
			Code           int
			ResponseString string
		}
	}{
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodGet,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusOK, ResponseString: `{"Hello":"World"}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodPost,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusCreated, ResponseString: `{"Hello":"World"}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodDelete,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusNoContent, ResponseString: ``},
		},
	}

	for _, item := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ctx.Request = &http.Request{Method: item.Input.HTTPMethod}
		R(ctx, item.Input.Body)
		resp.Flush()

		assert.Equal(t, item.Want.Code, resp.Code)
		assert.Equal(t, item.Want.ResponseString, resp.Body.String())
	}
}

func TestRJsonP(t *testing.T) {

	tests := [] struct {
		Input struct {
			Body       interface{}
			URL        string
			HTTPMethod string
		}
		Want struct {
			Code           int
			ResponseString string
		}
	}{
		{
			Input: struct {
				Body       interface{}
				URL        string
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				URL:        "http://localhost/?callback=c",
				HTTPMethod: http.MethodGet,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusOK, ResponseString: `c({"Hello":"World"})`},
		},
		{
			Input: struct {
				Body       interface{}
				URL        string
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				URL:        "http://localhost/?callback=c",
				HTTPMethod: http.MethodPost,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusCreated, ResponseString: `c({"Hello":"World"})`},
		},
		{
			Input: struct {
				Body       interface{}
				URL        string
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				URL:        "http://localhost/?callback=c",
				HTTPMethod: http.MethodDelete,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusNoContent, ResponseString: ``},
		},
	}

	for _, item := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		requestURL, err := url.Parse(item.Input.URL)
		assert.Nil(t, err)
		ctx.Request = &http.Request{Method: item.Input.HTTPMethod, URL: requestURL}
		RJsonP(ctx, item.Input.Body)
		resp.Flush()

		assert.Equal(t, item.Want.Code, resp.Code)
		assert.Equal(t, item.Want.ResponseString, resp.Body.String())
	}
}
