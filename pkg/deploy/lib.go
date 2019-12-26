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

package deploy

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	defaultApiServerPort = 6443
	defaultHAProxyPort   = 4443
)

var (
	AllFilesNeeded = func(path string) bool {
		return true
	}
)

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetControlPlaneEndpoint return control plane endpoint address
func GetControlPlaneEndpoint(clusterConfig *pb.ClusterConfig, masterNodes []*pb.Node) (addr string, err error) {
	conn := clusterConfig.KubeAPIServerConnect
	if conn == nil {
		return "", fmt.Errorf("nil %T encountered", conn)
	}

	// type could be ["firstMasterIP", "keepalived", "loadbalancer"]
	switch conn.Type {
	case "firstMasterIP":
		ip := masterNodes[0].Ip
		if ip == "" {
			err = fmt.Errorf("failed to get first master ip")
			return
		}
		addr = fmt.Sprintf("%v:%v", ip, defaultApiServerPort)
	case "keepalived":
		addr = fmt.Sprintf("%v:%v", conn.Keepalived.Vip, defaultHAProxyPort)
	case "loadbalancer":
		addr = fmt.Sprintf("%v:%v", conn.Loadbalancer.Ip, conn.Loadbalancer.Port)
	default:
		err = fmt.Errorf("unrecognized apiserver connect type: %v", conn.Type)
	}

	return
}

// PBErrLogger creates a new logging entry with the content of a pb.Error added as struct info,
// the new entry is set based on the passed in logging entry.
func PBErrLogger(pbErr *pb.Error, entry *logrus.Entry) *logrus.Entry {
	if pbErr == nil {
		// create an empty pb.Error to avoid return an error to caller.
		pbErr = new(pb.Error)
	}
	fields := logrus.Fields{
		"reason":    pbErr.Reason,
		"detail":    pbErr.Detail,
		"fixMethod": pbErr.FixMethods,
	}

	if entry == nil {
		return logrus.WithFields(fields)
	}
	return entry.WithFields(fields)
}

// return size with unit
func ReturnWithUnit(size float64) (valueWithUnit string) {
	if size < 1 {
		valueWithUnit = fmt.Sprintf("%.0f Byte", size)
	} else if size < 1024 {
		valueWithUnit = fmt.Sprintf("%.0f Bytes", size)
	} else if size < 1048576 {
		valueWithUnit = fmt.Sprintf("%.0f KiB", size/1024)
	} else if size < 1073741824 {
		valueWithUnit = fmt.Sprintf("%.0f MiB", size/1024/1024)
	} else {
		valueWithUnit = fmt.Sprintf("%.0f GiB", size/1024/1024/1024)
	}
	return
}

//func FetchEtcdCertAndKey(etcdNode *pb.Node, baseName string) (*x509.Certificate, crypto.Signer, error) {
//	certPath := fmt.Sprintf("%v/%v.crt", localEtcdCADir, baseName)
//	keyPath := fmt.Sprintf("%v/%v.key", localEtcdCADir, baseName)
//
//	localCert, err := os.Create(certPath)
//	localKey, err := os.Create(keyPath)
//
//	if err != nil {
//		return nil, nil, fmt.Errorf("failed to create local %v cert path:%v, error:%v", baseName, certPath, err)
//	}
//
//	m, err := machine.NewMachine(etcdNode)
//	if err != nil {
//		return nil, nil, fmt.Errorf("failed to create exec client for first etcd node:%v, error:%v", m.GetName(), err)
//	}
//
//	if err := m.FetchFile(localCert, etcd.DefaultEtcdCACertPath); err != nil {
//		return nil, nil, fmt.Errorf("failed to fetch etcd %v cert, error:%v", baseName, err)
//	}
//
//	if err := m.FetchFile(localKey, etcd.DefaultEtcdCAKeyPath); err != nil {
//		return nil, nil, fmt.Errorf("failed to fetch etcd %v key, error:%v", baseName, err)
//	}
//
//	etcdCACrt, etcCAKey, err := pkiutil.TryLoadCertAndKeyFromDisk(localEtcdCADir, baseName)
//	if err != nil {
//		err = fmt.Errorf("failed to load etcd %v cert and key from:%v, error:%v", baseName, localEtcdCADir, err)
//		logrus.Error(err)
//		return nil, nil, err
//	}
//
//	return etcdCACrt, etcCAKey, nil
//}
