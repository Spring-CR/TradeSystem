#!/bin/bash

docker rm -vf trade-engine-stars-fut; docker run --restart=always --name=trade-engine-stars-fut -it -d \
  -v /opt/rhino/stars-fut/trade-engine/log:/opt/log \
  -v /opt/rhino/stars-fut/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/stars-fut/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/stars-fut/trade-engine/config/:/config/ \
  -p 8094:8094 \
  dockertest.gf.com.cn/rhino/trade-engine-stars-fut:b20251010-p1.0.0.2-latest


