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

#HAPROXY_VERSION=1.9.6
#KEEPALIVED_VERSION=1.4.5
# latest version of haproxy and keepalived:
# |              | haproxy | keepalived |
# | ------------ | ------- | ---------- |
# | ubuntu 16.04 | 1.6.3   | 1.2.24     |
# | centos 7     | 1.5.18  | 1.3.5      |
# | rhel         | ----    | ----       |

HAPROXY_CONFIG_DIR=/etc/haproxy
HAPROXY_CONFIG=${HAPROXY_CONFIG_DIR}/haproxy.cfg
KEEPALIVED_CONFIG_DIR=/etc/keepalived
KEEPALIVED_CONFIG=${KEEPALIVED_CONFIG_DIR}/keepalived.conf

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd ${ROOT}

source "${ROOT}/init_deploy_haproxy_keepalived/lib.sh"

PKG_MGR=
PKG_INSTALL_OPTIONS=
PKG_VERSION_SYMBOL=

init() {
    local dist=$(. /etc/os-release && echo $ID)
    case ${dist} in
    ubuntu)
        PKG_MGR=apt
        PKG_INSTALL_OPTIONS='-y --allow-unauthenticated'
        PKG_VERSION_SYMBOL='='
    ;;
    centos)
        PKG_MGR=yum
        PKG_INSTALL_OPTIONS='-y --setopt=obsoletes=0 --nogpgcheck'
        PKG_VERSION_SYMBOL='-'
    ;;
    rhel)
        PKG_MGR=yum
        PKG_INSTALL_OPTIONS='-y --setopt=obsoletes=0 --nogpgcheck'
        PKG_VERSION_SYMBOL='-'
    ;;
    esac
}

keepalived::install() {
    ${PKG_MGR} install -y -q ${PKG_INSTALL_OPTIONS} keepalived &>/dev/null
}

keepalived::uninstall() {
    ${PKG_MGR} autoremove ${PKG_INSTALL_OPTIONS} keepalived
}

haproxy::install() {
    ${PKG_MGR} install -y -q ${PKG_INSTALL_OPTIONS} haproxy &>/dev/null
}

haproxy::uninstall() {
    ${PKG_MGR} autoremove ${PKG_INSTALL_OPTIONS} haproxy
}

haproxy::start() {
    systemctl start haproxy
}

haproxy::enable() {
    systemctl enable haproxy
}

haproxy::disable() {
    systemctl disable haproxy
}

haproxy::stop() {
    systemctl stop haproxy
}

haproxy::reload() {
    systemctl reload haproxy
}

haproxy::restart() {
    systemctl restart haproxy
}

haproxy::status() {
    systemctl status haproxy
}

haproxy::run() {
    haproxy::config
    haproxy::install
    haproxy::start
    haproxy::enable
}

haproxy::clean() {
    haproxy::stop
    haproxy::disable
    haproxy::uninstall
    haproxy::unconfig
}

keepalived::start() {
    systemctl start keepalived
}

keepalived::enable() {
    systemctl enable -q keepalived
}

keepalived::disable() {
    systemctl disable keepalived
}

keepalived::stop() {
    systemctl stop keepalived
}

keepalived::restart() {
    systemctl restart keepalived
}

keepalived::status() {
    systemctl status keepalived
}

keepalived::run() {
    keepalived::config
    keepalived::install
    keepalived::start
    keepalived::enable
}

keepalived::clean() {
    keepalived::stop
    keepalived::disable
    keepalived::uninstall
    keepalived::unconfig
}

init

