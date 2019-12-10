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

			metaFile := meta.NewFile(astFile.Name.Name)

			ast.Inspect(astFile, func(n ast.Node) bool {
				var expr ast.Expr
				var typNm string
				switch nt := n.(type) {
				case *ast.TypeSpec:
					expr = nt.Type
					typNm = nt.Name.Name
				}

				if expr == nil {
					return true
				}

				st, ok := expr.(*ast.StructType)
				if !ok {
					return true
				}

				log.Printf("Found struct: %s\n", typNm)

				rcv, _ := first(typNm)
				rcv = strings.ToLower(rcv)

				for _, f := range st.Fields.List {
					if f.Tag == nil {
						continue
					}

					metaTag := metaTagRegEx.FindString(f.Tag.Value)
					if metaTag == "" {
						return true
					}

					log.Printf("Found meta tag %s\n", metaTag)
					metaTag = strings.TrimPrefix(metaTag, "meta:\"")
					metaTag = strings.TrimSuffix(metaTag, "\"")

					var fldType string
					switch ft := f.Type.(type) {
					case *ast.Ident:
						fldType = ft.Name
					case *ast.ArrayType:
						fldType = "[]"
						fldType += ft.Elt.(*ast.Ident).Name
					case *ast.MapType:
						fldType = fmt.Sprintf("map[%s]%s", ft.Key.(*ast.Ident).Name, ft.Value.(*ast.Ident).Name)
					default:
						log.Printf("Unsupported field type: %v\n", ft)
						return true
					}

					elemType := strings.TrimPrefix(fldType, "[]")

					rcvType := typNm

					directives := strings.Split(metaTag, ";")
					for _, d := range directives {
						switch d {
						case "ptr":
							rcvType = directive.Ptr(typNm)
						case "getter":
							directive.Getter(&metaFile, rcv, rcvType, fldType, f)
						case "setter":
							directive.Setter(&metaFile, rcv, rcvType, elemType, fldType, f)
						case "find":
							directive.Find(&metaFile, rcv, rcvType, elemType, fldType, typNm, f)
						case "filter":
							directive.Filter(&metaFile, rcv, rcvType, elemType, fldType, typNm, f)
						default:
							log.Printf("Unknown directive: %s\n", d)
						}
					}
				}

				return true
			})

			if len(metaFile.Methods) < 1 {
				return nil
			}

			osFile := initFile(cleanPath)
			defer func() {
				if err := osFile.Close(); err != nil {
					log.Printf("File.Close() failed: %v\n", err)
				}
			}()

			if _, err := osFile.WriteString(metaFile.String()); err != nil {
				log.Fatalf("File.WriteString() failed: %v\n", err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
