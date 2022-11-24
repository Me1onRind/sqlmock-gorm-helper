package sqlmock_rows_helper

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type BaseModel struct {
	ID         uint64 `gorm:"column:id;primary"`
	CreateTime uint32 `gorm:"column:create_time"`
}

type TestModel struct {
	BaseModel
	Name        string     `gorm:"column:name"`
	Foobar      string     `gorm:"-"`
	CustomField Uint8Slice `gorm:"column:custom_field"`
}

type Uint8Slice []uint8

func (u *Uint8Slice) Scan(value interface{}) error {
	valueBytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal Uint8Slice value:", valueBytes))
	}

	if len(valueBytes) > 0 {
		bytesSlice := bytes.Split(valueBytes, []byte{','})
		for _, v := range bytesSlice {
			intValue, _ := strconv.ParseUint(string(v), 10, 32)
			*u = append(*u, uint8(intValue))
		}
	} else {
		*u = []uint8{}
	}

	return nil
}

func (u Uint8Slice) Value() (driver.Value, error) {
	if len(u) == 0 {
		return "", nil
	}

	var tmp []string
	for _, v := range u {
		tmp = append(tmp, strconv.Itoa(int(v)))
	}

	return strings.Join(tmp, ","), nil
}

func Test_ColumnsFromGormModleField(t *testing.T) {
	columns := columnsFromModelType(reflect.TypeOf(TestModel{}))
	assert.Equal(t, []string{"id", "create_time", "name", "custom_field"}, columns)
}

func Test_ValueFromField(t *testing.T) {
	model := TestModel{
		BaseModel: BaseModel{
			ID:         10,
			CreateTime: 1630248918,
		},
		Name:        "test",
		CustomField: []uint8{1, 2},
	}
	value := valuesFromModel(reflect.TypeOf(model), reflect.ValueOf(model))
	assert.Equal(t, []driver.Value{uint64(10), uint32(1630248918), "test", []byte("1,2")}, value)
}

func Test_SingleRow(t *testing.T) {
	model := TestModel{
		BaseModel: BaseModel{
			ID:         12,
			CreateTime: 1630248920,
		},
		Name:        "test_abc",
		CustomField: []uint8{1, 2},
	}

	targetRows := sqlmock.NewRows([]string{"id", "create_time", "name", "custom_field"}).
		AddRow(uint64(12), uint32(1630248920), "test_abc", []byte("1,2"))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)

	rows = ModelToRows(&model)
	assert.Equal(t, targetRows, rows)
}

func Test_MultiStructValueRows(t *testing.T) {
	model := []TestModel{
		{
			BaseModel: BaseModel{
				ID:         12,
				CreateTime: 1630248920,
			},
			Name:        "test_abc",
			CustomField: []uint8{1, 2},
		},
		{
			BaseModel: BaseModel{
				ID:         13,
				CreateTime: 1630248922,
			},
			Name: "test_efg",
		},
	}

	targetRows := sqlmock.NewRows([]string{"id", "create_time", "name", "custom_field"}).
		AddRow(uint64(12), uint32(1630248920), "test_abc", []byte("1,2")).
		AddRow(uint64(13), uint32(1630248922), "test_efg", []byte(""))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)
}

func Test_MultiStructPtrRows(t *testing.T) {
	model := []*TestModel{
		{
			BaseModel: BaseModel{
				ID:         12,
				CreateTime: 1630248920,
			},
			Name:        "test_abc",
			CustomField: []uint8{},
		},
		{
			BaseModel: BaseModel{
				ID:         13,
				CreateTime: 1630248922,
			},
			Name: "test_efg",
		},
	}

	targetRows := sqlmock.NewRows([]string{"id", "create_time", "name", "custom_field"}).
		AddRow(uint64(12), uint32(1630248920), "test_abc", []byte("")).
		AddRow(uint64(13), uint32(1630248922), "test_efg", []byte(""))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)
}

func Test_InsertSql(t *testing.T) {
	sql := InsertSql(TestModel{}, "test_table")
	wantSql := "INSERT INTO `test_table` \\(`create_time`,`name`,`custom_field`\\) VALUES \\(\\?,\\?,\\?\\)"
	assert.Equal(t, sql, wantSql)

	sql = InsertSql(&TestModel{}, "test_table")
	assert.Equal(t, sql, wantSql)
}
