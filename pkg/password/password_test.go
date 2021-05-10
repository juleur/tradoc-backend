package password

import (
	"testing"
)

func TestHashedPassword(t *testing.T) {
	pwd := "troubadour21"
	hash, err := HashPassword(pwd)
	if err != nil {
		t.Errorf(err.Message)
	}
	t.Log(hash)
}
