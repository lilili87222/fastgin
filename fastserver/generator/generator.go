package generator

import (
	"fastgin/database"
	"fmt"
	"gorm.io/gen"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

type TableConfig struct {
	TableName       string
	PrefixTableName string
	OutDir          string
	OutDirFront     string
	Module          string
	ModelName       string
	LowModelName    string
	TableComment    string
	GenerateFront   bool
	Columns         []TableColumn
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
		OutDirFront:     genConfig.OutDirFront,
		GenerateFront:   genConfig.GenerateFront,
	}

	for _, tableName := range genConfig.Tables {
		tableConfig.PrefixTableName = tableName
		tableConfig.TableName = strings.TrimPrefix(tableName, genConfig.TablePrefix)
		tableConfig.ModelName = ToCamelCase(tableConfig.TableName)
		tableConfig.LowModelName = strings.ToLower(tableConfig.ModelName)
		tableConfig.TableComment, _ = database.GetTableComment(DB, tableName)

		tableInfos, e := database.GetTableInfo(DB, tableName)
		if e != nil {
			panic(e)
		}
		tableConfig.Columns = ToTableColumns(tableInfos)

		err := Generate(tableConfig)
		if err != nil {
			panic(err)
		}
		if genConfig.GenerateFront {
			err := GenerateView(tableConfig)
			if err != nil {
				panic(err)
			}
		}
	}
}

func Generate(tc TableConfig) error {
	GenerateModel(tc)
	if err := GenerateGoFromTemplate(tc, "dao"); err != nil {
		return err
	}
	if err := GenerateGoFromTemplate(tc, "service"); err != nil {
		return err
	}
	if err := GenerateGoFromTemplate(tc, "controller"); err != nil {
		return err
	}
	if err := GenerateGoFromTemplate(tc, "route"); err != nil {
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
func GenerateGoFromTemplate(tc TableConfig, templateName string) error {
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
func GenerateView(tc TableConfig) error {
	if err := GenerateViewFromTemplate(tc, "api.ts", "api/app", tc.LowModelName+".ts"); err != nil {
		return err
	}
	if err := GenerateViewFromTemplate(tc, "types.ts", "types/app", tc.LowModelName+".ts"); err != nil {
		return err
	}
	if err := GenerateViewFromTemplate(tc, "dialog.vue", "views/app/"+tc.LowModelName, "dialog.vue"); err != nil {
		return err
	}
	if err := GenerateViewFromTemplate(tc, "index.vue", "views/app/"+tc.LowModelName, "index.vue"); err != nil {
		return err
	}
	return nil
}
func GenerateViewFromTemplate(tc TableConfig, templateName string, outputDir string, fileName string) error {
	outputPath := filepath.Join(tc.OutDirFront, outputDir, fileName)
	os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)

	tmpl, err := template.New(templateName+".gohtml").Funcs(template.FuncMap{
		"ToCamelCase": ToCamelCase,
		"ToLower":     strings.ToLower,
	}).Delims("{{{", "}}}").ParseFiles("generator/templates/" + templateName + ".gohtml")
	if err != nil {
		return err
	}
	file, err := os.Create(outputPath)
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
func ToTableColumns(infos []database.ColumnInfo) []TableColumn {
	excludeColumns := []string{"created_at", "updated_at", "deleted_at"}
	var columns []TableColumn
	for _, info := range infos {
		typeInfos := strings.Split(info.ColumnType, "(")
		typeSize := 0
		if len(typeInfos) > 1 {
			typeSize, _ = strconv.Atoi(strings.Split(typeInfos[1], ")")[0])
		}
		column := TableColumn{
			Name:         info.ColumnName,
			Type:         typeInfos[0],
			TypeSize:     typeSize,
			IsPriKey:     info.IsPriKey(),
			IsNullable:   info.IsNullableField(),
			DefaultValue: info.ColumnDefault,
			Comment:      info.ColumnComment,
		}
		column.Type = mysqlDataType2TypescriptType(column.Type)
		if slices.Contains(excludeColumns, column.Name) {
			continue
		}
		columns = append(columns, column)
	}
	return columns
}
func mysqlDataType2TypescriptType(mysqlType string) string {
	switch mysqlType {
	case "int", "tinyint", "smallint", "mediumint", "bigint":
		return "number"
	case "float", "double", "decimal":
		return "number"
	case "char", "varchar", "text", "tinytext", "mediumtext", "longtext":
		return "string"
	case "date", "time", "datetime", "timestamp":
		return "Date"
	default:
		return "any"
	}
}

type TableColumn struct {
	Name         string
	Type         string
	TypeSize     int
	IsPriKey     bool
	IsNullable   bool
	DefaultValue string
	Comment      string
}
