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
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type Schema struct {
	Project       string
	PackagePath   string
	OutputPackage string
	Name          string
	SchemaName    string
	Cols          []*Cols
	Primary       *Cols
	Index         map[string][]*Cols
	Shard         int
	ShardCols     *Cols
}

type Cols struct {
	Name       string
	SchemaName string
	Type       string
	Tag        string
	IsIndex    bool
}

const (
	tagKey = "db"

	tagPrimary = "primary"
	tagIndex   = "index"
	tagShard   = "shard"
)

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
	var schemaList []*Schema

	for _, v := range f {
		for kk, ff := range v.Files {
			for _, a := range ff.Scope.Objects {
				if a.Kind != ast.Typ {
					continue
				}
				tmpSchema := &Schema{}
				tmpSchema.Project = projectName
				tmpSchema.Name = a.Name
				tmpSchema.PackagePath = filepath.Dir(kk)
				tmpSchema.SchemaName = Camel2Snake(a.Name)
				tmpSchema.Index = make(map[string][]*Cols)
				if len(ff.Comments) > 0 {
					for _, docLine := range ff.Comments[0].List {
						tmpDoc := strings.Trim(strings.TrimLeft(docLine.Text, "//"), " ")
						if !strings.HasPrefix(tmpDoc, "@") {
							continue
						}
						tmpDoc = strings.Trim(strings.TrimLeft(tmpDoc, "@"), " ")
						sepIndex := strings.Index(tmpDoc, ":")
						docKey, docValue := strings.Trim(tmpDoc[:sepIndex], " "), tmpDoc[sepIndex+1:]
						switch docKey {
						case "Name":
							tmpSchema.SchemaName = strings.Trim(docValue, " ")
						default:
						}
					}
				}
				for _, field := range a.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
					cols := &Cols{
						Name:       field.Names[0].Name,
						SchemaName: CapLow(field.Names[0].Name),
						Type:       field.Type.(*ast.Ident).Name,
					}
					if field != nil {
						if field.Tag == nil {
							continue
						}
						cols.Tag = strings.TrimRight(field.Tag.Value, " ")
						tagStr, err := strconv.Unquote(cols.Tag)
						if err != nil {
							panic(err)
						}
						tv := reflect.StructTag(tagStr)
						value, ok := tv.Lookup(tagKey)
						if !ok || value == "_" {
							continue
						}
						for _, valueStr := range strings.Split(value, ";") {
							var vk, vv string
							spIndex := strings.Index(valueStr, ":")
							if spIndex != -1 {
								vk = valueStr[:spIndex]
								vv = valueStr[spIndex+1:]
							} else {
								vk = valueStr
							}
							if vk == tagPrimary {
								tmpSchema.Primary = cols
							}
							if vk == tagIndex {
								cols.IsIndex = true
								tmpSchema.Index[vv] = append(tmpSchema.Index[vv], cols)
							}
							if vk == tagShard {
								tmpSchema.ShardCols = cols
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
		schemaObj.OutputPackage = filepath.Base(outputPath)
		fpath := fmt.Sprintf("%s/schema-%s-generated.go", outputPath, schemaObj.Name)
		if !Exists(outputPath) {
			err := os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		fileGen, err := os.Create(fpath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		err = temp.Execute(fileGen, schemaObj)
		if err != nil {
			fmt.Println("err when template. ", err)
			panic(err)
		}

		_ = exec.Command("gofmt", "-w", fpath).Run()
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

func Camel2Snake(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
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
