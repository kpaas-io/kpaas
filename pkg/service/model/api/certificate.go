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

package api

import (
	"golang.org/x/crypto/ssh"

	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

type (
	SSHCertificate struct {
		Name    string `json:"name" binding:"required" minimum:"1" maximum:"20"`
		Content string `json:"content" binding:"required"`
	}

	GetSSHCertificateListResponse struct {
		Names []string `json:"names"`
	}
)

const (
	CertificateNameLimit    = 20
	CertificateContentLimit = 10000
)

func (cert *SSHCertificate) Validate() error {

	return validator.NewWrapper(
		validator.ValidateString(cert.Name, "name", validator.ItemNotEmptyLimit, CertificateNameLimit),
		validator.ValidateString(cert.Content, "content", validator.ItemNotEmptyLimit, CertificateContentLimit),
		func() error {
			return verifyPrivateKeyContent(cert.Content)
		},
	).Validate()
}

func verifyPrivateKeyContent(content string) (err error) {

	_, err = ssh.ParsePrivateKey([]byte(content))
	return
}
