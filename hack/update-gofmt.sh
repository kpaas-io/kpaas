#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

gofmt=$(which gofmt)

PKGS=$(go list ./...) 

xargs -n 1 -I pkg "$gofmt" -w -s ${GOPATH}/src/pkg <<<"${PKGS}"
