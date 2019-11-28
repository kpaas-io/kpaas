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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/common"
	"github.com/kpaas-io/kpaas/pkg/service/model/sshcertificate"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
)

func TestConvertModelAuthenticationTypeToAPIAuthenticationType(t *testing.T) {

	assert.Equal(t, api.AuthenticationTypePassword, convertModelAuthenticationTypeToAPIAuthenticationType(wizard.AuthenticationTypePassword))
	assert.Equal(t, api.AuthenticationTypePrivateKey, convertModelAuthenticationTypeToAPIAuthenticationType(wizard.AuthenticationTypePrivateKey))
	assert.Equal(t, api.AuthenticationType("unknown(OtherType)"), convertModelAuthenticationTypeToAPIAuthenticationType("OtherType"))
}

func TestConvertModelLabelToAPILabel(t *testing.T) {

	assert.EqualValues(t, api.Label{
		Key:   "key",
		Value: "value",
	}, convertModelLabelToAPILabel(&wizard.Label{
		Key:   "key",
		Value: "value",
	}))
}

func TestConvertModelAnnotationToAPIAnnotation(t *testing.T) {

	assert.EqualValues(t,
		api.Annotation{
			Key:   "key",
			Value: "value",
		}, convertModelAnnotationToAPIAnnotation(&wizard.Annotation{
			Key:   "key",
			Value: "value",
		}),
	)
}

func TestConvertModelTaintToAPITaint(t *testing.T) {

	assert.EqualValues(t,
		api.Taint{
			Key:    "key",
			Value:  "value",
			Effect: api.TaintEffectNoSchedule,
		},
		convertModelTaintToAPITaint(&wizard.Taint{
			Key:    "key",
			Value:  "value",
			Effect: wizard.TaintEffectNoSchedule,
		}),
	)
}

func TestConvertModelTaintEffectToAPITaintEffect(t *testing.T) {

	assert.Equal(t, api.TaintEffectNoSchedule, convertModelTaintEffectToAPITaintEffect(wizard.TaintEffectNoSchedule))
	assert.Equal(t, api.TaintEffectNoExecute, convertModelTaintEffectToAPITaintEffect(wizard.TaintEffectNoExecute))
	assert.Equal(t, api.TaintEffectPreferNoSchedule, convertModelTaintEffectToAPITaintEffect(wizard.TaintEffectPreferNoSchedule))
	assert.Equal(t, api.TaintEffect("unknown(OtherType)"), convertModelTaintEffectToAPITaintEffect("OtherType"))
}

func TestConvertModelErrorToAPIError(t *testing.T) {

	assert.EqualValues(t,
		&api.Error{
			Reason:     "reason",
			Detail:     "detail",
			FixMethods: "fixMethods",
			LogId:      1234,
		},
		convertModelErrorToAPIError(&common.FailureDetail{
			Reason:     "reason",
			Detail:     "detail",
			FixMethods: "fixMethods",
			LogId:      1234,
		}),
	)
}

func TestConvertModelDeployClusterStatusToAPIDeployClusterStatus(t *testing.T) {

	assert.Equal(t, api.DeployClusterStatusNotRunning, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusNotRunning))
	assert.Equal(t, api.DeployClusterStatusRunning, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusRunning))
	assert.Equal(t, api.DeployClusterStatusSuccessful, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusSuccessful))
	assert.Equal(t, api.DeployClusterStatusFailed, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusFailed))
	assert.Equal(t, api.DeployClusterStatusWorkedButHaveError, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusWorkedButHaveError))
	assert.Equal(t, api.DeployClusterStatus("unknown(OtherType)"), convertModelDeployClusterStatusToAPIDeployClusterStatus("OtherType"))
}

func TestConvertModelDeployStatusToAPiDeployStatus(t *testing.T) {

	assert.Equal(t, api.DeployStatusPending, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusPending))
	assert.Equal(t, api.DeployStatusDeploying, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusDeploying))
	assert.Equal(t, api.DeployStatusCompleted, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusCompleted))
	assert.Equal(t, api.DeployStatusFailed, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusFailed))
	assert.Equal(t, api.DeployStatusAborted, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusAborted))
	assert.Equal(t, api.DeployStatus("unknown(OtherType)"), convertModelDeployStatusToAPIDeployStatus("OtherType"))
}

