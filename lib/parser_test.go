package polygenlib

import (
	"testing"
	"reflect"
)

var example1 = `package foolib

type Result struct {
	Success	bool
	Code	int
	Note	string
}

type Person struct {
	Id int
	name string
	email string "pattern: \\S+@\\S+.\\S+"
	title string
    age float
}

type SampleService interface {
	Create(p Person) Result
	Add(a int, b int) int
    StoreName(name string)
    Say_Hi() string
    getPeople(params map[string] string) []Person
}`

func TestParseExample(t *testing.T) {
	pkg, err := Parse(example1)
	if err != nil {
		t.Fatal(err)
	}

	intType := PolyType{"int", "", false, false, false}
	floatType := PolyType{"float", "", false, false, false}
	stringType := PolyType{"string", "", false, false, false}
	boolType := PolyType{"bool", "", false, false, false}
	personType := PolyType{"Person", "", false, false, false}
	resultType := PolyType{"Result", "", false, false, false}

	structs := []Struct{
		Struct{"Result", []Property{
				Property{"Success", boolType},
				Property{"Code", intType},
				Property{"Note", stringType},
		}},
		Struct{"Person", []Property{
				Property{"Id", intType},
				Property{"name", stringType},
				Property{"email", stringType},
				Property{"title", stringType},
				Property{"age", floatType},
		}},
	}
	ifaces := []Interface{
		Interface{"SampleService", []Method{
				Method{"Create", []Property{Property{"p", personType}},
					resultType},
				Method{"Add", []Property{Property{"a", intType}, 
						Property{"b", intType},}, intType},
				Method{"StoreName", []Property{Property{"name", stringType}},
					NewVoidPolyType()},
				Method{"Say_Hi", []Property{}, stringType},
				Method{"getPeople", []Property{
						Property{"params", 
							PolyType{"string", "string", false, true, false}},
					}, PolyType{"Person", "", false, false, true}},
		}},
	}
	expected := Package{"foolib", structs, ifaces}
	if !reflect.DeepEqual(expected.Structs, pkg.Structs) {
		t.Errorf("%v != %v", expected.Structs, pkg.Structs)
	}
	if !reflect.DeepEqual(expected.Interfaces, pkg.Interfaces) {
		t.Errorf("%v != %v", expected.Interfaces, pkg.Interfaces)
	}
}   

// To try:
//    map[string] Foo
//    map[int] []Foo

// Validation tests
//
// Disallow func 
// Disallow var
// Disallow const
// Disallow types that are not int/float/string/bool or in file
// Disallow multiple return values on interface methods
// Disallow variadics
// Disallow nested maps
// Disallow interface args without names