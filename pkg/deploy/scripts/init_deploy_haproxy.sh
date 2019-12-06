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

# This script is aim to deploy haproxy

. lib.sh

OS=GetOS
HAPROXYADDR=$1

test -n $HAPROXYADDR || (echo "haproxy addr is not specific" && exit 1)

if [[ $OS -eq "ubuntu" ]]; then
    apt-get install -y haproxy
    dpkg -l | grep -q haproxy || (echo "haproxy install failed" && exit 1)
elif [[ $OS -eq "centos" ]] || [[ $OS -eq "rhel" ]]; then
    yum install -y haproxy
    rpm -qa | grep -q haproxy || (echo "haproxy install failed" && exit 1)
else
    echo "system not support"
    exit 1
fi

/bin/bash /tmp/init_deploy_haproxy_keepalived/setup.sh -u "$HAPROXYADDR" haproxy run 2>&1
