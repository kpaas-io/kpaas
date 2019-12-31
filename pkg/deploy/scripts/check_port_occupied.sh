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

portSet=(2379, 2380, 10249, 10250)
occupiedSet=

function detectPort() {
    for port in ${portSet[@]}; do
        count = `netstat -nltup | grep $port | wc -l`
        if [[ $count -ne 0 ]]; then
            occupiedSet+="$port, "
        fi
    done

    echo $occupiedSet
}

detectPort
