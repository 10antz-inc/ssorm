package tests

type Singers struct {
	SingerId  int64 `key:"primary"`
	FirstName string
	LastName  string
}

type Singer struct {
	SingerId  int64 `key:"primary"`
	FirstName string
	LastName  string
	Albums    []*Albums
	Concerts  []*Concerts
}

type Albums struct {
	SingerId int64
	AlbumId  int64
	Title    string
}
type Concerts struct {
	SingerId  int64
	ConcertId int64
	Price     int64
}
