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

type TestModel struct {
	ID          uint64     `gorm:"column:id"`
	Name        string     `gorm:"column:name"`
	Foobar      string     `gorm:"-"`
	CustomField Uint8Slice `gorm:"column:custom_field"`
	CreateTime  uint32     `gorm:"column:create_time"`
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
	assert.Equal(t, []string{"id", "name", "custom_field", "create_time"}, columns)
}

func Test_ValueFromField(t *testing.T) {
	model := TestModel{
		ID:          10,
		Name:        "test",
		CustomField: []uint8{1, 2},
		CreateTime:  1630248918,
	}
	value := valuesFromModel(reflect.TypeOf(model), reflect.ValueOf(model))
	assert.Equal(t, []driver.Value{uint64(10), "test", []byte("1,2"), uint32(1630248918)}, value)
}

func Test_SingleRow(t *testing.T) {
	model := TestModel{
		ID:          12,
		Name:        "test_abc",
		CustomField: []uint8{1, 2},
		CreateTime:  1630248920,
	}

	targetRows := sqlmock.NewRows([]string{"id", "name", "custom_field", "create_time"}).
		AddRow(uint64(12), "test_abc", []byte("1,2"), uint32(1630248920))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)

	rows = ModelToRows(&model)
	assert.Equal(t, targetRows, rows)
}

func Test_MultiStructValueRows(t *testing.T) {
	model := []TestModel{
		{
			ID:          12,
			Name:        "test_abc",
			CustomField: []uint8{1, 2},
			CreateTime:  1630248920,
		},
		{
			ID:         13,
			Name:       "test_efg",
			CreateTime: 1630248922,
		},
	}

	targetRows := sqlmock.NewRows([]string{"id", "name", "custom_field", "create_time"}).
		AddRow(uint64(12), "test_abc", []byte("1,2"), uint32(1630248920)).
		AddRow(uint64(13), "test_efg", []byte(""), uint32(1630248922))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)
}

func Test_MultiStructPtrRows(t *testing.T) {
	model := []*TestModel{
		{
			ID:          12,
			Name:        "test_abc",
			CreateTime:  1630248920,
			CustomField: []uint8{},
		},
		{
			ID:         13,
			Name:       "test_efg",
			CreateTime: 1630248922,
		},
	}

	targetRows := sqlmock.NewRows([]string{"id", "name", "custom_field", "create_time"}).
		AddRow(uint64(12), "test_abc", []byte(""), uint32(1630248920)).
		AddRow(uint64(13), "test_efg", []byte(""), uint32(1630248922))

	rows := ModelToRows(model)
	assert.Equal(t, targetRows, rows)
}
