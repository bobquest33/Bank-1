package gmdb_test

import (
	"testing"
	"util/gmdb"
	"fmt"
)

func TestTxo(t *testing.T) {
	db, err := gmdb.DbConnect()
	if err != nil {
		t.Error(err)
	}
	tt, err := db.Mdb.Begin()
	if err != nil {
		t.Error(err)
	}
	tx := gmdb.Transaction{	Tx : tt	}
	defer tx.Tx.Rollback()
	var name []string
	fv := make(map[string]interface{})
	fvw := make(map[string]interface{})
	name = append(name, "name", "digit")
	fv["name"] = "tangs"
	fv["digit"] = 99
	do := gmdb.DbOpera{
		Name:name,
		Table:gmdb.D_T,
		FV:fv,
		FVW:fvw,
	}

	//-------------- one -------------------------------
	//use table and fv
	_, err = tx.Insert(do)
	if err != nil {
		t.Error(err)
	}

	//use table name and fvw
	rows, err := tx.Query(do)
	if err != nil {
		t.Error(err)
	}
	var v1 string
	var v2 int
	//var id int
	for rows.Next() {
		err = rows.Scan(&v1, &v2)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("tx test! v1 = %s,v2 = %d\n", v1, v2)
	}

	//use table fv and fvw
	do.FV["name"] = "lily"
	do.FVW["name"] = "tangs"
	_, err = tx.Update(do)
	if err != nil {
		t.Error(err)
	}

	//use table and fvw
	do.FVW["name"] = "lily"
	_,err = tx.Delete(do, false)
	if err != nil {
		t.Error(err)
	}
	tx.Tx.Commit()
}