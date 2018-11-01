package models

import (
	"testing"
)

func TestUserName(t *testing.T) {
	var u User

	u = User{MyName: "", DisplayName: "DisplayName1"}
	if u.Name() != "DisplayName1" {
		t.Fatalf("WRONG: Name, current: \"%v\"", u.Name())
	}

	u = User{MyName: "MyName2", DisplayName: "DisplayName2"}
	if u.Name() != "MyName2" {
		t.Fatalf("WRONG: Name, current: \"%v\"", u.Name())
	}
}
