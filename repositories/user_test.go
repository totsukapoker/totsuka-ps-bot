package repositories

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
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
	//db.LogMode(true)
	return
}

func TestNewUserRepository(t *testing.T) {
	db, _, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("got unexpected error '%s'", err)
	}
	r := NewUserRepository(db)
	if r.db != db {
		t.Errorf("got: %v, expected: %v", r.db, db)
	}
}

func TestUserRepository_FirstOrCreate(t *testing.T) {
	t.Run("it creates new record", func(t *testing.T) {
		t.Skip("implement me")
	})
	t.Run("it returns exist record", func(t *testing.T) {
		t.Skip("implement me")
	})
}

func TestUserRepository_FindByIDs(t *testing.T) {
	db, mock, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("got unexpected error '%s'", err)
	}
	r := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND ((ID in (?,?,?)))",
	)).WithArgs(1, 2, 3).WillReturnError(sql.ErrNoRows)

	r.FindByIDs([]uint{1, 2, 3})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_SetMyName(t *testing.T) {
	t.Skip("implement me")
}
