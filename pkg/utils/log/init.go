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

package log

import (
	"flag"
	"fmt"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	// 	some vendor will use glog as logger, which will create logfile under /tmp when error is logged.
	// 	this will cause the programm exit, if the /tmp directory is not writable.
	// 	so we disable Glog, and prevent glog to create logfile
	logtostderr := flag.Lookup("logtostderr")
	if logtostderr != nil && logtostderr.Value != nil {
		logtostderr.Value.Set("true")
	}

	// Set log formatter to ouput source code file name, line number and function name.
	var re = regexp.MustCompile(`^github.com/kpaas-io/kpaas/`)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fileName := path.Base(f.File)
			return fmt.Sprintf("%s()", re.ReplaceAllString(f.Function, "")), fmt.Sprintf("%s:%d", fileName, f.Line)
		},
	})
}
