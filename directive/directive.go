package directive

import (
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
	for _, fldNm := range f.Names {
		if elemType == fldType {
			log.Printf("'filter' not valid for field %s.%s - must be a slice\n", typNm, fldNm)
			continue
		}

		method := "Filter" + upperFirst(fldNm.Name)

		arg, _ := first(elemType)
		arg = strings.ToLower(arg)

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewFilter(rcv, rcvType, method, elemType, fldNm.Name))
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
