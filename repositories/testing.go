package repositories

import (
	"github.com/jinzhu/gorm"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getDBMock() (db *gorm.DB, mock sqlmock.Sqlmock, err error) {
	rdb, mock, err := sqlmock.New()
	if err != nil {
		return
	}
	db, err = gorm.Open("mysql", rdb)
	if err != nil {
		return
	}
	return
}
