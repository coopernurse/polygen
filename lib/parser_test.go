package polygenlib

import (
	"reflect"
	"strings"
	"testing"
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
	pkg, err := Parse("test.go", example1)
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
				Property{"b", intType}}, intType},
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

// Validation tests
//
// Verify filename in output is correct

func TestErrFilename(t *testing.T) {
	idl := "package foo\nvar blah"
	pkg, err := Parse("example1.go", idl)
	if err == nil {
		t.Error("expected err")
	}
	if pkg != nil {
		t.Error("pkg should be nil")
	}
	if strings.Index(err.Error(), "example1.go:") != 0 {
		t.Errorf("err didn't start with example1.go: - %s", err.Error())
	}
}

// Verify line number in error output is correct
func TestErrLineNum(t *testing.T) {
	idl := "package foo\nimport (\"fmt\")"
	_, err := Parse("example1.go", idl)
	if err == nil {
		t.Fatal("expected err")
	}
	if strings.Index(err.Error(), "example1.go:2") != 0 {
		t.Errorf("err didn't start with example1.go:2: - %s", err.Error())
	}
}

var illegalIdl = []string{
	"const huge = 1 << 100\ntype foo struct {\n foo huge\n}",
	"func foo() int { return 1 }",
	"var foo string\n",
	"type foo interface { }",
	"type foo struct {\n a int64\n}",
	"type foo struct {\n a uint64\n}",
	"type foo struct {\n a float64\n}",
	"type foo struct {\n a map[int] string\n}",
	"type foo struct {\n a map[string] map[string] int\n}",
	"type foo interface {\n doSomething() (int, int)\n}",
	"type foo interface {\n doSomething(int, int) int\n}",
	"type foo interface {\n doSomething(foo ...int) int\n}",
}

func TestIllegalIdl(t *testing.T) {
	for _, s := range illegalIdl {
		s = "package foopkg\n" + s
		_, err := Parse("example1.go", s)
		//fmt.Printf("err=%v\n", err)
		if err == nil {
			t.Errorf("expected err for: %s", s)
		}
	}
}
