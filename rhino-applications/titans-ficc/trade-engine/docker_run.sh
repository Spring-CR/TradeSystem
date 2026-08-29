#!/bin/bash

docker rm -vf trade-engine-titans-ficc; docker run --restart=always --name=trade-engine-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/trade-engine/log:/opt/log \
  -v /opt/rhino/titans-ficc/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/titans-ficc/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/titans-ficc/trade-engine/config/:/config/ \
  -v /opt/rhino/titans-ficc/data/config/:/config_data/ \
  -p 8093:8093 \
  -p 8293:8293 \
  dockertest.gf.com.cn/rhino/trade-engine-titans-ficc:b20251010-p1.4.6.9-latest


docker rm -vf trade-engine-titans-ficc; docker run --restart=always --name=trade-engine-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/trade-engine/log:/opt/log \
  -v /opt/rhino/titans-ficc/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/titans-ficc/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/titans-ficc/trade-engine/config/:/config/ \
  -v /opt/rhino/titans-ficc/data/config/:/config_data/ \
  -p 8093:8093 \
  -p 8293:8293 \
  docker2.gf.com.cn/rhino/trade-engine-titans-ficc:b20251010-p1.3.4.4


docker rm -vf trade-engine-titans-ficc; docker run  --name=trade-engine-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/trade-engine/log:/opt/log \
  -v /opt/rhino/titans-ficc/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/titans-ficc/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/titans-ficc/trade-engine/config/:/config/ \
  -v /opt/rhino/titans-ficc/data/config/:/config_data/ \
  -p 8093:8093 \
  -p 8293:8293 \
  dockertest.gf.com.cn/rhino/trade-engine-titans-ficc:b20251010-p1.3.0.1-latest


docker rm -vf trade-engine-titans-ficc; docker run --restart=always --name=trade-engine-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/trade-engine/log:/opt/log \
  -v /opt/rhino/titans-ficc/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/titans-ficc/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/titans-ficc/trade-engine/config/:/config/ \
  -v /opt/rhino/titans-ficc/data/config/:/config_data/ \
  -p 8093:8093 \
  dockertest.gf.com.cn/rhino/trade-engine-titans-ficc:b20251010-p1.1.6.7-latest

docker run --restart=always --name=trade-engine-titans-ficc -it -d \
  -v /opt/rhino/titans-ficc/trade-engine/log:/opt/log \
  -v /opt/rhino/titans-ficc/trade-engine/tradeclient:/opt/tradeclient \
  -v /opt/rhino/titans-ficc/trade-engine/working_dir:/opt/working_dir \
  -v /opt/rhino/titans-ficc/trade-engine/config/config.json:/config.json \
  -p 8093:8093 \
  --dns=10.51.167.99 \
  dockertest.gf.com.cn/rhino/trade-engine-titans-ficc:b20251010-p1.0.2.1-latest

