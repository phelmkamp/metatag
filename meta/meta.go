package meta

import (
	"fmt"
	"strings"
)

const (
	fileTemplate   = "package %s\n%s%s"
	methodTemplate = "func (%s %s) %s(%s)%s{\n%s\n}"
	funcTemplate   = "func(%s)%s"

	getterTemplate   = "\treturn %s.%s"
	setterTemplate   = "\t%s.%s = %s"
	finderTemplate   = "\tfor i := range %s.%s {\n\t\tif reflect.DeepEqual(%s.%s[i], %s) {\n\t\t\treturn i\n\t\t}\n\t}\n\treturn -1"
	filtererTemplate = "\tfound := make([]%s, 0, len(%s.%s))\n\tfor i := range %s.%s {\n\t\tif fn(%s.%s[i]) {\n\t\t\tfound = append(found, %s.%s[i])\n\t\t}\n\t}\n\treturn found"
)

type File struct {
	Package string
	Imports Imports
	Methods Methods
}

func NewFile(pkg string) File {
	return File{
		Package: pkg,
		Imports: make(map[string]struct{}),
	}
}

func (f File) String() string {
	return fmt.Sprintf(fileTemplate, f.Package, f.Imports, f.Methods)
}

type Imports map[string]struct{}

func (is Imports) String() string {
	if len(is) < 1 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("\nimport (")
	for k := range is {
		sb.WriteString("\n\t\"")
		sb.WriteString(k)
		sb.WriteString("\"\n")
	}
	sb.WriteString(")\n")
	return sb.String()
}

type Method struct {
	RcvName string
	RcvType string
	Name    string
	RetVals string
	Body    string
	args    string
}

func NewGetter(rcvName, rcvType, name, retType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " " + retType + " ",
		Body:    fmt.Sprintf(getterTemplate, rcvName, field),
	}
}

func NewSetter(rcvName, rcvType, name, argName, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " ",
		Body:    fmt.Sprintf(setterTemplate, rcvName, field, argName),
		args:    argName + " " + argType,
	}
}

func NewFinder(rcvName, rcvType, name, argName, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " int ",
		Body:    fmt.Sprintf(finderTemplate, rcvName, field, rcvName, field, argName),
		args:    argName + " " + argType,
	}
}

func NewFilterer(rcvName, rcvType, name, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " []" + argType + " ",
		Body:    fmt.Sprintf(filtererTemplate, argType, rcvName, field, rcvName, field, rcvName, field, rcvName, field),
		args:    "fn " + fmt.Sprintf(funcTemplate, argType, " bool"),
	}
}

func (m Method) String() string {
	return fmt.Sprintf(methodTemplate, m.RcvName, m.RcvType, m.Name, m.args, m.RetVals, m.Body)
}

type Methods []Method

func (ms Methods) String() string {
	sb := strings.Builder{}
	for i := range ms {
		sb.WriteString("\n")
		sb.WriteString(ms[i].String())
		sb.WriteString("\n")
	}
	return sb.String()
}
