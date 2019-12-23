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

	"github.com/sirupsen/logrus"
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

func NextID() uint64 {
	id, err := idCreator.NextID()
	if err != nil {
		// Based on the readme from Sonyflake: NextID can continue to generate IDs for about 174 years from StartTime.
		// After that time, an error will return. In our case, we can ingore this error.
		// So we eat the error here.
		logrus.Error(err)
	}

	return id
}

func NextString() string {
	return fmt.Sprintf("%x", NextID())
}
