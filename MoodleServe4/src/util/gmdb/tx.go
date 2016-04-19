package gmdb

import (
	"fmt"
	"database/sql"
	"errors"
	"strings"
)

type Transaction struct {
	Tx		*sql.Tx
}

func (tx Transaction) Insert (d DbOpera) (sql.Result, error) {
	if d.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}
	var sels string
	var vals string
	length := len(d.FV)
	count := 0
	ms, err := MapI2MapS(d.FV)
	if err != nil {
		return nil, err
	}
	for k, v := range ms {
		sels = sels + fmt.Sprintf("%s%s%s", "`", k, "`")
		vals = vals + fmt.Sprintf("%s%s%s", "'", v, "'")
		if count != (length - 1) {
			sels += ","
			vals += ","
		}
		count++
	}

	cmd := fmt.Sprintf("Insert into %s ( %s ) value ( %s )", d.Table, sels, vals)
	//fmt.Println(cmd)
	stmt, err := tx.Tx.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (tx Transaction) Delete (d DbOpera, sure bool) (sql.Result, error) {
	if d.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}
	ifWhere := true
	length := len(d.FVW)
	if length< 1 && !sure {
		return nil, errors.New("Delete from " + d.Table + " is not safe, are you sure do that?")
	} else if len(d.FVW) < 1 && sure {
		ifWhere = false
	}
	whereCols := ""
	cmd := ""
	if ifWhere {
		count := 0
		msw, err := MapI2MapS(d.FVW)
		if err != nil {
			return nil, err
		}
		for k, v := range msw {
			whereCols += fmt.Sprintf("%s%s%s = %s%s%s", "`", k, "`", "'", v, "'")
			if count != (length - 1){
				whereCols += " and "
			}
			count ++
		}
		cmd = fmt.Sprintf("Delete from %s where %s", d.Table, whereCols)
	} else {
		cmd = fmt.Sprintf("Delete from %s", d.Table)
	}
	stmt, err := tx.Tx.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (tx Transaction) Update (d DbOpera) (sql.Result, error) {
	if d.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}
	if len(d.FV) < 1 {
		return nil, errors.New("FV fields is invalid or it is null")
	}
	var fv []string
	ms, err := MapI2MapS(d.FV)
	if err != nil {
		return nil, err
	}
	for k, v := range ms {
		str := fmt.Sprintf("%s%s%s = %s%s%s", "`", k, "`", "'", v, "'")
		fv = append(fv, str)
	}
	setCols := strings.Join(fv, " , ")

	if len(d.FVW) < 1 {
		return nil, errors.New("FVW field is invalid, maybe forget ?")
	}
	var fvw []string
	msw, err := MapI2MapS(d.FVW)
	if err != nil {
		return nil, err
	}
	for k, v := range msw {
		str := fmt.Sprintf("%s%s%s = %s%s%s", "`", k, "`", "'", v, "'")
		fvw = append(fvw, str)
	}
	whereCols := strings.Join(fvw, " and ")

	cmd := fmt.Sprintf("Update %s Set %s where %s", d.Table, setCols, whereCols)
	//fmt.Println(cmd)
	stmt, err := tx.Tx.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (tx Transaction) InsertMulti() {
	//InsertMulti
}

func (tx Transaction) Query(do DbOpera) (*sql.Rows, error) {
	if do.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}

	sels := ""
	if len(do.Name) < 1 {
		sels = "*"
	} else {
		selt := strings.Join(do.Name, "`,`")
		sels += fmt.Sprintf("%s%s%s", "`", selt, "`")
	}
	wherecols := ""
	var cmd string
	length := len(do.FVW)
	if length < 1 {
		cmd = fmt.Sprintf("select %s from %s", sels, do.Table)
	} else {
		count := 0
		msw, err := MapI2MapS(do.FVW)
		if err != nil {
			return nil, err
		}
		for k, v := range msw {
			wherecols += fmt.Sprintf("%s%s%s=%s%s%s", "`", k, "`", "'", v, "'")
			if count != (length - 1) {
				wherecols += " and "
			}
			count ++
		}
		cmd = fmt.Sprintf("select %s from %s where %s", sels, do.Table, wherecols)
	}
	rows, err := tx.Tx.Query(cmd)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
func (tx Transaction) QueryRow(do DbOpera) *sql.Row {
	//rows, err := tx.Query(do)
	//return &sql.Row{rows: rows, err: err}
	return &sql.Row{}
}