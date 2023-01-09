package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "gorm.io/driver/mysql"
	"testing"
	"time"
)

type sqlTestSuite struct {
	suite.Suite

	// 配置字段
	driver string
	dsn    string

	db *sql.DB
}

func (s *sqlTestSuite) TearDownTest() {
	_, err := s.db.Exec("delete from test_model;")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *sqlTestSuite) SetupSuite() {
	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_model(
			id INTEGER PRIMARY KEY ,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			age INTEGER
		)
	`)

	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *sqlTestSuite) TestCRUD() {

	var t = s.T()

	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, "INSERT INTO `test_model` (`id`,`first_name`,`last_name`,`age`) values (1,'Tom','Jerry',18)")
	if err != nil {
		t.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}

	if affected != 1 {
		t.Fatal(errors.New("影响行数为0"))
	}

	// 查询集合
	rows, err := db.QueryContext(ctx, "SELECT * FROM test_model LIMIT ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		tm := &TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.LastName, &tm.Age)
		if err != nil {
			rows.Close()
			t.Fatal(err)
		}
		assert.Equal(t, "Tom", tm.FirstName)
	}
	rows.Close()

	res, err = db.ExecContext(ctx, "update test_model set `first_name`='changed' where `id`=?", 1)
	if err != nil {
		t.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if rowsAffected != 1 {
		t.Fatal(errors.New("没有影响到数据"))
	}

	// 查询单行
	row := db.QueryRowContext(ctx, "select * from test_model limit 1")
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	tm := &TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.LastName, &tm.Age)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "changed", tm.FirstName)
}

type TestModel struct {
	Id        int64  `gorm:"auto_increment,primary_key"`
	FirstName string `gorm:"varchar(20)"`
	LastName  string `gorm:"varchar(50)"`
	Age       int8
}

func TestMysql(t *testing.T) {
	suite.Run(t, &sqlTestSuite{
		driver: "mysql",
		dsn:    "root:123456@tcp(192.168.1.183:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local",
	})
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(0)
	fmt.Println(timer.Stop())
	<-timer.C
}
