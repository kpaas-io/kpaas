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

# This script provide common function for scripts use

function::exists () {
  declare -f -F $1 > /dev/null
  return $?
}

echored() {
    echo -e "\033[31m$@\033[0m"
}

echogreen() {
    echo -e "\033[32m$@\033[0m"
}

echoyellow() {
    echo -e "\033[33m$@\033[0m"
}

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

error_exit() {
    echored $@
    exit 1
}

log() {
        local level=$1
                test -z "$level" && error_exit "[FATAL] log level not specified"

        case "$level" in
                INFO)
                        echogreen "[$level] ${@:2}"
                        ;;
                WARN)
                        echoyellow "[$level] ${@:2}"
                        ;;
                ERR)
                        error_exit "[$level] ${@:2}"
                        ;;
                *)
                        error_exit "[FATAL] unknown log level $level"
                        ;;
                esac
}

In() {
    local n=
    for n in $@; do
        if [[ $1 == $n ]]; then return 0; fi
    done
    return 1
}

RoleMatch() {
    local array=$(echo $1 | tr "," "\n")
    for i in $array; do
        if [[ "$i" == $2 ]]; then
            return 0
        fi
    done
    return 1
}

GetOS() {
    OS=`cat /etc/*-release | grep -w "ID" | awk '/ID/{print $1}' | awk -F "=" '{print $2}'`
    echo $OS
}
