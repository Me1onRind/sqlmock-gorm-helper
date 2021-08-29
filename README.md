# sqlmock-rows-helper
transfrom single or multi struct value to (https://github.com/DATA-DOG/go-sqlmock) sqlmock.rows

# usage
<b>support struct value, struct point, struct's slice, struct point's slice </b>

# example
```go
import (
    "database/sql"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    sqlmock_rows_helper "github.com/Me1onRind/sqlmock-rows-helper"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type TestTab struct {
    ID    uint64 `gorm:"column:id"`
    Name  string `gorm:"column:name"`
    CTime uint32 `gorm:"column:ctime"`
    MTime uint32 `gorm:"column:mtime"`
}

func newTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
    db, mock, err := sqlmock.New()
    if err != nil {
        return nil, nil, err
    }
    mock.ExpectQuery("SELECT VERSION").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("5.7.32"))
    gormDB, err := NewDBConnectPoolFRromDB(db)
    if err != nil {
        return nil, nil, err
    }
    return gormDB, mock, nil
}

func NewDBConnectPoolFRromDB(db *sql.DB) (*gorm.DB, error) {
    return doCreateDBConnectPool(mysql.New(mysql.Config{
        Conn: db, 
    })) 
}

func doCreateDBConnectPool(dial gorm.Dialector) (*gorm.DB, error) {
    db, err := gorm.Open(dial, &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })  
    if err != nil {
        return nil, err 
    }   

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err 
    }

    registerPlugin(db)

    return db, err 
}

func Test_Select(t *testing.T) {
    db, mock, err := newTestDB()
    if !assert.Empty(t, err) {
        return
    }
    mock.ExpectQuery("SELECT").WillReturnRows(sqlmock_rows_helper.ModelToRows(
        &TestTab{
            ID:    1,
            Name:  "test",
            CTime: 1630250445,
            MTime: 1630250445,
        },
    ))
    if !assert.Empty(t, err) {
        return
    }
    testTab := &TestTab{}
    if err := db.WithContext(context.Background()).Where("id = ?", 1).Find(testTab).Error; err != nil {
        if !assert.Empty(t, err) {
            return
        }
    }
    assert.Equal(t, uint64(1), testTab.ID)
    assert.Equal(t, "test", testTab.Name)
    assert.Equal(t, uint32(1630250445), testTab.CTime)
}
```
