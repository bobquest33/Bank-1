package gmdb_test

import (
	"testing"
	"fmt"
	"util/gmdb"
)

func TestColl(t *testing.T) {
	coll := make(map[string]interface{})
	coll["a"] = 1
	coll["b"] = "tangs"
	coll["c"] = false
	coll["d"] = 3.2
	res, err := gmdb.MapI2MapS(coll)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
