# NetXYScope

一个用于搜索NETCONF配置文件和YANG模型文件的命令行工具。支持任意XML根元素格式，适用于所有NETCONF消息类型。

## 功能特性

- 🔍 支持在XML、YANG和YIN文件中搜索
- 🎨 高亮显示搜索结果
- 📁 递归搜索目录结构
- 🏷️ 支持多种文件格式（.xml, .yang, .yin）
- 📊 按文件分组显示搜索结果
- 🌐 **通用XML解析**：支持任意XML根元素（hello, rpc, rpc-reply, config等）
- ⚡ **并发搜索**：多文件并行处理，提升搜索效率

## 支持的文件类型

- `.xml` - 所有类型的XML文件，包括：
  - NETCONF消息（hello, rpc, rpc-reply）
  - 配置文件（config）
  - 通知消息（notification）
  - 其他任意XML格式
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

# 搜索NETCONF消息类型
NetXYScope hello ./netconf-messages

# 搜索RPC调用
NetXYScope rpc-reply ./netconf-logs
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
│   ├── test.xml                 # 标准配置文件
│   ├── hello.xml               # NETCONF hello消息
│   ├── rpc-reply.xml           # NETCONF RPC响应
│   └── test.yang               # YANG模型文件
├── scripts/                     # 发布脚本
│   └── release.sh
├── go.mod
└── README.md
```

## 代码架构

### 核心组件

1. **models.SearchResult** - 搜索结果数据结构
2. **models.GenericXMLContent** - 通用XML内容解析结构（支持任意根元素）
3. **models.NETCONFConfig** - 向后兼容的别名（已废弃）
4. **search.Engine** - 搜索引擎核心逻辑，支持并发处理

### 关键特性

- **模块化设计**: 代码按功能分离到不同包中
- **通用XML解析**: 支持任意XML根元素，不限制于特定的NETCONF消息格式
- **并发处理**: 使用worker pool模式提升大目录搜索性能
- **智能XML验证**: 先验证XML格式，再进行内容搜索
- **结构体标签**: 通过XML标签指导解析过程
- **错误处理**: 完善的错误处理和用户友好的错误信息

## XML根元素支持

NetXYScope现在支持所有类型的XML根元素，包括但不限于：

```xml
<!-- NETCONF Hello消息 -->
<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <capabilities>...</capabilities>
</hello>

<!-- NETCONF RPC调用 -->
<rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <get-config>...</get-config>
</rpc>

<!-- NETCONF RPC响应 -->
<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <data>...</data>
</rpc-reply>

<!-- 配置数据 -->
<config>
  <interfaces>...</interfaces>
</config>

<!-- 通知消息 -->
<notification xmlns="urn:ietf:params:xml:ns:netconf:notification:1.0">
  <eventTime>...</eventTime>
</notification>
```

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
./netconf-search hello ./testdata        # 搜索NETCONF hello消息
./netconf-search rpc-reply ./testdata    # 搜索RPC响应

# 使用 go install 安装的版本
NetXYScope interface ./testdata
NetXYScope hostname ./testdata
NetXYScope capabilities ./testdata       # 搜索能力声明
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

项目现在使用更通用的XML解析结构：

```go
// 新的通用XML结构（推荐）
type GenericXMLContent struct {
    XMLName xml.Name `xml:""`          // 空字符串匹配任何根元素
    Content string   `xml:",innerxml"` // 获取所有内部XML内容
}

// 旧的结构（已废弃，保持向后兼容）
type NETCONFConfig struct {
    XMLName xml.Name `xml:"config"`    // 只能匹配 <config> 根元素
    Content string   `xml:",innerxml"` 
}
```

### 改进说明：

- **新版本**：`xml:""` - 空字符串表示匹配任何XML根元素
- **旧版本**：`xml:"config"` - 只能匹配特定的`<config>`元素
- `xml:",innerxml"` - 获取元素内部的所有XML内容（包括子元素）

这种改进使得工具能够处理所有类型的NETCONF消息和XML文件。

## 示例输出

```
Found 25 matches for 'interface':

=== testdata/test.xml ===
   2 | <interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
   3 | <interface>
   5 | <description>Management interface</description>

=== testdata/hello.xml ===
   5 | <capability>urn:ietf:params:netconf:capability:interface:1.0</capability>

=== testdata/rpc-reply.xml ===
   4 | <interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
   5 | <interface>
   9 | <description>Management interface for router</description>

=== testdata/test.yang ===
  17 | container interfaces {
  20 | list interface {
  26 | description "Interface name";
```

## 性能特性

- **并发处理**: 使用worker pool模式，支持多文件并行搜索
- **内存优化**: 流式处理大文件，避免内存溢出
- **智能解析**: 先验证XML格式，失败时自动降级到文本搜索
- **结果去重**: 避免XML解析和文本搜索产生的重复结果
  15 | container interfaces {
  20 |   list interface {
```

## 贡献

欢迎提交Issue和Pull Request来改进这个工具。

## 许可证

[MIT License](LICENSE)

## 更新日志

### v1.1.0 (最新)
- ✨ **重大改进**: 支持任意XML根元素，不再限制于`<config>`
- 🚀 **性能优化**: 增加并发搜索支持，提升大目录处理性能
- 🔧 **架构重构**: 引入`GenericXMLContent`替代固定的`NETCONFConfig`
- 📁 **测试扩展**: 添加多种NETCONF消息类型的测试文件
- 🛡️ **健壮性**: 改进错误处理和XML格式验证
- 📚 **文档完善**: 更新README和代码注释

### v1.0.0
- 初始版本
- 支持基本的文件搜索功能
- 支持XML、YANG、YIN文件格式
- 添加搜索结果高亮显示
