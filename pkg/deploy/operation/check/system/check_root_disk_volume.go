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

package system

import (
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

// compare root disk volume with desire root disk volume
func CheckRootDiskVolume(rootDiskVolume string, desiredDiskVolume float64) error {
	err := operation.CheckEntity(rootDiskVolume, desiredDiskVolume)
	if err != nil {
		return err
	}
	return nil

	//var diskCapacity float64
	//
	//diskVolumeFloat, err := strconv.ParseFloat(rootDiskVolume, 64)
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{
	//		"error_reason": operation.ErrParaInput,
	//		"actual_amount": rootDiskVolume,
	//		"desired_amount": desiredDiskVolume,
	//	})
	//	logrus.Error("parameter error")
	//	return fmt.Errorf("%v, desired disk volume: %v, input disk volume: %v", operation.ErrParaInput, desiredDiskVolume, rootDiskVolume)
	//}
	//
	//diskCapacity = diskVolumeFloat / 1024 / 1024
	//
	//if diskCapacity < float64(0) {
	//	logrus.WithFields(logrus.Fields{
	//		"error_reason": operation.ErrParaInput,
	//		"actual_amount": rootDiskVolume,
	//		"desired_amount": desiredDiskVolume,
	//	})
	//	logrus.Error("root disk volume can not be negative")
	//	return fmt.Errorf("%v, root disk volume can not be negative, input actual volume: %v", operation.ErrParaInput, rootDiskVolume)
	//}
	//
	//if diskVolumeFloat >= desiredDiskVolume {
	//	return nil
	//}
	//
	//logrus.WithFields(logrus.Fields{
	//	"error_reason": "root disk not enough",
	//	"actual_amount": rootDiskVolume,
	//	"desired_amount": desiredDiskVolume,
	//})
	//logrus.Errorf("node root disk not satisfied")
	//return fmt.Errorf("root disk not enough, desired disk volume: (%.1f), actual disk volume: (%v)", desiredDiskVolume, rootDiskVolume)
}
