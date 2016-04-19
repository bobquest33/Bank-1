package gmdb_test

import (
	"testing"
	. "util/gmdb"
	"fmt"
)

func TestDbo(t *testing.T) {
	db, err := DbConnect()
	if err != nil {
		t.Error(err)
	}
	var name []string
	fv := make(map[string]interface{})
	fvw := make(map[string]interface{})
	nequal := make(map[string]string)
	name = append(name, "name", "digit")
	fv["name"] = "tangs"
	fv["digit"] = 99
	do := DbOpera{
		Name:name,
		Table:D_T,
		FV:fv,
		FVW:fvw,
		NEqual:nequal,
	}
	//--------------------- one ----------------------------------
	// use table and fv
	_, err = db.Insert(do)
	if err != nil {
		t.Error(err)
	}

	//use table ,name and fvw
	do.FVW["name"] =  "tangs"
	rows, err := db.Query(do)
	if err != nil {
		t.Error(err)
	}
	var v1 string
	var v2 int
	for rows.Next() {
		err = rows.Scan(&v1, &v2)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("v1 value is %s, v2 value is %d\n", v1, v2)
	}

	//use table fv and fvw
	do.FV["name"] = "lily"
	do.FVW["id"] = 168
	_, err = db.Update(do)
	if err != nil {
		t.Error(err)
	}

	//use table, name and fvw
	do.FVW["name"] = "lily"
	rows, err = db.Query(do)
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		err = rows.Scan(&v1, &v2)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("v1 value is %s, v2 value is %d\n", v1, v2)
	}

	//use table and fvw
	do.FVW["name"] = "lily"
	_, err = db.Delete(do, false)
	if err != nil {
		t.Error(err)
	}

	//--------------------- two ----------------------------------
	//use table and fv
	do.FV["name"] = "lily"
	_, err = db.Insert(do)
	if err != nil {
		t.Error(err)
	}

	//use table, name and fvw
	var id int
	do.FVW["name"] = "lily"
	do.Name = []string{}
	rows, err = db.Query(do)
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &v1, &v2)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("id value is %d, v1 value is %s, v2 value is %d\n", id, v1, v2)
	}

	//use table and fvw
	do.FVW["name"] = "lily"
	_, err = db.Delete(do,false)
	if err != nil {
		t.Error(err)
	}
	//--------------------- three ----------------------------------
	do.FV["name"] = "tangs"
	_, err = db.Insert(do)
	if err != nil {
		t.Error(err)
	}
	do.FVW = make(map[string]interface{})
	do.FVW["id"] = 1
	do.NEqual["id"] = ">"
	do.Name = []string{"name", "digit"}
	rows, err = db.Query(do)
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		err = rows.Scan(&v1, &v2)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("v1 and v2 is : ", v1, v2)
	}
	do.FVW["name"] = "tangs"
	_, err = db.Delete(do, false)
	if err != nil {
		t.Error(err)
	}
}