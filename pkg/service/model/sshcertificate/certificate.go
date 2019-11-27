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

package sshcertificate

import (
	"sync"
)

type (
	Certificate struct {
		Name       string
		PrivateKey string
	}
)

var (
	list *sync.Map
)

func NewCertificate() *Certificate {
	return &Certificate{}
}

func init() {
	ClearList()
}

func ClearList() {
	list = new(sync.Map)
}

func AddCertificate(name, privateKey string) {

	list.Store(name, privateKey)
}

func GetNameList() []string {

	names := make([]string, 0, 0)
	list.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

func GetPrivateKey(name string) string {

	privateKey, exist := list.Load(name)
	if exist {
		return privateKey.(string)
	}
	return ""
}
