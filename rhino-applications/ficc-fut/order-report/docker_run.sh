#!/bin/bash

docker rm -vf order-report-ficc-fut; docker run --restart=always --name=order-report-ficc-fut -it -d \
  -v /opt/rhino/ficc-fut/order-report/log:/opt/log \
  -v /opt/rhino/ficc-fut/trade-engine/config/config.json:/config.json \
  -v /opt/rhino/ficc-fut/trade-engine/config/:/config/ \
  -v /opt/rhino/ficc-fut/data/config/:/config_data/ \
  -e ServerPort=9095 \
  -p 9095:9095 \
  dockertest.gf.com.cn/rhino/order-report-ficc-fut:b20251010-p1.0.3.8-latest


docker rm -vf order-report-ficc-fut; docker run --restart=always --name=order-report-ficc-fut -it -d \
  -v /opt/rhino/ficc-fut/order-report/log:/opt/log \
  -v /opt/rhino/ficc-fut/trade-engine/config/config.json:/config.json \
  -v /opt/rhino/ficc-fut/data/config/:/config_data/ \
  -e ServerPort=9093 \
  -p 9093:9093 \
  docker2.gf.com.cn/rhino/order-report-ficc-fut:b20251010-p1.0.0.0
