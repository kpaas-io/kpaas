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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/sshcertificate"
)

func TestAddSSHCertificate(t *testing.T) {

	sshcertificate.ClearList()

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	keyName := "id_rsa"
	privateKey := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEArWHF1jw6tVmUlFexam9F+BSFDQuIBsJ/Qo/yH8gl8P5HGJKL56nQ
C2OAJ9+Tks6Ay3cO40yvKBRhgZtEro+kkgBUiY1vzOfMPXoHa2U94lDa83Dkp0FHkdCxUy
INIk7fKnmX/x43/VTmbXcInM8oLUVf//MzM+Qqz5XhlF0Y1eUvNEXHPBV38oqXe6kbTGg2
TiBSfO08s2+ZctyWrlgY2hmS5MbkvkgcCnXfSVU6m74IVU6dTjShrCYS+AT/LFE+KseXf6
PeSH8KmGMIJApTfKT2y6JfT7MH4a+Wts+Kc2vNo6kQN1M7qA8237VMirUHCGXJZ5C/H7Mw
DqscRp9/CQAAA9DlbaVP5W2lTwAAAAdzc2gtcnNhAAABAQCtYcXWPDq1WZSUV7Fqb0X4FI
UNC4gGwn9Cj/IfyCXw/kcYkovnqdALY4An35OSzoDLdw7jTK8oFGGBm0Suj6SSAFSJjW/M
58w9egdrZT3iUNrzcOSnQUeR0LFTIg0iTt8qeZf/Hjf9VOZtdwiczygtRV//8zMz5CrPle
GUXRjV5S80Rcc8FXfyipd7qRtMaDZOIFJ87Tyzb5ly3JauWBjaGZLkxuS+SBwKdd9JVTqb
vghVTp1ONKGsJhL4BP8sUT4qx5d/o95IfwqYYwgkClN8pPbLol9Pswfhr5a2z4pza82jqR
A3UzuoDzbftUyKtQcIZclnkL8fszAOqxxGn38JAAAAAwEAAQAAAQBz5jTaZf6UtaIVm500
Wde64vSh6MBwTFnHg/PFfQSn2UJrUaMGJES3KDdF8DV04GfGGvsvxFYeA6m+eq1pxwmqs1
/PZ2WB4r1rpwQIrW+1tnj2XNPsXj3aYlf3C38eHP0fJpMNbgTdaoByUizGrc/cm1B2BvuG
R5K6myVlCOqOKAA7gaa6IwWbjZ40vfakIk+j08PNOAXAoKT6kDEmbBmdg22UAwbgTcjM2i
9K6PTKk6uLaagEMZcSPtUkbZo/j4beaW9RejQb/gSAomM0P+Iw4GFe2DJuig8+jPQOSHj/
8ryHNtnL/0/yihOsaAn5MWMOEzjjzKYAERXBso7PBu+VAAAAgQDI/1UHQtr1vUX9dVAItJ
19QMblA7lalRe18sd0WCxGY9+5Kf2wHQ6iWhrwHSRYZN6j65r78Uiw2JNFIE7UcteJQSnS
EddivJWFbYbffw9yrV0+Sb24N2Ckm9Uhf1bkwp+aNMIk7v9sBhv6sUmByC4DUG4E+Ypz+g
3pn/Np8yXneQAAAIEA5WDIQSTZYxNGai4w2bbxWMFynLiAhQqg0s0ISeQTcIWc4nrB28mK
Kh6dWFfjUMLm50M9LB7UndzYTdiuY8fTBUSpBGapWBTbbDoBolYlQZtFL7+Rj3dqMC1lOt
gKY6UUTMFIM9c7yGAK7zhfds57V5LALIPUNkTQU8kyePVhSBcAAACBAMGBP7iXiX+T6xTw
UpYgsuXDxma0RK713UCzZCP8o1ucwoNdGJaRDxIEIzTyKAwC2k6vURxlY6djoZHv79sfD7
NiyK6OkjUmiwIwsL4IQ/dsFD+Lrfp1Ilo3Yirz1UE3Zg6UNP5GUKiys8WnvvtC28uv4dGy
ls3Q/5aeF7hB2MXfAAAAGEx1Y2t5Ym95c0BMdWNreU1hYy5sb2NhbAEC
-----END OPENSSH PRIVATE KEY-----`
	body := api.SSHCertificate{
		Name:    keyName,
		Content: privateKey,
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/ssh_certificates", bodyReader)

	AddSSHCertificate(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.SuccessfulOption)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.True(t, responseData.Success)
}

func TestGetCertificateList(t *testing.T) {

	sshcertificate.ClearList()
	keyName := "id_rsa"
	privateKey := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEArWHF1jw6tVmUlFexam9F+BSFDQuIBsJ/Qo/yH8gl8P5HGJKL56nQ
C2OAJ9+Tks6Ay3cO40yvKBRhgZtEro+kkgBUiY1vzOfMPXoHa2U94lDa83Dkp0FHkdCxUy
INIk7fKnmX/x43/VTmbXcInM8oLUVf//MzM+Qqz5XhlF0Y1eUvNEXHPBV38oqXe6kbTGg2
TiBSfO08s2+ZctyWrlgY2hmS5MbkvkgcCnXfSVU6m74IVU6dTjShrCYS+AT/LFE+KseXf6
PeSH8KmGMIJApTfKT2y6JfT7MH4a+Wts+Kc2vNo6kQN1M7qA8237VMirUHCGXJZ5C/H7Mw
DqscRp9/CQAAA9DlbaVP5W2lTwAAAAdzc2gtcnNhAAABAQCtYcXWPDq1WZSUV7Fqb0X4FI
UNC4gGwn9Cj/IfyCXw/kcYkovnqdALY4An35OSzoDLdw7jTK8oFGGBm0Suj6SSAFSJjW/M
58w9egdrZT3iUNrzcOSnQUeR0LFTIg0iTt8qeZf/Hjf9VOZtdwiczygtRV//8zMz5CrPle
GUXRjV5S80Rcc8FXfyipd7qRtMaDZOIFJ87Tyzb5ly3JauWBjaGZLkxuS+SBwKdd9JVTqb
vghVTp1ONKGsJhL4BP8sUT4qx5d/o95IfwqYYwgkClN8pPbLol9Pswfhr5a2z4pza82jqR
A3UzuoDzbftUyKtQcIZclnkL8fszAOqxxGn38JAAAAAwEAAQAAAQBz5jTaZf6UtaIVm500
Wde64vSh6MBwTFnHg/PFfQSn2UJrUaMGJES3KDdF8DV04GfGGvsvxFYeA6m+eq1pxwmqs1
/PZ2WB4r1rpwQIrW+1tnj2XNPsXj3aYlf3C38eHP0fJpMNbgTdaoByUizGrc/cm1B2BvuG
R5K6myVlCOqOKAA7gaa6IwWbjZ40vfakIk+j08PNOAXAoKT6kDEmbBmdg22UAwbgTcjM2i
9K6PTKk6uLaagEMZcSPtUkbZo/j4beaW9RejQb/gSAomM0P+Iw4GFe2DJuig8+jPQOSHj/
8ryHNtnL/0/yihOsaAn5MWMOEzjjzKYAERXBso7PBu+VAAAAgQDI/1UHQtr1vUX9dVAItJ
19QMblA7lalRe18sd0WCxGY9+5Kf2wHQ6iWhrwHSRYZN6j65r78Uiw2JNFIE7UcteJQSnS
EddivJWFbYbffw9yrV0+Sb24N2Ckm9Uhf1bkwp+aNMIk7v9sBhv6sUmByC4DUG4E+Ypz+g
3pn/Np8yXneQAAAIEA5WDIQSTZYxNGai4w2bbxWMFynLiAhQqg0s0ISeQTcIWc4nrB28mK
Kh6dWFfjUMLm50M9LB7UndzYTdiuY8fTBUSpBGapWBTbbDoBolYlQZtFL7+Rj3dqMC1lOt
gKY6UUTMFIM9c7yGAK7zhfds57V5LALIPUNkTQU8kyePVhSBcAAACBAMGBP7iXiX+T6xTw
UpYgsuXDxma0RK713UCzZCP8o1ucwoNdGJaRDxIEIzTyKAwC2k6vURxlY6djoZHv79sfD7
NiyK6OkjUmiwIwsL4IQ/dsFD+Lrfp1Ilo3Yirz1UE3Zg6UNP5GUKiys8WnvvtC28uv4dGy
ls3Q/5aeF7hB2MXfAAAAGEx1Y2t5Ym95c0BMdWNreU1hYy5sb2NhbAEC
-----END OPENSSH PRIVATE KEY-----`
	sshcertificate.AddCertificate(keyName, privateKey)

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/ssh_certificates", nil)

	GetCertificateList(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.GetSSHCertificateListResponse)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, []string{keyName}, responseData.Names)
}
