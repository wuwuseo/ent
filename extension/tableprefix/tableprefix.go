package tableprefix

import (
	"embed"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed templates/dynamic/*.tmpl
var templates embed.FS

// TablePrefixExtension 是一个 Ent 扩展，用于为生成的数据库表名添加统一前缀
type TablePrefixExtension struct {
	entc.DefaultExtension
}

func NewTablePrefix() *TablePrefixExtension {
	return &TablePrefixExtension{}
}

// Templates 返回用于代码生成的模板列表
// 该方法加载 tableprefix.tmpl 模板，并将扩展实例注入到模板上下文中
// 使模板可以访问前缀配置
func (e *TablePrefixExtension) Templates() []*gen.Template {
	// 创建模板并解析所有扩展模板文件
	tmpl := gen.MustParse(
		gen.NewTemplate("").
			ParseFS(templates, "templates/dynamic/*.tmpl"),
	)

	return []*gen.Template{tmpl}
}
