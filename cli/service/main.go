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

// @title kpaasRestfulApi
// @version 0.1
// @description KPaaS RESTful API service for frontend and using Deploy service API to deployment kubernetes cluster.

// @contact.name Support
// @contact.url http://github.com/kpaas-io/kpaas/issues
// @contact.email support@kpaas.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

package main

import (
	"os"

	"github.com/kpaas-io/kpaas/pkg/service/application"
)

func main() {

	if err := application.GetCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
