package ssorm

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/10antz-inc/ssorm/v2/instrumentation/ssormotel"
	"github.com/10antz-inc/ssorm/v2/logger"
	"github.com/10antz-inc/ssorm/v2/utils"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

var tracing ssormotel.Tracing
var ssormLogger *logger.Logger

func init() {
	ssormLogger = logger.NewLogger()
}

func UseTrace(opts ...ssormotel.Option) {
	tracing = ssormotel.NewTracing(opts...)
}

func LoggerConfig(opts ...logger.Option) {
	ssormLogger = logger.NewLogger(opts...)
}

type Option func(*DB)

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

func (db *DB) QueryOptions(opt *spanner.QueryOptions) *DB {
	db.builder.queryOptions = opt
	return db
}

func (db *DB) ToRefresh() *DB {
	db.builder.refresh = true
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

func (db *DB) Find(ctx context.Context, spannerTransaction interface{}) error {
	cmd := db.find(ctx, spannerTransaction)
	if tracing != nil {
		return tracing.StartForRead(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) First(ctx context.Context, spannerTransaction interface{}) error {
	cmd := db.first(ctx, spannerTransaction)
	if tracing != nil {
		return tracing.StartForRead(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) Count(ctx context.Context, spannerTransaction interface{}, cnt interface{}) error {
	cmd := db.count(ctx, spannerTransaction, cnt)
	if tracing != nil {
		return tracing.StartForRead(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) Insert(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	cmd := db.insert(ctx, spannerTransaction)
	if tracing != nil {
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) Update(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	cmd := db.update(ctx, spannerTransaction)
	if tracing != nil {
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) UpdateColumns(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in []string) (int64, error) {
	cmd := db.updateColumns(ctx, spannerTransaction, in)
	if tracing != nil {
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) UpdateOmit(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in []string) (int64, error) {
	cmd := db.updateOmit(ctx, spannerTransaction, in)
	if tracing != nil {
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) UpdateParams(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in map[string]interface{}) (int64, error) {
	cmd := db.updateParams(ctx, spannerTransaction, in)
	if tracing != nil {
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) DeleteModel(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	cmd := db.deleteModel(ctx, spannerTransaction)
	if tracing != nil {
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) DeleteWhere(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) (int64, error) {
	cmd := db.deleteWhere(ctx, spannerTransaction)
	if tracing != nil {
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func SimpleQueryRead(ctx context.Context, spannerTransaction interface{}, query string, params map[string]interface{}, result interface{}) error {
	return SimpleQueryReadWithOptions(ctx, spannerTransaction, query, params, result, nil)
}

func SimpleQueryReadWithOptions(ctx context.Context, spannerTransaction interface{}, query string, params map[string]interface{}, result interface{}, queryOpts *spanner.QueryOptions) error {
	cmd := simpleQueryRead(ctx, spannerTransaction, query, params, result, queryOpts)
	if tracing != nil {
		statement := fmt.Sprintf("%s, params: %+v", query, params)
		tracing.SetStatement(statement)
		return tracing.StartForRead(ctx, cmd)
	}

	return cmd(ctx)
}

func SimpleQueryWrite(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, query string, params map[string]interface{}) (int64, error) {
	return SimpleQueryWriteWithOptions(ctx, spannerTransaction, query, params, nil)
}

func SimpleQueryWriteWithOptions(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, query string, params map[string]interface{}, queryOpts *spanner.QueryOptions) (int64, error) {
	cmd := simpleQueryWrite(ctx, spannerTransaction, query, params, queryOpts)
	if tracing != nil {
		statement := fmt.Sprintf("%s, params: %+v", query, params)
		tracing.SetStatement(statement)
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func (db *DB) find(ctx context.Context, spannerTransaction interface{}) func(context.Context) error {
	return func(ctx context.Context) error {
		var (
			err error
		)

		var query string
		if db.builder.subBuilder.subModels != nil {
			query = db.builder.buildSubQuery()
		} else {
			query = db.builder.selectQuery()
		}

		err = SimpleQueryReadWithOptions(ctx, spannerTransaction, query, db.builder.params, db.builder.model, db.builder.queryOptions)
		return err
	}
}

func (db *DB) first(ctx context.Context, spannerTransaction interface{}) func(context.Context) error {
	return func(ctx context.Context) error {
		var (
			err error
		)
		db.builder.limit = 1
		var query string
		if db.builder.subBuilder.subModels != nil {
			query = db.builder.buildSubQuery()
		} else {
			query = db.builder.selectQuery()
		}

		err = SimpleQueryReadWithOptions(ctx, spannerTransaction, query, db.builder.params, db.builder.model, db.builder.queryOptions)
		return err
	}
}

func (db *DB) count(ctx context.Context, spannerTransaction interface{}, cnt interface{}) func(context.Context) error {
	return func(ctx context.Context) error {
		var (
			err  error
			iter *spanner.RowIterator
			row  *spanner.Row
		)
		if db.builder.tableName == "" {
			return errors.New("Undefined table name. please set ssorm.Model(&struct{})")
		}
		query := db.Select([]string{"COUNT(1) AS CNT"}).builder.selectQuery()

		stmt := spanner.Statement{SQL: query, Params: db.builder.params}

		ssormLogger.ReadLog(ctx, "Select Query: %s Param: %+v", stmt.SQL, db.builder.params)

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
			if err := row.ColumnByName("CNT", cnt); err != nil {
				ssormLogger.ErrorLog(ctx, err, "")
				return err
			}
			break
		}

		return err
	}
}

func (db *DB) insert(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		query, err := db.builder.buildInsertModelQuery()
		if err != nil {
			return 0, err
		}
		ssormLogger.WriteLog(ctx, "Insert Query: %s Param: %+v", db.builder.query, db.builder.params)
		return execQueryWriter(ctx, spannerTransaction, query, db.builder)
	}
}

func (db *DB) update(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		query, err := db.builder.buildUpdateModelQuery()
		if err != nil {
			return 0, err
		}
		ssormLogger.WriteLog(ctx, "Update Query: %s Param: %+v", db.builder.query, db.builder.params)
		return execQueryWriter(ctx, spannerTransaction, query, db.builder)
	}
}

func (db *DB) updateColumns(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in []string) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		query, err := db.builder.buildUpdateColumnQuery(in, false)
		if err != nil {
			return 0, err
		}
		ssormLogger.WriteLog(ctx, "Update Query: %s Param: %+v", db.builder.query, db.builder.params)
		return execQueryWriter(ctx, spannerTransaction, query, db.builder)
	}
}

func (db *DB) updateOmit(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in []string) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		query, err := db.builder.buildUpdateColumnQuery(in, true)
		if err != nil {
			return 0, err
		}
		ssormLogger.WriteLog(ctx, "Update Query: %s Param: %+v", db.builder.query, db.builder.params)
		return execQueryWriter(ctx, spannerTransaction, query, db.builder)
	}
}

func (db *DB) updateParams(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, in map[string]interface{}) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		query, err := db.builder.buildUpdateParamsQuery(in)
		if err != nil {
			return 0, err
		}
		ssormLogger.WriteLog(ctx, "Update Query: %s Param: %+v", db.builder.query, db.builder.params)
		return execQueryWriter(ctx, spannerTransaction, query, db.builder)
	}
}

func (db *DB) deleteModel(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		var (
			err   error
			query string
		)

		if db.builder.softDelete {
			query, err := db.builder.buildSoftDeleteModelQuery()
			if err != nil {
				return 0, errors.New("no primary key set")
			}
			rowCount, err := execQueryWriter(ctx, spannerTransaction, query, db.builder)

			ssormLogger.WriteLog(ctx, "Update Query: %s Param: %+v", db.builder.query, db.builder.params)
			return rowCount, err
		}

		query, err = db.builder.buildDeleteModelQuery()
		if err != nil {
			return 0, err
		}
		ssormLogger.WriteLog(ctx, "DELETE Query: %s Param: %+v", db.builder.query, db.builder.params)
		return execQueryWriter(ctx, spannerTransaction, query, db.builder)
	}
}

func (db *DB) deleteWhere(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		var (
			err   error
			query string
		)
		if db.builder.softDelete {
			query, err = db.builder.buildSoftDeleteWhereQuery()
			if err != nil {
				return 0, err
			}
			ssormLogger.WriteLog(ctx, "Update Query: %s Param: %+v", db.builder.query, db.builder.params)
			return execQueryWriter(ctx, spannerTransaction, query, db.builder)
		}

		query, err = db.builder.buildDeleteWhereQuery()
		if err != nil {
			return 0, err
		}
		ssormLogger.WriteLog(ctx, "DELETE Query: %s Param: %+v", db.builder.query, db.builder.params)
		return execQueryWriter(ctx, spannerTransaction, query, db.builder)
	}
}

func simpleQueryRead(ctx context.Context, spannerTransaction interface{}, query string, params map[string]interface{}, result interface{}, queryOpts *spanner.QueryOptions) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var (
			iter *spanner.RowIterator
			row  *spanner.Row
		)

		stmt := spanner.Statement{SQL: query, Params: params}
		ssormLogger.ReadLog(ctx, "Select Query: %s Param: %+v", stmt.SQL, params)

		switch {
		case queryOpts != nil:
			switch txn := spannerTransaction.(type) {
			case interface {
				QueryWithOptions(context.Context, spanner.Statement, spanner.QueryOptions) *spanner.RowIterator
			}:
				iter = txn.QueryWithOptions(ctx, stmt, *queryOpts)
			}
		default:
			switch txn := spannerTransaction.(type) {
			case interface {
				Query(context.Context, spanner.Statement) *spanner.RowIterator
			}:
				iter = txn.Query(ctx, stmt)
			}
		}

		if iter == nil {
			panic("An unexpected transaction.")
		}

		defer iter.Stop()

		return reflectValues(ctx, result, row, iter)
	}
}

func reflectValues(ctx context.Context, result interface{}, row *spanner.Row, iter *spanner.RowIterator) error {
	var (
		err   error
		isPtr bool
	)
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
				ssormLogger.ErrorLog(ctx, err, "")
				return err
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
					return nil
				}
				return err
			}

			if err := row.ToStruct(result); err != nil {
				ssormLogger.ErrorLog(ctx, err, "")
				return err
			}
			break
		}
	}

	return err
}

