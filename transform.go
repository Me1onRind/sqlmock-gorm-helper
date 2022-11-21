package sqlmock_rows_helper

import (
	"database/sql/driver"
	"reflect"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
)

type GetColumnFromField func(field reflect.StructField) string

var (
	getColumnFromField GetColumnFromField = GetColumnFromGromModelField
)

func ModelToRows(dst interface{}) *sqlmock.Rows {
	dstType := reflect.TypeOf(dst)
	dstValue := reflect.ValueOf(dst)

	var allValues [][]driver.Value

	if dstType.Kind() == reflect.Slice {
		dstType = dstType.Elem()
		for i := 0; i < dstValue.Len(); i++ {
			dstValueItem := dstValue.Index(i)
			allValues = append(allValues, valuesFromModel(dstType, dstValueItem))
		}
	} else {
		allValues = append(allValues, valuesFromModel(dstType, dstValue))
	}

	rows := sqlmock.NewRows(columnsFromModelType(dstType))
	for _, row := range allValues {
		rows.AddRow(row...)
	}

	return rows
}

func valuesFromModel(dstType reflect.Type, dstValue reflect.Value) []driver.Value {
	if dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
		dstType = dstType.Elem()
	}
	var values []driver.Value
	for j := 0; j < dstValue.NumField(); j++ {
		fieldValue := dstValue.Field(j)
		dstTypeField := dstType.Field(j)
		if fieldValue.Kind() == reflect.Struct && dstTypeField.Anonymous {
			values = append(values, valuesFromModel(fieldValue.Type(), fieldValue)...)
		}
		columnValue := valueFromField(dstTypeField, fieldValue)
		if columnValue != nil {
			switch cv := columnValue.(type) {
			case driver.Valuer:
				v, _ := cv.Value()
				if strV, ok := v.(string); ok {
					values = append(values, []byte(strV))
				} else {
					values = append(values, v)
				}
			default:
				values = append(values, columnValue)
			}
		}
	}
	return values
}

func columnsFromModelType(dstType reflect.Type) []string {
	if dstType.Kind() == reflect.Ptr {
		dstType = dstType.Elem()
	}
	var columns []string
	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			columns = append(columns, columnsFromModelType(field.Type)...)
		} else {
			column := getColumnFromField(field)
			if len(column) > 0 {
				columns = append(columns, column)
			}
		}
	}
	return columns
}

func valueFromField(field reflect.StructField, value reflect.Value) driver.Value {
	if len(getColumnFromField(field)) > 0 {
		return value.Interface()
	}
	return nil
}

func GetColumnFromGromModelField(field reflect.StructField) string {
	tag := field.Tag.Get("gorm")
	items := strings.Split(tag, ";")
	for _, v := range items {
		if strings.HasPrefix(v, "column:") {
			column := strings.TrimPrefix(v, "column:")
			return column
		}
	}
	return ""
}
