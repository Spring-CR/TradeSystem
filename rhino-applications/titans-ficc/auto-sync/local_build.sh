#!/bin/bash

BASE_VERSION=20251010  # 基础版本
PLUG_VERSION=1.0.1.3   # 插件版本
IS_PROD=1              # 1 表示 true，即生产环境；0 表示 false，即测试环境

COMPONENT=auto-sync

if [ "$IS_PROD" -eq 1 ]; then
    REGISTRY="docker2"
    SUFFIX=""
else
    REGISTRY="dockertest"
    SUFFIX="-latest"
fi


VERSION=b$BASE_VERSION-p$PLUG_VERSION$SUFFIX

# ---------- 传输项目 ----------
# echo "开始传输项目到远程服务器..."
# rsync -avz --delete /Users/spring/Documents/workspace_go/src/ eagle@10.51.136.23:/usr/local/gopath/src/
# echo "项目已成功传输到服务器！"

ssh eagle@10.51.136.23 "rm -rf /tmp/$COMPONENT"
scp -r ./docker eagle@10.51.136.23:/tmp/$COMPONENT
ssh root@10.51.136.23 "cd /tmp/$COMPONENT; docker build -t $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION"

ssh root@10.128.12.31 "rm -rf /tmp/$COMPONENT"
scp -r ./docker root@10.128.12.31:/tmp/$COMPONENT
ssh root@10.128.12.31 "cd /tmp/$COMPONENT; docker build -t $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION"

echo "sudo DOCKER_CLI_EXPERIMENTAL=enabled docker --config ~/.docker/config-spring manifest create $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION --amend  $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION --amend $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION; sudo DOCKER_CLI_EXPERIMENTAL=enabled  docker --config ~/.docker/config-spring manifest push $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION"
ssh eagle@10.51.136.23 "sudo DOCKER_CLI_EXPERIMENTAL=enabled docker --config ~/.docker/config-spring manifest create $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION --amend  $REGISTRY.gf.com.cn/rhino/$COMPONENT:amd64_$VERSION --amend $REGISTRY.gf.com.cn/rhino/$COMPONENT:arm64_$VERSION; sudo DOCKER_CLI_EXPERIMENTAL=enabled  docker --config ~/.docker/config-spring manifest push $REGISTRY.gf.com.cn/rhino/$COMPONENT:$VERSION"
