package sqlmock_rows_helper

import (
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID         uint64 `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Foobar     string `gorm:"-"`
	CreateTime uint32 `gorm:"column:create_time"`
}

func Test_ColumnsFromGormModleField(t *testing.T) {
	columns := columnsFromModelType(reflect.TypeOf(TestModel{}))
	assert.Equal(t, []string{"id", "name", "create_time"}, columns)
}

func Test_ValueFromField(t *testing.T) {
	model := TestModel{
		ID:         10,
		Name:       "test",
		CreateTime: 1630248918,
	}
	value := valuesFromModel(reflect.TypeOf(model), reflect.ValueOf(model))
	assert.Equal(t, []driver.Value{uint64(10), "test", uint32(1630248918)}, value)
}

func Test_SingleRow(t *testing.T) {
	model := TestModel{
		ID:         12,
		Name:       "test_abc",
		CreateTime: 1630248920,
	}

	targetRows := sqlmock.NewRows([]string{"id", "name", "create_time"}).
		AddRow(uint64(12), "test_abc", uint32(1630248920))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)

	rows = ModelToRows(&model)
	assert.Equal(t, targetRows, rows)
}

func Test_MultiStructValueRows(t *testing.T) {
	model := []TestModel{
		{
			ID:         12,
			Name:       "test_abc",
			CreateTime: 1630248920,
		},
		{
			ID:         13,
			Name:       "test_efg",
			CreateTime: 1630248922,
		},
	}

	targetRows := sqlmock.NewRows([]string{"id", "name", "create_time"}).
		AddRow(uint64(12), "test_abc", uint32(1630248920)).
		AddRow(uint64(13), "test_efg", uint32(1630248922))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)
}

func Test_MultiStructPtrRows(t *testing.T) {
	model := []*TestModel{
		{
			ID:         12,
			Name:       "test_abc",
			CreateTime: 1630248920,
		},
		{
			ID:         13,
			Name:       "test_efg",
			CreateTime: 1630248922,
		},
	}

	targetRows := sqlmock.NewRows([]string{"id", "name", "create_time"}).
		AddRow(uint64(12), "test_abc", uint32(1630248920)).
		AddRow(uint64(13), "test_efg", uint32(1630248922))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)
}
