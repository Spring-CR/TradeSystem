#!/bin/bash

docker run --rm --restart=always --name=trade-engine-titans-cbf -it -d \
  -v /opt/rhino/titans-cbf/trade-engine/log:/opt/log \
  -v /opt/rhino/titans-cbf/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/titans-cbf/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/titans-cbf/trade-engine/config/config.json:/config.json \
  -p 8092:8092 \
  dockertest.gf.com.cn/rhino/trade-engine-titans-cbf:b20250718-p1.0.0.7-latest

docker run --rm --restart=always --name=trade-engine-titans-cbf -it -d \
  -v /opt/rhino/titans-cbf/trade-engine/log:/opt/log \
  -v /opt/rhino/titans-cbf/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/titans-cbf/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/titans-cbf/trade-engine/config/config.json:/config.json \
  -p 8092:8092 \
  --dns=10.51.167.99 \
  dockertest.gf.com.cn/rhino/trade-engine-titans-cbf:b20250718-p1.0.0.7-latest

