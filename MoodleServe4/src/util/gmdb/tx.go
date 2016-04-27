package gmdb

import (
	"fmt"
	"database/sql"
	"errors"
	"strings"
	"strconv"
)

type Transaction struct {
	Tx		*sql.Tx
}

func GetTx() (Transaction, error) {
	trans, err := GetDb().Mdb.Begin()
	return Transaction{ Tx:trans }, err
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
func (tx Transaction) CopyRow(do DbOpera) (sql.Result, error) {
	if len(do.Name) < 1 {
		return nil, errors.New("invalid field to insert")
	}
	inf := strings.Join(do.Name, "`,`")
	inf = fmt.Sprintf("%s%s%s", "`", inf, "`")
	var sf string
	if len(do.SF) < 1 {
		sf = " * "
	} else {
		sf := strings.Join(do.SF, "`,`")
		sf = fmt.Sprintf("%s%s%s", "`", sf, "`")
	}
	var wherecols string
	var length int = len(do.FVW)
	var count int = 0
	msw, err := MapI2MapS(do.FVW)
	if err != nil {
		return nil, err
	} else {
		for k, v := range msw {
			wherecols += fmt.Sprintf("%s%s%s=%s%s%s", "`", k, "`", "`", v, "``")
			if count != (length - 1) {
				wherecols += " and "
			}
			count ++
		}
	}
	cmd := fmt.Sprintf("insert into %s ( %s ) select %s from %s where %s", do.Table, inf, sf, do.FTable, wherecols)
	fmt.Println(cmd)
	if stmt, err := tx.Tx.Prepare(cmd); err != nil {
		return nil, err
	} else {
		if res, err := stmt.Exec(); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}
	return nil, err
}

func (tx Transaction) Copy(do []DbOpera) (sql.Result, error) {
	var res sql.Result
	var err error
	for _, v := range do {
		if res, err = tx.CopyRow(v); err != nil {
			return res, err
		}
	}
	return res, err
}

// database opera multi, only insert, update isn't allow to update multi
type DbOM struct {
	Table		string
	Name		[]string
	Value		[][]interface{}
}
func (tx Transaction) InsertMulti(dbom DbOM) (sql.Result, error) {
	var inf, inv string
	fields := strings.Join(dbom.Name, "`,`")
	inf = fmt.Sprintf("%s%s%s", "`", fields, "`")
	vls, err := I2S(dbom.Value)
	if err != nil {
		return nil, err
	}
	length := len(vls)
	for ti, tv := range vls {
		inv = strings.Join(tv, "','")
		inv = fmt.Sprintf("%s%s%s", "'", inv, "'")
		if ti != (length - 1) {
			inv += "),("
		}
	}
	cmd := fmt.Sprintf("insert into %s ( %s ) values ( %s )", dbom.Table, inf, inv)
	fmt.Println(cmd)
	if stmt, err := tx.Tx.Prepare(cmd); err != nil {
		return nil, err
	} else {
		if res, err := stmt.Exec(); err != nil {
			return res, err
		} else {
			return res, nil
		}
	}
}

func I2S(v [][]interface{}) ([][]string, error) {
	res := [][]string{}
	for _, v := range v {
		rt := []string{}
		for _, vv := range v {
			switch vv.(type) {
			case string:
				rt = append(rt, vv.(string))
			case int:
				tyv := strconv.Itoa(vv.(int))
				rt = append(rt, tyv)
				break
			case float32:
				tyv := strconv.FormatFloat(float64(vv.(float32)), 'e', -1, 32)
				rt = append(rt, tyv)
				break
			case float64:
				tyv := strconv.FormatFloat(vv.(float64), 'e', -1, 32)
				rt = append(rt, tyv)
				break
			case bool:
				tyv := strconv.FormatBool(vv.(bool))
				rt = append(rt, tyv)
				break
			default:
				return nil, errors.New("Unknow type of the value")
			}
		}
		res = append(res, rt)
	}
	return res, nil
}