#!/bin/bash

# Gitea 插件系统 Docker 构建和发布脚本

set -e

# 配置
DOCKER_REGISTRY="registry.cn-chengdu.aliyuncs.com"
IMAGE_NAME="gitea-plugins"
NAMESPACE="harti123"
VERSION="1.0.0"
LATEST_TAG="latest"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

echo_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo_error "Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查是否登录 Docker 仓库
echo_info "检查 Docker 登录状态..."
if ! docker info | grep -q "Username"; then
    echo_warn "未登录 Docker 仓库，请先登录"
    echo "运行: docker login ${DOCKER_REGISTRY}"
    exit 1
fi

# 清理旧的构建
echo_info "清理旧的构建文件..."
rm -f gitea
rm -rf dist/

# 构建 Docker 镜像
echo_info "开始构建 Docker 镜像..."
docker build \
    -f Dockerfile.custom \
    -t ${IMAGE_NAME}:${VERSION} \
    -t ${IMAGE_NAME}:${LATEST_TAG} \
    .

if [ $? -ne 0 ]; then
    echo_error "Docker 镜像构建失败"
    exit 1
fi

echo_info "Docker 镜像构建成功"

# 标记镜像
echo_info "标记镜像..."
docker tag ${IMAGE_NAME}:${VERSION} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${VERSION}
docker tag ${IMAGE_NAME}:${LATEST_TAG} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${LATEST_TAG}

# 推送到仓库
echo_info "推送镜像到 ${DOCKER_REGISTRY}..."
docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${VERSION}
docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${LATEST_TAG}

if [ $? -ne 0 ]; then
    echo_error "镜像推送失败"
    exit 1
fi

echo_info "镜像推送成功"
echo_info "镜像地址: ${DOCKER_REGISTRY}/${IMAGE_NAME}:${VERSION}"
echo_info "镜像地址: ${DOCKER_REGISTRY}/${IMAGE_NAME}:${LATEST_TAG}"

# 显示镜像信息
echo_info "镜像信息:"
docker images | grep ${IMAGE_NAME}

echo_info "完成！"
