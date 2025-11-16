package tableprefix

import (
	"embed"
	"fmt"
	"regexp"
	"unicode"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed templates/*.tmpl
var templates embed.FS

// TablePrefixExtension 是一个 Ent 扩展，用于为生成的数据库表名添加统一前缀
type TablePrefixExtension struct {
	entc.DefaultExtension
	prefix string // 存储表前缀配置
}

// NewTablePrefix 创建一个新的表前缀扩展实例
// 参数 prefix 是要应用到所有表名的前缀字符串
// 如果 prefix 为空字符串，则不会修改原始表名
// 如果 prefix 包含非法字符，将返回错误
func NewTablePrefix(prefix string) (*TablePrefixExtension, error) {
	// 如果前缀为空，允许但不修改表名
	if prefix == "" {
		return &TablePrefixExtension{
			prefix: "",
		}, nil
	}

	// 验证前缀是否为合法的 SQL 标识符
	if err := validatePrefix(prefix); err != nil {
		return nil, fmt.Errorf("无效的表前缀: %w", err)
	}

	return &TablePrefixExtension{
		prefix: prefix,
	}, nil
}

// Prefix 返回配置的表前缀
// 该方法用于在模板中访问前缀配置
func (e *TablePrefixExtension) Prefix() string {
	return e.prefix
}

// validatePrefix 验证前缀是否为合法的 SQL 标识符
func validatePrefix(prefix string) error {
	// 检查前缀长度限制（建议不超过 20 字符）
	if len(prefix) > 20 {
		return fmt.Errorf("前缀长度不能超过 20 个字符，当前长度: %d", len(prefix))
	}

	// 检查前缀是否为合法的 SQL 标识符
	// 规则：只能包含字母、数字、下划线，且不能以数字开头
	sqlIdentifierPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !sqlIdentifierPattern.MatchString(prefix) {
		return fmt.Errorf("前缀必须是合法的 SQL 标识符（字母、数字、下划线，不能以数字开头）")
	}

	// 额外检查：确保第一个字符不是数字
	firstChar := rune(prefix[0])
	if unicode.IsDigit(firstChar) {
		return fmt.Errorf("前缀不能以数字开头")
	}

	return nil
}

// Templates 返回用于代码生成的模板列表
// 该方法加载 tableprefix.tmpl 模板，并将扩展实例注入到模板上下文中
// 使模板可以访问前缀配置
func (e *TablePrefixExtension) Templates() []*gen.Template {
	// 创建模板并解析所有扩展模板文件
	tmpl := gen.MustParse(
		gen.NewTemplate("").
			Funcs(map[string]any{
				"getPrefix": func() string { return e.prefix },
			}).
			ParseFS(templates, "templates/*.tmpl"),
	)

	return []*gen.Template{tmpl}
}
