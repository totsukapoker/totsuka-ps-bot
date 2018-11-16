package models

import (
	"testing"
)

func TestUser_Name(t *testing.T) {
	tests := []struct {
		myName, displayName, want string
	}{
		{"", "", ""},
		{"", "DisplayName1", "DisplayName1"},
		{"MyName2", "", "MyName2"},
		{"MyName3", "DisplayName3", "MyName3"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			u := &User{MyName: tt.myName, DisplayName: tt.displayName}
			if u.Name() != tt.want {
				t.Errorf("user.Name() = %#v; want: %#v", u.Name(), tt.want)
			}
		})
	}
}
