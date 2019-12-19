package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/phelmkamp/metatag/directive"
	"github.com/phelmkamp/metatag/meta"
)

const (
	accessTemplate = "%s.%s"
)

var (
	goFileRegEx  = regexp.MustCompile(`.+\.go$`)
	metaTagRegEx = regexp.MustCompile(`meta:".+"`)
)

func initFile(origPath string) *os.File {
	filename := strings.Replace(origPath, ".go", "_meta.go", 1)
	log.Printf("Creating file: %s\n", filename)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("os.Create() failed: %v\n", err)
	}
	return f
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

func main() {
	var root string
	flag.StringVar(&root, "path", ".", "directory path to scan for *.go files")
	flag.Parse()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && goFileRegEx.MatchString(info.Name()) {
			cleanPath := filepath.Clean(path)

			log.Printf("Parsing file: %s\n", cleanPath)
			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, cleanPath, nil, 0)
			if err != nil {
				return fmt.Errorf("parser.ParseFile() failed: %w", err)
			}

			tgt := directive.Target{
				MetaFile: meta.NewFile(astFile.Name.Name),
			}
			importPaths := make([]string, 0, 8)

			ast.Inspect(astFile, func(n ast.Node) bool {
				var expr ast.Expr
				switch nt := n.(type) {
				case *ast.TypeSpec:
					expr = nt.Type
					tgt.RcvType = nt.Name.Name
				case *ast.ImportSpec:
					p := strings.Trim(nt.Path.Value, `"`)
					importPaths = append(importPaths, p)
				}

				if expr == nil {
					return true
				}

				st, ok := expr.(*ast.StructType)
				if !ok {
					return true
				}

				log.Printf("Found struct: %s\n", tgt.RcvType)

				tgt.RcvName, _ = first(tgt.RcvType)
				tgt.RcvName = strings.ToLower(tgt.RcvName)

				for _, f := range st.Fields.List {
					if f.Tag == nil {
						continue
					}

					metaTag := metaTagRegEx.FindString(f.Tag.Value)
					if metaTag == "" {
						continue
					}

					log.Printf("Found meta tag %s\n", metaTag)
					metaTag = strings.TrimPrefix(metaTag, "meta:\"")
					metaTag = strings.TrimSuffix(metaTag, "\"")

					// some directives modify target, use a local copy
					fldTgt := tgt

					var fldPkg string
					switch ft := f.Type.(type) {
					case *ast.Ident:
						fldTgt.FldType = ft.Name
					case *ast.SelectorExpr:
						// package.type
						fldPkg = ft.X.(*ast.Ident).Name
						fldTgt.FldType = fmt.Sprintf(accessTemplate, fldPkg, ft.Sel.Name)
					case *ast.ArrayType:
						switch elt := ft.Elt.(type) {
						case *ast.Ident:
							fldTgt.FldType = "[]" + elt.Name
						case *ast.SelectorExpr:
							// package.type
							fldPkg = elt.X.(*ast.Ident).Name
							fldTgt.FldType = fmt.Sprintf(accessTemplate, "[]"+fldPkg, elt.Sel.Name)
						}
					case *ast.MapType:
						fldTgt.FldType = fmt.Sprintf("map[%s]%s", ft.Key.(*ast.Ident).Name, ft.Value.(*ast.Ident).Name)
					default:
						log.Printf("Unsupported field type: %v\n", ft)
						continue
					}

					fldTgt.FldNames = make([]string, len(f.Names))
					for i := range f.Names {
						fldTgt.FldNames[i] = f.Names[i].Name
					}

					directives := strings.Split(metaTag, ";")
					for _, d := range directives {
						dSubs := strings.SplitN(d, ":", 2)
						if len(dSubs) < 1 {
							continue
						}

						opts := strings.Split(dSubs[len(dSubs)-1], ",")
						dSubs[len(dSubs)-1] = opts[0]
						opts = opts[1:]

						switch dSubs[0] {
						case "ptr":
							directive.Ptr(&fldTgt)
						case "getter":
							directive.Getter(&fldTgt)
						case "setter":
							directive.Setter(&fldTgt)
						case "filter":
							directive.Filter(&fldTgt, opts)
						case "map":
							if len(dSubs) < 2 {
								continue
							}
							directive.Map(&fldTgt, dSubs[1], opts)
						case "sort":
							directive.Sort(&fldTgt, opts)
						case "stringer":
							directive.Stringer(&fldTgt)
						case "new":
							directive.New(&fldTgt)
						default:
							log.Printf("Unknown directive: %s\n", d)
						}
					}

					if fldPkg != "" {
						var importPath string
						for _, s := range importPaths {
							subs := strings.Split(s, "/")
							last := subs[len(subs)-1]
							if last == fldPkg {
								importPath = s
								break
							}
						}
						log.Printf("Adding import: \"%s\"\n", importPath)
						tgt.MetaFile.Imports[importPath] = struct{}{}
					}
				}

				return true
			})

			if len(tgt.MetaFile.Methods) < 1 {
				return nil
			}

			osFile := initFile(cleanPath)
			defer func() {
				if err := osFile.Close(); err != nil {
					log.Printf("File.Close() failed: %v\n", err)
				}
			}()

			if _, err := osFile.WriteString(tgt.MetaFile.String()); err != nil {
				log.Fatalf("File.WriteString() failed: %v\n", err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
