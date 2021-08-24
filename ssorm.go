package ssorm

import (
	"cloud.google.com/go/spanner"
	"context"
	"errors"
	"fmt"
	"github.com/10antz-inc/ssorm/utils"
	"google.golang.org/api/iterator"
	"os"
	"reflect"
	"sync"
)

type Option func(*DB)

func Logger(l ILogger) Option {
	return func(d *DB) {
		d.logger = l
	}
}

type DB struct {
	mu      sync.RWMutex
	builder *Builder
	logger  ILogger
}

func CreateDB(opts ...Option) *DB {
	db := &DB{mu: sync.RWMutex{}}
	for _, opt := range opts {
		opt(db)
	}
	if db.logger == nil {
		db.logger = NewLogger(os.Stdout)
	}
	return db
}

func (db *DB) Build() *DB {
	db.builder = &Builder{}
	return db
}

func (db *DB) Model(model interface{}) *DB {
	db.builder.setModel(model)
	return db
}

func (db *DB) Select(query []string, args ...interface{}) *DB {
	db.builder.setSelects(query, args)
	return db
}

func (db *DB) Offset(offset int64) *DB {
	db.builder.setOffset(offset)
	return db
}

func (db *DB) Order(order string) *DB {
	db.builder.setOrder(order)
	return db
}

func (db *DB) Where(query interface{}, values ...interface{}) *DB {
	db.builder.setWhere(query, values)
	return db
}

func (db *DB) Limit(limit int64) *DB {
	db.builder.limit = limit
	return db
}

func (db *DB) DeleteModel(model interface{}, ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	query, err := db.Model(model).builder.deleteModelQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Delete Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) DeleteWhere(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	query, err := db.builder.deleteWhereQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Delete Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) First(model interface{}, ctx context.Context, spannerTransaction interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.builder.limit = 1
	query, _ := db.Model(model).builder.selectQuery()

	var (
		err  error
		iter *spanner.RowIterator
		row  *spanner.Row
	)

	stmt := spanner.Statement{SQL: query}
	db.logger.Infof("Select Query: %s", stmt.SQL)

	rot, readOnly := spannerTransaction.(*spanner.ReadOnlyTransaction)
	rwt, readWrite := spannerTransaction.(*spanner.ReadWriteTransaction)
	if readOnly {
		iter = rot.Query(ctx, stmt)
	}
	if readWrite {
		iter = rwt.Query(ctx, stmt)
	}

	defer iter.Stop()
	for {
		if row, err = iter.Next(); err != nil {
			if err == iterator.Done {
				fmt.Printf("Result: %+v", model)
				return nil
			}
			return err
		}
		row.ToStruct(model)
		break
	}

	return err
}

func (db *DB) Count(cnt interface{}, ctx context.Context, spannerTransaction interface{}) error {
	var (
		err  error
		iter *spanner.RowIterator
		row  *spanner.Row
	)
	if db.builder.tableName == "" {
		return errors.New("Undefined table name. please set db.Model(&struct{})")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	query, err := db.Select([]string{"COUNT(1) AS CNT"}).builder.selectQuery()

	stmt := spanner.Statement{SQL: query}
	db.logger.Infof("Select Query: %s", stmt.SQL)

	rot, readOnly := spannerTransaction.(*spanner.ReadOnlyTransaction)
	rwt, readWrite := spannerTransaction.(*spanner.ReadWriteTransaction)
	if readOnly {
		iter = rot.Query(ctx, stmt)
	}
	if readWrite {
		iter = rwt.Query(ctx, stmt)
	}

	defer iter.Stop()
	for {
		if row, err = iter.Next(); err != nil {
			if err == iterator.Done {
				return nil
			}
			return err
		}
		row.ColumnByName("CNT", cnt)
		break
	}

	return err
}

func (db *DB) Find(model interface{}, ctx context.Context, spannerTransaction interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	var (
		err  error
		iter *spanner.RowIterator
		row  *spanner.Row
	)

	var (
		isSlice, isPtr bool
	)

	query, err := db.Model(model).builder.selectQuery()
	stmt := spanner.Statement{SQL: query}
	db.logger.Infof("Select Query: %s", stmt.SQL)

	rot, readOnly := spannerTransaction.(*spanner.ReadOnlyTransaction)
	rwt, readWrite := spannerTransaction.(*spanner.ReadWriteTransaction)
	if readOnly {
		iter = rot.Query(ctx, stmt)
	}
	if readWrite {
		iter = rwt.Query(ctx, stmt)
	}

	defer iter.Stop()

	results := utils.Indirect(reflect.ValueOf(model))
	var resultType reflect.Type
	if kind := results.Kind(); kind == reflect.Slice {
		isSlice = true
		resultType = results.Type().Elem()

		results.Set(reflect.MakeSlice(results.Type(), 0, 0))

		if resultType.Kind() == reflect.Ptr {
			isPtr = true
			resultType = resultType.Elem()
		}
	}

	for {
		if row, err = iter.Next(); err != nil {
			if err == iterator.Done {
				return nil
			}
			return err
		}
		results := utils.Indirect(reflect.ValueOf(model))
		elem := reflect.New(resultType).Interface()
		row.ToStruct(elem)

		if isSlice {
			if isPtr {
				results.Set(reflect.Append(results, reflect.ValueOf(elem).Elem().Addr()))
			} else {
				results.Set(reflect.Append(results, reflect.ValueOf(elem).Elem()))
			}
		}
	}
}

func (db *DB) Insert(model interface{}, ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Model(model)
	query, err := db.Model(model).builder.insertModelQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Update Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) Update(model interface{}, ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	query, err := db.Model(model).builder.updateModelQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Update Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) UpdateMap(model interface{}, ctx context.Context, in map[string]interface{}, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	query, err := db.Model(model).builder.updateMapQuery(in)
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Update Query: %s", db.builder.query)
	return rowCount, err
}
