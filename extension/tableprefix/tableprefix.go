package tableprefix

import (
	"embed"
	"strings"

	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed templates/dynamic/*.tmpl
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
		gen.NewTemplate("").
			ParseFS(templates, "templates/dynamic/*.tmpl"),
	)

	return []*gen.Template{tmpl}
}

type TablePrefixGenExtension struct {
	entc.DefaultExtension
	Prefix string
}

func NewTablePrefixGen(prefix string) (*TablePrefixGenExtension, error) {
	return &TablePrefixGenExtension{Prefix: prefix}, nil
}

func (e *TablePrefixGenExtension) Hooks() []gen.Hook {
	return []gen.Hook{
		e.generateTablePrefix(e.Prefix),
	}
}

func (e *TablePrefixGenExtension) generateTablePrefix(prefix string) gen.Hook {
	return AddTablePrefix(prefix)
}

func AddTablePrefix(prefix string) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {

			for _, n := range g.Nodes {
				if n.Annotations == nil {
					n.Annotations = make(map[string]any)
				}

				ann := &entsql.Annotation{}
				if a := n.EntSQL(); a != nil {
					*ann = *a
				}
				if ann.Table == "" {
					ann.Table = prefix + strings.ToLower(n.Name)
				}
				n.Annotations[ann.Name()] = ann
			}
			if err := next.Generate(g); err != nil {
				return err
			}
			return nil
		})
	}
}
