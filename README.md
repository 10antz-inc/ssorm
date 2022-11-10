SSORM
=========

SSORM is a simple spanner orm for Golang

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
    * SimpleQueryRead
    * SimpleQueryWrite

* Supported data type
    * STRING
    * INT64
    * FLOAT64
    * ARRAY
    * BOOL
    * DATE
    * TIMESTAMP

* SSORM Tag
    * primary
    * create_time
    * update_time
    * delete_time
    * ignore_write
    * nullable_write

Test
=========

* Configure spanner-emulator && create instance && create database && insert record
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
ssorm.Logger({custom logger})
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
WithContext(ctx context.Context) *logrus.Entry
}
```

Tracing
=========

```go
ssorm.UseTracing()
```

Additional Tracing Option
=========
```go
import (
  "go.opentelemetry.io/otel/trace"

  "github.com/10antz-inc/ssorm"
  "github.com/10antz-inc/ssorm/ssormotel"
)

func main() {
  ssorm.UseTrace(
    // add attribute 'db.statement'
    ssormotel.WithQueryStatement(),
    // set the created TraceProvider.
    ssormotel.WithTraceProvider(provider trace.TraceProvider),
    // add any attribute
    ssormotel.WithAttributes(
      semconv.DBConnectionStringKey.String("projects/.../instances/.../databases/..."),
      semconv.DBSystemKey.String("Google Cloud Spanner"),
    )
  )
}

```

## License

Copyright (c) 2021 10ANTZ, Inc.

SSORM is released under the [MIT License](https://github.com/10antz-inc/ssorm/blob/master/LICENSE)
