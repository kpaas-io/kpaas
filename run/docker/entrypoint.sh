#!/usr/bin/env sh
## Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
##
## Licensed under the Apache License, Version 2.0 (the "License");
## you may not use this file except in compliance with the License.
## You may obtain a copy of the License at
##
##      http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software
## distributed under the License is distributed on an "AS IS" BASIS,
## WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and
## limitations under the License.

logLevel=${LOG_LEVEL:-info}
serviceId=${SERVICE_ID:-0}
dashboardPort="${DASHBOARD_PORT:-80}"

dashboardOptions=""
if [ "${logLevel}" = "debug" ]; then
  dashboardOptions="--debug=true"
fi

/app/service --log-level="${logLevel}" --service-id="${serviceId}" >>/var/log/kpaas.log 2>&1 &
/app/deploy --log-level="${logLevel}" >>/var/log/kpaas.log 2>&1 &
/app/dashboard --config-file=/app/config/dashboard.conf --server-addr="0.0.0.0:${dashboardPort}" "${dashboardOptions}" >>/var/log/kpaas.log 2>&1 &

tail -f /var/log/kpaas.log
