package foolib

// Result is a thing that does something
type Result struct {
	// 
	Success bool
	Code    int
	Note    string
}

type Person struct {
	Id    int
	Name  string
	Email string "pattern: \\S+@\\S+.\\S+"
	Title string
}

type SampleService interface {
	Create(p Person) Result
	Add(a int, b int) int
	StoreName(name string)
	Say_Hi() string
}
