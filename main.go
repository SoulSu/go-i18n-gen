package main

import (
	"fmt"
	"go/build"
	"os"
	"reflect"
	"text/template"
)

type I18nErr struct {
	VariableName string
	Code         string
	ZhCnMsg      string
	EnUsMsg      string

	GenPackage string `ignore:"-"`
}

var errTpl = `// Code auto generated, DO NOT EDIT.
package {{.GenPackage}}

type ErrorMsg string

func (e ErrorMsg) Error() string {
	return " msg: " + string(e)
}

type ErrorCode string

func (e ErrorCode) Error() string {
	return "error code: " + string(e)
}

var I18nError = make(map[string]map[ErrorCode]ErrorMsg)

func registerError(lang string, err ErrorCode, msg ErrorMsg) {
	_, ok := I18nError[lang]
	if !ok {
		I18nError[lang] = make(map[ErrorCode]ErrorMsg)
	}
	I18nError[lang][err] = msg
}

func GetErrorMsg(lang string, err error) string {
	errCodeErr, ok := err.(ErrorCode)
	if !ok {
		return err.Error()
	}

	_, ok = I18nError[lang]
	if !ok {
		return "Not register lang: " + lang
	}

	_, ok = I18nError[lang][errCodeErr]
	if !ok {
		return "Not register " + errCodeErr.Error()
	}

	return I18nError[lang][errCodeErr].Error()
}

`

var tpl = `// Code auto generated, DO NOT EDIT.
package {{.GenPackage}}

const {{.VariableName}} ErrorCode = "{{.Code}}"

func init(){
	registerError("zh_cn", {{.VariableName}}, ErrorMsg("{{.ZhCnMsg}}"))
	registerError("en_us", {{.VariableName}}, ErrorMsg("{{.EnUsMsg}}"))
}
`

func panicf(format string, a ...interface{}) {
	panic(fmt.Errorf(format, a))
}

func findPackageName() string {
	p, err := build.Default.Import(".", ".", build.ImportMode(0))
	if err != nil {
		panicf("Failed to read packages in current directory: %s", err.Error())
	}
	return p.Name
}

func createErrorTpl(i18nErr *I18nErr) {
	const errorTplFileName = "error.go"
	f, err := os.Create(errorTplFileName)
	if err != nil {
		panicf("create file: %s err: %s", errorTplFileName, err.Error())
	}
	defer f.Close()
	tmp, err := template.New("error").Parse(errTpl)
	if err != nil {
		panicf("parse template err:%s", err.Error())
	}
	err = tmp.Execute(f, i18nErr)
	if err != nil {
		panicf("execute template err:%s", err.Error())
	}
	_ = f.Sync()
}

func main() {
	i18nErr := new(I18nErr)
	vf := reflect.Indirect(reflect.ValueOf(i18nErr))
	numField := vf.NumField() - 1
	args := os.Args
	if len(args)-1 != numField {
		panicf("file num length should is %d", numField)
	}
	for i := 0; i < numField; i++ {
		vf.Field(i).SetString(args[1+i])
	}

	i18nErr.GenPackage = findPackageName()

	tmp, err := template.New("tpl").Parse(tpl)
	if err != nil {
		panicf("new template err: %v", err)
	}

	outFileName := fmt.Sprintf("%s.go", i18nErr.Code)
	f, _ := os.Create(outFileName)
	defer f.Close()

	err = tmp.Execute(f, i18nErr)
	if err != nil {
		panicf("template excute err: %v", err)
	}

	createErrorTpl(i18nErr)

}
