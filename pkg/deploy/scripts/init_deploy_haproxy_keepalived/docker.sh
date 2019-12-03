#!/usr/bin/env bash
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

set -o errexit
set -o nounset
set -o pipefail

HAPROXY_IMAGE=index-dev.qiniu.io/kelibrary/haproxy
#HAPROXY_IMAGE=haproxy
KEEPALIVED_IMAGE=index-dev.qiniu.io/kelibrary/keepalived
#KEEPALIVED_IMAGE=osixia/keepalived
HAPROXY_VERSION=1.9.6
KEEPALIVED_VERSION=1.4.5

KUBERNETES_CONFIG_DIR=/etc/kubernetes
HAPROXY_CONFIG_DIR=${KUBERNETES_CONFIG_DIR}/haproxy
HAPROXY_CONFIG=${HAPROXY_CONFIG_DIR}/haproxy.cfg
KEEPALIVED_CONFIG_DIR=${KUBERNETES_CONFIG_DIR}/keepalived
KEEPALIVED_CONFIG=${KEEPALIVED_CONFIG_DIR}/keepalived.conf

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd ${ROOT}

source "${ROOT}/init_deploy_haproxy_keepalived/lib.sh"

docker::check() {
    docker version > /dev/null
}

haproxy::start() {
    docker run -dit --restart always --name kubernetes-ha-haproxy -l svc=kubernetes-ha -l app=haproxy -v ${HAPROXY_CONFIG_DIR}:/usr/local/etc/haproxy:ro,slave -p ${HAPROXY_PORT}:${HAPROXY_PORT} -p ${HAPROXY_STATS_PORT}:${HAPROXY_STATS_PORT} ${HAPROXY_IMAGE}:${HAPROXY_VERSION}
}

haproxy::stop() {
    docker rm -f kubernetes-ha-haproxy
}

haproxy::reload() {
    docker kill -s HUP kubernetes-ha-haproxy
}

haproxy::restart() {
    haproxy::stop
    haproxy::start
}

haproxy::status() {
    docker inspect -f '{{.State.Status}}' kubernetes-ha-haproxy
}

haproxy::run() {
    haproxy::config
    haproxy::start
}

haproxy::clean() {
    haproxy::stop
    haproxy::unconfig
}

keepalived::start() {
    docker run -dit --restart always --name kubernetes-ha-keepalived -l svc=kubernetes-ha -l app=keepalived -v ${KEEPALIVED_CONFIG}:/container/service/keepalived/assets/keepalived.conf:ro,slave --cap-add NET_ADMIN --net host ${KEEPALIVED_IMAGE}:${KEEPALIVED_VERSION} --copy-service
}

keepalived::stop() {
    docker rm -f kubernetes-ha-keepalived
}

keepalived::restart() {
    keepalived::stop
    keepalived::start
}

keepalived::status() {
    docker inspect -f '{{.State.Status}}' kubernetes-ha-keepalived
}

keepalived::run() {
    keepalived::config
    keepalived::start
}

keepalived::clean() {
    keepalived::stop
    keepalived::unconfig
}
