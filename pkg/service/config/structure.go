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

package config

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	WebServiceModeDebug            WebServiceMode = "debug"
	WebServiceModeRelease          WebServiceMode = "release"
	DefaultListenPort                             = uint16(8080)
	DefaultLogLevel                               = "info"
	DefaultReadWriteTimeout                       = time.Minute
	DefaultDeployControllerAddress                = "127.0.0.1:8081"
	DefaultServiceId                              = 0
)

type (
	// Gin Service Mode
	WebServiceMode string

	configuration struct {
		Service          serviceSetting          `json:"service"`
		Log              logSetting              `json:"log"`
		DeployController deployControllerSetting `json:"deployController"`
	}

	serviceSetting struct {
		Port             uint16         `json:"port"`
		Mode             WebServiceMode `json:"mode"`
		ReadWriteTimeout time.Duration  `json:"readWriteTimeout"`
		ServiceId        uint16         // used to distinguish between different services when highly available. no parse from configuration file, because services will use the same configuration file.
	}

	logSetting struct {
		Level string `json:"level"` // level: trace, debug, info, warn|warning, error, fatal, panic
	}

	deployControllerSetting struct {
		Address string        `json:"address"`
		Timeout time.Duration `json:"timeout"`
	}
)

var (
	Config = new(configuration)
)

func (service *serviceSetting) GetPort() uint16 {
	if service.Port <= 0 {
		return DefaultListenPort
	}
	return service.Port
}

func (service *serviceSetting) GetMode() WebServiceMode {
	if service.Mode == "" {
		return gin.ReleaseMode
	}
	return service.Mode
}

func (service *serviceSetting) GetReadWriteTimeout() time.Duration {
	if service.ReadWriteTimeout == 0 {
		return DefaultReadWriteTimeout
	}
	return service.ReadWriteTimeout
}

func (service *serviceSetting) GetServiceId() uint16 {
	if service.ServiceId == 0 {
		return DefaultServiceId
	}
	return service.ServiceId
}

func (log *logSetting) GetLevel() string {
	if log.Level == "" {
		return DefaultLogLevel
	}
	return log.Level
}

func (controller *deployControllerSetting) GetAddress() string {
	if controller.Address == "" {
		return DefaultDeployControllerAddress
	}
	return controller.Address
}

func (controller *deployControllerSetting) GetTimeout() time.Duration {

	if controller.Timeout == 0 {
		return DefaultReadWriteTimeout
	}
	return controller.Timeout
}
