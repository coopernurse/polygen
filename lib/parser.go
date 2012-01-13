package polygenlib

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type scanState int
const (
	STRUCT scanState = iota
	INTERFACE
	METHOD
)

type Package struct {
	Name string
	Structs  []Struct
	Interfaces []Interface
}

type Struct struct {
	Name string
	Props []Property
}

type Interface struct {
	Name string
	Methods []Method
}

type Method struct {
	Name string
	Args []Property
	ReturnType string // should be something cooler
}

type Property struct {
	Name  string
	Type  string // should be something more sophisticated..
}

type Visitor struct { 
	pkg *Package
	lastName string
	state scanState
}

func (v *Visitor) Visit(n ast.Node) ast.Visitor {
	//fmt.Printf("  node: type=%v, value=%v\n", reflect.TypeOf(n), n)
	switch t := n.(type) {
	case *ast.TypeSpec:
		v.lastName = t.Name.Name
	case *ast.StructType:
		s := Struct{v.lastName, []Property{}}
		v.pkg.Structs = append(v.pkg.Structs, s)
		v.state = STRUCT
	case *ast.InterfaceType:
		i := Interface{v.lastName, []Method{}}
		v.pkg.Interfaces = append(v.pkg.Interfaces, i)
		v.state = INTERFACE
	case *ast.FieldList:
		if v.state == INTERFACE {
			i := &v.pkg.Interfaces[len(v.pkg.Interfaces)-1]
			if len(i.Methods) > 0 {
				meth := &i.Methods[len(i.Methods)-1]
				fields := t.List
				if meth.Args == nil {
					meth.Args = []Property{}
					for x := 0; x < len(fields); x++ {
						if len(fields[x].Names) > 0 {
							tname := fmt.Sprintf("%v", fields[x].Type)
							prop := Property{fields[x].Names[0].Name, tname}
							meth.Args = append(meth.Args, prop)
						}
					}
				} else {
					if len(fields) > 0 {
						meth.ReturnType = fmt.Sprintf("%v", fields[0].Type)
					}
				}
				
				return nil
			}
		}
	case *ast.Field:
		switch v.state {
		case STRUCT:
			if len(t.Names) > 0 {
				tmp := &v.pkg.Structs[len(v.pkg.Structs)-1]
				tname := fmt.Sprintf("%v", t.Type)
				tmp.Props = append(tmp.Props, Property{t.Names[0].Name, tname})
				return nil
			}
		case INTERFACE:
			if len(t.Names) > 0 {
				tmp := &v.pkg.Interfaces[len(v.pkg.Interfaces)-1]
				m := Method{t.Names[0].Name, nil, ""}
				tmp.Methods = append(tmp.Methods, m)
			}

		}
	}
	return v
}

func Parse(code string) (*Package, error) {
	//lines = strings.Split(code, "\n")
	//for i := 0; i < len(lines); i++ {
	//	f.AddLine(i+1, 
	//}
	fs := &token.FileSet{}
	af, err := parser.ParseFile(fs, "myfile.go", code, 0)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("err=%v\n", err)
	//fmt.Printf("ast=%v\n", af)

	v := &Visitor{&Package{}, "", STRUCT}
	v.pkg.Name = af.Name.Name
	v.pkg.Structs = []Struct{}
	v.pkg.Interfaces = []Interface{}
	ast.Walk(v, af)

	return v.pkg, nil
}

