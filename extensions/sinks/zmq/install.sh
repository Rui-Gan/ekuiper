#!/bin/sh
#
# Copyright 2023 EMQ Technologies Co., Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set +e -x -u

DISTRO='unknow'

Get_Dist_Name()
{
    if grep -Eqii "CentOS" /etc/issue || grep -Eq "CentOS" /etc/*-release; then
        DISTRO='CentOS'
    elif grep -Eqi "Red Hat Enterprise Linux Server" /etc/issue || grep -Eq "Red Hat Enterprise Linux Server" /etc/*-release; then
        DISTRO='RHEL'
    elif grep -Eqi "Aliyun" /etc/issue || grep -Eq "Aliyun" /etc/*-release; then
        DISTRO='Aliyun'
    elif grep -Eqi "Fedora" /etc/issue || grep -Eq "Fedora" /etc/*-release; then
        DISTRO='Fedora'
    elif grep -Eqi "Debian" /etc/issue || grep -Eq "Debian" /etc/*-release; then
        DISTRO='Debian'
    elif grep -Eqi "Ubuntu" /etc/issue || grep -Eq "Ubuntu" /etc/*-release; then
        DISTRO='Ubuntu'
    elif grep -Eqi "Raspbian" /etc/issue || grep -Eq "Raspbian" /etc/*-release; then
        DISTRO='Raspbian'
    elif grep -Eqi "Alpine" /etc/issue || grep -Eq "Alpine" /etc/*-release; then
        DISTRO='Alpine'
    else
        DISTRO='unknow'
    fi
    echo $DISTRO;
}


Get_Dist_Name

case $DISTRO in \
    Debian|Ubuntu|Raspbian ) \
	apt update \
	&& apt upgrade \
        && apt install -y libczmq-dev findutils 2> /dev/null \
    ;; \
    Alpine ) \
        apk add libzmq findutils \
    ;; \
    *) \
        yum install -y zeromq 2> /dev/null \
    ;; \
esac

LIBZMQ_PATH=$(find / -name "libzmq.so.5" -print -quit 2>/dev/null)
if [ -n "$LIBZMQ_PATH" ]; then
  echo "找到 libzmq.so.5 文件，路径为：$LIBZMQ_PATH"
  echo "$LIBZMQ_PATH" > /etc/ld.so.conf
  echo "已将路径添加到 /etc/ld.so.conf 文件中"
else
  echo "未找到 libzmq.so.5 文件"
  case $DISTRO in \
      Debian|Ubuntu|Raspbian ) \
  	apt-get update \
          && apt-get install -y libzmq5 2> /dev/null \
          && apt-get clean \
          && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
      ;; \
      Alpine ) \
          apk add zeromq-dev \
      ;; \
      *) \
          yum install -y zeromq 2> /dev/null \
      ;; \
  esac
fi
    
echo "install success";
