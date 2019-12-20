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

package utils

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

// ExecuteLogItem a log item about executing a command.
type ExecuteLogItem struct {
	StartTime   time.Time
	EndTime     time.Time
	Command     string // string or []string ?
	Stdout      []byte // ssh package uses []byte, so []byte is used here
	Stderr      []byte
	Err         error // other error messages
	Description string
}

// WriteExecuteLog write an item into writer.
func WriteExecuteLog(w io.Writer, item *ExecuteLogItem) {
	if w == nil || item == nil {
		return
	}
	buf := &bytes.Buffer{}
	buf.Write([]byte(consts.DashLine + "\n"))
	// write description
	if item.Description != "" {
		buf.Write([]byte("[description] "))
		buf.Write([]byte(item.Description + "\n"))
	}
	// write start time
	startTimeMsg := fmt.Sprintf("[start time] %s\n", item.StartTime.Format(
		"2006-01-02 15:04:05"))
	buf.Write([]byte(startTimeMsg))
	// write command
	buf.Write([]byte(fmt.Sprintf("[command] %s\n", item.Command)))
	// write stderr
	buf.Write([]byte("[stderr]\n"))
	buf.Write(item.Stderr)
	buf.Write([]byte("\n"))
	// write stdout
	buf.Write([]byte("[stdout]\n"))
	buf.Write(item.Stdout)
	buf.Write([]byte("\n"))
	// write error message
	if item.Err != nil {
		buf.Write([]byte("[error]\n"))
		buf.Write([]byte(item.Err.Error()))
		buf.Write([]byte("\n"))
	}
	// write end time
	endTimeMsg := fmt.Sprintf("[end time] %s\n", item.EndTime.Format(
		"2006-01-02 15:04:05"))
	buf.Write([]byte(endTimeMsg))
	// write an extra empty line
	buf.Write([]byte("\n"))
	w.Write(buf.Bytes())
}
