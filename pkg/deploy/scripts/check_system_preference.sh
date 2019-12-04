#!/bin/bash
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

function check_sysctl {
	sysctl_key=$1
	expected_value=$2
	value=$(sysctl ${sysctl_key} | awk '{print $3}')
	if [[ ${value} != ${expected_value} ]]; then
		echo "check failed: ${sysctl_key} is ${value}, expected ${expected_value}"
		if [[ ${fix_sysctl} = "true" ]]; then
			echo "try to set ${sysctl_key} to ${expected_value}"
			sysctl -w ${sysctl_key}=${expected_value}
			# check the value again.
			value=$(sysctl ${sysctl_key} | awk '{print $3}')
			if [[ ${value} != ${expected_value} ]]; then
				echo "failed to set ${sysctl_key} to ${expected_value}"
				return 1
			else
				echo "successfully set ${sysctl_key} to ${expected_value}"
			fi
		else
			return 1
		fi
	fi
	return 0
}


err=""
newline=$'\n'

# check net.ipv4.ip_forward
echo "check net.ipv4.ip_forward..."
check_sysctl "net.ipv4.ip_forward" 1
ip_forward_ok=$?
if [ ${ip_forward_ok} -ne 0 ]; then
	echo "net.ipv4.ip_forward is not opened"
	err="${err}${newline}ip_forward not opened"
fi

# check net.ipv4.conf.all.forwarding
echo "check net.ipv4.conf.all.forwarding..."
check_sysctl "net.ipv4.conf.all.forwarding" 1
all_forwarding_ok=$?
if [ ${all_forwarding_ok} -ne 0 ]; then
	echo "net.ipv4.conf.all.forwarding not opened"
	err="${err}${newline}net.ipv4.conf.all.forwarding not opened"
fi

# check proxy_arp and forwarding of cali* interfaces
cali_interfaces=$(ifconfig | grep ^cali | awk '{print $1}')
for interface in ${cali_interfaces}; do
	echo "check configuration of interface ${interface}..."
	check_sysctl "net.ipv4.conf.${interface}.proxy_arp" 1
	proxy_arp_ok=$?
	if [ ${proxy_arp_ok} -ne 0 ]; then
		echo "error: proxy_arp of interface ${interface} is not opened"
		err="${err}${newline}proxy_arp of interface ${interface} is not opened"
	fi
	check_sysctl "net.ipv4.conf.${interface}.forwarding" 1
	forwarding_ok=$?
	if [ ${forwarding_ok} -ne 0 ] ;then
		echo "error: forwarding of interface ${interface} is not opened"
		err="${err}${newline}forwarding of interface ${interface} is not opened"
	fi
done

if [ "${err}" == "" ]; then
	echo "sysctl check OK"
	exit 0
else
	echo "sysctl errors:${err}"
	exit 1
fi