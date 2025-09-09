# 发布和安装指南

## 为第三方用户提供 go install 安装

### 1. 发布到 GitHub

首先需要将项目推送到 GitHub：

```bash
# 初始化 git 仓库（如果还没有）
git init

# 添加所有文件
git add .

# 提交代码
git commit -m "Initial release of netconf-search-tool"

# 添加远程仓库（替换为您的实际仓库地址）
git remote add origin https://github.com/jerryhanjj/netconf-search-tool.git

# 推送到 GitHub
git push -u origin main
```

### 2. 创建发布版本

在 GitHub 上创建一个发布版本（release）：

1. 进入您的 GitHub 仓库
2. 点击 "Releases"
3. 点击 "Create a new release"
4. 设置标签版本（如 v1.0.0）
5. 填写发布说明
6. 点击 "Publish release"

### 3. 用户安装方式

用户可以通过以下方式安装您的工具：

#### 方式一：直接安装最新版本
```bash
go install github.com/your-username/netconf-search-tool/cmd/netconf-search@latest
```

#### 方式二：安装特定版本
```bash
go install github.com/your-username/netconf-search-tool/cmd/netconf-search@v1.0.0
```

#### 方式三：从源码安装
```bash
git clone https://github.com/your-username/netconf-search-tool.git
cd netconf-search-tool
go install ./cmd/netconf-search
```

### 4. 验证安装

用户安装后，可以这样验证：

```bash
# 检查是否安装成功
which netconf-search

# 查看帮助信息
netconf-search

# 运行测试搜索
netconf-search interface /path/to/config/files
```

### 5. 更新模块路径

请确保在发布前，将 `go.mod` 中的模块路径更新为您的实际 GitHub 仓库地址：

```go
module github.com/your-actual-username/netconf-search-tool
```

并相应更新所有 import 语句。

### 6. 最佳实践

1. **版本标签**: 使用语义化版本号（如 v1.0.0, v1.1.0）
2. **文档**: 保持 README.md 更新，包含安装和使用说明
3. **变更日志**: 在发布说明中描述新功能和修复
4. **兼容性**: 遵循 Go 模块的兼容性保证

### 7. 示例安装命令

一旦发布到 GitHub，用户只需要运行：

```bash
go install github.com/your-username/netconf-search-tool/cmd/netconf-search@latest
```

工具将自动安装到 `$GOPATH/bin` 或 `$HOME/go/bin` 目录。
