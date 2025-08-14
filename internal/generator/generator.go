package generator

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/rnikrozoft/orm2hex/internal/parser"
)

//go:embed templates/*.tmpl
var tmplFS embed.FS

type GeneratorConfig struct {
	ORM     string
	RawSQL  bool
	WithCtx bool
}

func GenerateHexCRUD(s parser.StructInfo, config GeneratorConfig, outDir string) error {
	absOutDir, _ := filepath.Abs(outDir)
	absStructFile, _ := filepath.Abs(s.FilePath)
	if strings.HasPrefix(absStructFile, absOutDir) {
		return nil
	}

	tplName := "helper"
	if config.RawSQL {
		tplName = "raw"
	}

	tplData, err := tmplFS.ReadFile("templates/" + tplName + ".tmpl")
	if err != nil {
		return err
	}

	tpl, err := template.New(tplName).Funcs(template.FuncMap{
		"lowerFirst":  lowerFirst,
		"toSnakeCase": toSnakeCase,
	}).Parse(string(tplData))
	if err != nil {
		return err
	}

	baseName := strings.TrimSuffix(s.Name, "Repository")
	fileName := filepath.Join(outDir, toSnakeCase(baseName)+"_repository.go")

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	data := map[string]interface{}{
		"Struct":  s,
		"ORM":     config.ORM,
		"RawSQL":  config.RawSQL,
		"WithCtx": config.WithCtx,
	}

	return tpl.Execute(f, data)
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
