package tests

type Singers struct {
	SingerId  int64 `key:"primary"`
	FirstName string
	LastName  string
}

