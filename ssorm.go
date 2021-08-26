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
)

type Option func(*DB)

func Logger(l ILogger) Option {
	return func(d *DB) {
		d.logger = l
	}
}

type DB struct {
	builder *Builder
	logger  ILogger
}

func CreateDB(opts ...Option) *DB {
	db := &DB{}
	for _, opt := range opts {
		opt(db)
	}
	if db.logger == nil {
		db.logger = NewLogger(os.Stdout)
	}
	return db
}

func (db *DB) Model(model interface{}) *DB {
	db.builder = &Builder{
		subBuilder: &SubBuilder{},
		model:      model,
		tableName:  utils.GetTableName(model),
		softDelete: false,
	}
	return db
}

func (db *DB) SoftDeleteModel(model interface{}) *DB {
	db.builder = &Builder{
		subBuilder: &SubBuilder{},
		model:      model,
		tableName:  utils.GetTableName(model),
		softDelete: true,
	}
	return db
}

func (db *DB) TableName(tableName string) *DB {
	db.builder.tableName = tableName
	return db
}

func (db *DB) AddSub(model interface{}, query interface{}, values ...interface{}) *DB {
	db.builder.addSub(model, query, values)
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

func (db *DB) DeleteModel(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	var (
		rowCount int64
		err      error
		query    string
	)

	if db.builder.softDelete {
		query, err = db.builder.buildSoftDeleteModelQuery()
		if err != nil {
			return 0, errors.New("no primary key set")
		}
		stmt := spanner.Statement{SQL: query}
		rowCount, err := spannerTransaction.Update(ctx, stmt)
		db.logger.Infof("Update Query: %s", db.builder.query)
		return rowCount, err
	}

	query, err = db.builder.deleteModelQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err = spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Delete Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) DeleteWhere(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	var (
		rowCount int64
		err      error
		query    string
	)
	if db.builder.softDelete {
		query, err = db.builder.buildDeleteWhereQuery()
		if err != nil {
			return 0, err
		}
		stmt := spanner.Statement{SQL: query}
		rowCount, err := spannerTransaction.Update(ctx, stmt)
		db.logger.Infof("Update Query: %s", db.builder.query)
		return rowCount, err
	}

	query, err = db.builder.deleteWhereQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err = spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Delete Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) First(ctx context.Context, spannerTransaction interface{}) error {

	db.builder.limit = 1
	var query string
	if db.builder.subBuilder.subModels != nil {
		query, _ = db.builder.buildSubQuery()
	} else {
		query, _ = db.builder.selectQuery()
	}

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
				fmt.Printf("Result: %+v", db.builder.model)
				return nil
			}
			return err
		}

		row.ToStruct(db.builder.model)
		break
	}

	return err
}

func (db *DB) Count(ctx context.Context, spannerTransaction interface{}, cnt interface{}) error {
	var (
		err  error
		iter *spanner.RowIterator
		row  *spanner.Row
	)
	if db.builder.tableName == "" {
		return errors.New("Undefined table name. please set db.Model(&struct{})")
	}
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

func (db *DB) Find(ctx context.Context, spannerTransaction interface{}) error {

	var (
		err  error
		iter *spanner.RowIterator
		row  *spanner.Row
	)

	var (
		isSlice, isPtr bool
	)
	var query string
	if db.builder.subBuilder.subModels != nil {
		query, _ = db.builder.buildSubQuery()
	} else {
		query, _ = db.builder.selectQuery()
	}
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

	results := reflect.Indirect(reflect.ValueOf(db.builder.model))
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
		results := reflect.Indirect(reflect.ValueOf(db.builder.model))
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

func (db *DB) Insert(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	query, err := db.builder.buildInsertModelQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Insert Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) Update(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	query, err := db.builder.buildUpdateModelQuery()
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Update Query: %s", db.builder.query)
	return rowCount, err
}

func (db *DB) UpdateMap(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in []string) (int64, error) {
	query, err := db.builder.buildUpdateMapQuery(in)
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Update Query: %s", db.builder.query)
	return rowCount, err
}
func (db *DB) UpdateWhere(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in map[string]interface{}) (int64, error) {
	query, err := db.builder.buildUpdateWhereQuery(in)
	if err != nil {
		return 0, errors.New("no primary key set")
	}
	stmt := spanner.Statement{SQL: query}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	db.logger.Infof("Update Query: %s", db.builder.query)
	return rowCount, err
}
