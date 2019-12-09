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

package idcreator

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

const (
	DefaultServiceID = 0
)

var idCreator *sonyflake.Sonyflake

func init() {

	InitCreator(DefaultServiceID)
}

func InitCreator(serviceId uint16) {

	idCreator = sonyflake.NewSonyflake(
		sonyflake.Settings{
			StartTime: time.Now(),
			MachineID: func() (u uint16, e error) {
				return serviceId, nil
			},
		})
}

func NextID() (uint64, error) {

	return idCreator.NextID()
}

func NextString() (string, error) {

	uid, err := idCreator.NextID()
	return fmt.Sprintf("%x", uid), err
}
