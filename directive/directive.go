package directive

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/phelmkamp/metatag/meta"
)

const (
	optOmitField = "omitfield"
	optStringer  = "stringer"
	optReflect   = "reflect"
)

var (
	runFuncs = map[string]runFunc{
		"ptr":      ptr,
		"getter":   getter,
		"setter":   setter,
		"filter":   filter,
		"mapper":   mapper,
		"sort":     sort,
		"stringer": stringer,
		"new":      runNew,
		"equal":    equal,
	}
)

// Target represents the target of the directive.
type Target struct {
	MetaFile         *meta.File
	RcvName, RcvType string
	FldNames         []string
	FldType          string
}

type runFunc func(*Target, []string)

// RunAll runs all of the given directives.
func RunAll(ds []string, tgt *Target) {
	for i := range ds {
		Run(ds[i], tgt)
	}
}

// Run runs the given directive.
func Run(d string, tgt *Target) {
	opts := strings.Split(d, ",")
	d = opts[0]
	opts = opts[1:]

	run, ok := runFuncs[d]
	if !ok {
		log.Printf("Unknown directive: %s\n", d)
		return
	}

	run(tgt, opts)
}

// ptr converts the receiver to a pointer for all subsequent directives.
func ptr(tgt *Target, opts []string) {
	tgt.RcvType = "*" + tgt.RcvType
	log.Printf("Using pointer receiver: %s\n", tgt.RcvType)
}

// getter generates a getter method for each name of the given field.
func getter(tgt *Target, opts []string) {
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

// setter generates a setter method for each name of the given field.
func setter(tgt *Target, opts []string) {
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

// filter generates a filter method for each name of the given field.
func filter(tgt *Target, opts []string) {
	elemType := strings.TrimPrefix(tgt.FldType, "[]")

	var isOmitField bool
	for i := range opts {
		if opts[i] == optOmitField {
			isOmitField = true
			break
		}
	}

	for _, fldNm := range tgt.FldNames {

		method := "Filter"
		if !isOmitField {
			method += upperFirst(fldNm)
		}

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

// mapper generates a mapper method for each name of the given field.
func mapper(tgt *Target, opts []string) {
	if len(opts) < 1 {
		log.Print("skipping 'mapper' - must specify target type as first option\n")
		return
	}

	result := opts[0]
	opts = opts[1:]

	sel := result
	if resSubs := strings.SplitN(result, ".", 2); len(resSubs) > 1 {
		sel = resSubs[1]
	}

	elemType := strings.TrimPrefix(tgt.FldType, "[]")

	var isOmitField bool
	for i := range opts {
		if opts[i] == optOmitField {
			isOmitField = true
			break
		}
	}

	for _, fldNm := range tgt.FldNames {
		var fldPart string
		if !isOmitField {
			fldPart = upperFirst(fldNm)
		}
		method := fmt.Sprintf("Map%sTo%s", fldPart, upperFirst(sel))

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

// sort generates sort methods for the first name of the given field.
func sort(tgt *Target, opts []string) {
	log.Print("Adding import: \"sort\"\n")
	tgt.MetaFile.Imports["sort"] = struct{}{}

	fldNm := tgt.FldNames[0]

	log.Println("Adding method: Len")
	log.Println("Adding method: Swap")
	log.Println("Adding method: Sort")
	sort := meta.Method{
		RcvName: tgt.RcvName,
		RcvType: tgt.RcvType,
		FldName: fldNm,
		Tmpl:    "sort",
	}
	tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, &sort)

	var isStringer bool
	for i := range opts {
		if opts[i] == optStringer {
			isStringer = true
			break
		}
	}
	if isStringer {
		log.Println("Adding method: Less")
		less := meta.Method{
			RcvName: tgt.RcvName,
			RcvType: tgt.RcvType,
			FldName: fldNm,
			Tmpl:    "less",
		}
		tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, &less)
	}
}

// stringer adds each name of the given field to the String() implementation.
func stringer(tgt *Target, opts []string) {
	log.Print("Adding import: \"fmt\"\n")
	tgt.MetaFile.Imports["fmt"] = struct{}{}

	for _, fldNm := range tgt.FldNames {
		log.Print("Adding to method: String\n")
		found := tgt.MetaFile.FilterMethods(
			func(m *meta.Method) bool {
				return m.RcvName == tgt.RcvName && m.RcvType == tgt.RcvType && m.Name == "String"
			},
			1,
		)
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

// runNew adds each name of the given field to the New() implementation.
func runNew(tgt *Target, opts []string) {
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
		new.Misc["Fields"] = fmt.Sprintf("%s%s: %s", fields, fldNm, arg) + ","
	}
}

// equal adds each name of the given field to the Equal() implementation.
func equal(tgt *Target, opts []string) {
	for _, fldNm := range tgt.FldNames {
		log.Print("Adding to method: Equal\n")
		found := tgt.MetaFile.FilterMethods(
			func(m *meta.Method) bool {
				return m.RcvName == tgt.RcvName && m.RcvType == tgt.RcvType && m.Name == "Equal"
			},
			1,
		)
		var cmps string
		var equal *meta.Method
		if len(found) > 0 {
			equal = found[0]
			cmps = equal.Misc["Cmps"].(string) + "\n\t"
		} else {
			equal = &meta.Method{
				RcvName: tgt.RcvName,
				RcvType: tgt.RcvType,
				Name:    "Equal",
				Misc:    make(map[string]interface{}),
				Tmpl:    "equal",
			}
			tgt.MetaFile.Methods = append(tgt.MetaFile.Methods, equal)
		}
		var cmp string
		if len(opts) > 0 && opts[0] == optReflect {
			log.Print("Adding import: \"reflect\"\n")
			tgt.MetaFile.Imports["reflect"] = struct{}{}
			cmp = fmt.Sprintf(
				"if !reflect.DeepEqual(%s.%s, %s2.%s) {\n\t\treturn false\n\t}",
				tgt.RcvName, fldNm, tgt.RcvName, fldNm,
			)
		} else {
			cmp = fmt.Sprintf(
				"if %s.%s != %s2.%s {\n\t\treturn false\n\t}",
				tgt.RcvName, fldNm, tgt.RcvName, fldNm,
			)
		}
		equal.Misc["Cmps"] = cmps + cmp
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
