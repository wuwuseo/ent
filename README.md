# Ent 表前缀扩展

一个用于 [Ent](https://entgo.io/) ORM 框架的扩展插件，可以为生成的数据库表名添加统一的前缀。在多租户系统、多应用共享数据库或需要表名隔离的场景中，表前缀功能可以帮助开发者更好地组织和管理数据库表结构。

## ✨ 特性

- 🎯 **简单易用**：通过一行代码即可为所有表添加前缀
- 🔧 **灵活配置**：支持空前缀、动态前缀、多环境配置
- 🛡️ **安全验证**：自动验证前缀是否为合法的 SQL 标识符
- 🔌 **无缝集成**：与 Ent 框架完美集成，不影响其他功能
- 🚀 **零性能开销**：代码生成时确定表名，运行时无额外计算
- 🤝 **扩展兼容**：可与其他 Ent 扩展同时使用
- 📦 **嵌入式模板**：使用 Go embed 特性，无需担心模板文件路径问题

## ⚙️ 模板覆盖说明

- 覆盖 `meta.tmpl`（SQL 方言）：
  - 将实体包中的 `Table` 由常量改为变量：`var Table = TablePrefix + "<原表名>"`
  - 在生成文件中引入 `github.com/wuwuseo/cmf/orm` 并定义 `var TablePrefix = orm.GetTablePrefix()`
- 覆盖 `meta/order`（SQL 方言）：
  - 邻接查询与排序步骤使用前缀化表名，确保跨实体关系时表名一致
- 覆盖 `migrate/schema.tmpl`：
  - 迁移生成的 `schema.Table{Name: ...}` 使用前缀化表名，`Tables` 列表保持一致

> 注意：生成的代码会导入 `github.com/wuwuseo/cmf/orm`，请确保该依赖在你的 `go.mod` 中（见安装章节）。

## 🛠️ 命令行用法

- 如果不通过扩展注册，也可以使用 CLI 覆盖模板：

```bash
ent generate --template ./extension/tableprefix/templates ./ent/schema
```

> 等价于注册扩展后自动加载模板；两种方式选其一即可。

## 📦 安装

使用 `go get` 命令安装：

```bash
go get github.com/wuwuseo/ent/extension/tableprefix
go get github.com/wuwuseo/cmf/orm
```

## 🚀 快速开始

### 1. 创建扩展实例

在你的 Ent 代码生成入口文件（通常是 `ent/generate.go`）中：

```go
//go:build ignore

package main

import (
    "log"
    
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/wuwuseo/ent/extension/tableprefix"
)

func main() {
    // 创建表前缀扩展，为所有表添加 "myapp_" 前缀
    ext, err := tableprefix.NewTablePrefix("myapp_")
    if err != nil {
        log.Fatalf("创建表前缀扩展失败: %v", err)
    }

    // 执行代码生成
    err = entc.Generate("./schema", &gen.Config{
        Target: "./ent",
    }, entc.Extensions(
        ext, // 注册表前缀扩展
    ))
    if err != nil {
        log.Fatalf("运行 ent 代码生成失败: %v", err)
    }
}
```

### 2. 运行代码生成

```bash
go generate ./ent
```

### 3. 效果

假设你有一个 `User` 实体，原始表名为 `users`，使用前缀 `myapp_` 后，生成的表名将变为 `myapp_users`。

## 📖 使用示例

### 基本用法

```go
// 为所有表添加固定前缀
ext, err := tableprefix.NewTablePrefix("myapp_")
if err != nil {
    log.Fatalf("创建扩展失败: %v", err)
}

err = entc.Generate("./schema", &gen.Config{
    Target: "./ent",
}, entc.Extensions(ext))
```

### 多环境配置

根据环境变量使用不同的表前缀：

```go
// 从环境变量读取前缀
prefix := os.Getenv("TABLE_PREFIX")
if prefix == "" {
    // 根据环境设置默认前缀
    env := os.Getenv("APP_ENV")
    switch env {
    case "production":
        prefix = "prod_"
    case "staging":
        prefix = "stg_"
    case "test":
        prefix = "test_"
    default:
        prefix = "dev_"
    }
}

ext, err := tableprefix.NewTablePrefix(prefix)
if err != nil {
    log.Fatalf("创建扩展失败: %v", err)
}

err = entc.Generate("./schema", &gen.Config{
    Target: "./ent",
}, entc.Extensions(ext))
```

设置环境变量：

```bash
# Linux/Mac
export TABLE_PREFIX="dev_"

# Windows CMD
set TABLE_PREFIX=dev_

# Windows PowerShell
$env:TABLE_PREFIX="dev_"
```

### 与其他扩展组合使用

```go
// 创建表前缀扩展
tablePrefixExt, err := tableprefix.NewTablePrefix("app_")
if err != nil {
    log.Fatalf("创建扩展失败: %v", err)
}

// 同时使用多个扩展
err = entc.Generate("./schema", &gen.Config{
    Target: "./ent",
    Features: []gen.Feature{
        gen.FeaturePrivacy,  // 启用隐私层
        gen.FeatureEntQL,    // 启用 EntQL
        gen.FeatureSnapshot, // 启用快照
    },
}, entc.Extensions(
    tablePrefixExt,
    // 其他扩展...
))
```

### 空前缀处理

如果需要在某些情况下不使用前缀：

```go
// 传入空字符串，不修改原始表名
ext, err := tableprefix.NewTablePrefix("")
if err != nil {
    log.Fatalf("创建扩展失败: %v", err)
}
```

### 多租户场景

为不同租户使用不同的表前缀：

```go
tenantID := "tenant123"
prefix := fmt.Sprintf("%s_", tenantID)

ext, err := tableprefix.NewTablePrefix(prefix)
if err != nil {
    log.Fatalf("创建扩展失败: %v", err)
}
```

## 📚 API 文档

### NewTablePrefix

```go
func NewTablePrefix(prefix string) (*TablePrefixExtension, error)
```

创建一个新的表前缀扩展实例。

**参数：**
- `prefix string`：要应用到所有表名的前缀字符串

**返回值：**
- `*TablePrefixExtension`：配置好的扩展实例
- `error`：如果前缀无效，返回错误信息

**前缀验证规则：**
- 只能包含字母、数字、下划线
- 不能以数字开头
- 长度不能超过 20 个字符
- 空字符串是合法的（表示不使用前缀）

**示例：**

```go
// 有效的前缀
ext, _ := ent.NewTablePrefix("myapp_")
ext, _ := ent.NewTablePrefix("v2_")
ext, _ := ent.NewTablePrefix("tenant_123_")
ext, _ := ent.NewTablePrefix("")  // 空前缀

// 无效的前缀（会返回错误）
ext, err := ent.NewTablePrefix("123_invalid")  // 以数字开头
ext, err := ent.NewTablePrefix("my-app_")      // 包含非法字符 '-'
ext, err := ent.NewTablePrefix("very_long_prefix_name_")  // 超过 20 字符
```

### TablePrefixExtension

```go
type TablePrefixExtension struct {
    entc.DefaultExtension
    // 私有字段
}
```

表前缀扩展结构体，实现了 Ent 的 Extension 接口。

**方法：**

#### Prefix

```go
func (e *TablePrefixExtension) Prefix() string
```

返回配置的表前缀。该方法主要用于在模板中访问前缀配置。

#### Templates

```go
func (e *TablePrefixExtension) Templates() []*gen.Template
```

返回用于代码生成的模板列表。该方法由 Ent 框架在代码生成时自动调用。

## ❓ 常见问题解答（FAQ）

### Q1: 如何修改已有项目的表前缀？

**A:** 修改表前缀需要以下步骤：

1. 更新代码生成配置中的前缀
2. 重新运行代码生成：`go generate ./ent`
3. 创建数据库迁移脚本，重命名现有表
4. 执行迁移脚本

**注意**：修改表前缀会影响现有数据库，请务必备份数据并在测试环境验证。

### Q2: 前缀会影响性能吗？

**A:** 不会。表前缀在代码生成时确定，生成的 `Table()` 方法返回常量字符串，运行时没有任何计算开销。

### Q3: 可以为不同的实体使用不同的前缀吗？

**A:** 当前版本不支持。所有实体使用相同的前缀。如果需要为不同实体使用不同前缀，可以考虑：
- 使用多个 Schema 目录，分别生成代码
- 或者在实体的 `Annotations` 中自定义表名

### Q4: 前缀可以包含哪些字符？

**A:** 前缀必须是合法的 SQL 标识符：
- ✅ 字母（a-z, A-Z）
- ✅ 数字（0-9，但不能作为首字符）
- ✅ 下划线（_）
- ❌ 连字符（-）
- ❌ 空格
- ❌ 特殊字符

### Q5: 如何在测试中使用不同的表前缀？

**A:** 可以通过环境变量或配置文件动态设置前缀：

```go
func getTablePrefix() string {
    if os.Getenv("GO_ENV") == "test" {
        return "test_"
    }
    return "prod_"
}

ext, _ := ent.NewTablePrefix(getTablePrefix())
```

### Q6: 表前缀会影响外键关系吗？

**A:** 不会。Ent 框架会自动处理外键关系，表前缀只影响表名，不影响实体之间的关联。

### Q7: 可以在运行时动态修改前缀吗？

**A:** 不可以。前缀在代码生成时确定，运行时无法修改。如果需要动态表名，建议使用数据库的 schema 或其他隔离机制。

### Q8: 前缀长度有限制吗？

**A:** 是的，前缀长度限制为 20 个字符。这是为了确保加上实体表名后，总长度不超过大多数数据库的表名限制（通常为 64 字符）。

### Q9: 如何验证生成的表名是否正确？

**A:** 代码生成后，可以查看生成的代码文件（通常在 `ent/` 目录下），搜索 `Table()` 方法查看生成的表名。

### Q10: 这个扩展与 Ent 的哪些版本兼容？

**A:** 
- 最低版本：Ent v0.10.0（引入扩展机制）
- 推荐版本：Ent v0.14.0+
- Go 版本要求：1.19+

## 🎯 使用场景

### 1. 多应用共享数据库

当多个应用共享同一个数据库时，使用表前缀可以避免表名冲突：

```go
// 应用 A
ext, _ := ent.NewTablePrefix("app_a_")

// 应用 B
ext, _ := ent.NewTablePrefix("app_b_")
```

### 2. 多租户系统

为不同租户使用不同的表前缀实现数据隔离：

```go
tenantID := getCurrentTenantID()
prefix := fmt.Sprintf("tenant_%s_", tenantID)
ext, _ := ent.NewTablePrefix(prefix)
```

### 3. 版本控制

在进行数据库架构升级时，使用版本前缀：

```go
// 旧版本表
ext, _ := ent.NewTablePrefix("v1_")

// 新版本表
ext, _ := ent.NewTablePrefix("v2_")
```

### 4. 环境隔离

为不同环境使用不同的表前缀：

```go
env := os.Getenv("APP_ENV")
var prefix string
switch env {
case "production":
    prefix = "prod_"
case "staging":
    prefix = "stg_"
default:
    prefix = "dev_"
}
ext, _ := ent.NewTablePrefix(prefix)
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. Fork 本仓库
2. 创建你的特性分支：`git checkout -b feature/amazing-feature`
3. 提交你的更改：`git commit -m '添加某个很棒的特性'`
4. 推送到分支：`git push origin feature/amazing-feature`
5. 提交 Pull Request

### 代码规范

- 遵循 Go 语言官方代码规范
- 所有公开函数和类型必须有中文注释
- 提交前运行 `go fmt` 和 `go vet`
- 添加必要的单元测试
- 更新相关文档

### 报告问题

如果你发现了 bug 或有功能建议，请在 GitHub Issues 中提交。提交时请包含：

- 问题的详细描述
- 复现步骤
- 预期行为和实际行为
- 环境信息（Go 版本、Ent 版本等）

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

### MIT License 摘要

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

## 🔗 相关链接

- [Ent 官方文档](https://entgo.io/)
- [Ent 扩展开发指南](https://entgo.io/docs/extensions/)
- [Go 模板语法](https://pkg.go.dev/text/template)

## 💬 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 提交 GitHub Issue
- 发送邮件至项目维护者

---

**感谢使用 Ent 表前缀扩展！** 🎉
