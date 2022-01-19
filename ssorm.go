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

var logger ILogger

func getLogger() ILogger {
	if logger == nil {
		logger = NewLogger(os.Stdout)
	}
	return logger
}

func Logger(l ILogger) Option {
	return func(d *DB) {
		logger = l
	}
}

type DB struct {
	builder *Builder
}

func Model(model interface{}, opts ...Option) *DB {
	db := &DB{}
	db.builder = &Builder{
		subBuilder: &SubBuilder{},
		model:      model,
		tableName:  utils.GetTableName(model),
		params:     make(map[string]interface{}),
		softDelete: false,
	}
	return db
}

func SoftDeleteModel(model interface{}, opts ...Option) *DB {
	db := &DB{}
	db.builder = &Builder{
		subBuilder: &SubBuilder{},
		model:      model,
		tableName:  utils.GetTableName(model),
		softDelete: true,
		params:     make(map[string]interface{}),
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

func (db *DB) Select(query []string) *DB {
	db.builder.setSelects(query)
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
		err   error
		query string
	)

	if db.builder.softDelete {
		query, err := db.builder.buildSoftDeleteModelQuery()
		if err != nil {
			return 0, errors.New("no primary key set")
		}
		rowCount, err := SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)

		getLogger().Infof("Update Query: %s Param: %+v", db.builder.query, db.builder.params)
		return rowCount, err
	}

	query, err = db.builder.buildDeleteModelQuery()
	if err != nil {
		return 0, err
	}
	getLogger().Infof("DELETE Query: %s Param: %+v", db.builder.query, db.builder.params)
	return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)
}

func (db *DB) DeleteWhere(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	var (
		err   error
		query string
	)
	if db.builder.softDelete {
		query, err = db.builder.buildSoftDeleteWhereQuery()
		if err != nil {
			return 0, err
		}
		getLogger().Infof("Update Query: %s Param: %+v", db.builder.query, db.builder.params)
		return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)
	}

	query, err = db.builder.buildDeleteWhereQuery()
	if err != nil {
		return 0, err
	}
	getLogger().Infof("DELETE Query: %s Param: %+v", db.builder.query, db.builder.params)
	return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)
}

func (db *DB) First(ctx context.Context, spannerTransaction interface{}) error {
	var (
		err error
	)
	db.builder.limit = 1
	var query string
	if db.builder.subBuilder.subModels != nil {
		query, _ = db.builder.buildSubQuery()
	} else {
		query, _ = db.builder.selectQuery()
	}

	err = SimpleQueryRead(ctx, spannerTransaction, query, db.builder.params, db.builder.model)
	return err
}

func (db *DB) Count(ctx context.Context, spannerTransaction interface{}, cnt interface{}) error {
	var (
		err  error
		iter *spanner.RowIterator
		row  *spanner.Row
	)
	if db.builder.tableName == "" {
		return errors.New("Undefined table name. please set ssorm.Model(&struct{})")
	}
	query, err := db.Select([]string{"COUNT(1) AS CNT"}).builder.selectQuery()

	stmt := spanner.Statement{SQL: query, Params: db.builder.params}

	getLogger().Infof("Select Query: %s Param: %+v", stmt.SQL, db.builder.params)

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
		err error
	)

	var query string
	if db.builder.subBuilder.subModels != nil {
		query, _ = db.builder.buildSubQuery()
	} else {
		query, _ = db.builder.selectQuery()
	}

	err = SimpleQueryRead(ctx, spannerTransaction, query, db.builder.params, db.builder.model)
	return err
}

func (db *DB) Insert(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	query, err := db.builder.buildInsertModelQuery()
	if err != nil {
		return 0, err
	}
	getLogger().Infof("Insert Query: %s Param: %+v", db.builder.query, db.builder.params)
	return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)

}

func (db *DB) Update(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	query, err := db.builder.buildUpdateModelQuery()
	if err != nil {
		return 0, err
	}
	getLogger().Infof("Update Query: %s Param: %+v", db.builder.query, db.builder.params)
	return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)

}

func (db *DB) UpdateColumns(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in []string) (int64, error) {
	query, err := db.builder.buildUpdateColumnQuery(in, false)
	if err != nil {
		return 0, err
	}
	getLogger().Infof("Update Query: %s Param: %+v", db.builder.query, db.builder.params)
	return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)
}

func (db *DB) UpdateOmit(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in []string) (int64, error) {
	query, err := db.builder.buildUpdateColumnQuery(in, true)
	if err != nil {
		return 0, err
	}
	getLogger().Infof("Update Query: %s Param: %+v", db.builder.query, db.builder.params)
	return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)
}

func (db *DB) UpdateParams(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in map[string]interface{}) (int64, error) {
	query, err := db.builder.buildUpdateParamsQuery(in)
	if err != nil {
		return 0, err
	}
	getLogger().Infof("Update Query: %s Param: %+v", db.builder.query, db.builder.params)
	return SimpleQueryWrite(ctx, spannerTransaction, query, db.builder.params)
}

func SimpleQueryRead(ctx context.Context, spannerTransaction interface{}, query string, params map[string]interface{}, result interface{}) error {
	var (
		err  error
		iter *spanner.RowIterator
		row  *spanner.Row
	)

	var (
		isPtr bool
	)
	stmt := spanner.Statement{SQL: query, Params: params}
	getLogger().Infof("Select Query: %s Param: %+v", stmt.SQL, params)

	rot, readOnly := spannerTransaction.(*spanner.ReadOnlyTransaction)
	rwt, readWrite := spannerTransaction.(*spanner.ReadWriteTransaction)
	if readOnly {
		iter = rot.Query(ctx, stmt)
	}
	if readWrite {
		iter = rwt.Query(ctx, stmt)
	}

	defer iter.Stop()

	results := reflect.Indirect(reflect.ValueOf(result))
	var resultType reflect.Type
	if kind := results.Kind(); kind == reflect.Slice {
		resultType = results.Type().Elem()

		results.Set(reflect.MakeSlice(results.Type(), 0, 0))

		if resultType.Kind() == reflect.Ptr {
			isPtr = true
			resultType = resultType.Elem()
		}
		for {
			if row, err = iter.Next(); err != nil {
				if err == iterator.Done {
					return nil
				}
				return err
			}
			results := reflect.Indirect(reflect.ValueOf(result))
			elem := reflect.New(resultType).Interface()
			if err := row.ToStruct(elem); err != nil {
				return fmt.Errorf("failed to struct: %w", err)
			}

			if isPtr {
				results.Set(reflect.Append(results, reflect.ValueOf(elem).Elem().Addr()))
			} else {
				results.Set(reflect.Append(results, reflect.ValueOf(elem).Elem()))
			}
		}
	} else {
		for {
			if row, err = iter.Next(); err != nil {
				if err == iterator.Done {
					fmt.Printf("Result: %+v", result)
					return nil
				}
				return err
			}

			if err := row.ToStruct(result); err != nil {
				return fmt.Errorf("failed to struct: %w", err)
			}
			break
		}
	}

	return err
}

func SimpleQueryWrite(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, query string, params map[string]interface{}) (int64, error) {
	stmt := spanner.Statement{SQL: query, Params: params}
	rowCount, err := spannerTransaction.Update(ctx, stmt)
	return rowCount, err
}
