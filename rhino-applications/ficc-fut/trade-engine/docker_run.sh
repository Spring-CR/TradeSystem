#!/bin/bash

docker rm -vf trade-engine-ficc-fut; docker run --restart=always --name=trade-engine-ficc-fut -it -d \
  -v /opt/rhino/ficc-fut/trade-engine/log:/opt/log \
  -v /opt/rhino/ficc-fut/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/ficc-fut/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/ficc-fut/trade-engine/config/:/config/ \
  -v /opt/rhino/ficc-fut/data/config/:/config_data/ \
  -p 8095:8095 \
  -p 8295:8295 \
  dockertest.gf.com.cn/rhino/trade-engine-ficc-fut:b20251010-p1.0.6.5-latest


