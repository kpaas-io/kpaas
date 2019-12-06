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

# This script is aim to set hostaliases

test -f ~/.bashrc || touch ~/.bashrc
test -f ~/.bash_aliases || touch ~/.bash_aliases

grep "alias" ~/.bash_aliases -q || {
    cat >~/.bash_aliases<<EOF
alias k='kubectl'
alias kc='kubectl -n default'
alias kcks='kubectl -n kube-system'
alias kce='kubectl -n qce'
EOF
}

grep "bash_aliases" ~/.bashrc || {
    echo `
if [ -f ~/.bash_aliases ]; then
    . ~/.bash_aliases
fi` >> ~/.bashrc
}

source ~/.bash_aliases
source ~/.bashrc
