package gmdb_test

import (
	"testing"
	. "util/gmdb"
)


func TestDbO(t *testing.T) {
	_, err := DbConnect()
	if err != nil {
		t.Error(err)
	}
}