#!/bin/bash

docker rm -vf data-ficc-fut; docker run --restart=always --name=data-ficc-fut -it -d \
  -v /opt/rhino/ficc-fut/data/log:/opt/log \
  -v /opt/rhino/ficc-fut/data/config/:/config/ \
  -e ServerPort=6095 \
  -p 6095:6095 \
  dockertest.gf.com.cn/rhino/data-ficc-fut:b20251010-p1.0.0.6-latest

docker rm -vf data-ficc-fut; docker run --restart=always --name=data-ficc-fut -it -d \
  -v /opt/rhino/ficc-fut/data/log:/opt/log \
  -v /opt/rhino/ficc-fut/data/config/:/config/ \
  -e ServerPort=6095 \
  -p 6095:6095 \
  docker2.gf.com.cn/rhino/data-ficc-fut:b20251010-p1.0.2.1


--修复dns的问题
修改 Docker 的守护进程配置，设置默认的 DNS 服务器。编辑 /etc/docker/daemon.json 文件（如果不存在则创建），添加以下内容：
json
{
  "dns": ["10.51.167.99"]
}
然后重启 Docker 服务：sudo systemctl restart docker。这样以后启动的容器都会使用这个 DNS。