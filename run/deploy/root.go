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
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kpaas-io/kpaas/pkg/deploy/server"
	_ "github.com/kpaas-io/kpaas/pkg/utils/log"
)

var (
	cfgFile    string
	port       uint16
	logLevel   string
	logFileLoc string
)

const (
	defaultPort       uint16 = 8081
	defaultLogLevel   string = "info"
	defaultLogFileLoc string = "/app/log/deploy"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kpaas",
	Short: "the kpaas deploy controller",
	Long:  `The kpass deploy controller provides gRPC API services to deploy a k8s cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLogLevel()
		options := server.ServerOptions{
			Port:       port,
			LogFileLoc: logFileLoc,
		}
		server.New(options).Run(SetupSignalHandler())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&cfgFile, "config-file", "", "config file")
	rootCmd.Flags().Uint16VarP(&port, "port", "p", defaultPort, "gRPC service listening port")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", defaultLogLevel, "log level(options: trace, debug, info, warn|warning, error, fatal, panic)")
	rootCmd.Flags().StringVar(&logFileLoc, "log-file-location", defaultLogFileLoc, "the location to store the detail logs")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setupLogLevel() {
	logLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Parse log level error")
	} else {
		logrus.SetLevel(logLevel)
	}
}

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
