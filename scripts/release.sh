#!/bin/bash

# Git 发布工作流示例脚本
# 使用前请先设置您的GitHub仓库地址

set -e

echo "🚀 NETCONF Search Tool 发布脚本"
echo "================================"

# 检查是否有未提交的更改
if [[ -n $(git status -s) ]]; then
    echo "❌ 检测到未提交的更改，请先提交所有更改"
    git status
    exit 1
fi

# 检查当前分支
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]; then
    echo "⚠️  当前不在主分支，确认要在 '$CURRENT_BRANCH' 分支发布吗? (y/N)"
    read -r confirm
    if [[ $confirm != [yY] ]]; then
        echo "❌ 发布已取消"
        exit 1
    fi
fi

# 构建测试
echo "🔨 构建测试..."
go build -o netconf-search ./cmd/netconf-search
if [[ $? -eq 0 ]]; then
    echo "✅ 构建成功"
    rm -f netconf-search
else
    echo "❌ 构建失败"
    exit 1
fi

# 运行测试
echo "🧪 运行测试..."
go test ./...
if [[ $? -eq 0 ]]; then
    echo "✅ 测试通过"
else
    echo "❌ 测试失败"
    exit 1
fi

# 获取版本号
echo "📝 请输入版本号 (例如: v1.0.0):"
read -r VERSION

if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "❌ 版本号格式错误，应该是 vX.Y.Z 格式"
    exit 1
fi

# 检查版本是否已存在
if git tag -l | grep -q "^$VERSION$"; then
    echo "❌ 版本 $VERSION 已存在"
    exit 1
fi

# 推送到远程
echo "📤 推送到远程仓库..."
git push origin $CURRENT_BRANCH

# 创建标签
echo "🏷️  创建版本标签 $VERSION..."
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION

echo "🎉 发布完成！"
echo ""
echo "用户现在可以通过以下命令安装："
echo "go install github.com/your-username/netconf-search-tool/cmd/netconf-search@$VERSION"
echo ""
echo "或安装最新版本："
echo "go install github.com/your-username/netconf-search-tool/cmd/netconf-search@latest"
echo ""
echo "请不要忘记在GitHub上创建正式的Release！"
