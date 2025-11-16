# Ent TablePrefix 扩展

为 Ent 生成的 SQL 表名统一添加可配置前缀，适用于多租户、模块化系统、或需要区分组件的场景。本扩展通过自定义模板覆盖 Ent 的 SQL 方言元模板与迁移模板，使生成代码在运行时从 `github.com/wuwuseo/cmf/orm` 读取表前缀，并拼接到所有表名上。

## 主要特性
- 统一为实体表名添加前缀：`Table = TablePrefix + "<原始表名>"`
- 邻接查询与排序使用前缀化的表名，避免跨包依赖循环
- 迁移时的表定义同样应用前缀，保证 DDL 与运行时一致
- 与 `ent v0.14.x` 兼容

## 安装
- 使用 `go get` 安装扩展：

```
go get github.com/wuwuseo/ent@latest
```

## 使用方法
建议通过一个生成器程序调用 `entc.Generate` 并注册扩展。

示例：创建 `cmd/entgen/main.go`：

```go
package main

import (
    "log"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/wuwuseo/ent/extension/tableprefix"
)

func main() {
    ext, err := tableprefix.NewTablePrefix()
    if err != nil {
        log.Fatal(err)
    }

    // 根据你的项目实际路径调整 Target 与 Package
    if err := entc.Generate("./schema", &gen.Config{
        Target:  "./ent",
        //Package: "your/module/ent",
    }, entc.Extensions(ext)); err != nil {
        log.Fatal(err)
    }
}
```

运行生成：

```
go run ./cmd/entgen
```

## 配置表前缀
扩展依赖 `github.com/wuwuseo/cmf/orm` 提供的 `GetTablePrefix()` 在运行时返回前缀。你需要在应用启动时设置前缀，例如：

```go
package orm

var tablePrefix string

func GetTablePrefix() string  { return tablePrefix }
```

## 生成效果
- 在实体代码中会出现：

```go
var TablePrefix = orm.GetTablePrefix()
var Table = TablePrefix + "users" // 以 users 为例
```

- 在迁移代码中会出现：

```go
Name: TablePrefix + "users"
```

## 注意事项
- 仅覆盖 SQL 方言相关模板；非 SQL 方言未做适配
- 生成代码会依赖 `github.com/wuwuseo/cmf/orm`，请确保该包在你的模块中可用
- 不更改实体名、字段名及 API，仅影响底层表名与相关 SQL 构造
- 与 `ent v0.14.x` 配套测试，升级 Ent 版本时请验证生成结果

## 目录结构
- `extension/tableprefix/`：扩展实现与模板
- `extension/tableprefix/templates/`：覆盖的模板文件（meta、order、migrate_schema）

## 许可证
本项目基于 MIT 许可证开源，详见 `LICENSE`。