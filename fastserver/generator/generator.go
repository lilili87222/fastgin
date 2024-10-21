package generator

import (
	"fastgin/database"
	"fmt"
	"gorm.io/gen"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TableConfig struct {
	TableName       string
	PrefixTableName string
	OutDir          string
	Module          string
	ModelName       string
	LowModelName    string
	TableComment    string
}

func GenerateAll() {
	InitConfig()
	InitDatabase()
	genConfig := Instance.Generator
	dirs := []string{"model", "dao", "service", "controller", "route"}
	for _, dir := range dirs {
		fileDir := filepath.Join(genConfig.OutDir, dir)
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			os.MkdirAll(fileDir, os.ModePerm)
		}
	}
	tableConfig := TableConfig{
		PrefixTableName: "",
		TableName:       "",
		OutDir:          genConfig.OutDir,
		Module:          genConfig.Module,
	}

	for _, tableName := range genConfig.Tables {
		tableConfig.PrefixTableName = tableName
		tableConfig.TableName = strings.TrimPrefix(tableName, genConfig.TablePrefix)
		tableConfig.ModelName = ToCamelCase(tableConfig.TableName)
		tableConfig.LowModelName = strings.ToLower(tableConfig.ModelName)
		tableConfig.TableComment, _ = database.GetTableComment(DB, tableName)
		err := Generate(tableConfig)
		if err != nil {
			panic(err)
		}
	}
}
func Generate(tc TableConfig) error {
	GenerateModel(tc)
	if err := GenerateFromTemplate(tc, "dao"); err != nil {
		return err
	}
	if err := GenerateFromTemplate(tc, "service"); err != nil {
		return err
	}
	if err := GenerateFromTemplate(tc, "controller"); err != nil {
		return err
	}
	if err := GenerateFromTemplate(tc, "route"); err != nil {
		return err
	}
	return nil
}
func GenerateModel(tc TableConfig) {
	g := gen.NewGenerator(gen.Config{
		OutPath:          filepath.Join(tc.OutDir, "model"),
		FieldWithTypeTag: true,
	})
	g.UseDB(DB)
	g.GenerateModelAs(tc.PrefixTableName, tc.ModelName)
	g.Execute()
}
func GenerateFromTemplate(tc TableConfig, templateName string) error {
	tmpl, err := template.New(templateName + ".tpl").Funcs(template.FuncMap{
		"ToCamelCase": ToCamelCase,
		"ToLower":     strings.ToLower,
	}).ParseFiles("generator/templates/" + templateName + ".tpl")
	if err != nil {
		return err
	}
	file, err := os.Create(fmt.Sprintf("%s/%s/%s_%s.go", tc.OutDir, templateName, tc.TableName, templateName))
	if err != nil {
		return err
	}
	defer file.Close()
	return tmpl.Execute(file, tc)
}

func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}
