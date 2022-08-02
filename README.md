# Stream Unlock
[![](https://img.shields.io/badge/Telegram-Group-blue?style=flat-square)](https://t.me/aioCloud)
[![](https://img.shields.io/badge/Telegram-Channel-green?style=flat-square)](https://t.me/aioCloud_channel) 
[![](https://img.shields.io/github/downloads/aiocloud/stream/total.svg?style=flat-square)](https://github.com/aiocloud/stream/releases)
[![](https://img.shields.io/github/v/release/aiocloud/stream?style=flat-square)](https://github.com/aiocloud/stream/releases)

流媒体解锁后端

## 推荐系统
- Debian 11
- Debian 10
- Ubuntu 20.04
- CentOS 8 Stream

## 一键部署
```bash
curl -fsSL https://raw.githubusercontent.com/loadinghtml/stream/loadinghtml-stream/scripts/kickstart.sh | bash
```

## 更新
```bash
curl -fsSL https://raw.githubusercontent.com/loadinghtml/stream/loadinghtml-stream/scripts/upgrade.sh | bash
```

## 卸载
```bash
curl -fsSL https://raw.githubusercontent.com/loadinghtml/stream/loadinghtml-stream/scripts/remove.sh | bash
```

## 配置文件
Stream
```
/etc/stream.json
```

## 控制命令

##启动
```
systemctl start stream
``` 
##停止
```
systemctl stop stream
```
##重启
```
systemctl restart stream
``` 
##状态
```
systemctl status stream
``` 
##具体日志
```
journalctl -f -u stream
``` 
