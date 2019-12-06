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

#KUBE_APISERVER_ADDRS=(
#  "10.0.0.9:6443"
#)
#VIP="10.0.0.88"
#INTERFACE="ens3"

HAPROXY_PORT=4443
HAPROXY_STATS_PORT=1936
# DAEMON: "systemd" or "docker"
DAEMON=systemd

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd ${ROOT}

source "${ROOT}/init_deploy_haproxy_keepalived/lib.sh"

source "${ROOT}/init_deploy_haproxy_keepalived/${DAEMON}.sh"

usage() {
    cat <<EOF
Usage: $0 [flags] <action> [app]

Flags:
    -u  kube-apiserver proxy upstream addrs
    -n  ha vip
    -i  ha bind interface

Apps:
    haproxy
    keepalived

Actions:
    run     config and start app
    clean   stop and unconfig app

    config  generate config file to config dir
    start   start app
    reload  reload app config from config file
    restart stop and start app
    status  get running status of app
    stop    stop app
    unconfig remove config dir of app

Examples:
    $0 -u "10.0.0.10:6443 10.0.0.11:6443 10.0.0.12:6443" haproxy run
    $0 -n "10.0.0.88" -i eth0 keepalived run
    $0 haproxy clean|start|reload|status|stop|unconfig
    $0 keepalived clean|start|restart|reload|status|stop|unconfig
EOF
}

# Resetting OPTIND is necessary if getopts was used previously in the script.
# It is a good idea to make OPTIND local if you process options in a function.
OPTIND=1

while getopts ":hu:n:i:" opt; do
    case "$opt" in
    h)
        usage
        ;;
    u)
        KUBE_APISERVER_ADDRS=($OPTARG)
        ;;
    n)
        VIP=$OPTARG
        ;;
    i)
        INTERFACE=$OPTARG
        ;;
    \?)
        usage "Invalid option: -$OPTARG"
        ;;
    :)
        usage "Option -$OPTARG requires an argument."
        ;;
    esac
done
shift $((OPTIND-1))

if [[ "$#" < 2 ]]; then
    usage
    die
fi

APP=$1
ACTION=$2

if [[ "config|run" == *"${ACTION}"* ]];
then
    ${APP}::config::check
fi

${APP}::${ACTION}
