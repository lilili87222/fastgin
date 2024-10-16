package generator

import (
	"codegen/config"
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
}

func GenerateAll() {
	// 加载配置文件到全局配置结构体
	config.InitConfig()
	// 初始化数据库
	config.InitDatabase()
	genConfig := config.Instance.Generator
	dirs := []string{"model", "dao", "service", "controller"}
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
	//if err := GenerateDAO(tc); err != nil {
	//	return err
	//}
	//if err := GenerateService(tc); err != nil {
	//	return err
	//}
	//if err := GenerateController(tc); err != nil {
	//	return err
	//}
	return nil
}
func GenerateModel(tc TableConfig) {
	g := gen.NewGenerator(gen.Config{
		OutPath: filepath.Join(tc.OutDir, "model"),
	})
	g.UseDB(config.DB)
	g.GenerateModelAs(tc.PrefixTableName, ToCamelCase(tc.TableName))
	g.Execute()
}
func GenerateFromTemplate(tc TableConfig, templateName string) error {
	tmpl, err := template.New(templateName + ".tmpl").Funcs(template.FuncMap{
		"ToCamelCase": ToCamelCase,
	}).ParseFiles("templates/" + templateName + ".tmpl")
	if err != nil {
		return err
	}
	file, err := os.Create(fmt.Sprintf("%s/%s/%s_%s.go", tc.OutDir, templateName, templateName, tc.TableName))
	if err != nil {
		return err
	}
	defer file.Close()
	return tmpl.Execute(file, tc)
}

/*
func GenerateDAO(tc TableConfig) error {
	tmpl, err := template.New("dao.tmpl").Funcs(template.FuncMap{
		"ToCamelCase": ToCamelCase,
	}).ParseFiles("templates/dao.tmpl")
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/dao/%s_dao.go", tc.OutDir, tc.TableName))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, tc)
}

func GenerateService(tc TableConfig) error {
	tmpl, err := template.New("service.tmpl").Funcs(template.FuncMap{
		"ToCamelCase": ToCamelCase,
	}).ParseFiles("templates/service.tmpl")
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf(output+"/service/%s_service.go", tableName))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, struct {
		TableName string
	}{
		TableName: tableName,
	})
}

func GenerateController(tc TableConfig) error {
	tmpl, err := template.New("controller.tmpl").Funcs(template.FuncMap{
		"ToCamelCase": ToCamelCase,
	}).ParseFiles("templates/controller.tmpl")
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf(output+"/controller/%s_controller.go", tableName))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, struct {
		TableName string
	}{
		TableName: tableName,
	})
}*/

func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func ToGoType(sqlType string) string {
	switch sqlType {
	case "int", "integer", "smallint", "mediumint", "bigint":
		return "int"
	case "float", "double", "real":
		return "float64"
	case "decimal", "numeric":
		return "float64"
	case "char", "varchar", "text", "tinytext", "mediumtext", "longtext":
		return "string"
	case "date", "datetime", "timestamp", "time", "year":
		return "time.Time"
	case "boolean", "bool":
		return "bool"
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		return "[]byte"
	default:
		return "string"
	}
}
