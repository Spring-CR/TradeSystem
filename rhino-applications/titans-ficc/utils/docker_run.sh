#!/bin/bash

docker rm -vf utils-titans-ficc; docker run --restart=always --name=utils-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/utils/config/config.json:/config.json \
  -v /opt/rhino/titans-ficc/utils/log:/opt/log \
  -p 5093:5093 \
  dockertest.gf.com.cn/rhino/utils-titans-ficc:b20251010-p1.0.5.6-latest

docker rm -vf utils-titans-ficc; docker run --restart=always --name=utils-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/utils/config/config.json:/config.json \
  -p 5093:5093 \
  docker2.gf.com.cn/rhino/utils-titans-ficc:b20251010-p1.0.2.5