func TestConvertAPITaintEffectToModelTaintEffect(t *testing.T) {

	assert.Equal(t, wizard.TaintEffectNoSchedule, convertAPITaintEffectToModelTaintEffect(api.TaintEffectNoSchedule))
	assert.Equal(t, wizard.TaintEffectNoExecute, convertAPITaintEffectToModelTaintEffect(api.TaintEffectNoExecute))
	assert.Equal(t, wizard.TaintEffectPreferNoSchedule, convertAPITaintEffectToModelTaintEffect(api.TaintEffectPreferNoSchedule))
	assert.Equal(t, wizard.TaintEffect("unknown(OtherType)"), convertAPITaintEffectToModelTaintEffect("OtherType"))
}

func TestConvertAPIAuthenticationTypeToModelAuthenticationType(t *testing.T) {

	assert.Equal(t, wizard.AuthenticationTypePassword, convertAPIAuthenticationTypeToModelAuthenticationType(api.AuthenticationTypePassword))
	assert.Equal(t, wizard.AuthenticationTypePrivateKey, convertAPIAuthenticationTypeToModelAuthenticationType(api.AuthenticationTypePrivateKey))
	assert.Equal(t, wizard.AuthenticationType("unknown(OtherType)"), convertAPIAuthenticationTypeToModelAuthenticationType("OtherType"))
}

func TestConvertDeployControllerErrorToAPIError(t *testing.T) {

	var nilStruct *api.Error
	assert.Equal(t, nilStruct, convertDeployControllerErrorToAPIError(nil))
	assert.Equal(t, &api.Error{
		Reason:     "reason",
		Detail:     "detail",
		FixMethods: "fixMethods",
		LogId:      0,
	}, convertDeployControllerErrorToAPIError(&protos.Error{
		Reason:     "reason",
		Detail:     "detail",
		FixMethods: "fixMethods",
	}))
}

func TestConvertModelConnectionDataToDeployControllerSSHData(t *testing.T) {

	var nilStruct *protos.SSH
	assert.Equal(t, nilStruct, convertModelConnectionDataToDeployControllerSSHData(nil))
	assert.Equal(t, &protos.SSH{
		Port: 22,
		Auth: &protos.Auth{
			Type:       "password",
			Username:   "root",
			Credential: "123456",
		},
	}, convertModelConnectionDataToDeployControllerSSHData(&wizard.ConnectionData{
		Port:               uint16(22),
		Username:           "root",
		AuthenticationType: wizard.AuthenticationTypePassword,
		Password:           "123456",
	}))

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

	assert.Equal(t, &protos.SSH{
		Port: 22,
		Auth: &protos.Auth{
			Type:       "privatekey",
			Username:   "root",
			Credential: privateKey,
		},
	}, convertModelConnectionDataToDeployControllerSSHData(&wizard.ConnectionData{
		Port:               uint16(22),
		Username:           "root",
		AuthenticationType: wizard.AuthenticationTypePrivateKey,
		PrivateKeyName:     keyName,
	}))
}

func TestConvertDeployControllerCheckResultToModelCheckResult(t *testing.T) {

	assert.Equal(t, constant.CheckResultNotRunning, convertDeployControllerCheckResultToModelCheckResult(string(constant.CheckResultNotRunning)))
	assert.Equal(t, constant.CheckResultChecking, convertDeployControllerCheckResultToModelCheckResult(string(constant.CheckResultChecking)))
	assert.Equal(t, constant.CheckResultPassed, convertDeployControllerCheckResultToModelCheckResult(string(constant.CheckResultPassed)))
	assert.Equal(t, constant.CheckResultFailed, convertDeployControllerCheckResultToModelCheckResult(string(constant.CheckResultFailed)))
	assert.Equal(t, constant.CheckResult("unknown(OtherType)"), convertDeployControllerCheckResultToModelCheckResult("OtherType"))
}
