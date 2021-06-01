package gen

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type Schema struct {
	Project    string
	Name       string
	SchemaName string
	Cols       []*Cols
	Primary    *Cols
	Index      [][]*Cols
	Shard      int
	ShardCols  *Cols
}

type Cols struct {
	Name       string
	SchemaName string
	Type       string
	Tag        string
	IsIndex    bool
}

var templateFile string
var projectName string

func GenerateWithFlagScan() {
	var (
		scanPath   string
		outputPath string
	)
	flag.StringVar(&projectName, "projectName", "a", "projectName")
	flag.StringVar(&scanPath, "scanPath", "", "scanPath")
	flag.StringVar(&outputPath, "outputPath", "", "outputPath")
	flag.StringVar(&templateFile, "templateFile", "", "templateFile")
	flag.Parse()
	fmt.Println("start template")
	fmt.Printf("projectName: %s\n", projectName)
	fmt.Printf("scanPath: %s\n", scanPath)
	fmt.Printf("outputPath: %s\n", outputPath)
	schemaList := scan(scanPath)
	runGenerate(schemaList, outputPath)
}

func Generate(project, scanPath, outputPath string) {
	projectName = project
	fmt.Println("start template")
	fmt.Printf("projectName: %s\n", projectName)
	fmt.Printf("scanPath: %s\n", scanPath)
	fmt.Printf("outputPath: %s\n", outputPath)
	schemaList := scan(scanPath)
	runGenerate(schemaList, outputPath)
}

func scan(scanPath string) []*Schema {
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, scanPath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}
	var files []*ast.File
	var schemaList []*Schema
	for _, v := range f {
		for _, ff := range v.Files {
			files = append(files, ff)
			for _, a := range ff.Scope.Objects {
				if a.Kind != ast.Typ {
					continue
				}
				tmpSchema := &Schema{}
				tmpSchema.Project = projectName
				tmpSchema.Name = a.Name
				sname := CapLow(a.Name)
				if !strings.HasSuffix(sname, "s") {
					sname += "s"
				}
				tmpSchema.SchemaName = sname
				for _, field := range a.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
					cols := &Cols{
						Name:       field.Names[0].Name,
						SchemaName: CapLow(field.Names[0].Name),
						Type:       field.Type.(*ast.Ident).Name,
					}
					if field != nil {
						cols.Tag = strings.TrimRight(field.Tag.Value, " ")
						tagStr, err := strconv.Unquote(cols.Tag)
						if err != nil {
							panic(err)
						}
						tv := reflect.StructTag(tagStr)
						value, ok := tv.Lookup("db")
						if ok {
							for _, valueStr := range strings.Split(value, ";") {
								if valueStr == "primary" {
									tmpSchema.Primary = cols
								}
								if valueStr == "index" {
									cols.IsIndex = true
								}
								if valueStr == "shard" {
									tmpSchema.ShardCols = cols
								}
							}
						}

					}
					tmpSchema.Cols = append(tmpSchema.Cols, cols)
				}
				if tmpSchema.Shard > 0 && tmpSchema.ShardCols == nil {
					panic(errors.New(fmt.Sprintf("schema %s set shard=%d but no set shard field", tmpSchema.Name, tmpSchema.Shard)))
				}
				schemaList = append(schemaList, tmpSchema)
			}
		}
	}
	return schemaList
}

func runGenerate(schemaList []*Schema, outputPath string) {
	var temp *template.Template
	if templateFile != "" {
		temp = template.Must(template.ParseFiles(templateFile))
	} else {
		temp = template.Must(template.New("db").Parse(dbTemplate))
	}
	for _, schemaObj := range schemaList {
		filepath := fmt.Sprintf("%s/%s-generated.go", outputPath, schemaObj.SchemaName)
		if !Exists(outputPath) {
			err := os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		fileGen, err := os.Create(filepath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		err = temp.Execute(fileGen, schemaObj)
		if err != nil {
			fmt.Println("err when template. ", err)
			panic(err)
		}

		_ = exec.Command("gofmt", "-w", filepath).Run()
	}
}

func CapLow(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 {
				vv[i] += 32
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
