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

package application

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FlagPort       = "port"
	FlagLogLevel   = "log-level"
	FlagConfigFile = "config-file"
	FlagServiceId  = "service-id"
)

func GetCommand() *cobra.Command {

	setFlags()
	return &cobra.Command{
		Use:   "restful",
		Short: "restful is a service that provides management RESTful APIs",
		Long:  `restful is a service that provides a management API that interfaces to the web front end and also calls the deployment service and the Kubernetes API service.`,
		Run:   NewApp().run,
	}
}

func setFlags() {

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Uint16(FlagPort, 8080, "restful service listening port.")
	pflag.String(FlagLogLevel, "info", "log level(options: trace, debug, info, warn|warning, error, fatal, panic)")
	pflag.String(FlagConfigFile, "", "config file for json format")
	pflag.Uint16(FlagServiceId, 0, "distinguish between different services when highly available.")
}
