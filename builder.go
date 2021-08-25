package ssorm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/10antz-inc/ssorm/utils"
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
	subBuilder      *SubBuilder
}

type SubBuilder struct {
	subModels  []interface{}
	conditions []map[string]interface{}
}

func (builder *Builder) addSub(model interface{}, query interface{}, values ...interface{}) {
	builder.subBuilder.subModels = append(builder.subBuilder.subModels, model)
	condition := map[string]interface{}{"query": query, "args": values}
	builder.subBuilder.conditions = append(builder.subBuilder.conditions, condition)
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

func (builder *Builder) buildCondition() error {
	condition := builder.buildWhereCondition(builder.whereConditions)
	if condition != "" {
		builder.query = fmt.Sprintf("%s WHERE %s", builder.query, condition)
	}
	builder.buildOrders()
	builder.buildLimit()
	builder.buildOffset()
	return nil
}

func (builder *Builder) selectQuery() (string, error) {
	err := builder.buildSelectQuery()
	err = builder.buildCondition()
	return builder.query, err
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

func (builder *Builder) buildSubQuery() (string, error) {
	if builder.selects != nil {
		selectQuery := strings.Join(builder.selects, ",")
		builder.query = fmt.Sprintf("SELECT %s", selectQuery)
	}
	builder.query = "SELECT *"

	var subQueries []string
	for i, v := range builder.subBuilder.subModels {
		//ARRAY(SELECT AS STRUCT * FROM Albums WHERE SingerId > 12) as Albums,
		tableName := utils.GetTableName(v)
		query := fmt.Sprintf("SELECT AS STRUCT * FROM %s", tableName)
		if builder.subBuilder.conditions[i] != nil {
			condition := builder.buildWhereCondition(builder.subBuilder.conditions[i])
			if condition != "" {
				query = fmt.Sprintf("%s WHERE %s", query, condition)
			}
		}
		query = fmt.Sprintf("ARRAY ( %s ) as %s", query, tableName)
		subQueries = append(subQueries, query)
	}
	builder.query = fmt.Sprintf("%s, %s, FROM %s", builder.query, strings.Join(subQueries, ","), builder.tableName)
	condition := builder.buildWhereCondition(builder.whereConditions)
	if condition != "" {
		builder.query = fmt.Sprintf("%s WHERE %s", condition)
	}
	err := builder.buildCondition()
	return builder.query, err
}



func (builder *Builder) deleteModelQuery() (string, error) {
	builder.query = fmt.Sprintf("DELETE FROM %s WHERE", builder.tableName)
	e := utils.Indirect(reflect.ValueOf(builder.model))
	value := reflect.TypeOf(e.Interface())

	var replacement []string
	for i := 0; i < e.NumField(); i++ {
		tag := value.Field(i).Tag
		format := "%s=%v"
		if tag.Get("key") == "primary" {
			varName := e.Type().Field(i).Name
			varType := e.Type().Field(i).Type
			varValue := e.Field(i).Interface()
			switch varType.Kind() {
			case reflect.String:
				format = "%s=\"%v\""
			}
			replacement = append(replacement, fmt.Sprintf(format, varName, varValue))
		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(replacement, " AND "))
	if len(replacement) == 0 {
		return "", errors.New("no primary key set")
	}
	return builder.query, nil
}

func (builder *Builder) deleteWhereQuery() (string, error) {
	builder.query = fmt.Sprintf("DELETE FROM %s", builder.tableName)
	condition := builder.buildWhereCondition(builder.whereConditions)
	if condition != "" {
		builder.query = fmt.Sprintf("%s WHERE %s", builder.query, condition)
	}
	return builder.query, nil
}

func (builder *Builder) buildInsertModelQuery() (string, error) {
	builder.query = fmt.Sprintf("INSERT INTO  %s", builder.tableName)
	e := utils.Indirect(reflect.ValueOf(builder.model))
	var (
		cols []string
		vals []string
	)
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		varValue := e.Field(i).Interface()
		format := "%v"
		cols = append(cols, fmt.Sprintf("%s", varName))
		switch varType.Kind() {
		case reflect.String:
			format = "\"%v\""
		}
		vals = append(vals, fmt.Sprintf(format, varValue))
	}
	builder.query = fmt.Sprintf("%s (%s) VALUES (%s)", builder.query, strings.Join(cols, ", "), strings.Join(vals, ", "))
	return builder.query, nil
}

func (builder *Builder) buildUpdateModelQuery() (string, error) {
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)
	e := utils.Indirect(reflect.ValueOf(builder.model))
	value := reflect.TypeOf(e.Interface())

	var (
		replacement []string
		updateData  []string
	)

	for i := 0; i < e.NumField(); i++ {
		tag := value.Field(i).Tag
		if tag.Get("key") == "primary" {
			varName := e.Type().Field(i).Name
			varType := e.Type().Field(i).Type
			varValue := e.Field(i).Interface()
			format := "%s=%v"
			switch varType.Kind() {
			case reflect.String:
				format = "%s=\"%v\""
			}
			replacement = append(replacement, fmt.Sprintf(format, varName, varValue))
		} else {

			varName := e.Type().Field(i).Name
			varType := e.Type().Field(i).Type
			varValue := e.Field(i).Interface()
			format := "%s=%v"
			switch varType.Kind() {
			case reflect.String:
				format = "%s=\"%v\""
			}
			updateData = append(updateData, fmt.Sprintf(format, varName, varValue))

		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	builder.query = fmt.Sprintf("%s WHERE %s ", builder.query, strings.Join(replacement, " AND "))
	if len(replacement) == 0 {
		return "", errors.New("no primary key set")
	}
	return builder.query, nil
}

func (builder *Builder) buildUpdateMapQuery(in map[string]interface{}) (string, error) {
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)
	e := utils.Indirect(reflect.ValueOf(builder.model))
	value := reflect.TypeOf(e.Interface())

	var (
		replacement []string
		updateData  []string
	)

	for i := 0; i < e.NumField(); i++ {
		tag := value.Field(i).Tag
		if tag.Get("key") == "primary" {
			varName := e.Type().Field(i).Name
			varType := e.Type().Field(i).Type
			varValue := e.Field(i).Interface()
			format := "%s=%v"
			switch varType.Kind() {
			case reflect.String:
				format = "%s=\"%v\""
			}
			replacement = append(replacement, fmt.Sprintf(format, varName, varValue))

		} else {
			varName := e.Type().Field(i).Name
			if val, ok := in[varName]; ok {
				varType := reflect.TypeOf(val)
				varValue := val
				format := "%s=%v"
				switch varType.Kind() {
				case reflect.String:
					format = "%s=\"%v\""
				}
				updateData = append(updateData, fmt.Sprintf(format, varName, varValue))
			}
		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	builder.query = fmt.Sprintf("%s WHERE %s ", builder.query, strings.Join(replacement, " AND "))
	if len(replacement) == 0 {
		return "", errors.New("no primary key set")
	}
	return builder.query, nil
}

func (builder *Builder) buildUpdateWhereQuery(in map[string]interface{}) (string, error) {
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)

	var (
		updateData []string
	)

	for k, v := range in {
		varType := reflect.TypeOf(v)
		format := "%s=%v"
		switch varType.Kind() {
		case reflect.String:
			format = "%s=\"%v\""
		}
		updateData = append(updateData, fmt.Sprintf(format, k, v))
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	condition := builder.buildWhereCondition(builder.whereConditions)
	if condition != "" {
		builder.query = fmt.Sprintf("%s WHERE %s", builder.query, condition)
	}
	return builder.query, nil
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

func (builder *Builder) buildWhereCondition(conditions map[string]interface{}) string {
	if conditions == nil || len(conditions) == 0 {
		return ""
	}
	clause := conditions
	query := clause["query"].(string)
	args := clause["args"].([]interface{})

	var replacements []string
	for _, arg := range args {
		format := "%v"
		if reflect.ValueOf(arg).Kind() == reflect.Slice {
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
		} else {
			if reflect.ValueOf(arg).Kind() == reflect.String {
				format = "\"%v\""
			}
			replacements = append(replacements, fmt.Sprintf(format, arg))
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

	return buff.String()
}
