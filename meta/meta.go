package meta

import (
	"fmt"
	"strings"
)

const (
	fileTemplate     = "package %s\n%s%s"
	methodTemplate   = "func (%s %s) %s(%s)%s{\n%s\n}"
	anonFuncTemplate = "func(%s)%s"

	getterBodyTemplate   = "\treturn %s.%s"
	setterBodyTemplate   = "\t%s.%s = %s"
	finderBodyTemplate   = "\tfor i := range %s.%s {\n\t\tif reflect.DeepEqual(%s.%s[i], %s) {\n\t\t\treturn i\n\t\t}\n\t}\n\treturn -1"
	filtererBodyTemplate = "\tfound := make([]%s, 0, len(%s.%s))\n\tfor i := range %s.%s {\n\t\tif fn(%s.%s[i]) {\n\t\t\tfound = append(found, %s.%s[i])\n\t\t}\n\t}\n\treturn found"
)

// File represents a generated code file
type File struct {
	Package string
	Imports Imports
	Methods Methods
}

// NewFile creates a new File with all fields initialized
func NewFile(pkg string) File {
	return File{
		Package: pkg,
		Imports: make(map[string]struct{}),
	}
}

// String generates the file content
func (f File) String() string {
	return fmt.Sprintf(fileTemplate, f.Package, f.Imports, f.Methods)
}

// Imports represents a set of import paths
type Imports map[string]struct{}

// String generates the import statement
func (is Imports) String() string {
	if len(is) < 1 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("\nimport (\n")
	for k := range is {
		sb.WriteString("\t\"")
		sb.WriteString(k)
		sb.WriteString("\"\n")
	}
	sb.WriteString(")\n")
	return sb.String()
}

// Method represents a generated method
type Method struct {
	RcvName string
	RcvType string
	Name    string
	RetVals string
	Body    string
	args    string
}

// NewGetter creates a new getter metthod
func NewGetter(rcvName, rcvType, name, retType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " " + retType + " ",
		Body:    fmt.Sprintf(getterBodyTemplate, rcvName, field),
	}
}

// NewSetter creates a new setter method
func NewSetter(rcvName, rcvType, name, argName, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " ",
		Body:    fmt.Sprintf(setterBodyTemplate, rcvName, field, argName),
		args:    argName + " " + argType,
	}
}

// NewFinder creates a new finder method
func NewFinder(rcvName, rcvType, name, argName, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " int ",
		Body:    fmt.Sprintf(finderBodyTemplate, rcvName, field, rcvName, field, argName),
		args:    argName + " " + argType,
	}
}

// NewFilterer creates a new filterer method
func NewFilterer(rcvName, rcvType, name, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: " []" + argType + " ",
		Body:    fmt.Sprintf(filtererBodyTemplate, argType, rcvName, field, rcvName, field, rcvName, field, rcvName, field),
		args:    "fn " + fmt.Sprintf(anonFuncTemplate, argType, " bool"),
	}
}

// String generates the method code
func (m Method) String() string {
	return fmt.Sprintf(methodTemplate, m.RcvName, m.RcvType, m.Name, m.args, m.RetVals, m.Body)
}

// Methods represents a collection of generated methods
type Methods []Method

// String generates the code for all methods
func (ms Methods) String() string {
	sb := strings.Builder{}
	for i := range ms {
		sb.WriteString("\n")
		sb.WriteString(ms[i].String())
		sb.WriteString("\n")
	}
	return sb.String()
}
