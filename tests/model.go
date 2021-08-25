package tests

import (
	"cloud.google.com/go/spanner"
	"time"
)

type Singers struct {
	SingerId        int64 `ssorm_key:"primary"`
	FirstName       string
	LastName        string
	TestTime        time.Time
	TestSpannerTime spanner.NullTime
	DeleteTime      spanner.NullTime `ssorm_key:"delete_time"`
	CreateTime      time.Time        `ssorm_key:"create_time"`
	UpdateTime      time.Time        `ssorm_key:"update_time"`
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
