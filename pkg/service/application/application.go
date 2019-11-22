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
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/kpaas-io/kpaas/pkg/service/config"
	configUtils "github.com/kpaas-io/kpaas/pkg/utils/config"
	"github.com/kpaas-io/kpaas/pkg/utils/logger"
)

type (
	app struct {
		httpHandler *gin.Engine
		httpServer  *http.Server
		isClosing   bool
	}
)

func NewApp() *app {
	return new(app)
}

func (a *app) run(cmd *cobra.Command, args []string) {

	a.parseFlags()
	a.initService()
	a.startService()
	a.waitSignal()
}

func (a *app) startService() {
	a.startRESTfulAPIListener()
}

func (a *app) initService() {

	a.initRandomSeed()
	a.initLogLevel()
	// TODO Lucky Init Memories Database
	// TODO Lucky Init Clients
	a.initRESTfulAPIHandler()
	a.initRequestLogger()
	a.setRoutes()
	a.initRESTfulListener()
}

func (a *app) parseFlags() {
	pflag.Parse()
	a.loadConfig()
	a.parseParameters()
}

func (a *app) initRandomSeed() {

	rand.Seed(time.Now().UnixNano())
}

func (a *app) initLogLevel() {

	logLevel, err := logrus.ParseLevel(config.Config.Log.GetLevel())
	if err != nil {
		logrus.Errorf("Parse log level error")
	} else {
		logrus.SetLevel(logLevel)
	}
}

func (a *app) waitSignal() {
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		select {
		case sig := <-chanSignal:
			logrus.Infof("Received signal: %d", sig)
			a.close()
			goto exit
		}
	}
exit:
	logrus.Infof("loop exited")
	return
}

func (a *app) close() {

	a.isClosing = true

	// TODO Lucky Clean Memory Database

	logrus.Infof("closing http server")
	err := a.httpServer.Close()
	if err != nil {
		logrus.Errorf("happened error at close http server: %v", err)
	}
	logrus.Infof("http server closed")
}

func (a *app) loadConfig() {
	configFile, _ := pflag.CommandLine.GetString("config-file")
	if configFile == "" {
		return
	}
	configUtils.MustLoadConf(config.Config, configFile)
	logrus.Infof("ConfigFile: %s\n%+v", configFile, config.Config)
}

func (a *app) parseParameters() {

	a.parseParameterListenPort()
	a.parseParameterLogLevel()
}

func (a *app) parseParameterListenPort() {

	var err error
	var listenPort uint16
	listenPort, err = pflag.CommandLine.GetUint16(FlagPort)
	if err == nil && listenPort > 0 {
		config.Config.Service.Port = listenPort
	}
	logrus.Infof("listen port: %v", config.Config.Service.GetPort())
}

func (a *app) parseParameterLogLevel() {

	var err error
	var logLevel string
	logLevel, err = pflag.CommandLine.GetString(FlagLogLevel)
	if err == nil && logLevel != "" {
		config.Config.Log.Level = logLevel
	}
	logrus.Infof("log level: %v", config.Config.Log.GetLevel())
}

func (a *app) initRESTfulAPIHandler() {

	logrus.Debug("start to init restful api service handler")
	a.httpHandler = gin.Default()
	switch config.Config.Service.GetMode() {
	case config.WebServiceModeDebug:
		gin.SetMode(gin.DebugMode)
	default: // Default Release Mode
		gin.SetMode(gin.ReleaseMode)
	}
	logrus.Debug("init restful api handler service succeed")
}

func (a *app) initRequestLogger() {
	logrus.Debug("start to register log middleware")
	a.httpHandler.Use(logger.ReqLoggerMiddleware())
	logrus.Debug("register log middleware succeed")
}

func (a *app) initRESTfulListener() {
	logrus.Debug("start to init restful api listener")
	// start the server，For services exposed on the public network, timeout must be set
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", config.Config.Service.GetPort()),
		Handler:      a.httpHandler,
		ReadTimeout:  config.Config.Service.GetReadWriteTimeout(),
		WriteTimeout: config.Config.Service.GetReadWriteTimeout(),
	}
	logrus.Debug("init restful api listener succeed")
}

func (a *app) startRESTfulAPIListener() {
	logrus.Infof("start server listening")
	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil && !a.isClosing {
			logrus.Errorf("listen error: %v", err)
		}
	}()
}
