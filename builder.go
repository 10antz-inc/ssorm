package ssorm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/10antz-inc/cp-service-go/ssorm/utils"
	"reflect"
	"strings"
)

type Builder struct {
	query           string
	limit           int64
	whereConditions map[string]interface{}
	selects         []string
	orders          string
	offset          int64
	tableName       string
	model           interface{}
}

func (builder *Builder) setModel(model interface{}) {
	builder.model = model
	builder.tableName = utils.GetTableName(builder.model)
}

func (builder *Builder) setSelects(query []string, args ...interface{}) {
	builder.selects = query
}

func (builder *Builder) setOffset(offset int64) {
	builder.offset = offset
}

func (builder *Builder) setOrder(order string) {
	builder.orders = order
}

func (builder *Builder) setWhere(query interface{}, values interface{}) {
	builder.whereConditions = map[string]interface{}{"query": query, "args": values}
}

func (builder *Builder) setLimit(limit int64) {
	builder.limit = limit
}

func (builder *Builder) buildTableName(model ...interface{}) {
	builder.tableName = utils.GetTableName(model)
}

func (builder *Builder) selectQuery() (string, error) {
	err := builder.buildSelectQuery()
	err = builder.buildCondition()
	return builder.query, err
}

func (builder *Builder) deleteModelQuery() (string, error) {
	builder.buildDeleteModelQuery()
	return builder.query, nil
}

func (builder *Builder) deleteWhereQuery() (string, error) {
	builder.buildDeleteConditionQuery()
	return builder.query, nil
}

func (builder *Builder) buildCondition() error {
	builder.buildWhereCondition()
	builder.buildOrders()
	builder.buildLimit()
	builder.buildOffset()
	return nil
}

func (builder *Builder) buildSelectQuery() error {
	if builder.selects != nil {
		selectQuery := strings.Join(builder.selects, ",")
		builder.query = fmt.Sprintf("SELECT %s FROM %s", selectQuery, builder.tableName)
		return nil
	}
	builder.query = fmt.Sprintf("SELECT * FROM %s", builder.tableName)
	return nil
}

func (builder *Builder) buildDeleteModelQuery() error {
	builder.query = fmt.Sprintf("DELETE FROM %s WHERE", builder.tableName)
	e := utils.Indirect(reflect.ValueOf(builder.model))
	value := reflect.TypeOf(e.Interface())

	var replacement []string
	for i := 0; i < e.NumField(); i++ {
		tag := value.Field(i).Tag
		if tag.Get("key") == "primary" {
			varName := e.Type().Field(i).Name
			varType := e.Type().Field(i).Type
			varValue := e.Field(i).Interface()
			switch varType.Kind() {
			case reflect.String:
				replacement = append(replacement, fmt.Sprintf("%s=\"%s\"", varName, varValue))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				replacement = append(replacement, fmt.Sprintf("%s=%v", varName, varValue))
			case reflect.Float32, reflect.Float64:
				replacement = append(replacement, fmt.Sprintf("%s=%v", varName, varValue))
			}
		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(replacement, " AND "))
	if len(replacement) == 0 {
		return errors.New("no primary key set")
	}
	return nil
}

func (builder *Builder) buildDeleteConditionQuery() error {
	builder.query = fmt.Sprintf("DELETE FROM %s", builder.tableName)
	builder.buildWhereCondition()
	return nil
}

func (builder *Builder) buildLimit() {
	if builder.limit > 0 {
		builder.query = fmt.Sprintf("%s %s %d", builder.query, "LIMIT", builder.limit)
	}
}

func (builder *Builder) buildOrders() {
	if builder.orders != "" {
		builder.query = fmt.Sprintf("%s %s %s", builder.query, "ORDER BY", builder.orders)
	}
}

func (builder *Builder) buildOffset() {
	if builder.offset > 0 {
		builder.query = fmt.Sprintf("%s %s %d", builder.query, "OFFSET", builder.offset)
	}
}

func (builder *Builder) buildWhereCondition() {
	if builder.whereConditions == nil {
		return
	}
	clause := builder.whereConditions
	query := clause["query"].(string)
	args := clause["args"].([]interface{})

	var replacements []string
	for _, arg := range args {

		switch reflect.ValueOf(arg).Kind() {
		case reflect.Slice: // For where("id in (?)", []int64{1,2})
			if values := reflect.ValueOf(arg); values.Len() > 0 {
				var tempMarks []string
				var isString bool
				for i := 0; i < values.Len(); i++ {
					isString = values.Index(i).Kind() == reflect.String
					if isString {
						strValue := fmt.Sprintf("\"%s\"", fmt.Sprint(values.Index(i)))
						tempMarks = append(tempMarks, strValue)
					} else {
						tempMarks = append(tempMarks, fmt.Sprint(values.Index(i)))
					}
				}

				replacements = append(replacements, strings.Join(tempMarks, ","))
			}
		case reflect.String:
			replacements = append(replacements, fmt.Sprintf("\"%s\"", arg.(string)))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			replacements = append(replacements, fmt.Sprintf("%v", arg))
		case reflect.Float32, reflect.Float64:
			replacements = append(replacements, fmt.Sprintf("%v", arg))

		}
	}

	buff := bytes.NewBuffer([]byte{})
	i := 0
	for _, s := range query {
		if s == '?' && len(replacements) > i {
			buff.WriteString(replacements[i])
			i++
		} else {
			buff.WriteRune(s)
		}
	}

	builder.query = fmt.Sprintf("%s WHERE %s", builder.query, buff.String())
}
