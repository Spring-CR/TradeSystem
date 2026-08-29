#!/bin/bash

docker rm -vf data-sync-ficc-fut; docker run --restart=always --name=data-sync-ficc-fut -it -d \
  -v /opt/rhino/ficc-fut/data-sync/log:/opt/log \
  -v /opt/rhino/ficc-fut/data-sync/config/:/config/ \
  -e ServerPort=7095 \
  -p 7095:7095 \
  dockertest.gf.com.cn/rhino/data-sync-ficc-fut:b20251010-p1.0.1.8-latest


docker rm -vf data-sync-ficc-fut; docker run --restart=always --name=data-sync-ficc-fut -it -d \
  -v /opt/rhino/ficc-fut/data-sync/log:/opt/log \
  -v /opt/rhino/ficc-fut/data-sync/config/:/config/ \
  -e ServerPort=7093 \
  -p 7093:7093 \
  docker2.gf.com.cn/rhino/data-sync-ficc-fut:b20251010-p1.0.6.2

--修复dns的问题
修改 Docker 的守护进程配置，设置默认的 DNS 服务器。编辑 /etc/docker/daemon.json 文件（如果不存在则创建），添加以下内容：
json
{
  "dns": ["10.51.167.99"]
}
然后重启 Docker 服务：sudo systemctl restart docker。这样以后启动的容器都会使用这个 DNS。