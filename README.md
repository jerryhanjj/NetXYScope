# NetXYScope

一个用于搜索NETCONF配置文件和YANG模型文件的命令行工具。

## 功能特性

- 🔍 支持在XML、YANG和YIN文件中搜索
- 🎨 高亮显示搜索结果
- 📁 递归搜索目录结构
- 🏷️ 支持多种文件格式（.xml, .yang, .yin）
- 📊 按文件分组显示搜索结果

## 支持的文件类型

- `.xml` - NETCONF XML配置文件
- `.yang` - YANG模型文件
- `.yin` - YANG的XML格式文件

## 安装

### 使用 go install（推荐）

```bash
go install github.com/jerryhanjj/NetXYScope@latest
```

安装完成后，可执行文件将被命名为 `NetXYScope`（与仓库名相同）。

### 从源码构建

```bash
git clone https://github.com/jerryhanjj/NetXYScope.git
cd NetXYScope
go build -o netconf-search .
```

### 使用 Makefile

```bash
git clone https://github.com/jerryhanjj/NetXYScope.git
cd NetXYScope
make install
```

## 使用方法

```bash
# 如果通过 go install 安装
NetXYScope <搜索词> <目录路径>

# 如果从源码构建
netconf-search <搜索词> <目录路径>
```

### 示例

```bash
# 搜索包含"interface"的配置（使用 go install 安装的版本）
NetXYScope interface /path/to/config/files

# 搜索YANG模型中的"container"定义（从源码构建的版本）
netconf-search container ./yang-models

# 搜索特定的配置节点
NetXYScope hostname ./testdata
```

## 项目结构

```
NetXYScope/
├── main.go                      # 主程序入口
├── internal/
│   ├── models/                  # 数据模型定义
│   │   └── models.go
│   └── search/                  # 搜索引擎实现
│       └── engine.go
├── testdata/                    # 测试数据
│   ├── test.xml
│   └── test.yang
├── scripts/                     # 发布脚本
│   └── release.sh
├── go.mod
└── README.md
```

## 代码架构

### 核心组件

1. **models.SearchResult** - 搜索结果数据结构
2. **models.NETCONFConfig** - NETCONF配置解析结构
3. **search.Engine** - 搜索引擎核心逻辑

### 关键特性

- **模块化设计**: 代码按功能分离到不同包中
- **XML解析**: 使用Go的encoding/xml包解析NETCONF配置
- **结构体标签**: 通过XML标签指导解析过程
- **错误处理**: 完善的错误处理和用户友好的错误信息

## 开发

### 构建

```bash
go build -o netconf-search .
```

### 测试

```bash
go test ./...
```

### 运行示例

```bash
# 使用项目中的测试数据（从源码构建的版本）
./netconf-search interface ./testdata
./netconf-search hostname ./testdata

# 使用 go install 安装的版本
NetXYScope interface ./testdata
NetXYScope hostname ./testdata
```

### 使用 Makefile

```bash
# 构建
make build

# 安装
make install

# 运行示例
make run-example

# 查看所有可用命令
make help
```

## XML结构体标签说明

项目中使用了Go的结构体标签来解析XML：

```go
type NETCONFConfig struct {
    XMLName xml.Name `xml:"config"`      // 指定XML根元素名称
    Content string   `xml:",innerxml"`   // 获取所有内部XML内容
}
```

### 标签含义：

- `xml:"config"` - 告诉解析器这个结构体对应XML中的`<config>`元素
- `xml:",innerxml"` - 获取元素内部的所有XML内容（包括子元素）

## 示例输出

```
Found 2 matches for 'interface':

=== ./testdata/test.xml ===
   2 | <interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
   3 |   <interface>

=== ./testdata/test.yang ===
  15 | container interfaces {
  20 |   list interface {
```

## 贡献

欢迎提交Issue和Pull Request来改进这个工具。

## 许可证

[MIT License](LICENSE)

## 更新日志

### v1.0.0
- 初始版本
- 支持基本的文件搜索功能
- 支持XML、YANG、YIN文件格式
- 添加搜索结果高亮显示
