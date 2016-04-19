package logic

import (
	"util/gmdb"
	"errors"
	"fmt"
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

