package tableprefix

import (
	"embed"

	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed templates/*.tmpl
var templates embed.FS

// TablePrefixExtension 是一个 Ent 扩展，用于为生成的数据库表名添加统一前缀
type TablePrefixExtension struct {
	entc.DefaultExtension
}

func NewTablePrefix() (*TablePrefixExtension, error) {
	return &TablePrefixExtension{}, nil
}

// Templates 返回用于代码生成的模板列表
// 该方法加载 tableprefix.tmpl 模板，并将扩展实例注入到模板上下文中
// 使模板可以访问前缀配置
func (e *TablePrefixExtension) Templates() []*gen.Template {
	// 创建模板并解析所有扩展模板文件
	tmpl := gen.MustParse(
		gen.NewTemplate("schema").
			ParseFS(templates, "templates/*.tmpl"),
	)

	return []*gen.Template{tmpl}
}

type TablePrefixGenerateExtension struct {
	entc.DefaultExtension
	prefix string
}

func NewTablePrefixGenerate(prefix string) (*TablePrefixGenerateExtension, error) {
	return &TablePrefixGenerateExtension{prefix: prefix}, nil
}

func Hooks(prefix string) []gen.Hook {
	return []gen.Hook{
		AddTablePrefix(prefix),
	}
}

func AddTablePrefix(prefix string) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, n := range g.Nodes {
				// 获取现有 entsql 注解（如果有）
				var sqlAnnot *entsql.Annotation
				annotName := (&entsql.Annotation{}).Name() // 先获取注解名称

				if n.Annotations != nil {
					if annot, ok := n.Annotations[annotName]; ok {
						if a, ok := annot.(*entsql.Annotation); ok {
							sqlAnnot = a
						}
					}
				}

				if sqlAnnot == nil {
					sqlAnnot = &entsql.Annotation{}
				}

				sqlAnnot.Table = prefix + n.Table()

				if n.Annotations == nil {
					n.Annotations = gen.Annotations{}
				}
				n.Annotations[sqlAnnot.Name()] = sqlAnnot
			}
			return next.Generate(g)
		})
	}
}

func (e *TablePrefixGenerateExtension) AddTablePrefix(prefix string) gen.Hook {
	return AddTablePrefix(prefix)
}
