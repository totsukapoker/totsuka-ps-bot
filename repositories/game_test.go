package repositories

import (
	"database/sql"
	"regexp"
	"testing"
)

func TestNewGameRepository(t *testing.T) {
	db, _, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("got unexpected error '%s'", err)
	}
	r := NewGameRepository(db)
	if r.db != db {
		t.Errorf("got: %v, expected: %v", r.db, db)
	}
}

func TestGameRepository_First(t *testing.T) {
	db, mock, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("got unexpected error '%s'", err)
	}
	r := NewGameRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `games` WHERE `games`.`deleted_at` IS NULL AND ((`games`.`id` = 89)) ORDER BY `games`.`id` ASC LIMIT 1",
	)).WillReturnError(sql.ErrNoRows)

	r.First(89)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGameRepository_Current(t *testing.T) {
	db, mock, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("got unexpected error '%s'", err)
	}
	r := NewGameRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `games` WHERE `games`.`deleted_at` IS NULL AND ((? BETWEEN started_at AND ended_at)) ORDER BY `games`.`id` ASC LIMIT 1",
	))

	r.Current()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
