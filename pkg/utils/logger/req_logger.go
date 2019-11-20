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

package logger

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	REQID_CTX_KEY = "X-Reqid"
	REQID_HEADER  = "X-Reqid"
)

func init() {
	/*
		some vendor will use glog as logger, which will create logfile under /tmp when error is logged.
		this will cause the programm exit, if the /tmp directory is not writable.
		so we disable Glog, and prevent glog to create logfile
	*/
	logtostderr := flag.Lookup("logtostderr")
	if logtostderr != nil && logtostderr.Value != nil {
		logtostderr.Value.Set("true")
	}
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano, FullTimestamp: true})
}

func genReqId() string {
	var b [12]byte
	io.ReadFull(rand.Reader, b[:])
	return base64.URLEncoding.EncodeToString(b[:])
}

func ReqLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339})

		reqid := c.Request.Header.Get(REQID_HEADER)
		if reqid == "" {
			reqid = genReqId()
			c.Request.Header.Set(REQID_HEADER, reqid)
		}
		c.Set(REQID_CTX_KEY, reqid)
		// Set request id into response header
		c.Writer.Header().Set(REQID_HEADER, reqid)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		entry := logrus.WithFields(logrus.Fields{
			"reqid":      reqid,
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL,
			"size":       c.Writer.Size(),
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}

// usage: ReqEntry(c).Debug(".....")
func ReqEntry(c context.Context) *logrus.Entry {
	reqid, _ := c.Value(REQID_CTX_KEY).(string)
	return logrus.WithField("reqid", reqid)
}
