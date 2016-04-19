package gmdb

import (
	"fmt"
	"strings"
	"errors"
	"database/sql"
)

type Row struct {
			   // One of these two will be non-nil:
	err  error // deferred error for easy chaining
	rows *Rows
}

func (d DbController) Insert(do DbOpera) (sql.Result, error) {
	if do.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}
	var sels string
	var vals string
	ms, err := MapI2MapS(do.FV)
	if err != nil {
		return nil, err
	}
	length := len(ms)
	count := 0
	for k, v := range ms {
		sels = sels + fmt.Sprintf("%s%s%s", "`", k, "`")
		vals = vals + fmt.Sprintf("%s%s%s", "'", v, "'")
		if count != (length - 1) {
			sels += ","
			vals += ","
		}
		count++
	}

	cmd := fmt.Sprintf("Insert into %s ( %s ) value ( %s )", do.Table, sels, vals)
	//fmt.Println(cmd)
	stmt, err := d.Mdb.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d DbController) Update(do DbOpera) (sql.Result, error) {
	if do.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}
	if len(do.FV) < 1 {
		return nil, errors.New("FV fields is invalid or it is null")
	}
	var fv []string
	ms, err := MapI2MapS(do.FV)
	if err != nil {
		return nil, err
	}
	for k, v := range ms {
		str := fmt.Sprintf("%s%s%s = %s%s%s", "`", k, "`", "'", v, "'")
		fv = append(fv, str)
	}
	setCols := strings.Join(fv, " , ")

	if len(do.FVW) < 1 {
		return nil, errors.New("FVW field is invalid, maybe forget ?")
	}
	msw, err := MapI2MapS(do.FVW)
	if err != nil {
		return nil, err
	}
	var fvw []string
	for k, v := range msw{
		str := fmt.Sprintf("%s%s%s = %s%s%s", "`", k, "`", "'", v, "'")
		fvw = append(fvw, str)
	}
	whereCols := strings.Join(fvw, " and ")

	cmd := fmt.Sprintf("Update %s Set %s where %s", do.Table, setCols, whereCols)
	//fmt.Println(cmd)
	stmt, err := d.Mdb.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d DbController) Delete(do DbOpera, sure interface{}) (sql.Result, error) {
	if do.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}
	ifWhere := true
	length := len(do.FVW)
	var su bool
	var ok bool
	if sure == nil {
		su = false
	} else {
		su, ok = sure.(bool)
		if !ok {
			return nil, errors.New("Unknow argument type")
		}
	}
	if length< 1 && !su {
		return nil, errors.New("Delete from " + do.Table + " is not safe, are you sure do that?")
	} else if len(do.FVW) < 1 && su {
		ifWhere = false
	}
	whereCols := ""
	cmd := ""
	msw, err := MapI2MapS(do.FVW)
	if err != nil {
		return nil, err
	}
	if ifWhere {
		count := 0
		for k, v := range msw {
			whereCols += fmt.Sprintf("%s%s%s = %s%s%s", "`", k, "`", "'", v, "'")
			if count != (length - 1){
				whereCols += " and "
			}
			count ++
		}
		cmd = fmt.Sprintf("Delete from %s where %s", do.Table, whereCols)
	} else {
		cmd = fmt.Sprintf("Delete from %s", do.Table)
	}
	stmt, err := d.Mdb.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d DbController) Query(do DbOpera) (*sql.Rows, error) {
	if do.Table == "" {
		return nil, errors.New("table is invalid or it is null")
	}
	sels := ""
	if len(do.Name) < 1 {
		sels = "*"
	} else {
		selt := strings.Join(do.Name, "`,`")
		sels = fmt.Sprintf("%s%s%s", "`", selt, "`")
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
			opType, ok := do.NEqual[k]
			if !ok {
				wherecols += fmt.Sprintf("%s%s%s=%s%s%s", "`", k, "`", "'", v, "'")
			} else {
				wherecols += fmt.Sprintf("%s%s%s%s%s%s%s", "`", k, "`", opType, "'", v, "'")
			}
			if count != (length - 1) {
				wherecols += " and "
			}
			count ++
		}
		cmd = fmt.Sprintf("select %s from %s where %s", sels, do.Table, wherecols)
	}
	//fmt.Println(cmd)
	rows, err := d.Mdb.Query(cmd)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

//write form sql.go
//rows, err := db.Query(query, args...)
//return &Row{rows: rows, err: err}
func (d DbController) QueryRow(do DbOpera) (*sql.Row) {
	//rows, err := d.Query(do)
	//return &sql.Row{rows:rows, err:err}
	return &sql.Row{}
}

//select from where select
//嵌套查询 最好不要超过三层,否则数据库解析会非常的复杂,甚至会解析出错误的数据
func (d DbController) SS(do DbOpera) {
}