func execQueryWriter(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, query string, builder *Builder) (int64, error) {
	var cmd ssormotel.Write

	switch {
	case builder.refresh:
		cmd = simpleQueryWriteForRefresh(ctx, spannerTransaction, query, builder)
	default:
		cmd = simpleQueryWrite(ctx, spannerTransaction, query, builder.params, builder.queryOptions)
	}

	if tracing != nil {
		statement := fmt.Sprintf("%s, params: %+v", query, builder.params)
		tracing.SetStatement(statement)
		return tracing.StartForWrite(ctx, cmd)
	}

	return cmd(ctx)
}

func simpleQueryWrite(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, query string, params map[string]interface{}, queryOpts *spanner.QueryOptions) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		ssormLogger.WriteLog(ctx, "Simple Write Query: %s Param: %+v", query, params)
		switch {
		case queryOpts != nil:
			return spannerTransaction.UpdateWithOptions(ctx, spanner.Statement{SQL: query, Params: params}, *queryOpts)
		default:
			return spannerTransaction.Update(ctx, spanner.Statement{SQL: query, Params: params})
		}
	}
}

func simpleQueryWriteForRefresh(ctx context.Context, spannerTransaction *spanner.ReadWriteTransaction, query string, builder *Builder) func(context.Context) (int64, error) {
	var (
		iter *spanner.RowIterator
		row  *spanner.Row
		err  error
	)

	return func(ctx context.Context) (int64, error) {
		ssormLogger.WriteLog(ctx, "Simple Write Query: %s Param: %+v", query, builder.params)
		switch {
		case builder.queryOptions != nil:
			iter = spannerTransaction.QueryWithOptions(ctx, spanner.Statement{SQL: query, Params: builder.params}, *builder.queryOptions)
		default:
			iter = spannerTransaction.Query(ctx, spanner.Statement{SQL: query, Params: builder.params})
		}
		defer iter.Stop()
		err = reflectValues(ctx, builder.model, row, iter)
		return iter.RowCount, err
	}
}
