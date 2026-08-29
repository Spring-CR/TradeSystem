#!/bin/bash

BASE_VERSION=20251010  # 基础版本
PLUG_VERSION=1.0.3.3   # 插件版本
IS_PROD=0              # 1 表示 true，即生产环境；0 表示 false，即测试环境

COMPONENT=data-titans-ficc

if [ "$IS_PROD" -eq 1 ]; then
    REGISTRY="docker2"
    SUFFIX=""
else
    REGISTRY="dockertest"
    SUFFIX="-latest"
fi


VERSION=b$BASE_VERSION-p$PLUG_VERSION$SUFFIX

# ---------- 传输项目 ----------
echo "开始传输项目到远程服务器..."
rsync -avz --delete /Users/spring/Documents/workspace_go/src/ eagle@10.51.136.23:/usr/local/gopath/src/
echo "项目已成功传输到服务器！"
echo "开始在远程服务器上编译..."

# ---------- 远程编译 ----------
ssh eagle@10.51.136.23 << 'EOF'
source ~/.bashrc
cd /usr/local/gopath/src/rhino-applications/titans-ficc/data

# 创建编译输出目录
mkdir -p artifacts/{x86,arm}

export GOPATH=/usr/local/gopath

# 编译 x86 版本
echo "开始编译 x86 版本..."
GO111MODULE=off CGO_ENABLED=1 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o artifacts/x86/ cmd/*.go
echo "x86 版本编译完成！"

# 编译 arm64 版本
echo "开始编译 arm64 版本..."
GO111MODULE=off CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc /usr/local/go/bin/go build -o artifacts/arm/ cmd/*.go
echo "arm64 版本编译完成！"
EOF

echo "远程编译全部完成！"



cat > ./docker/Dockerfile <<EOF
FROM docker2.gf.com.cn/library/ubuntu:18.04-latest
#FROM docker2.gf.com.cn/library/java:v1.0.0-oraclejdk8-maven-3.3.9
ADD ./$COMPONENT /

# 配置目录并安装日志轮转工具
RUN mkdir -p /opt/tradeclient; mkdir -p /opt/working_dir; mkdir -p /opt/log; apt-get update && \\
    apt-get install -y logrotate cron && \\
    rm -rf /var/lib/apt/lists/*

# 配置日志轮转规则（每日00:00触发）
RUN echo "/opt/log/data.log {"  > /etc/logrotate.d/data && \\
    echo "  daily"                     >> /etc/logrotate.d/data && \\
    echo "  rotate 7"                  >> /etc/logrotate.d/data && \\
    echo "  compress"                  >> /etc/logrotate.d/data && \\
    echo "  delaycompress"             >> /etc/logrotate.d/data && \\
    echo "  missingok"                 >> /etc/logrotate.d/data && \\
    echo "  notifempty"                >> /etc/logrotate.d/data && \\
    echo "  dateext"                   >> /etc/logrotate.d/data && \\
    echo "  copytruncate"              >> /etc/logrotate.d/data && \\
    echo "  create 644 root root"      >> /etc/logrotate.d/data && \\
    echo "}"

# 设置定时任务（每日00:00执行）
RUN echo "0 0 * * * root /usr/sbin/logrotate -f /etc/logrotate.d/data" > /etc/cron.d/trade-logrotate && \\
    chmod 644 /etc/cron.d/trade-logrotate


WORKDIR /opt/log
VOLUME /opt/log
CMD ["/bin/sh", "-c", "service cron start;/$COMPONENT 2>&1 | tee -a /opt/log/data.log"]
EOF

scp eagle@10.51.136.23:/usr/local/gopath/src/rhino-applications/titans-ficc/data/artifacts/x86/data ./docker/$COMPONENT
ssh eagle@10.51.136.23 "rm -rf /tmp/$COMPONENT"
scp -r ./docker eagle@10.51.136.23:/tmp/$COMPONENT
ssh root@10.51.136.23 "cd /tmp/$COMPONENT; docker build -t $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION"

scp eagle@10.51.136.23:/usr/local/gopath/src/rhino-applications/titans-ficc/data/artifacts/arm/data ./docker/$COMPONENT
ssh root@10.128.12.31  "rm -rf /tmp/$COMPONENT"
scp -r ./docker root@10.128.12.31:/tmp/$COMPONENT
ssh root@10.128.12.31 "cd /tmp/$COMPONENT; docker build -t $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION"

echo "sudo DOCKER_CLI_EXPERIMENTAL=enabled docker --config ~/.docker/config-spring manifest create $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION --amend  $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION --amend $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION; sudo DOCKER_CLI_EXPERIMENTAL=enabled  docker --config ~/.docker/config-spring manifest push $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION"
ssh eagle@10.51.136.23 "sudo DOCKER_CLI_EXPERIMENTAL=enabled docker --config ~/.docker/config-spring manifest create $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION --amend  $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION --amend $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION; sudo DOCKER_CLI_EXPERIMENTAL=enabled  docker --config ~/.docker/config-spring manifest push $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION"

rm ./docker/$COMPONENT