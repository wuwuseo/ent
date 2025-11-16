# 需求文档

## 简介

本项目是一个 Ent ORM 框架的扩展插件，用于为生成的数据库表名添加统一的前缀。在多租户系统或需要表名隔离的场景中，表前缀功能可以帮助开发者更好地组织和管理数据库表结构。

## 术语表

- **Ent**: Go 语言的实体框架（Entity Framework），用于构建和查询图结构的数据库模式
- **Extension**: Ent 框架的扩展机制，允许开发者自定义代码生成行为
- **Table Prefix**: 数据库表名前缀，例如将 `users` 表改为 `app_users`
- **Template**: Go 模板文件，用于自定义 Ent 代码生成的输出
- **Schema**: 数据库模式定义，描述实体的结构和关系

## 需求

### 需求 1

**用户故事：** 作为一个使用 Ent 框架的开发者，我希望能够为所有生成的数据库表添加统一的前缀，以便在多应用共享数据库时避免表名冲突。

#### 验收标准

1. THE TablePrefixExtension SHALL 实现 Ent 框架的扩展接口
2. WHEN 开发者调用 WithTablePrefix 函数并传入前缀字符串，THE TablePrefixExtension SHALL 存储该前缀配置
3. THE TablePrefixExtension SHALL 提供模板注册机制以自定义表名生成逻辑
4. THE TablePrefixExtension SHALL 确保前缀应用于所有实体的表名

### 需求 2

**用户故事：** 作为一个开发者，我希望能够通过简单的 API 配置表前缀，而不需要修改每个实体的定义。

#### 验收标准

1. THE WithTablePrefix 函数 SHALL 接受一个字符串参数作为表前缀
2. THE WithTablePrefix 函数 SHALL 返回一个可用于 Ent 代码生成配置的扩展实例
3. WHEN 前缀为空字符串，THE TablePrefixExtension SHALL 不修改原始表名
4. THE WithTablePrefix 函数 SHALL 支持任意合法的表名前缀字符串

### 需求 3

**用户故事：** 作为一个开发者，我希望扩展能够与 Ent 的代码生成流程无缝集成，不影响其他功能。

#### 验收标准

1. THE TablePrefixExtension SHALL 继承 DefaultExtension 以保持与 Ent 框架的兼容性
2. THE Templates 方法 SHALL 返回包含表前缀逻辑的模板列表
3. WHEN Ent 执行代码生成，THE TablePrefixExtension SHALL 自动应用表前缀模板
4. THE TablePrefixExtension SHALL 不干扰 Ent 框架的其他扩展功能

### 需求 4

**用户故事：** 作为一个开发者，我希望模板文件能够正确处理表名转换逻辑，生成符合预期的代码。

#### 验收标准

1. THE tableprefix.tmpl 模板 SHALL 定义完整的模板块结构
2. THE tableprefix.tmpl 模板 SHALL 访问配置中的包名和前缀信息
3. THE tableprefix.tmpl 模板 SHALL 生成为每个实体添加表前缀的代码
4. WHEN 模板被解析，THE tableprefix.tmpl SHALL 不产生语法错误

### 需求 5

**用户故事：** 作为一个开发者，我希望能够在 Ent 的代码生成入口文件中轻松集成此扩展。

#### 验收标准

1. THE TablePrefixExtension SHALL 可通过 entc.Generate 函数的 Extensions 选项注册
2. WHEN 扩展被注册，THE Ent 代码生成器 SHALL 在生成过程中调用扩展的模板
3. THE 扩展 SHALL 支持与其他 Ent 扩展同时使用
4. THE 扩展 SHALL 在代码生成完成后生成包含表前缀逻辑的辅助代码
