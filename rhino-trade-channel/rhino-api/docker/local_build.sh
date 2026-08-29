#REGISTRY=dockertest
VERSION=1.0.2.0
REGISTRY=docker2

lux_go_build -o rhino-api-server ../cmd/*.go
ssh eagle@10.51.136.22 "rm -rf /tmp/rhino-api"
scp -r ../docker eagle@10.51.136.22:/tmp/rhino-api
ssh root@10.51.136.22 "cd /tmp/rhino-api; docker build -t $REGISTRY.gf.com.cn/rhino/api-server:amd64_v$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/api-server:amd64_v$VERSION"

arm_go_build -o rhino-api-server ../cmd/*.go
ssh root@10.128.12.31  "rm -rf /tmp/rhino-api"
scp -r ../docker root@10.128.12.31:/tmp/rhino-api
ssh root@10.128.12.31 "cd /tmp/rhino-api; mv sources-arm.list sources.list; docker build -t $REGISTRY.gf.com.cn/rhino/api-server:arm64_v$VERSION .; docker --config ~/.docker/config-spring push $REGISTRY.gf.com.cn/rhino/api-server:arm64_v$VERSION"

ssh eagle@10.51.136.22 "sudo DOCKER_CLI_EXPERIMENTAL=enabled docker --config ~/.docker/config-spring manifest create $REGISTRY.gf.com.cn/rhino/api-server:v$VERSION --amend  $REGISTRY.gf.com.cn/rhino/api-server:amd64_v$VERSION --amend $REGISTRY.gf.com.cn/rhino/api-server:arm64_v$VERSION; sudo DOCKER_CLI_EXPERIMENTAL=enabled  docker --config ~/.docker/config-spring manifest push $REGISTRY.gf.com.cn/rhino/api-server:v$VERSION"
