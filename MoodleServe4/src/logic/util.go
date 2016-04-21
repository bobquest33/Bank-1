package logic

import (
	"net/http"
	"encoding/json"
	"util/log"
	"database/sql"
	"os"
	"fmt"
	"util/gmdb"
	"errors"
)

func OutPut(w http.ResponseWriter, status int, msg string, data interface{}) error {
	res := Result{ Status:status, Msg:msg, Data:data }
	rd, err := json.Marshal(res)
	if err != nil {
		log.AddError(err)
		return err
	}
	w.Write([]byte(rd))
	return nil
}

func FindId(db *sql.DB, table string, uid string, v ...*interface{}) (string, error) {
	dbCmd := "select Id from " + table + " where Id = '" + uid + "';"
	rows, err := db.Query(dbCmd)
	if err != nil {
		return "", err
	}
	id := ""
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return "", err
		}
	}
	return id, nil
}


func Struct2Map(v interface{}) (map[string]interface{}, error) {
	dt, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	err = json.Unmarshal(dt, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Map2Struct(data map[string]interface{}, v interface{}) (error) {
	dt, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dt, &v)
	if err != nil {
		return err
	}
	return nil
}

//Json struct convert to struct, the purpose is create a map to db.
func JS2S(v interface{}, k interface{}) (error) {
	vs, err := Struct2Map(v)
	if err != nil {
		return err
	}
	ks, err := Struct2Map(k)
	if err != nil {
		return err
	}
	for key, _ := range ks {
		ks[key] = vs[key]
	}
	err = Map2Struct(ks, &k)
	if err != nil {
		return err
	}
	return nil
}

//Json struct convert to struct, the purpose is create a map to db.
func JS2M(v interface{}, k interface{}) (map[string]interface{}, error) {
	vs, err := Struct2Map(v)
	if err != nil {
		return nil, err
	}
	ks, err := Struct2Map(k)
	if err != nil {
		return nil, err
	}
	for key, _ := range ks {
		ks[key] = vs[key]
	}
	return ks, nil
}

// map string key interface value to slice
func MskIv2S(v map[string]interface{}) []string {
	var names []string
	for k, _ := range v {
		names = append(names, k)
	}
	return names
}

func Guid() (string, error) {
	f, err := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	defer f.Close()
	if err != nil {
		return "", err
	}
	b := make([]byte, 16)
	f.Read(b)
	//uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	uuid := fmt.Sprintf("%x", b)
	return uuid, nil
}

func UGuid(db gmdb.DbController, table string) (string, error) {
	for i := 0 ;i < 3; i ++ {
		uuid, err := Guid()
		if err != nil {
			return "", err
		}
		do := gmdb.DbOpera{
			Table:table,
			Name:[]string{"id"},
		}
		do.FVW = make(map[string]interface{})
		do.FVW["id"] = uuid
		rows, err := db.Query(do)
		id := ""
		for rows.Next() {
			err = rows.Scan(&id)
			if err != nil {
				return "", err
			}
		}
		if id == "" {
			return uuid, nil
		}
	}
	return "", errors.New("Repetitive id appearence.")
}