SSORM
=========

SSORM is a simple spanner orm

Inspired by [GORM](https://github.com/go-gorm/gorm)

Overview
=========

* Feature ORM
  * Insert 
  * Update (Model and Map)
  * Find
  * First
  * Count
  * Delete (Model and Condition)


Run Test
=========
Config spanner-emulator && create instance && create database && insert record
```
. ./tests/ddl/create_datbase.sh
```

Run test
=========
```
go test -v ./tests/...
```


Custom Logger
=========

```go
ssorm.CreateDB(sorm.Logger({custom logger}))
```

Logger Interface
=========

```go
type ILogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}
```

## License

@m-mine @k-haga, 2021~time.Now

Released under the [MIT License](https://github.com/go-gorm/gorm/blob/master/License)
