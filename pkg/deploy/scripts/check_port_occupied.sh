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

# This script is aim to check if port occupied

# 80, 443 ingress
# 2379, 2380 etcd
# 6443 kube-apiserver
# 10248 kubelet (healthz-port)
# 10249 kube-proxy(healthz-port)
# 10250 kubelet (local-listen-port)
# 10251 kube-scheduler
# 10252 kube-manager

ingressNode=(80 443 10248 10249 10250)
masterNode=(6443 10248 10249 10250 10251 10252)
etcdNode=(2379 2380 10248 10249 10250)
role=()

function main() {
    copyArray $1
    detectPort
}

function copyArray() {
    test -n $1 || (echo "role can not be non-exists" 2>&1 && exit 1)
    if [[ $1 == "" ]]; then
        echo "role can not be empty string" 2>&1 && exit 1
    fi

    if [[ $1 == "master" ]]; then
        role=(${masterNode[@]})
    elif [[ $1 == "ingress" ]]; then
        role=(${ingressNode[@]})
    elif [[ $1 == "etcd" ]]; then
        role=(${etcdNode[@]})
    else
        echo "role can not be recognise, exit..." && exit 1
    fi
}

function detectPort() {
    occupiedSet=
    for port in ${role[@]}; do

        countPort=`netstat -nltp | grep -v "Active" | grep -v "Proto" | awk '{print $4}' | awk -F ":" '{print $NF}' | grep -w $port | wc -l`
        if [[ $countPort -ne 0 ]]; then
            occupiedSet=$occupiedSet"$port, "
        fi
    done

    echo $occupiedSet
}

main $@
