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

die() {
  if test $# -gt 0; then
    echo "$@" 1>&2
  fi
  exit 1
}

haproxy::config() {
    mkdir -p ${HAPROXY_CONFIG_DIR}

    local upstreams=""
    for i in "${!KUBE_APISERVER_ADDRS[@]}";
    do
        printf -v line "  server server$[i+1] ${KUBE_APISERVER_ADDRS[$i]} maxconn 2048 check fall 3 rise 2\n"
        upstreams+=${line}
    done

    cat > ${HAPROXY_CONFIG} <<EOF
global
  #user haproxy
  #group haproxy
  daemon
  maxconn 4096
defaults
  mode    tcp
  balance leastconn
  timeout client      30s
  timeout server      30s
  timeout connect      3s
  retries 3
listen stats
  bind 0.0.0.0:1936
  mode http
  stats enable
  stats uri /
frontend kubernetes
  bind 0.0.0.0:${HAPROXY_PORT}
  default_backend kube-apiserver
backend kube-apiserver
${upstreams}
EOF
}

haproxy::config::check() {
    if [[ -z "${KUBE_APISERVER_ADDRS+x}" ]]; then
        usage
        die "KUBE_APISERVER_ADDRS is required"
    fi
}

haproxy::unconfig() {
    rm -rf ${HAPROXY_CONFIG_DIR}
}

keepalived::config() {
    mkdir -p ${KEEPALIVED_CONFIG_DIR}

    local virtual_router_id=${VIP##*.}
    cat > ${KEEPALIVED_CONFIG} <<EOF
vrrp_script chk_proxy {
  script "nc -w 3 -z 127.0.0.1 ${HAPROXY_PORT}"
  interval 2
  weight -100
  fall 3
  rise 2
}

vrrp_instance kubernetes_ha {
  interface ${INTERFACE}
  virtual_router_id ${virtual_router_id}
  #nopreempt
  priority 100
  advert_int 1
  authentication {
    auth_type MD5
  }
  virtual_ipaddress {
    ${VIP}/32
  }
  track_script {
    chk_proxy
  }
}
EOF
}

keepalived::config::check() {
    if [[ -z "${VIP+x}" ]]; then
        usage
        die "VIP is required"
    fi
    if [[ -z "${INTERFACE+x}" ]]; then
        usage
        die "INTERFACE is required"
    fi
}

keepalived::unconfig() {
    rm -rf ${KEEPALIVED_CONFIG_DIR}
}
