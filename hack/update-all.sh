#!/bin/bash

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

set -o errexit
set -o nounset
set -o pipefail

hack/update-gofmt.sh
