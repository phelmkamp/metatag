package meta

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("package %s\n%s%s", f.Package, f.Imports, f.Methods)
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
	// Field   string
	args string
	// ret     string
	// assign  string
}

func NewGetter(rcvName, rcvType, name, retType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " " + retType + " ",
		Body:    fmt.Sprintf("\treturn %s.%s", rcvName, field),
		// Field:   field,
		// ret:     "return ",
	}
}

func NewSetter(rcvName, rcvType, name, argName, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " ",
		Body:    fmt.Sprintf("\t%s.%s = %s", rcvName, field, argName),
		// Field:   field,
		args: argName + " " + argType,
		// assign:  " = " + argName,
	}
}

func NewFinder(rcvName, rcvType, name, argName, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " int ",
		Body: fmt.Sprintf("\tfor i := range %s.%s {\n\t\tif reflect.DeepEqual(%s.%s[i], %s) {\n\t\t\treturn i\n\t\t}\n\t}\n\treturn -1",
			rcvName, field, rcvName, field, argName),
		// Field:   field,
		args: argName + " " + argType,
		// assign:  " = " + argName,
	}
}

func NewFilterer(rcvName, rcvType, name, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " []" + argType + " ",
		Body: fmt.Sprintf("\tfound := make([]%s, 0, len(%s.%s))\n\tfor i := range %s.%s {\n\t\tif fn(%s.%s[i]) {\n\t\t\tfound = append(found, %s.%s[i])\n\t\t}\n\t}\n\treturn found",
			argType, rcvName, field, rcvName, field, rcvName, field, rcvName, field),
		// Field:   field,
		args: fmt.Sprintf("fn func(%s) bool", argType),
		// assign:  " = " + argName,
	}
}

func (m Method) String() string {
	return fmt.Sprintf("func (%s %s) %s(%s)%s{\n%s\n}",
		m.RcvName, m.RcvType, m.Name, m.args, m.RetVals, m.Body)
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
