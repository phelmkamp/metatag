package meta

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
)

const (
	topComment   = "// GENERATED BY metatag, DO NOT EDIT\n// (or edit away - I'm a comment, not a cop)\n"
	fileTemplate = "package %s\n%s%s"
)

var (
	tmplBox = rice.MustFindBox("templates")
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
	return topComment + fmt.Sprintf(fileTemplate, f.Package, f.Imports, f.Methods)
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
	ArgName string
	ArgType string
	RetVals string
	FldName string
	FldType string
	tmpl    string
}

// NewGetter creates a new getter metthod
func NewGetter(rcvName, rcvType, name, retType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		RetVals: retType,
		FldName: field,
		tmpl:    "getter",
	}
}

// NewSetter creates a new setter method
func NewSetter(rcvName, rcvType, name, argName, argType, field string) Method {
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		ArgName: argName,
		ArgType: argType,
		FldName: field,
		tmpl:    "setter",
	}
}

// NewFilter creates a new filter method
func NewFilter(rcvName, rcvType, name, argType, field string) Method {
	fldType := "[]" + argType
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		ArgType: argType,
		RetVals: fldType,
		FldName: field,
		FldType: fldType,
		tmpl:    "filter",
	}
}

// NewMapper creates a new mapper method
func NewMapper(rcvName, rcvType, name, argType, field, target string) Method {
	fldType := "[]" + argType
	return Method{
		RcvName: rcvName,
		RcvType: rcvType,
		Name:    name,
		ArgType: fmt.Sprintf("func(%s) %s", argType, target),
		RetVals: "[]" + target,
		FldName: field,
		FldType: fldType,
		tmpl:    "mapper",
	}
}

// String generates the method code
func (m Method) String() string {
	tmplTxt := tmplBox.MustString(m.tmpl + ".tmpl")
	tmplMessage, err := template.New(m.tmpl).Parse(tmplTxt)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := tmplMessage.Execute(&buf, m); err != nil {
		log.Fatal(err)
	}

	return buf.String()
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
