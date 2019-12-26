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

package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/coreos/etcd/clientv3"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func newCertPool(caCrt *x509.Certificate) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	_, pemByte, err := ToByte(caCrt, nil)
	if err != nil {
		return nil, err
	}

	for {
		var block *pem.Block
		block, pemByte = pem.Decode(pemByte)
		if block == nil {
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certPool.AddCert(cert)
	}

	return certPool, nil
}

func newCert(d *deployEtcdOperation) (*tls.Certificate, error) {
	tlsCert, err := tls.X509KeyPair(d.encodedPeerCert, d.encodedPeerKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate for etcd client:%v, error: %v", d.machine.GetName(), err)
	}

	return &tlsCert, nil
}

func getClientV3TLS(d *deployEtcdOperation) (*tls.Config, error) {
	var err error

	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: d.machine.GetName(),
	}

	cfg.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return newCert(d)
	}
	cfg.GetClientCertificate = func(unused *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return newCert(d)
	}

	cfg.RootCAs, err = newCertPool(d.caCrt)
	if err != nil {
		return nil, fmt.Errorf("failed to get cert pool for etcd client:%v, error:%v", d.machine.GetName(), err)
	}

	return cfg, nil
}

func composeEndpoints(nodes []*pb.Node) (endpoints []string) {
	for i := range nodes {
		endpoints = append(endpoints, fmt.Sprintf("https://%v:%v", nodes[i].Ip, defaultEtcdServerPort))
	}

	return
}

func newEtcdV3SecureClient(d *deployEtcdOperation) (*clientv3.Client, error) {

	tlsCfg, err := getClientV3TLS(d)
	if err != nil {
		return nil, err
	}

	// check etcd cluster health
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   composeEndpoints(d.clusterNodes),
		DialTimeout: defaultEtcdDialTimeout,
		TLS:         tlsCfg,
	})

	if err != nil {
		return nil, err
	}

	return cli, nil
}
