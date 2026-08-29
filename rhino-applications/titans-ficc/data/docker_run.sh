#!/bin/bash

docker rm -vf data-titans-ficc; docker run --restart=always --name=data-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/data/log:/opt/log \
  -v /opt/rhino/titans-ficc/data/config/:/config/ \
  -e ServerPort=6093 \
  -p 6093:6093 \
  dockertest.gf.com.cn/rhino/data-titans-ficc:b20251010-p1.0.3.0-latest

docker rm -vf data-titans-ficc; docker run --restart=always --name=data-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/data/log:/opt/log \
  -v /opt/rhino/titans-ficc/data/config/:/config/ \
  -e ServerPort=6093 \
  -p 6093:6093 \
  docker2.gf.com.cn/rhino/data-titans-ficc:b20251010-p1.0.2.1


--修复dns的问题
修改 Docker 的守护进程配置，设置默认的 DNS 服务器。编辑 /etc/docker/daemon.json 文件（如果不存在则创建），添加以下内容：
json
{
  "dns": ["10.51.167.99"]
}
然后重启 Docker 服务：sudo systemctl restart docker。这样以后启动的容器都会使用这个 DNS。