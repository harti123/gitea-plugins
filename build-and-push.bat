@echo off
REM Gitea 插件系统 Docker 构建和发布脚本 (Windows)

setlocal enabledelayedexpansion

REM 配置
set DOCKER_REGISTRY=registry.cn-hangzhou.aliyuncs.com
set NAMESPACE=harti
set IMAGE_NAME=gitea-plugins
set VERSION=1.0.0
set LATEST_TAG=latest

echo [INFO] 检查 Docker 是否安装...
docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker 未安装，请先安装 Docker Desktop
    exit /b 1
)

echo [INFO] 检查 Docker 登录状态...
docker info | findstr /C:"Username" >nul 2>&1
if errorlevel 1 (
    echo [WARN] 未登录 Docker 仓库，请先登录
    echo 运行: docker login %DOCKER_REGISTRY%
    exit /b 1
)

echo [INFO] 清理旧的构建文件...
if exist gitea.exe del gitea.exe
if exist dist rmdir /s /q dist

echo [INFO] 开始构建 Docker 镜像...
docker build -f Dockerfile.custom -t %IMAGE_NAME%:%VERSION% -t %IMAGE_NAME%:%LATEST_TAG% .
if errorlevel 1 (
    echo [ERROR] Docker 镜像构建失败
    exit /b 1
)

echo [INFO] Docker 镜像构建成功

echo [INFO] 标记镜像...
docker tag %IMAGE_NAME%:%VERSION% %DOCKER_REGISTRY%/%IMAGE_NAME%:%VERSION%
docker tag %IMAGE_NAME%:%LATEST_TAG% %DOCKER_REGISTRY%/%IMAGE_NAME%:%LATEST_TAG%

echo [INFO] 推送镜像到 %DOCKER_REGISTRY%...
docker push %DOCKER_REGISTRY%/%IMAGE_NAME%:%VERSION%
docker push %DOCKER_REGISTRY%/%IMAGE_NAME%:%LATEST_TAG%
if errorlevel 1 (
    echo [ERROR] 镜像推送失败
    exit /b 1
)

echo [INFO] 镜像推送成功
echo [INFO] 镜像地址: %DOCKER_REGISTRY%/%IMAGE_NAME%:%VERSION%
echo [INFO] 镜像地址: %DOCKER_REGISTRY%/%IMAGE_NAME%:%LATEST_TAG%

echo [INFO] 镜像信息:
docker images | findstr %IMAGE_NAME%

echo [INFO] 完成！
pause
