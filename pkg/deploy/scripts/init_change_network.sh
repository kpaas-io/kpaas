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

# This script is aim to change basic network env

grep "net.ipv4.ip_forward" /etc/sysctl.conf && {
    sed -i -E 's/(.*)net.ipv4.ip_forward(.*)/net.ipv4.ip_forward\2/1' /etc/sysctl.conf
}

sysctl -w net.ipv4.ip_forward=1 1>/dev/null
sysctl -p /etc/sysctl.conf 1>/dev/null
