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

# This script is aim to deploy kubectl, kubelet, kubeadm

set -Eeuo pipefail

DEBUG=false
LSB_DIST=
DIST_VERSION=
ACTION=
COMPONENT=
VERSION=
IMAGE_REPOSITORY=index-dev.qiniu.io/kelibrary
DEVICE_MOUNTS=

# kubelet specific
KUBELET_VERSION=
CLUSTER_DNS=
KUBELET_PKG=

# kubeadm specific
JOIN_CONTROL_PLANE=
INIT_CONFIG=/etc/kubernetes/kubeadm_config.yaml

# package specific
PKG_MGR=
INSTALL_OPTIONS=
VERSION_SYMBOL=
LOCALREPO_ADDR=
EXTRAPKGS=
PKG_MIRROR=mirrors.aliyun.com

# library
echored() {
    echo -e "\033[31m$@\033[0m"
}

echogreen() {
    echo -e "\033[32m$@\033[0m"
}

echoyellow() {
    echo -e "\033[33m$@\033[0m"
}

error_exit() {
    echored "$@" >&2
    exit 2
}

command::exists() {
    command -v "$@" > /dev/null 2>&1
}

command::exec() {
    $DEBUG && {
        log D $@
        eval $@
        return
    }

    #log I $@
    err=$( { eval $@ > /dev/null; } 2>&1) || log E "exec $@ failed, err: $err"
}

usage_exit() {
    log I "$@" >&2
    usage
    exit 1
}

log() {
    local level=$1
    local step=$2
    local logts=$(date '+%m%d %T')
    local hostname=$(hostname)

    test -z "$level" && error_exit "F$LOGTS [$hostname] log level not specified"

    case "$level" in
        D)
            echogreen "D$logts [$hostname] ${@:2}"
        ;;
        I)
            echo "I$logts [$hostname] ${@:2}"
        ;;
        W)
            echoyellow "W$logts [$hostname] ${@:2}"
        ;;
        E)
            error_exit "E$logts [$hostname] ${@:2}"
        ;;
        F)
            error_exit "F$logts [$hostname] ${@:2}"
        ;;
        *)
            error_exit "F$logts [$hostname] unknown log level $level"
        ;;
    esac
}

# reference from https://stackoverflow.com/questions/4023830/how-to-compare-two-strings-in-dot-separated-version-format-in-bash
vercomp() {
    if [[ $1 == $2 ]]
    then
        return 0
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

# setup specific
# trap unexpected error
trap 'log F "unexpected error occured at line: $LINENO, command: $BASH_COMMAND"' ERR

clean() {
    command::exists kubeadm && command::exec kubeadm reset -f
    command::exec mount -a
    kubelet::setup::undo
    repos::setup::undo
}

repos::setup() {
    log I "setup package repos and update cache"
    repos::setup::${LSB_DIST}
}

repos::setup::undo() {
    log I "clean up package repos"
    repos::setup::${LSB_DIST}::undo
}

repos::setup::ubuntu() {
    local sourcedir=/etc/apt/sources.list.d
    cat > /etc/apt/sources.list <<EOF
deb http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION} main restricted universe multiverse
deb http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-security main restricted universe multiverse
deb http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-updates main restricted universe multiverse
deb http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-proposed main restricted universe multiverse
deb http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-backports main restricted universe multiverse
deb-src http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION} main restricted universe multiverse
deb-src http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-security main restricted universe multiverse
deb-src http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-updates main restricted universe multiverse
deb-src http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-proposed main restricted universe multiverse
deb-src http://$PKG_MIRROR/ubuntu/ ${DIST_VERSION}-backports main restricted universe multiverse
EOF
    cat > $sourcedir/kubernetes.list <<EOF
deb https://$PKG_MIRROR/kubernetes/apt/ kubernetes-${DIST_VERSION} main
EOF
    command::exec apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 6A030B21BA07F4FB E84AC2C0460F3994 7EA0A9C3F273FCD8 F76221572C52609D
    command::exec apt clean
    command::exec apt update
}

repos::setup::ubuntu::undo() {
    local sourcedir=/etc/apt/sources.list.d
    rm -f $sourcedir/{kirk,kubernetes,local}.list

    command::exec apt-key del 6A030B21BA07F4FB E84AC2C0460F3994 7EA0A9C3F273FCD8 F76221572C52609D
    command::exec apt clean
}

