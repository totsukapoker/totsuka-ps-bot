package repositories

import (
	"testing"
)

func TestNewTransactionRepository(t *testing.T) {
	db, _, err := getDBMock()
	defer db.Close()
	if err != nil {
		t.Fatalf("%#v", err)
	}
	r := NewTransactionRepository(db)
	if r.db != db {
		t.Errorf("TransactionRepository.db = %#v; want: %#v", r.db, db)
	}
}

func TestTransactionRepository_FindByGame(t *testing.T) {
	t.Skip("implement me")
}

func TestTransactionRepository_LastBy(t *testing.T) {
	t.Skip("implement me")
}

func TestTransactionRepository_CurrentAmountBy(t *testing.T) {
	t.Skip("implement me")
}

func TestTransactionRepository_CurrentAmountBuyinBy(t *testing.T) {
	t.Skip("implement me")
}

func TestTransactionRepository_Create(t *testing.T) {
	t.Skip("implement me")
}

func TestTransactionRepository_Delete(t *testing.T) {
	t.Skip("implement me")
}
