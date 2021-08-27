SSORM
=========

SSORM is a simple spanner orm

Overview
=========

* Feature
    * Insert (Model)
    * Update (Model and Columns and Params)
    * Find (Model)
    * First (Model)
    * Count (Model)
    * Delete (Model and Where)
    * SubQuery (Model)
    * SoftDelete (Insert Update Find First Count Delete SubQuery)

Test
=========

* Config spanner-emulator && create instance && create database && insert record
    ```
    . ./tests/ddl/create_datbase.sh
    ```

* Run test
    ```
    go test -v ./tests/...
    ```

Custom Logger
=========

```go
ssorm.CreateDB(ssorm.Logger({custom logger}))
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

Copyright 2021 10antz, Inc.

Released under the [MIT License](https://opensource.org/licenses/MIT)