repos::setup::centos() {
    local repodir=/etc/yum.repos.d

    if [[ -z $LOCALREPO_ADDR ]]; then
        command::exists yum-config-manager || {
            command::exec "$PKG_MGR install $INSTALL_OPTIONS yum-utils"
        }

        cat > $repodir/epel.repo <<EOF
[epel]
name=Extra Packages for Enterprise Linux 7 - \$basearch
baseurl=http://$PKG_MIRROR/epel/7/\$basearch
failovermethod=priority
enabled=1
gpgcheck=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7
EOF
        cat > $repodir/k8s.repo <<EOF
[kubernetes]
name = k8s
baseurl = http://$PKG_MIRROR/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled = 1
gpgcheck = 0
EOF
    else
        test -d /etc/yum.repos.d/bak || mkdir /etc/yum.repos.d/bak
        mv -f /etc/yum.repos.d/*.repo /etc/yum.repos.d/bak &> /dev/null || true
        cat > /etc/yum.repos.d/local.repo <<EOF
[kirk-poc]
name=local-yum
baseurl=$LOCALREPO_ADDR
enabled=1
gpgcheck=0
EOF
    fi
}

repos::setup::centos::undo() {
    local repodir=/etc/yum.repos.d
    command::exec yum clean all
    rm -f $repodir/{epel,kirk,k8s,local}.repo
    command::exec rm -rf /var/cache/yum/*
    # command::exec yum makecache fast
}

repos::setup::rhel() {
    local repodir=/etc/yum.repos.d

    if [[ -z $LOCALREPO_ADDR ]]; then
        command::exists yum-config-manager || {
            command::exec "$PKG_MGR install $INSTALL_OPTIONS yum-utils"
        }

        cat > $repodir/epel.repo <<EOF
[epel]
name=Extra Packages for Enterprise Linux 7 - \$basearch
baseurl=http://$PKG_MIRROR/epel/7/\$basearch
failovermethod=priority
enabled=1
gpgcheck=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7
EOF
        cat > $repodir/k8s.repo <<EOF
[kubernetes]
name = k8s
baseurl = http://$PKG_MIRROR/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled = 1
gpgcheck = 0
EOF
    else
        test -d /etc/yum.repos.d/bak || mkdir /etc/yum.repos.d/bak
        mv -f /etc/yum.repos.d/*.repo /etc/yum.repos.d/bak &> /dev/null || true
        cat > /etc/yum.repos.d/local.repo <<EOF
[kirk-poc]
name=local-yum
baseurl=$LOCALREPO_ADDR
enabled=1
gpgcheck=0
EOF
    fi
}

repos::setup::rhel::undo() {
    local repodir=/etc/yum.repos.d
    command::exec yum clean all
    rm -f $repodir/{epel,kirk,k8s,local}.repo
    command::exec rm -rf /var/cache/yum/*
    # command::exec yum makecache fast
}

kubelet::validate() {
    log I "validate kubelet installation"
    local kubelet_version=
    KUBELET_INSTALLED=false

    command::exists kubelet && {
        kubelet_version=$(kubelet --version | sed 's/-/_/' | awk -Fv '{print $2}')

        [[ $KUBELET_VERSION != $kubelet_version ]] && log E "a different version kubelet already installed on host, please uninstall it and try again" || KUBELET_INSTALLED=true
    } || true
}

kubelet::install() {
    log I "installing kubelet${VERSION_SYMBOL}${KUBELET_VERSION}"
    local kubeadm_version=$(echo $KUBELET_VERSION | awk -F'[_-]' '{print $1}')

    $KUBELET_INSTALLED || {
        command::exec "$PKG_MGR install ${INSTALL_OPTIONS} kubelet${VERSION_SYMBOL}${KUBELET_VERSION}*"

        $DEBUG && log D "installing kubectl${VERSION_SYMBOL}${kubeadm_version} and kubeadm${VERSION_SYMBOL}${kubeadm_version}"
        command::exec "$PKG_MGR install ${INSTALL_OPTIONS} kubectl${VERSION_SYMBOL}${kubeadm_version}* kubeadm${VERSION_SYMBOL}${kubeadm_version}*"
    }
}

kubelet::install::undo() {
    log I "uninstall kubelet kubectl kubeadm"
    if command::exists kubelet; then command::exec "$PKG_MGR autoremove ${INSTALL_OPTIONS} kubelet"; fi
    if command::exists kubectl; then command::exec "$PKG_MGR autoremove ${INSTALL_OPTIONS} kubectl"; fi
    if command::exists kubeadm; then command::exec "$PKG_MGR autoremove ${INSTALL_OPTIONS} kubeadm"; fi
}

kubelet::config() {
    log I "generate config for kubelet${VERSION_SYMBOL}${KUBELET_VERSION}"
    [[ -d /etc/systemd/system/kubelet.service.d/ ]] || mkdir /etc/systemd/system/kubelet.service.d/

    echo '[Service]
    Environment="KUBELET_CGROUP_DRIVER=--cgroup-driver=cgroupfs"
    Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
    Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
    Environment="KUBELET_NETWORK_ARGS=--network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin"
    Environment="KUBELET_DNS_ARGS=--cluster-dns='$CLUSTER_DNS' --cluster-domain=cluster.local"
    Environment="KUBELET_AUTHZ_ARGS=--authorization-mode=Webhook --client-ca-file=/etc/kubernetes/pki/ca.crt"
    #Environment="KUBELET_CADVISOR_ARGS=--cadvisor-port=0"
    Environment="KUBELET_CERTIFICATE_ARGS=--rotate-certificates=true --cert-dir=/var/lib/kubelet/pki"
    Environment="KUBELET_POD_INFRA_ARGS=--pod-infra-container-image='${IMAGE_REPOSITORY%*/}'/pause-amd64:3.0"
    Environment="KUBELET_FEATURE_GATES=--feature-gates=DevicePlugins=true,MountPropagation=true"
    Environment="KUBELET_LOG_LEVEL=-v=4"
    ExecStart=
    ExecStart=/usr/bin/kubelet $KUBELET_CGROUP_DRIVER $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_NETWORK_ARGS $KUBELET_DNS_ARGS $KUBELET_AUTHZ_ARGS $KUBELET_CADVISOR_ARGS $KUBELET_CERTIFICATE_ARGS $KUBELET_EXTRA_ARGS $KUBELET_POD_INFRA_ARGS $KUBELET_NODE_IP_ARGS $KUBELET_FEATURE_GATES $KUBELET_LOG_LEVEL $KUBELET_RESERVE_COMPUTE_RESOURCE_ARGS
    ' > /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
}

