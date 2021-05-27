#!/usr/bin/env bash
if [[ -f /etc/redhat-release ]]; then
    release="centos"
elif cat /etc/issue | grep -q -E -i "debian"; then
    release="debian"
elif cat /etc/issue | grep -q -E -i "ubuntu"; then
    release="ubuntu"
elif cat /etc/issue | grep -q -E -i "centos|red hat|redhat"; then
    release="centos"
elif cat /proc/version | grep -q -E -i "raspbian|debian"; then
    release="debian"
elif cat /proc/version | grep -q -E -i "ubuntu"; then
    release="ubuntu"
elif cat /proc/version | grep -q -E -i "centos|red hat|redhat"; then
    release="centos"
else
    exit 1
fi

if [[ "$release" == "centos" ]]; then
    yum install openssl-devel pkgconfig make gcc git -y || exit $?
else
    apt update || exit $?
    apt install build-essential pkg-config libssl-dev make git -y || exit $?
fi

cd /opt && rm -fr smartdns
git clone https://github.com/pymumu/smartdns --depth 1 || exit $?

cd smartdns && make -j$(nproc) || exit $?
cd src && cp -f smartdns /usr/bin/smartdns

rm -fr /etc/smartdns && mkdir /etc/smartdns
wget -O /etc/smartdns/smartdns.conf          https://raw.githubusercontent.com/aiocloud/stream/master/smartdns/smartdns.conf    || exit $?
wget -O /etc/systemd/system/smartdns.service https://raw.githubusercontent.com/aiocloud/stream/master/smartdns/smartdns.service || exit $?

cd /opt && rm -fr smartdns
exit 0
