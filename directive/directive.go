package directive

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/phelmkamp/metatag/meta"
)

// Ptr converts the receiver to a pointer for all subsequent directives
func Ptr(typNm string) string {
	ptrTyp := "*" + typNm
	log.Printf("Using pointer receiver: %s\n", ptrTyp)
	return ptrTyp
}

// Getter generates a getter method for each name of the given field
func Getter(metaFile *meta.File, rcv, rcvType, fldType string, f *ast.Field) {
	for _, fldNm := range f.Names {
		method := upperFirst(fldNm.Name)
		if method == fldNm.Name {
			method = "Get" + method
		}

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewGetter(rcv, rcvType, method, fldType, fldNm.Name))
	}
}

// Setter generates a setter method for each name of the given field
func Setter(metaFile *meta.File, rcv, rcvType, elemType, fldType string, f *ast.Field) {
	arg := argName(rcv, elemType)

	for _, fldNm := range f.Names {
		method := "Set" + upperFirst(fldNm.Name)

		ptrRcvType := rcvType
		if !strings.HasPrefix(rcvType, "*") {
			ptrRcvType = "*" + rcvType
		}

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewSetter(rcv, ptrRcvType, method, arg, fldType, fldNm.Name))
	}
}

// Filter generates a filter method for each name of the given field
func Filter(metaFile *meta.File, rcv, rcvType, elemType, fldType, typNm string, f *ast.Field) {
	arg, _ := first(elemType)
	arg = strings.ToLower(arg)

	for _, fldNm := range f.Names {

		method := "Filter" + upperFirst(fldNm.Name)

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewFilter(rcv, rcvType, method, elemType, fldNm.Name, fldType))
	}
}

// Map generates a mapper method for each name of the given field
func Map(metaFile *meta.File, rcv, rcvType, elemType, fldType, typNm, target string, f *ast.Field) {
	for _, fldNm := range f.Names {
		if elemType == fldType {
			log.Printf("'map' not valid for field %s.%s - must be a slice\n", typNm, fldNm)
			continue
		}

		targetType := target
		if tgtSubs := strings.SplitN(target, ".", 2); len(tgtSubs) > 1 {
			targetType = tgtSubs[1]
		}
		method := fmt.Sprintf("Map%sTo%s", upperFirst(fldNm.Name), upperFirst(targetType))

		arg, _ := first(elemType)
		arg = strings.ToLower(arg)

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewMapper(rcv, rcvType, method, elemType, fldNm.Name, target))
	}
}

// Stringer adds each name of the given field to the String() implementation
func Stringer(metaFile *meta.File, rcv, rcvType, fldType string, f *ast.Field) {
	log.Print("Adding import: \"fmt\"\n")
	metaFile.Imports["fmt"] = struct{}{}

	for _, fldNm := range f.Names {
		log.Print("Adding to method: String\n")
		found := metaFile.FilterMethods(func(m *meta.Method) bool { return m.Name == "String" }, 1)
		var format, a string
		var stringer *meta.Method
		if len(found) > 0 {
			stringer = found[0]
			format = stringer.Misc["Format"].(string) + ", "
			a = stringer.Misc["A"].(string) + ", "
		} else {
			stringer = meta.NewStringer(rcv, rcvType)
			metaFile.Methods = append(metaFile.Methods, stringer)
		}
		stringer.Misc["Format"] = fmt.Sprintf("%s%s: %%v", format, fldNm)
		stringer.Misc["A"] = fmt.Sprintf("%s%s.%s", a, rcv, fldNm)
	}
}

// New adds each name of the given field to the New() implementation
func New(metaFile *meta.File, rcvType, fldType string, f *ast.Field) {
	method := "New" + upperFirst(rcvType)
	for _, fldNm := range f.Names {
		log.Printf("Adding to method: %s\n", method)
		found := metaFile.FilterMethods(func(m *meta.Method) bool { return m.Name == method }, 1)
		var args, fields string
		var new *meta.Method
		if len(found) > 0 {
			new = found[0]
			args = new.Misc["Args"].(string) + ", "
			fields = new.Misc["Fields"].(string) + "\n\t\t"
		} else {
			new = meta.NewNew(rcvType, method)
			metaFile.Methods = append(metaFile.Methods, new)
		}

		arg := lowerFirst(fldNm.Name)
		new.Misc["Args"] = fmt.Sprintf("%s%s %s", args, arg, fldType)
		new.Misc["Fields"] = fmt.Sprintf("%s%s: %s", fields, fldNm.Name, arg) + ", "
	}
}

func first(s string) (string, int) {
	if s == "" {
		return "", 0
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(r), n
}

func lowerFirst(s string) string {
	f, n := first(s)
	return strings.ToLower(f) + s[n:]
}

func upperFirst(s string) string {
	f, n := first(s)
	return strings.ToUpper(f) + s[n:]
}

func argName(rcv, argType string) string {
	subs := strings.Split(argType, ".")
	arg, _ := first(subs[len(subs)-1])
	arg = strings.ToLower(arg)
	if arg == rcv {
		// just double up
		arg += arg
	}
	return arg
}
