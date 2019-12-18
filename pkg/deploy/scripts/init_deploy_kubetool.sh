#! /usr/bin/env bash

set -Eeuo pipefail

DEBUG=false
LSB_DIST=
DIST_VERSION=
ACTION=
COMPONENT=
VERSION=
IMAGE_REPOSITORY=index-dev.qiniu.io/kelibrary
DEVICE_MOUNTS=

# docker specific
DOCKER_INSTALLED=
DOCKER_THINPOOL_STATE=
DOCKER_PKG=
DOCKER_VERSION=
DOCKER_DEVICE_MAPPER_DEVICE=
DOCKER_DEVICE_MAPPER_PKGS=
DOCKER_STORAGE_DRIVER=
INSECURE_REGS=

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

# etcd specific
ETCD_IP=
APISERVERYAML=/etc/kubernetes/manifests/kube-apiserver.yaml
ETCD_YAML=/etc/kubernetes/manifests/etcd.yaml

## source lib.sh
. lib.sh

command::exists() {
    command -v "$@" > /dev/null 2>&1
}

command::exec() {
    $DEBUG && {
        log::deploy D $@
        eval $@
        return
    }

    #log::deploy I $@
    err=$( { eval $@ > /dev/null; } 2>&1) || log::deploy E "exec $@ failed, err: $err"
}

usage_exit() {
    log::deploy I "$@" >&2
    usage
    exit 1
}

