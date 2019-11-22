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

package wizard

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogs(t *testing.T) {

	logId, err := logIdCreator.NextID()
	assert.Nil(t, err)

	logs[logId] = []byte("test")

	currentLogs := logs
	currentLogIdCreator := logIdCreator

	InitLogs()

	assert.NotEqual(t, currentLogs, logs)
	assert.NotEqual(t, currentLogIdCreator, logIdCreator)
}

func TestSetLogByString(t *testing.T) {

	logId, err := SetLogByString("testString")
	assert.Nil(t, err)
	assert.Greater(t, logId, uint64(0))

	assert.Equal(t, logs[logId], []byte("testString"))
}

func TestSetLogByReader(t *testing.T) {

	logId, err := SetLogByReader(bytes.NewReader([]byte("testBytes")))
	assert.Nil(t, err)
	assert.Greater(t, logId, uint64(0))

	assert.Equal(t, logs[logId], []byte("testBytes"))
}

func TestGetLog(t *testing.T) {

	logId := uint64(9644)
	logContent := []byte("testGetLog")
	logs[logId] = logContent
	assert.Equal(t, logContent, GetLog(logId))
}

func TestGetLogReader(t *testing.T) {

	logId := uint64(44944)
	logContent := []byte("testGetLogReader")
	logs[logId] = logContent
	reader := GetLogReader(logId)
	assert.NotNil(t, reader)
	readContent, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, logContent, readContent)
}
