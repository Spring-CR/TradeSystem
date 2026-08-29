#!/bin/bash

BASE_VERSION=20251010  # 基础版本
PLUG_VERSION=1.4.6.9   # 插件版本
IS_PROD=0              # 1 表示 true，即生产环境；0 表示 false，即测试环境

COMPONENT=trade-engine-titans-ficc

if [ "$IS_PROD" -eq 1 ]; then
    REGISTRY="docker2"
    SUFFIX=""
else
    REGISTRY="dockertest"
    SUFFIX="-latest"
fi


VERSION=b$BASE_VERSION-p$PLUG_VERSION$SUFFIX

#echo "FROM docker2.gf.com.cn/library/ubuntu:18.04-latest" > ./docker/Dockerfile
#echo "ADD ./$COMPONENT /" >> ./docker/Dockerfile
#echo "RUN mkdir -p /opt/tradeclient; mkdir -p /opt/working_dir; mkdir -p /opt/log" >> ./docker/Dockerfile
#echo "WORKDIR /opt/log" >> ./docker/Dockerfile
#echo "VOLUME /opt/log" >> ./docker/Dockerfile
#echo "CMD [\"/bin/sh\", \"-c\", \"/trade-engine-titans-ficc 2>&1 | tee -a /opt/log/trade-engine.log\"]" >> ./docker/Dockerfile

cat > ./docker/Dockerfile <<EOF
FROM docker2.gf.com.cn/library/ubuntu:18.04-latest
#FROM docker2.gf.com.cn/library/java:v1.0.0-oraclejdk8-maven-3.3.9
ADD ./$COMPONENT /

# 配置目录并安装日志轮转工具
RUN mkdir -p /opt/tradeclient; mkdir -p /opt/working_dir; mkdir -p /opt/log; apt-get update && \\
    apt-get install -y logrotate cron && \\
    rm -rf /var/lib/apt/lists/*

# 配置日志轮转规则（每日00:00触发）
RUN echo "/opt/log/trade-engine.log {"  > /etc/logrotate.d/trade-engine && \\
    echo "  daily"                     >> /etc/logrotate.d/trade-engine && \\
    echo "  rotate 800"                >> /etc/logrotate.d/trade-engine && \\
    echo "  compress"                  >> /etc/logrotate.d/trade-engine && \\
    echo "  delaycompress"             >> /etc/logrotate.d/trade-engine && \\
    echo "  missingok"                 >> /etc/logrotate.d/trade-engine && \\
    echo "  notifempty"                >> /etc/logrotate.d/trade-engine && \\
    echo "  dateext"                   >> /etc/logrotate.d/trade-engine && \\
    echo "  copytruncate"              >> /etc/logrotate.d/trade-engine && \\
    echo "  create 644 root root"      >> /etc/logrotate.d/trade-engine && \\
    echo "}"

# 设置定时任务（每日00:00执行）
RUN echo "0 0 * * * root /usr/sbin/logrotate -f /etc/logrotate.d/trade-engine" > /etc/cron.d/trade-logrotate && \\
    chmod 644 /etc/cron.d/trade-logrotate


WORKDIR /opt/log
VOLUME /opt/log
CMD ["/bin/sh", "-c", "service cron start;/$COMPONENT 2>&1 | tee -a /opt/log/trade-engine.log"]
EOF

lux_go_build -o ./docker/$COMPONENT ./cmd/*.go
ssh eagle@10.51.136.23 "rm -rf /tmp/$COMPONENT"
scp -r ./docker eagle@10.51.136.23:/tmp/$COMPONENT
ssh root@10.51.136.23 "cd /tmp/$COMPONENT; docker build -t $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION"

arm_go_build -o ./docker/$COMPONENT ./cmd/*.go
ssh root@10.128.12.31  "rm -rf /tmp/$COMPONENT"
scp -r ./docker root@10.128.12.31:/tmp/$COMPONENT
ssh root@10.128.12.31 "cd /tmp/$COMPONENT; docker build -t $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION"

echo "sudo DOCKER_CLI_EXPERIMENTAL=enabled docker --config ~/.docker/config-spring manifest create $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION --amend  $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION --amend $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION; sudo DOCKER_CLI_EXPERIMENTAL=enabled  docker --config ~/.docker/config-spring manifest push $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION"
ssh eagle@10.51.136.23 "sudo DOCKER_CLI_EXPERIMENTAL=enabled docker --config ~/.docker/config-spring manifest create $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION --amend  $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION --amend $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION; sudo DOCKER_CLI_EXPERIMENTAL=enabled  docker --config ~/.docker/config-spring manifest push $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION"

rm ./docker/$COMPONENT