log::deploy() {
    local level=$1
    local step=$2
    local logts=$(date '+%m%d %T')
    local hostname=$(hostname)

    test -z "$level" && error_exit "F$LOGTS [$hostname] log::deploy level not specified"

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
            error_exit "F$logts [$hostname] unknown log::deploy level $level"
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

repos::setup() {
    log::deploy I "setup package repos and update cache"
    repos::setup::${LSB_DIST}
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
    cat > $sourcedir/ceph.list <<EOF
deb http://$PKG_MIRROR/ceph/debian-luminous ${DIST_VERSION} main
EOF
    command::exec apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 6A030B21BA07F4FB E84AC2C0460F3994 7EA0A9C3F273FCD8 F76221572C52609D
    command::exec apt clean
    command::exec apt update
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
        cat > $repodir/ceph.repo <<EOF
[ceph]
name = ceph-luminous
baseurl = http://$PKG_MIRROR/ceph/rpm-luminous/el7/\$basearch
enabled = 1
gpgcheck = 1
gpgkey = https://download.ceph.com/keys/release.asc
EOF
    else
        test -d /etc/yum.repos.d/bak || mkdir /etc/yum.repos.d/bak
        mv -f /etc/yum.repos.d/*.repo /etc/yum.repos.d/bak &> /dev/null || true
        cat > /etc/yum.repos.d/local.repo <<EOF
[kpaas-deploy]
name=local-yum
baseurl=$LOCALREPO_ADDR
enabled=1
gpgcheck=0
EOF
    fi
}

repos::setup::rhel() {
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
    cat > $sourcedir/ceph.list <<EOF
deb http://$PKG_MIRROR/ceph/debian-luminous ${DIST_VERSION} main
EOF
    command::exec apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 6A030B21BA07F4FB E84AC2C0460F3994 7EA0A9C3F273FCD8 F76221572C52609D
    command::exec apt clean
    command::exec apt update
}

kubelet::validate() {
    log::deploy I "validate kubelet installation"
    local kubelet_version=
    KUBELET_INSTALLED=false

    command::exists kubelet && {
        kubelet_version=$(kubelet --version | sed 's/-/_/' | awk -Fv '{print $2}')

        [[ $KUBELET_VERSION != $kubelet_version ]] && log::deploy E "a different version kubelet already installed on host, please uninstall it and try again" || KUBELET_INSTALLED=true
    } || true
}

kubelet::install() {
    log::deploy I "installing kubelet${VERSION_SYMBOL}${KUBELET_VERSION}"
    local kubeadm_version=$(echo $KUBELET_VERSION | awk -F'[_-]' '{print $1}')

    $KUBELET_INSTALLED || {
        command::exec "$PKG_MGR install ${INSTALL_OPTIONS} kubelet${VERSION_SYMBOL}${KUBELET_VERSION}*"

        $DEBUG && log::deploy D "installing kubectl${VERSION_SYMBOL}${kubeadm_version} and kubeadm${VERSION_SYMBOL}${kubeadm_version}"
        command::exec "$PKG_MGR install ${INSTALL_OPTIONS} kubectl${VERSION_SYMBOL}${kubeadm_version}* kubeadm${VERSION_SYMBOL}${kubeadm_version}*"
    }
}

kubelet::config() {
    log::deploy I "generate config for kubelet${VERSION_SYMBOL}${KUBELET_VERSION}"
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

kubelet::run() {
    log::deploy I "run kubelet"
    command::exec systemctl daemon-reload
    command::exec systemctl enable -f kubelet
    command::exec systemctl restart kubelet
}

kubelet::setup() {
    kubelet::validate
    kubelet::install
    kubelet::config
    kubelet::run
}

join() {
    log::deploy I "join node to cluster"
    local skip_ca=

    if kubeadm join --help | grep -q '^\s*--discovery-token-unsafe-skip-ca-verification'
    then
        export skip_ca=--discovery-token-unsafe-skip-ca-verification
    fi

    #$PKG_MGR remove -y cri-tools &> /dev/null || true

    #kubeadm join --token $TOKEN $MASTERIP --discovery-token-unsafe-skip-ca-verification [--experimental-control-plane]
    command::exec kubeadm join --token $TOKEN $MASTER $skip_ca $JOIN_CONTROL_PLANE
}

usage() {
cat <<EOF
Usage:
    $0 setup repos [--local-repo-addr http://10.10.0.1:8880/localrepo --pkg-mirror mirrors.aliyun.com] [--debug]
    $0 setup kubelet --cluster-dns 169.169.0.10 --version 1.11.0 --image-repository index.qiniu.com/library [--debug]
    $0 join --token 845e36.bc466480ab621387 --master 10.10.0.1:6443 [--control-plane] [--debug]
    $0 clean [--debug]
EOF
}

main() {;./
    [[ -r /etc/os-release ]] && LSB_DIST=$(. /etc/os-release && echo $ID)
    [[ -z $LSB_DIST ]] && log::deploy F "failed to detect linux distro"

    case $LSB_DIST in
    ubuntu)
        PKG_MGR=apt
        INSTALL_OPTIONS=' -y --allow-unauthenticated'
        VERSION_SYMBOL='='
        DIST_VERSION=$(. /etc/os-release && echo $UBUNTU_CODENAME)
        DOCKER_DEVICE_MAPPER_PKGS="thin-provisioning-tools lvm2"
    ;;
    centos)
        PKG_MGR=yum
        INSTALL_OPTIONS=' -y --setopt=obsoletes=0 --nogpgcheck'
        VERSION_SYMBOL='-'
        DIST_VERSION=$(. /etc/os-release && echo $VERSION_ID)
        DOCKER_DEVICE_MAPPER_PKGS="device-mapper-persistent-data lvm2"
    ;;
    rhel)
        PKG_MGR=yum
        INSTALL_OPTIONS=' -y --setopt=obsoletes=0 --nogpgcheck'
        VERSION_SYMBOL='-'
        DIST_VERSION=$(. /etc/os-release && echo $VERSION_ID)
        DOCKER_DEVICE_MAPPER_PKGS="device-mapper-persistent-data lvm2"

    ;;
    *)
        log::deploy F "unrecognized Linux distro: $LSB_DIST, currently only support centos and ubuntu"
    ;;
    esac

    parse "$@"

    $DEBUG && log::deploy D "Linux distro: $LSB_DIST $DIST_VERSION detected"
    $ACTION
}

parse() {
    while [ $# -gt 0 ]; do
        case "$1" in
            setup)
                ACTION=setup
            ;;
            join)
                ACTION=join
            ;;
            clean)
                ACTION=clean
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
            --etcd-ip)
                [[ -n ${2+x} ]] && ! echo $2 | grep -q ^- && {
                    ETCD_IP="$2"
                    shift
                } || {
                    usage_exit "no etcd ip given for --etcd-ip"
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
                docker)
                    ACTION=docker::setup
                    DOCKER_VERSION=$VERSION
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
        init|join|clean)
        ;;
        *)
            usage_exit "invalid action: $ACTION"
        ;;
    esac
}

# main
main "$@"
