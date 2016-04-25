package logic

import (
	"util/gmdb"
	"errors"
	"fmt"
	"util/log"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// equal true sign if != { return }
// equal false sign if == { return }
type Mapping struct {
	equal		bool
	mapping		map[string]interface{}
	name 		[]string
}

//mapping string to db name, (map[string]interface) to fvw
func F(prefix string, db gmdb.DbController, mapping map[string](Mapping), status int) (int, error) {
	for dbn, mpi := range mapping {
		do := gmdb.DbOpera{ Table:dbn, Name:mpi.name, FVW:mpi.mapping }
		rows, err := db.Query(do)
		if err != nil {
			return status, err
		}
		id := ""
		for rows.Next() {
			err = rows.Scan(&id)
			if err != nil {
				status ++
				return status, err
			}
		}
		if mpi.equal {
			if id != "" {
				status ++
				var kv string
				for k, v := range mpi.mapping {
					if v == id {
						kv = k
						break
					}
				}
				return status, errors.New(fmt.Sprintf("%s but %s is already exist", prefix, kv))
			}
		} else {
			if id == "" {
				status ++
				var kv string
				for k, v := range mpi.mapping {
					if v == id {
						kv = k
						break
					}
				}
				return status, errors.New(fmt.Sprintf("%s but %s is not exist", prefix, kv))
			}
		}
	}
	return 200, nil
}

func mapMapping(ump *UMP) {
	var err error
	db := gmdb.GetDb()
	if ump.EbMap, err = GID(db, gmdb.D_1); err != nil {
		log.AddError(err)
		panic (err)
	}
	if ump.PgMap, err = GID(db, gmdb.D_2); err != nil {
		log.AddError(err)
		panic (err)
	}
	if ump.PMap, err = GID(db, gmdb.D_3); err != nil {
		log.AddError(err)
		panic (err)
	}
	if ump.QgMap, err = GID(db, gmdb.D_4); err != nil {
		log.AddError(err)
		panic (err)
	}
	if ump.QMap, err = GID(db, gmdb.D_5); err != nil {
		log.AddError(err)
		panic (err)
	}
	if ump.PqMap, err = GID(db, gmdb.D_6); err != nil {
		log.AddError(err)
		panic (err)
	}
}

//get id mapping to database
func GID(db gmdb.DbController, table string) (map[string]CI, error) {
	var m map[string]CI
	m = make(map[string]CI)
	do := gmdb.DbOpera{ Table:table, Name:[]string{"id"} }
	if rows, err := db.Query(do); err == nil {
		for rows.Next() {
			var id string
			if err = rows.Scan(&id); err == nil {
				m[id] = CI{ C:false, I:true }
			} else {
				return nil, err
			}
		}
	}
	return m, nil
}

func UnmarshalJ(r *http.Request, v interface{}) ([]byte, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return data, err
	}
	if err = json.Unmarshal([]byte(data), v); err != nil {
		return data, err
	}
	return data, nil
}