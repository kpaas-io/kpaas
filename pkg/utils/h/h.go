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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type H = map[string]interface{}

type AppErr struct {
	Status  int         `json:"-"`
	Msg     string      `json:"msg"`
	Payload interface{} `json:"payload"`
}

var (
	ENotFound             = &AppErr{http.StatusNotFound, "NotFound", nil}
	EParamsError          = &AppErr{http.StatusBadRequest, "ParamsError", nil}
	EUnknown              = &AppErr{http.StatusInternalServerError, "Unknown", nil}
	EExists               = &AppErr{http.StatusConflict, "Exists", nil}
)

const ()

func (e *AppErr) WithPayload(payload interface{}) *AppErr {
	return &AppErr{e.Status, e.Msg, payload}
}

func NewAppErr(code int, msg string, payload interface{}) *AppErr {
	return &AppErr{code, msg, payload}
}

func (e *AppErr) Error() string {
	data, err := json.Marshal(e)
	if err != nil {
		logrus.Errorln("AppErr.Error():", err)
		return ""
	}
	return string(data)
}

func E(c *gin.Context, e error, msg ...string) {
	appErr := WrapErr(e, msg...)
	c.JSON(appErr.Status, appErr)
}

func R(c *gin.Context, body interface{}) {
	if c.Request.Method == "POST" {
		c.JSON(http.StatusCreated, body)
		return
	}
	if c.Request.Method == "DELETE" {
		c.JSON(http.StatusNoContent, body)
		return
	}
	c.JSON(http.StatusOK, body)
}

func RJsonP(c *gin.Context, body interface{}) {
	if c.Request.Method == "POST" {
		c.JSONP(http.StatusCreated, body)
		return
	}
	if c.Request.Method == "DELETE" {
		c.JSONP(http.StatusNoContent, body)
		return
	}
	c.JSONP(http.StatusOK, body)
}

func WrapErr(err error, msg ...string) *AppErr {

	payload := fmt.Sprintf("%s %s", err.Error(), strings.Join(msg, " "))

	switch appErr := err.(type) {
	case *AppErr:
		return appErr.WithPayload(payload)
	default:
		return EUnknown.WithPayload(payload)
	}

}
