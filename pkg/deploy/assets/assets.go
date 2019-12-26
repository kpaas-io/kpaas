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

// +build dev

package assets

import (
	"net/http"
	"os"
	"strings"

	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/union"
)

// assign scriptsPath as directory for static assembled
// we can add more static file as config in the future
const (
	relativeScriptsPath string = "../scripts"
)

// define scripts as filesystem contains entries in relative path
var scripts http.FileSystem = filter.Keep(
	http.Dir(relativeScriptsPath),
	func(path string, fi os.FileInfo) bool {
		return path == "/" ||
			strings.HasPrefix(path, "/init_deploy_haproxy_keepalived") ||
			strings.HasSuffix(path, ".sh")
	},
)

// var Assets contains the deploy's scripts
var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/scripts": scripts,
})
