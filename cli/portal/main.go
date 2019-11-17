// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	goflag "flag"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := getCommand()
	decorateFlags(command)

	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func decorateFlags(command *cobra.Command) {

	command.Flags().Int16("port", 8080, "web service listening port")
}

func getCommand() *cobra.Command {

	return &cobra.Command{
		Use:   "portal",
		Short: "portal is a tool for deploying Kubernetes",
		Long:  `portal is a tool for deploying Kubernetes clusters using web interface`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
}
