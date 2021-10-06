package ssorm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/10antz-inc/ssorm/utils"
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
	softDelete      bool
	softDeleteQuery string
}

type SubBuilder struct {
	subModels  []interface{}
	conditions []map[string]interface{}
}

func (builder *Builder) addSub(model interface{}, query interface{}, values interface{}) {
	builder.subBuilder.subModels = append(builder.subBuilder.subModels, model)
	condition := map[string]interface{}{"query": query, "args": values}
	builder.subBuilder.conditions = append(builder.subBuilder.conditions, condition)
}

func (builder *Builder) setSelects(query []string) {
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

	var deleteColumn string
	if conditions == nil || len(conditions) == 0 {
		if !builder.softDelete {
			return ""
		}
	}
	if builder.softDelete {
		deleteColumn = utils.GetDeleteColumnName(builder.model)
		if deleteColumn != "" {
			if conditions == nil || len(conditions) == 0 {
				conditions = map[string]interface{}{}
				conditions["query"] = fmt.Sprintf("%s IS NULL", deleteColumn)
				conditions["args"] = []interface{}{}
			} else {
				if conditions["query"] == "" {
					conditions["query"] = fmt.Sprintf("%s IS NULL", deleteColumn)
				} else {
					conditions["query"] = fmt.Sprintf("%s AND %s IS NULL", conditions["query"], deleteColumn)
				}

			}
		}
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
				for i := 0; i < values.Len(); i++ {
					if utils.IsTypeString(values.Index(i).Type()) {
						strValue := fmt.Sprintf("\"%s\"", fmt.Sprint(values.Index(i)))
						tempMarks = append(tempMarks, strValue)
					} else {
						tempMarks = append(tempMarks, fmt.Sprint(values.Index(i)))
					}
				}

				replacements = append(replacements, strings.Join(tempMarks, ","))
			}
		} else {
			if utils.IsTypeString(reflect.ValueOf(arg).Type()) {
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

	builder.buildCondition()

	return builder.query, nil
}

func (builder *Builder) buildInsertModelQuery() (string, error) {
	builder.query = fmt.Sprintf("INSERT INTO  %s", builder.tableName)
	e := reflect.Indirect(reflect.ValueOf(builder.model))
	var (
		cols []string
		vals []string
	)
	for i := 0; i < e.NumField(); i++ {
		addColumn := true
		tag, varName, varValue, varType := utils.ReflectValues(e, i)

		if tag.Get(utils.SSORM_TAG_KEY) == utils.SSORM_TAG_IGNORE_WRITE {
			continue
		}

		if utils.IsNullable(varValue) && !utils.IsValid(varValue) {
			addColumn = false
		}

		format := "%v"
		if utils.IsTypeString(varType) {
			format = "\"%v\""
		}
		switch tag.Get(utils.SSORM_TAG_KEY) {
		case utils.SSORM_TAG_CREATE_TIME:
			varValue = "CURRENT_TIMESTAMP()"
			break
		case utils.SSORM_TAG_UPDATE_TIME:
			varValue = "CURRENT_TIMESTAMP()"
			break
		case utils.SSORM_TAG_DELETE_TIME:
			addColumn = false
			break
		default:
			if utils.IsTime(varValue) {
				varValue = utils.GetTimestampStr(varValue)
			}

			if varType.Kind() == reflect.Slice || varType.Kind() == reflect.Array {
				varValue = utils.GetArrayStr(varValue, varType)
			}
		}

		if addColumn {
			cols = append(cols, fmt.Sprintf("%s", varName))
			if utils.IsTime(varType) {
				format = "\"%v\""
			}
			vals = append(vals, fmt.Sprintf(format, varValue))
		}

	}
	builder.query = fmt.Sprintf("%s (%s) VALUES (%s)", builder.query, strings.Join(cols, ", "), strings.Join(vals, ", "))
	return builder.query, nil
}

func (builder *Builder) buildUpdateModelQuery() (string, error) {
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)
	e := reflect.Indirect(reflect.ValueOf(builder.model))

	var (
		conditions []string
		updateData []string
	)

	for i := 0; i < e.NumField(); i++ {
		tag, varName, varValue, varType := utils.ReflectValues(e, i)

		if tag.Get(utils.SSORM_TAG_KEY) == utils.SSORM_TAG_IGNORE_WRITE {
			continue
		}

		if utils.IsNullable(varValue) && !utils.IsValid(varValue) {
			continue
		}
		format := "%s=%v"
		if utils.IsTypeString(varType) {
			format = "%s=\"%v\""
		}

		switch tag.Get(utils.SSORM_TAG_KEY) {
		case utils.SSORM_TAG_PRIMARY:
			conditions = append(conditions, fmt.Sprintf(format, varName, varValue))
			break
		case utils.SSORM_TAG_UPDATE_TIME:
			updateData = append(updateData, fmt.Sprintf(format, varName, "CURRENT_TIMESTAMP()"))
			break
		case utils.SSORM_TAG_CREATE_TIME:
			break
		case utils.SSORM_TAG_DELETE_TIME:
			break
		default:
			if utils.IsTime(varValue) {
				varValue = utils.GetTimestampStr(varValue)
			}

			if varType.Kind() == reflect.Slice || varType.Kind() == reflect.Array {
				varValue = utils.GetArrayStr(varValue, varType)
			}
			updateData = append(updateData, fmt.Sprintf(format, varName, varValue))
		}

	}

	if builder.softDelete {
		deleteColumn := utils.GetDeleteColumnName(builder.model)
		if deleteColumn != "" {
			conditions = append(conditions, fmt.Sprintf("%s IS NULL", deleteColumn))
		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	builder.query = fmt.Sprintf("%s WHERE %s ", builder.query, strings.Join(conditions, " AND "))
	if len(conditions) == 0 {
		return "", errors.New("no primary key set")
	}

	return builder.query, nil
}

func (builder *Builder) buildUpdateColumnQuery(in []string) (string, error) {
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)
	e := reflect.Indirect(reflect.ValueOf(builder.model))

	var (
		conditions []string
		updateData []string
	)

	for i := 0; i < e.NumField(); i++ {
		tag, varName, varValue, varType := utils.ReflectValues(e, i)
		if tag.Get(utils.SSORM_TAG_KEY) == utils.SSORM_TAG_IGNORE_WRITE {
			continue
		}
		format := "%s=%v"
		if utils.IsTypeString(varType) {
			format = "%s=\"%v\""
		}

		switch tag.Get(utils.SSORM_TAG_KEY) {
		case utils.SSORM_TAG_PRIMARY:
			conditions = append(conditions, fmt.Sprintf(format, varName, varValue))
			break
		case utils.SSORM_TAG_CREATE_TIME:
			break
		case utils.SSORM_TAG_DELETE_TIME:
			break
		case utils.SSORM_TAG_UPDATE_TIME:
			updateData = append(updateData, fmt.Sprintf(format, varName, "CURRENT_TIMESTAMP()"))
			break
		default:
			if utils.ArrayContains(in, varName) {
				if utils.IsTime(varValue) {
					varValue = utils.GetTimestampStr(varValue)
				}
				if varType.Kind() == reflect.Slice || varType.Kind() == reflect.Array {
					varValue = utils.GetArrayStr(varValue, varType)
				}
				format := "%s=%v"
				if utils.IsTypeString(varType) {
					format = "%s=\"%v\""
				}
				updateData = append(updateData, fmt.Sprintf(format, varName, varValue))
			}
		}

	}
	if builder.softDelete {
		deleteColumn := utils.GetDeleteColumnName(builder.model)
		if deleteColumn != "" {
			conditions = append(conditions, fmt.Sprintf("%s IS NULL", deleteColumn))
		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	builder.query = fmt.Sprintf("%s WHERE %s ", builder.query, strings.Join(conditions, " AND "))
	if len(conditions) == 0 {
		return "", errors.New("no primary key set")
	}
	return builder.query, nil
}
func (builder *Builder) buildUpdateParamsQuery(in map[string]interface{}) (string, error) {
	if builder.whereConditions == nil || len(builder.whereConditions) == 0 {
		return "", errors.New("no update condition set")
	}
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)

	var (
		updateData []string
	)

	for k, v := range in {
		varType := reflect.TypeOf(v)
		if varType.Kind() == reflect.Slice || varType.Kind() == reflect.Array {
			v = utils.GetArrayStr(v, varType)
		}
		format := "%s=%v"
		if utils.IsTypeString(varType) {
			format = "%s=\"%v\""
		}
		updateData = append(updateData, fmt.Sprintf(format, k, v))
	}

	e := reflect.Indirect(reflect.ValueOf(builder.model))
	for i := 0; i < e.NumField(); i++ {
		tag, varName, _, varType := utils.ReflectValues(e, i)
		format := "%s=%v"
		if utils.IsTypeString(varType) {
			format = "%s=\"%v\""
		}

		switch tag.Get(utils.SSORM_TAG_KEY) {
		case utils.SSORM_TAG_UPDATE_TIME:
			updateData = append(updateData, fmt.Sprintf(format, varName, "CURRENT_TIMESTAMP()"))
			break
		case utils.SSORM_TAG_DELETE_TIME:
			break
		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	condition := builder.buildWhereCondition(builder.whereConditions)
	if condition != "" {
		builder.query = fmt.Sprintf("%s WHERE %s", builder.query, condition)
	}
	return builder.query, nil
}

func (builder *Builder) buildDeleteModelQuery() (string, error) {
	builder.query = fmt.Sprintf("DELETE FROM %s WHERE", builder.tableName)
	e := reflect.Indirect(reflect.ValueOf(builder.model))
	var replacement []string
	for i := 0; i < e.NumField(); i++ {
		tag, varName, varValue, varType := utils.ReflectValues(e, i)
		format := "%s=%v"
		if tag.Get(utils.SSORM_TAG_KEY) == utils.SSORM_TAG_PRIMARY {
			if utils.IsTypeString(varType) {
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

func (builder *Builder) buildDeleteWhereQuery() (string, error) {
	if builder.whereConditions == nil || len(builder.whereConditions) == 0 {
		return "", errors.New("no delete condition set")
	}
	builder.query = fmt.Sprintf("DELETE FROM %s", builder.tableName)

	condition := builder.buildWhereCondition(builder.whereConditions)
	if condition != "" {
		builder.query = fmt.Sprintf("%s WHERE %s", builder.query, condition)
	}
	return builder.query, nil
}

func (builder *Builder) buildSoftDeleteModelQuery() (string, error) {
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)
	e := reflect.Indirect(reflect.ValueOf(builder.model))

	var (
		conditions []string
		updateData []string
	)

	for i := 0; i < e.NumField(); i++ {
		tag, varName, varValue, varType := utils.ReflectValues(e, i)
		format := "%s=%v"
		if utils.IsTypeString(varType) {
			format = "%s=\"%v\""
		}
		switch tag.Get(utils.SSORM_TAG_KEY) {
		case utils.SSORM_TAG_PRIMARY:
			conditions = append(conditions, fmt.Sprintf(format, varName, varValue))
			break
		case utils.SSORM_TAG_UPDATE_TIME:
			updateData = append(updateData, fmt.Sprintf(format, varName, "CURRENT_TIMESTAMP()"))
			break
		case utils.SSORM_TAG_CREATE_TIME:
			break
		case utils.SSORM_TAG_DELETE_TIME:
			updateData = append(updateData, fmt.Sprintf(format, varName, "CURRENT_TIMESTAMP()"))
			break
		}

	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	builder.query = fmt.Sprintf("%s WHERE %s ", builder.query, strings.Join(conditions, " AND "))
	return builder.query, nil
}

func (builder *Builder) buildSoftDeleteWhereQuery() (string, error) {
	if builder.whereConditions == nil || len(builder.whereConditions) == 0 {
		return "", errors.New("no delete condition set")
	}
	builder.query = fmt.Sprintf("UPDATE %s SET", builder.tableName)

	var (
		updateData []string
	)

	e := reflect.Indirect(reflect.ValueOf(builder.model))
	for i := 0; i < e.NumField(); i++ {
		tag, varName, _, varType := utils.ReflectValues(e, i)
		format := "%s=%v"
		if utils.IsTypeString(varType) {
			format = "%s=\"%v\""
		}
		switch tag.Get(utils.SSORM_TAG_KEY) {
		case utils.SSORM_TAG_UPDATE_TIME:
			updateData = append(updateData, fmt.Sprintf(format, varName, "CURRENT_TIMESTAMP()"))
			break
		case utils.SSORM_TAG_DELETE_TIME:
			updateData = append(updateData, fmt.Sprintf(format, varName, "CURRENT_TIMESTAMP()"))
			break
		}
	}
	builder.query = fmt.Sprintf("%s %s", builder.query, strings.Join(updateData, ","))
	condition := builder.buildWhereCondition(builder.whereConditions)
	if condition != "" {
		builder.query = fmt.Sprintf("%s WHERE %s", builder.query, condition)
	}
	return builder.query, nil
}
