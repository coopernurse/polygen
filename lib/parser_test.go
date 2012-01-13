package polygenlib

import (
	"testing"
	"reflect"
)

var example1 = `package foolib

type JobTitle string
const (
    BOSS JobTitle = "Boss"
    WORKER = "Worker"
)

type Result struct {
	Success	bool
	Code	int
	Note	string
}

type Person struct {
	Id int
	Name string
	Email string "pattern: \\S+@\\S+.\\S+"
	Title string
}

type SampleService interface {
	Create(p Person) Result
	Add(a int, b int) int
    StoreName(name string)
    Say_Hi() string
}`

func TestExample(t *testing.T) {
	pkg, err := Parse(example1)
	if err != nil {
		t.Fatal(err)
	}
	structs := []Struct{
		Struct{"Result", []Property{
				Property{"Success", "bool"},
				Property{"Code", "int"},
				Property{"Note", "string"},
		}},
		Struct{"Person", []Property{
				Property{"Id", "int"},
				Property{"Name", "string"},
				Property{"Email", "string"},
				Property{"Title", "string"},
		}},
	}
	ifaces := []Interface{
		Interface{"SampleService", []Method{
				Method{"Create", []Property{
						Property{"p", "Person"},
					}, "Result"},
				Method{"Add", []Property{
						Property{"a", "int"},
						Property{"b", "int"},
					}, "int"},
				Method{"StoreName", []Property{
						Property{"name", "string"}}, ""},
				Method{"Say_Hi", []Property{}, "string"},
		}},
	}
	expected := Package{"foolib", structs, ifaces}
	if !reflect.DeepEqual(expected.Interfaces, pkg.Interfaces) {
		t.Errorf("%v != %v", expected.Interfaces, pkg.Interfaces)
	}
}

// JSON Schema generation
//   

// Java generation
//
//  create a .java file per struct
//  create two .java files per interface:
//    - java interface 
//    - java json-rpc server class, which dispatches to an impl of interface,
//        but also validates inputs against the schema
//    - java json-rpc client class, which makes requests to server

// To try:
//    map[string] Foo
//    map[int] []Foo

// Validation tests
//
// Disallow func 
// Disallow types that are not int/float/string/bool or in file
// Disallow multiple return values on interface methods
// Disallow variadics
// Disallow nested maps
// Disallow interface args without names