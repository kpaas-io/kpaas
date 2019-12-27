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
	"crypto/x509"
	"fmt"
	"net"

	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
)

var (
	caCertConfig = &certutil.Config{
		CommonName: "etcd-ca",
	}

	apiserverClientCertConfig = &certutil.Config{
		CommonName:   kubeadmconstants.APIServerEtcdClientCertCommonName,
		Organization: []string{kubeadmconstants.SystemPrivilegedGroup},

		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
)

func GetCaCrtConfig() *certutil.Config {
	return caCertConfig
}

func GetServerCrtConfig(hostName, IP string) (*certutil.Config, error) {
	ip := net.ParseIP(IP)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse etcd node ip: %v, possiblly invalid ip", IP)
	}

	serverCertConfig := &certutil.Config{
		// TODO: etcd 3.2 introduced an undocumented requirement for ClientAuth usage on the
		// server cert: https://github.com/coreos/etcd/issues/9785#issuecomment-396715692
		// Once the upstream issue is resolved, this should be returned to only allowing
		// ServerAuth usage.
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	// create AltNames with defaults DNSNames/IPs
	serverCertConfig.AltNames = certutil.AltNames{
		DNSNames: []string{hostName, "localhost"},
		IPs:      []net.IP{ip, net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}

	// set server cert CommonName to etcd hostname
	serverCertConfig.CommonName = hostName

	return serverCertConfig, nil
}

func GetPeerCrtConfig(hostName, IP string) (*certutil.Config, error) {
	ip := net.ParseIP(IP)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse etcd node ip: %v, possiblly invalid ip", IP)
	}

	peerCertConfig := &certutil.Config{
		// TODO: etcd 3.2 introduced an undocumented requirement for ClientAuth usage on the
		// server cert: https://github.com/coreos/etcd/issues/9785#issuecomment-396715692
		// Once the upstream issue is resolved, this should be returned to only allowing
		// ServerAuth usage.
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	// create AltNames with defaults DNSNames/IPs
	peerCertConfig.AltNames = certutil.AltNames{
		DNSNames: []string{hostName, "localhost"},
		IPs:      []net.IP{ip, net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}

	// set server cert CommonName to etcd hostname
	peerCertConfig.CommonName = hostName

	return peerCertConfig, nil
}

func GetAPIServerClientCrtConfig() *certutil.Config {
	return apiserverClientCertConfig
}
