package directive

import (
	"go/ast"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/phelmkamp/metatag/meta"
)

// type Directive func()

// func New(name string) Directive {
// 	switch name {
// 	case "getter":
// 		return getter
// 	case "setter":
// 		return setter
// 	case "filter":
// 		return filter
// 	case "find":
// 		return find
// 	}
// 	return nil
// }

func Ptr(typNm string) string {
	ptrTyp := "*" + typNm
	log.Printf("Using pointer receiver: %s\n", ptrTyp)
	return ptrTyp
}

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

func Setter(metaFile *meta.File, rcv, rcvType, elemType, fldType string, f *ast.Field) {
	for _, fldNm := range f.Names {
		method := "Set" + upperFirst(fldNm.Name)

		arg, _ := first(elemType) // lowerFirst(fldNm.Name)
		arg = strings.ToLower(arg)

		ptrRcvType := rcvType
		if !strings.HasPrefix(rcvType, "*") {
			ptrRcvType = "*" + rcvType
		}

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewSetter(rcv, ptrRcvType, method, arg, fldType, fldNm.Name))
	}
}

func Filter(metaFile *meta.File, rcv, rcvType, elemType, fldType, typNm string, f *ast.Field) {
	for _, fldNm := range f.Names {
		if elemType == fldType {
			log.Printf("'filter' not valid for field %s.%s - must be a slice\n", typNm, fldNm)
			continue
		}

		method := "Filter" + upperFirst(fldNm.Name)

		arg, _ := first(elemType) // lowerFirst(fldNm.Name)
		arg = strings.ToLower(arg)

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewFilterer(rcv, rcvType, method, elemType, fldNm.Name))
	}
}

func Find(metaFile *meta.File, rcv, rcvType, elemType, fldType, typNm string, f *ast.Field) {
	log.Print("Adding import: \"reflect\"\n")
	metaFile.Imports["reflect"] = struct{}{}

	for _, fldNm := range f.Names {
		if elemType == fldType {
			log.Printf("'find' not valid for field %s.%s - must be a slice\n", typNm, fldNm)
			continue
		}

		method := "Find" + upperFirst(strings.TrimSuffix(fldNm.Name, "s"))

		arg, _ := first(elemType) // lowerFirst(fldNm.Name)
		arg = strings.ToLower(arg)

		log.Printf("Adding method: %s\n", method)
		metaFile.Methods = append(metaFile.Methods, meta.NewFinder(rcv, rcvType, method, arg, elemType, fldNm.Name))
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
