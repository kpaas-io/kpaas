// Copyright 2020 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/stretchr/testify/assert"
)

func TestSetNetwork(t *testing.T) {
	tests := []struct {
		inputBody      []byte
		wantOptions    *api.NetworkOptions
		wantStatusCode int
	}{
		{
			inputBody: []byte(`{
	"networkType":"calico",
	"calicoOptions": {
		"encapsulationMode":"vxlan",
		"vxlanPort": 1234,
		"initialPodIPs": "10.0.0.0/16",
		"vethMtu":1440,
		"ipDetectionMethod": "interface",
		"ipDetectionInterface": "eth0"
	}
	}`),
			wantOptions: &api.NetworkOptions{
				NetworkType: "calico",
				CalicoOptions: &api.CalicoOptions{
					EncapsulationMode:    api.EncapsulationVxlan,
					VxlanPort:            1234,
					InitialPodIPs:        "10.0.0.0/16",
					VethMtu:              1440,
					IPDetectionMethod:    api.IPDetectionMethodInterface,
					IPDetectionInterface: "eth0",
				},
			},
			wantStatusCode: http.StatusCreated,
		},

		{
			inputBody: []byte(`{"networkType":"calico",
			"calicoOptions":{"vxlanPort":"abcd"}}`),
			wantOptions:    &wizard.DefaultNetworkOptions,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		wizard.ClearCurrentWizardData()
		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		bodyReader := bytes.NewReader(testCase.inputBody)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/networks", bodyReader)

		SetNetwork(ctx)
		assert.Equal(t, testCase.wantStatusCode, resp.Code)
		assert.Equal(t, testCase.wantOptions, wizard.GetCurrentWizard().NetworkOptions)

	}
}

func TestGetNetwork(t *testing.T) {
	tests := []struct {
		inputOptions *api.NetworkOptions
		wantOptions  *api.NetworkOptions
	}{
		{
			inputOptions: nil,
			wantOptions:  &wizard.DefaultNetworkOptions,
		},
		{
			inputOptions: &api.NetworkOptions{
				NetworkType: api.NetworkTypeCalico,
				CalicoOptions: &api.CalicoOptions{
					EncapsulationMode: api.EncapsulationVxlan,
					VxlanPort:         1234,
					InitialPodIPs:     "10.0.0.0/16",
				},
			},

			wantOptions: &api.NetworkOptions{
				NetworkType: "calico",
				CalicoOptions: &api.CalicoOptions{
					EncapsulationMode: api.EncapsulationVxlan,
					VxlanPort:         1234,
					InitialPodIPs:     "10.0.0.0/16",
				},
			},
		},
	}
	for _, testCase := range tests {
		wizard.ClearCurrentWizardData()
		wizard.GetCurrentWizard().SetNetworkOptions(testCase.inputOptions)

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		bodyReader := bytes.NewReader([]byte{})
		ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/networks", bodyReader)

		GetNetwork(ctx)
		var options *api.NetworkOptions
		err := json.Unmarshal(resp.Body.Bytes(), &options)
		assert.Nil(t, err)
		assert.Equal(t, testCase.wantOptions, options)
	}
}
