# API Docs

## Table of Contents
* [Action](#action)
* [APIServer](#apiserver)
* [CheckMachines](#checkmachines)
* [CommandAction](#commandaction)
* [CommonNode](#commonnode)
* [ControlPlaneComponent](#controlplanecomponent)
* [ControlPlaneConfig](#controlplaneconfig)
* [DeployEtcd](#deployetcd)
* [DeployMaster](#deploymaster)
* [DeployNode](#deplounode)
* [DeployIngress](#deployingress)
* [DeviceMountPoints](#devicemountpoints)
* [Etcd](#etcd)
* [FetchLog](#fetchlog)
* [HostPathMount](#hostpathmount)
* [InitMachines](#initmachines)
* [KubeletConfig](#kubeletconfig)
* [LocalEtcd](#localetcd)
* [Machine](#machine)
* [MachineConfig](#machineconfig)
* [MachineGroup](#machinegroup)
* [MachineGroupList](#machinegrouplist)
* [Registry](#registry)
* [ScriptAction](#scriptaction)
* [Tasks](#tasks)
* [TaskManager](#taskmanager)

## Action

Action defines the exact action which is executed on target machines

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| script |  | *[ScriptAction](#scriptaction) | false |
| command |  | *[CommandAction](#commandaction) | false |

[Back to TOC](#table-of-contents)

## APIServer

APIServer holds settings necessary for APIServer deployments in the clusters

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| certSANs | CertSANs sets extra Subject Alternative Names for the API Server signing cert. | []string | false |
| timeoutForControlPlane | TimeoutForControlPlane controls the timeout that we use for API server to appear | *metav1.Duration | false |

[Back to TOC](#table-of-contents)

## CheckMachines

Checkmachines defines methods of checking machine groups

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| targetMachine | TargetMachine defines a bunch of entities including all roles in machine group list | *[machineGroupList](#machinegrouplist) | false |
| action | Action ensure the real command needs to be exec on nodes | [][Action](#action) | false |

[Back to TOC](#table-of-contents)

## CommandAction

CommandAction defines the action in the form of pure shell command

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| cmd | Cmd defines the raw command that executed on target machine | string | true |

[Back to TOC](#table-of-contents)

## CommonNode

CommonNode defines node excepts master and etcd nodes role

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| extraArgs | ExtraArgs is an extra set of flags to pass to the common node. | map[string]string | false |
| discoveryToken | DiscoveryToken is token for join node to control plane node. | string | false |
| discoveryFile | DiscoveryFile is file for join node to control plane node. | string | false |

[Back to TOC](#table-of-contents)

## ControlPlaneComponent

ControlPlaneComponent holds settings common to control plane component of the cluster

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| extraArgs | ExtraArgs is an extra set of flags to pass to the control plane component. use ComponentConfig + ConfigMaps. | map[string]string | false |
| extraVolume | ExtraVolume is an extra set of host volumes, mounted to the control plane component. | *[HostPathMount](#hostpathmount) | false |

[Back to TOC](#table-of-contents)

## ControlPlaneConfig

ControlPlaneConfig is control plane components config of a node whose role is master

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| bindInterface | BindInterface is network interface name which control plane components bind to | string | false |

[Back to TOC](#table-of-contents)

## DeployEtcd

DeployEtcd is mainly for deploying ETCD component on specific node

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| targetMachine | TargetMachine is the list of machine group which assigned specific roles | *[machineGroupList] | false |
| etcdInfo | EtcdInfo gathers all etcd setting for deploying | *[Etcd](#Etcd) | true |

[Back to TOC](#table-of-contents)

## DeployMaster

DeployMaster is mainly for deploying control plane on specific node

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| targetMachine | TargetMachine is the list of machine group which assigned specific roles | *[machineGroupList] | false |
| controlPlaneComponent | ControlPlaneComponent holds settings common to control plane component of the cluster | *[controlPlaneComponent](#controlplanecomponent) | true |

[Back to TOC](#table-of-contents)

## DeployNode

DeployNode is mainly for deploying common role on specific node

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| targetMachine | TargetMachine is the list of machine group which assigned specific roles | *[machineGroupList] | false |
| commonNode | CommonNode defines node excepts master and etcd nodes role | *[commandNode] | true |

[Back to TOC](#table-of-contents)

## DeployIngress

DeployIngress is mainly for deploying ingress roles on specific node

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| targetMachine | TargetMachine is the list of machine group which assigned specific roles | *[machineGroupList] | false |

[Back to TOC](#table-of-contents)

## DeviceMountPoints

DeviceMountPoints is a mapping from node's block device to mount point

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| devicePath | DevicePath is the path of a device like \"/dev/vdb1\" | string | true |
| mountPoint | MountPoint is the mount point of a specific device like \"/var/lib/etcd\", if MountPoint == \"\", which means device must not be mounted | string | true |
| fsType | FSType is the filesystem type of the device like \"xfs\" | string | false |

[Back to TOC](#table-of-contents)

## Etcd

Etcd contains elements describing Etcd configuration.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| local | Local provides configuration knobs for configuring the local etcd | *[LocalEtcd](#localetcd) | false |

[Back to TOC](#table-of-contents)

## FetchLog

FetchLog returns log of one actions

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| logPath | LogPath dedicated the path on the host which stores log files | string | false |
| readOnly | ReadOnly controls write access to the volume | bool | false |

[Back to TOC](#table-of-contents)

## HostPathMount

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Name of the volume inside pod template | string | true |
| hostPath | HostPath is the actual path exists on the host, will be mounted inside the pod | string | true |
| mountPath | MountPath is the actual path exists on the pod, which the mountpoint hostpath been mounted | string | true |
| readOnly | ReadOnly controls write access to the volume | bool | false |
| pathType | PathType is the type of the HostPath | v1.HostPathType | false |

[Back to TOC](#table-of-contents)

## KubeletConfig

KubeletConfig is kubelet config of a node

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| version | Version is the version of kubelet if version is nil, it will use cluster kubernetes version by default | string | false |
| bindInterface | BindInterface is the network interface name which kubelet bind to | string | false |

[Back to TOC](#table-of-contents)

## LocalEtcd

LocalEtcd describes that we should run an etcd cluster locally

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| dataDir | DataDir is the directory etcd will place its data. Defaults to \"/var/lib/etcd\". | string | true |
| extraArgs | ExtraArgs are extra arguments provided to the etcd binary when run inside a static pod. | map[string]string | false |
| serverCertSANs | ServerCertSANs sets extra Subject Alternative Names for the etcd server signing cert. | []string | false |
| peerCertSANs | PeerCertSANs sets extra Subject Alternative Names for the etcd peer signing cert. | []string | false |

[Back to TOC](#table-of-contents)

## Machine

Machine is a instance of a node/role

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Name is machine's hostname | string | true |
| ip | IP is machine's (internal) IP address | string | true |

[Back to TOC](#table-of-contents)

## MachineConfig

MachineConfig is the config of all machines in same machine group

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| labels | Labels is node labels while joining to cluster | map[string]string | false |
| taints | Taints is node taints while joining to cluster | []v1.Taint | false |
| deviceMountPoints | DeviceMountPoints provide a path from node's block device to mount point | [][DeviceMountPoints](#devicemountpoints) | false |
| controlPlaneConfig | ControlPlaneConfig is control plane components config of a specific node whose role is master | *[ControlPlaneConfig](#controlplaneconfig) | false |
| kubeletConfig | KubeletConfig is kubelet config of a specific node | *[KubeletConfig](#kubeletconfig) | false |

[Back to TOC](#table-of-contents)

## MachineGroup

MachineGroup contains bunch of machines with role

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Name is machinegroup's name | string | true |
| role | Role defines a set of machines | string | true |
| machines | Machines is a set of machines with common flags | [][Machine](#machine) | true |

[Back to TOC](#table-of-contents)

## MachineGroupList

MachineGroupList contains a list of machine group

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| items | MachineGroupList define a list of machine group | [][MachineGroup](#machinegroup) | true |

[Back to TOC](#table-of-contents)

## ScriptAction

ScriptAction defines the action in the form of script files

| Field | Description | Scheme | Required |
| ----- | ----------- | ----- | -------- |
| path | Path defines the location of script file | string | true |

[Back to TOC](#table-of-contents)

## Registry

Registry is a private docker registry for deploying cluster

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| version |  | string | true |

[Back to TOC](#table-of-contents)

## Tasks

Tasks contain a set of actions which would be run on the target machines

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| action | Action defines the exact action which is executed on target machines | *[Action](#action) | false |
| status | Status shows task's procedure eg.pending, complete, stopped. | *[TaskStatus](#taskstatus) | false |

[Back to TOC](#table-of-contents)

## TaskManager

TaskManager manage all tasks which would be run at nodes

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| taskID | TaskID is signed as unique for each unique actions executed on machine groups | int32 | false |
| task | task gives explanation of specific set of tasks and defination | [Task](#task) | false |
| targetMachine | TargetMachine is the list of machine group which assigned specific roles | *[machineGroupList] | false |

[Back to TOC](#table-of-contents)
