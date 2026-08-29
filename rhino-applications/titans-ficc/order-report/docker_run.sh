#!/bin/bash

docker rm -vf order-report-titans-ficc; docker run --restart=always --name=order-report-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/order-report/log:/opt/log \
  -v /opt/rhino/titans-ficc/trade-engine/config/config.json:/config.json \
  -v /opt/rhino/titans-ficc/trade-engine/config/:/config/ \
  -v /opt/rhino/titans-ficc/data/config/:/config_data/ \
  -e ServerPort=9093 \
  -p 9093:9093 \
  dockertest.gf.com.cn/rhino/order-report-titans-ficc:b20251010-p1.1.8.8-latest


docker rm -vf order-report-titans-ficc; docker run --restart=always --name=order-report-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/order-report/log:/opt/log \
  -v /opt/rhino/titans-ficc/trade-engine/config/config.json:/config.json \
  -v /opt/rhino/titans-ficc/data/config/:/config_data/ \
  -e ServerPort=9093 \
  -p 9093:9093 \
  docker2.gf.com.cn/rhino/order-report-titans-ficc:b20251010-p1.1.3.5
