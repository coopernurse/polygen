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
	ReturnType PolyType
}

type Property struct {
	Name  string
	Type  PolyType
}

type Visitor struct { 
	pkg *Package
	lastName string
	state scanState
	errors []PolyError
}

func (v *Visitor) AddErr(e *PolyError) {
	v.errors = append(v.errors, *e)
}

func (v Visitor) Error() string {
	err := ""
	for i := 0; i < len(v.errors); i++ {
		if i > 0 {
			err += "\n"
		}
		err += v.errors[i].Error()
	}
	return err
}

type PolyType struct {
	GoType string
    MapKeyType string
	IsVoid bool
	IsMap  bool
    IsList bool
}

func NewVoidPolyType() PolyType {
	return PolyType{"", "", true, false, false}
}

func NewPolyTypeFromField(f *ast.Field) (PolyType, *PolyError) {
	switch t := f.Type.(type) {
	case *ast.MapType:
		kname := fmt.Sprintf("%v", t.Key)
		vname := fmt.Sprintf("%v", t.Value)
		return PolyType{vname, kname, false, true, false}, nil
	case *ast.ArrayType:
		tname := fmt.Sprintf("%v", t.Elt)
		return PolyType{tname, "", false, false, true}, nil
	case *ast.SliceExpr:
		tname := fmt.Sprintf("%v", t.X)
		return PolyType{tname, "", false, false, true}, nil
	//default:
	//	fmt.Printf("NewPoly. type: %v\n", t)
	}

	tname := fmt.Sprintf("%v", f.Type)
	return PolyType{tname, "", false, false, false}, nil
}

func NewPolyTypeFromGoType(gotype string) (PolyType, *PolyError) {
	return PolyType{gotype, "", false, false, false}, nil
}

type PolyError struct {
	Line    int
	Message string
}

func (e PolyError) Error() string {
	return fmt.Sprintf("polygen: line: %d err: %s", e.Line, e.Message)
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
							ptype, err := NewPolyTypeFromField(fields[x])
							if err == nil {
								prop := Property{fields[x].Names[0].Name, ptype}
								meth.Args = append(meth.Args, prop)
							} else {
								v.AddErr(err)
							}
						}
					}
				} else {
					if len(fields) > 0 {
						rtype, err := NewPolyTypeFromField(fields[0])
						if err == nil {
							meth.ReturnType = rtype
						} else {
							v.AddErr(err)
						}
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
				ptype, err := NewPolyTypeFromField(t)
				if err == nil {
					tmp.Props = append(tmp.Props, Property{t.Names[0].Name, ptype})
				} else {
					v.AddErr(err)
				}
				return nil
			}
		case INTERFACE:
			if len(t.Names) > 0 {
				tmp := &v.pkg.Interfaces[len(v.pkg.Interfaces)-1]
				m := Method{t.Names[0].Name, nil, NewVoidPolyType()}
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

	v := &Visitor{&Package{}, "", STRUCT, make([]PolyError,0)}
	v.pkg.Name = af.Name.Name
	v.pkg.Structs = []Struct{}
	v.pkg.Interfaces = []Interface{}
	ast.Walk(v, af)

	if len(v.errors) > 0 {
		return nil, v
	}
	return v.pkg, nil
}

