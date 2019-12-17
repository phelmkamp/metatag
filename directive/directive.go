package directive

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/phelmkamp/metatag/meta"
)

// Target represents the target of the directive
type Target struct {
	MetaFile *meta.File
	RcvName  string
	RcvType  string
	FldNames []string
	FldType  string
}

// Ptr converts the receiver to a pointer for all subsequent directives
func Ptr(tgt *Target) {
	tgt.RcvType = "*" + tgt.RcvType
	log.Printf("Using pointer receiver: %s\n", tgt.RcvType)
}

// Getter generates a getter method for each name of the given field
func Getter(tgt *Target) {
	for _, fldNm := range tgt.FldNames {
		method := upperFirst(fldNm)
		if method == fldNm {
			method = "Get" + method
		}

		log.Printf("Adding method: %s\n", method)
		getter := meta.Method{
			RcvName: tgt.RcvName,
			RcvType: tgt.RcvType,
			Name:    method,
			RetVals: tgt.FldType,
			FldName: fldNm,
			Tmpl:    "getter",
		}
		tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, &getter)
	}
}

// Setter generates a setter method for each name of the given field
func Setter(tgt *Target) {
	elemType := strings.TrimPrefix(tgt.FldType, "[]")

	arg := argName(tgt.RcvName, elemType)

	ptrRcvType := tgt.RcvType
	if !strings.HasPrefix(tgt.RcvType, "*") {
		ptrRcvType = "*" + tgt.RcvType
	}

	for _, fldNm := range tgt.FldNames {
		method := "Set" + upperFirst(fldNm)

		log.Printf("Adding method: %s\n", method)
		setter := meta.Method{
			RcvName: tgt.RcvName,
			RcvType: ptrRcvType,
			Name:    method,
			ArgName: arg,
			ArgType: tgt.FldType,
			FldName: fldNm,
			Tmpl:    "setter",
		}
		tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, &setter)
	}
}

// Filter generates a filter method for each name of the given field
func Filter(tgt *Target) {
	elemType := strings.TrimPrefix(tgt.FldType, "[]")

	for _, fldNm := range tgt.FldNames {

		method := "Filter" + upperFirst(fldNm)

		log.Printf("Adding method: %s\n", method)
		filter := meta.Method{
			RcvName: tgt.RcvName,
			RcvType: tgt.RcvType,
			Name:    method,
			ArgType: elemType,
			RetVals: tgt.FldType,
			FldName: fldNm,
			FldType: tgt.FldType,
			Tmpl:    "filter",
		}
		tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, &filter)
	}
}

// Map generates a mapper method for each name of the given field
func Map(tgt *Target, result string) {
	elemType := strings.TrimPrefix(tgt.FldType, "[]")

	sel := result
	if resSubs := strings.SplitN(result, ".", 2); len(resSubs) > 1 {
		sel = resSubs[1]
	}

	for _, fldNm := range tgt.FldNames {
		if elemType == tgt.FldType {
			log.Printf("'map' not valid for field %s.%s - must be a slice\n", tgt.RcvName, fldNm)
			continue
		}

		method := fmt.Sprintf("Map%sTo%s", upperFirst(fldNm), upperFirst(sel))

		log.Printf("Adding method: %s\n", method)
		mapper := meta.Method{
			RcvName: tgt.RcvName,
			RcvType: tgt.RcvType,
			Name:    method,
			ArgType: fmt.Sprintf("func(%s) %s", elemType, result),
			RetVals: "[]" + result,
			FldName: fldNm,
			Tmpl:    "mapper",
		}
		tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, &mapper)
	}
}

// Stringer adds each name of the given field to the String() implementation
func Stringer(tgt *Target) {
	log.Print("Adding import: \"fmt\"\n")
	tgt.MetaFile.Imports["fmt"] = struct{}{}

	for _, fldNm := range tgt.FldNames {
		log.Print("Adding to method: String\n")
		found := tgt.MetaFile.FilterMethods(func(m *meta.Method) bool { return m.Name == "String" }, 1)
		var format, a string
		var stringer *meta.Method
		if len(found) > 0 {
			stringer = found[0]
			format = stringer.Misc["Format"].(string) + " "
			a = stringer.Misc["A"].(string) + ", "
		} else {
			stringer = &meta.Method{
				RcvName: tgt.RcvName,
				RcvType: tgt.RcvType,
				Name:    "String",
				RetVals: "string",
				Misc:    make(map[string]interface{}),
				Tmpl:    "stringer",
			}
			tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, stringer)
		}
		stringer.Misc["Format"] = fmt.Sprintf("%s%%v", format)
		stringer.Misc["A"] = fmt.Sprintf("%s%s.%s", a, tgt.RcvName, fldNm)
	}
}

// New adds each name of the given field to the New() implementation
func New(tgt *Target) {
	method := "New" + upperFirst(tgt.RcvType)
	for _, fldNm := range tgt.FldNames {
		log.Printf("Adding to method: %s\n", method)
		found := tgt.MetaFile.FilterMethods(func(m *meta.Method) bool { return m.Name == method }, 1)
		var args, fields string
		var new *meta.Method
		if len(found) > 0 {
			new = found[0]
			args = new.Misc["Args"].(string) + ", "
			fields = new.Misc["Fields"].(string) + "\n\t\t"
		} else {
			new = &meta.Method{
				RcvType: tgt.RcvType,
				Name:    method,
				RetVals: tgt.RcvType,
				Misc:    make(map[string]interface{}),
				Tmpl:    "new",
			}
			tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, new)
		}

		arg := lowerFirst(fldNm)
		new.Misc["Args"] = fmt.Sprintf("%s%s %s", args, arg, tgt.FldType)
		new.Misc["Fields"] = fmt.Sprintf("%s%s: %s", fields, fldNm, arg) + ", "
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
