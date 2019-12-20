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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

func TestWriteExecuteLogs(t *testing.T) {
	var tests = []struct {
		input *ExecuteLogItem
		want  string
	}{
		{
			input: &ExecuteLogItem{
				StartTime: time.Date(2019, 12, 19, 19, 0, 0, 0, time.Local),
				EndTime:   time.Date(2019, 12, 19, 19, 0, 1, 1024, time.Local),
				Command:   "ls",
				Stdout:    []byte("a b c\n"),
				Stderr:    []byte{},
			},
			want: consts.DashLine + `
[start time] 2019-12-19 19:00:00
[command] ls
[stderr]

[stdout]
a b c

[end time] 2019-12-19 19:00:01

`,
		},
		{
			input: &ExecuteLogItem{
				StartTime:   time.Date(2019, 12, 19, 19, 0, 0, 0, time.Local),
				EndTime:     time.Date(2019, 12, 19, 19, 0, 0, 1024, time.Local),
				Description: "create directory",
				Command:     "mkdir /etc/l78z/config",
				Stdout:      []byte{},
				Stderr:      []byte("mkdir: cannot create directory ‘/etc/l78z/config’: No such file or directory\n"),
				Err:         fmt.Errorf("exit code 1"),
			},
			want: consts.DashLine + `
[description] create directory
[start time] 2019-12-19 19:00:00
[command] mkdir /etc/l78z/config
[stderr]
mkdir: cannot create directory ‘/etc/l78z/config’: No such file or directory

[stdout]

[error]
exit code 1
[end time] 2019-12-19 19:00:00

`,
		},
	}
	for _, testCase := range tests {
		buf := &bytes.Buffer{}
		WriteExecuteLog(buf, testCase.input)
		assert.Equal(t, testCase.want, buf.String())
	}
}
