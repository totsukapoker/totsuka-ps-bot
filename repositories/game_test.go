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
		t.Fatalf("%#v", err)
	}
	if r := NewGameRepository(db); r.db != db {
		t.Errorf("GameRepository.db = %#v; want: %#v", r.db, db)
	}
}

func TestGameRepository_First(t *testing.T) {
	db, mock, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("%#v", err)
	}
	r := NewGameRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `games` WHERE `games`.`deleted_at` IS NULL AND ((`games`.`id` = 89)) ORDER BY `games`.`id` ASC LIMIT 1",
	)).WillReturnError(sql.ErrNoRows)

	r.First(89)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %#v", err)
	}
}

func TestGameRepository_Current(t *testing.T) {
	db, mock, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("%#v", err)
	}
	r := NewGameRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `games` WHERE `games`.`deleted_at` IS NULL AND ((? BETWEEN started_at AND ended_at)) ORDER BY `games`.`id` ASC LIMIT 1",
	))

	r.Current()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %#v", err)
	}
}
