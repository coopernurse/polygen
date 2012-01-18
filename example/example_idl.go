package foolib

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
}