kubelet::config::undo() {
    log I "remove config of kubelet"
    command::exec rm -rf /etc/systemd/system/kubelet.service.d
    command::exec rm -rf /var/lib/kubelet/*
}

kubelet::run() {
    log I "run kubelet"
    command::exec systemctl daemon-reload
    command::exec systemctl enable -f kubelet
    command::exec systemctl restart kubelet
}

kubelet::run::undo() {
    log I "stop kubelet"
    if systemctl list-unit-files --type=service | grep -q kubelet; then
        command::exec systemctl stop kubelet
        command::exec systemctl disable kubelet
    fi
}

kubelet::setup() {
    kubelet::validate
    kubelet::install
    kubelet::config
    kubelet::run
}

kubelet::setup::undo() {
    kubelet::run::undo
    kubelet::config::undo
    kubelet::install::undo
}

init() {
    log I "init master"
    #$PKG_MGR remove -y cri-tools &> /dev/null || true
    command::exec "kubeadm init --config $INIT_CONFIG"
    init::postrun
}

init::postrun() {
    local protocol=
    [[ -f $ETCD_YAML ]] || err_exit "etcd yaml not found"
    # update etcd to listen at lan ip
    [[ -z $ETCD_IP ]] && log E "etcd ip is empty"

    grep -q 'listen-client-urls=http://' $ETCD_YAML && protocol=http || protocol=https
    grep -q "listen-client-urls=.*$ETCD_IP:2379" $ETCD_YAML || \
    sed -i "s#listen-client-urls=$protocol://.*#&,$protocol://$ETCD_IP:2379#" $ETCD_YAML
}

package::install() {
    log I "installing $EXTRAPKGS"
    local pkgs=$(echo $EXTRAPKGS | tr ',' ' ')
    echo $pkgs | grep -q ceph-common && pip uninstall -y urllib3 &> /dev/null || true

    command::exec "$PKG_MGR install $INSTALL_OPTIONS $pkgs"
}

device::mounted() {
    local device=$1

    if df | grep -q "^$device\s"; then
        return 0
    else
        return 1
    fi
}

path::mounted() {
    local path=$1

    if df | grep -q "\s$path$"; then
        return 0
    else
        return 1
    fi
}

device::is_lvm() {
    IS_LVM=$(lvdisplay | grep -w "$1" -c || true)
    if [[ ${IS_LVM} == 1 ]]; then
        return 0
    else
        return 1
    fi
}

device::mount::writefstab() {
    local uuid=
    local mntops=defaults,noatime
    local freq=0
    local passno=2
    local fstab=/etc/fstab

    uuid=$(blkid $device | cut -d ' ' -f 2 | tr -d '"')

    if ! grep -q "$uuid.*$mountPath" $fstab; then
        $DEBUG && log D "write fstab $uuid\t$mountPath\t$fs\t$mntops\t$freq\t$passno"
        echo -e "$uuid\t$mountPath\t$fs\t$mntops\t$freq\t$passno" >> $fstab
    fi
}

device::mount() {
    log I "mount $DEVICE_MOUNTS"
    local currentPath=
    local currentDevice=
    local currentFs=
    local device=
    local mountPath=
    local fs=

    for dmp in $DEVICE_MOUNTS; do
        read device mountPath fs<<<$(echo $dmp | tr ':' ' ')

        [[ -z $device ]] && log F "device to mount should be specified"
        [[ -z $mountPath ]] && log F "mount path should be specified"
        [[ -z $fs ]] && log F "disk file system type not specified"

        [[ -b $device ]] || log F "invalid block device: $device"
        [[ -d $mountPath ]] || {
            log I "$mountPath do not exist, creating ..."
            mkdir -p $mountPath
        }

        if device::mounted $device; then
            currentPath=$(df | grep -w "$device" | awk '{print $NF}' | xargs)
            if [[ $mountPath != $currentPath ]]; then
                log F "$device already mounted to $currentPath"
            fi
        elif path::mounted $mountPath; then

            if device::is_lvm $device; then
                currentMountPathMapper=$(df | grep -w "$mountPath" | awk '{print $1}' | xargs -I '{}' ls -l '{}' | awk '{print $11}' | xargs)
                currentDeviceMapper=$(ls -l "$device" |awk '{print $11}')
                if [[ $currentDeviceMapper != $currentMountPathMapper ]]; then
                    currentDevice=$(df | grep -w "$mountPath" | awk '{print $1}' | xargs)
                    log F "$currentDevice already mounted to $mountPath"
                fi
            else
                currentDevice=$(df | grep -w "$mountPath" | awk '{print $1}' | xargs)
                if [[ $currentDevice != $device ]]; then
                    log F "$currentDevice already mounted to $mountPath"
                fi
            fi
        else
            $DEBUG && log D "mount $device to $mountPath"
            currentFs=$(lsblk -nf $device | cut -d' ' -f 2)
            [[ $currentFs != $fs ]] && log F "filesystem type of $device is: $currentFs, which is not the same as configured: $fs"
            mount -t $fs $device $mountPath
            device::mount::writefstab
        fi
    done
}

usage() {
cat <<EOF
Usage:
    $0 setup repos [--local-repo-addr http://10.10.0.1:8880/localrepo --pkg-mirror mirrors.aliyun.com] [--debug]
    $0 setup kubelet --cluster-dns 169.169.0.10 --version 1.11.0 --image-repository index.qiniu.com/library [--debug]
    $0 install --extra-pkgs pkg1,pkg2 [--debug]
    $0 mount --mounts "/dev/sda1:/var/lib/etcd:xfs /dev/sdb2:/var/lib/kubelet:xfs" [--debug]
    $0 init --config kubeadm_config.yaml --etcd-ip 10.10.0.1 [--debug]
    $0 join --token 845e36.bc466480ab621387 --master 10.10.0.1:6443 [--control-plane] [--debug]
    $0 clean [--debug]
EOF
}

main() {
    [[ -r /etc/os-release ]] && LSB_DIST=$(. /etc/os-release && echo $ID)
    [[ -z $LSB_DIST ]] && log F "failed to detect linux distro"

    case $LSB_DIST in
    ubuntu)
        PKG_MGR=apt
        INSTALL_OPTIONS=' -y --allow-unauthenticated'
        VERSION_SYMBOL='='
        DIST_VERSION=$(. /etc/os-release && echo $UBUNTU_CODENAME)
    ;;
    centos)
        PKG_MGR=yum
        INSTALL_OPTIONS=' -y --setopt=obsoletes=0 --nogpgcheck'
        VERSION_SYMBOL='-'
        DIST_VERSION=$(. /etc/os-release && echo $VERSION_ID)
    ;;
    rhel)
        PKG_MGR=yum
        INSTALL_OPTIONS=' -y --setopt=obsoletes=0 --nogpgcheck'
        VERSION_SYMBOL='-'
        DIST_VERSION=$(. /etc/os-release && echo $VERSION_ID)
    *)
        log F "unrecognized Linux distro: $LSB_DIST, currently only support centos and ubuntu"
    ;;
    esac

    parse "$@"

    $DEBUG && log D "Linux distro: $LSB_DIST $DIST_VERSION detected"
    $ACTION
}

parse() {
    while [ $# -gt 0 ]; do
        case "$1" in
            setup)
                ACTION=setup
            ;;
            init)
                ACTION=init
            ;;
            join)
                ACTION=join
            ;;
            mount)
                ACTION=mount
            ;;
            clean)
                ACTION=clean
            ;;
            install)
                ACTION=install
            ;;
            repos)
                COMPONENT=repos
            ;;
            kubelet)
                COMPONENT=kubelet
            ;;
            --cluster-dns)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    CLUSTER_DNS="$2"
                    shift
                } || {
                    usage_exit "no cluster dns ip given for --cluster-dns"
                }
            ;;
            --master)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    MASTER="$2"
                    shift
                } || {
                    usage_exit "no master address(IP:PORT) given for --master"
                }
            ;;
            --storage-driver)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    DOCKER_STORAGE_DRIVER="$2"
                    shift
                } || {
                    usage_exit "no docker storage driver given for --storage-driver"
                }
            ;;
            --token)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    TOKEN="$2"
                    shift
                } || {
                    usage_exit "no token given for --token"
                }
            ;;
            --version)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    VERSION="$2"
                    shift
                } || {
                    usage_exit "no version given for --version"
                }
            ;;
            --device-mapper-device)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    DOCKER_DEVICE_MAPPER_DEVICE="$2"
                    shift
                } || {
                    usage_exit "no docker devicemapper device given for --device-mapper-device"
                }
            ;;
            --image-repository)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    IMAGE_REPOSITORY="$2"
                    shift
                } || {
                    usage_exit "no docker image repository given for --image-repository"
                }
            ;;
            --config)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    INIT_CONFIG="$2"
                    shift
                } || {
                    usage_exit "no kubeadm init config given for --config"
                }
            ;;
            --local-repo-addr)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    LOCALREPO_ADDR="$2"
                    shift
                } || {
                    usage_exit "no local package repository addr given for --local-repo-addr"
                }
            ;;
            --insecure-registries)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    INSECURE_REGS="$2"
                    shift
                } || {
                    usage_exit "no insecure docker registry given for --insecure-registries"
                }
            ;;
            --extra-pkgs)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    EXTRAPKGS="$2"
                    shift
                } || {
                    usage_exit "no extra packages given for --extra-pkgs"
                }
            ;;
            --pkg-mirror)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    PKG_MIRROR="$2"
                    shift
                } || {
                    usage_exit "no package mirror given for --pkg-mirror"
                }
            ;;
            --mounts)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    DEVICE_MOUNTS="$2"
                    shift
                } || {
                    usage_exit "no mount info ip given for --mounts"
                }
            ;;
            --control-plane)
                JOIN_CONTROL_PLANE=--experimental-control-plane
            ;;
            --debug)
                DEBUG=true
            ;;
            -h|--help)
                usage
                exit 0
            ;;
            *)
                usage_exit "invalid option: $1"
            ;;
        esac
        shift $(( $# > 0 ? 1 : 0 ))
    done

    case "$ACTION" in
        setup)
            case "$COMPONENT" in
                repos)
                    ACTION=repos::setup
                ;;
                kubelet)
                    ACTION=kubelet::setup
                    KUBELET_VERSION=$VERSION
                ;;
                *)
                    usage_exit "invalid component"
                ;;
            esac
        ;;
        install)
            ACTION=package::install
        ;;
        mount)
            ACTION=device::mount
        ;;
        init|join|clean)
        ;;
        *)
            usage_exit "invalid action: $ACTION"
        ;;
    esac
}

# main
main "$@"
