package logic

import (
	"math"
	"util/gmdb"
	"util/log"
	"fmt"
)

func Analyze(dig int) []int {
	res := []int{}
	if dig < 1 {
		res = append(res, 0)
		return res
	}
	var digit int = dig
	m := math.Log2(float64(dig))
	max := int(m)
	for i := max; i > -1; i -- {
		d := math.Pow(2, float64(i))
		if digit >= int(d) {
			res = append(res, int(d))
			digit -= int(d)
		}
	}
	return res
}

func IfRelease(dig int, audit []int) bool {
	if len(audit) == 0 {
		return false
	}
	if dig == audit[0] {
		return true
	}
	return false
}
func IfRelPG(audit []int) bool {
	if len(audit) == 0 {
		return false
	} else {
		for _, v := range audit {
			if v == Rel_1 || v == Rel_2 || v == Rel_4 {
				return true
			}
		}
	}
	return false
}

func DupPg(tx gmdb.Transaction, id string) (string, error) {
	pg := Paper_Grp{}
	fvw := make(map[string]interface{})
	fvw["id"] = id
	var idPg string
	do := gmdb.DbOpera { Table:gmdb.D_2, Name:[]string{"name", "type", "exam_bank_id", "remark", "status"}, FVW:fvw }
	if rows, err := tx.Query(do); err != nil {
		log.AddError(err, do)
		return idPg, err
	} else {
		db := gmdb.GetDb()
		for rows.Next() {
			if err = rows.Scan(&pg.Name,&pg.Type,&pg.Exam_Bank_Id,&pg.Remark,&pg.Status); err != nil {
				log.AddError(err, rows, do)
				return idPg, err
			}
		}
		if pg.Id, err = UGuid(db, gmdb.D_2); err != nil {
			log.AddError(err)
			return idPg, err
		}
		idPg = pg.Id
	}
	if fv, err := JS2M(pg, pg); err != nil {
		log.AddError(err, pg)
		return idPg, err
	} else {
		do := gmdb.DbOpera {Table:gmdb.D_2, FV:fv }
		if _, err = tx.Insert(do); err != nil {
			log.AddError(err, do, pg)
			return idPg, err
		}
	}
	return idPg, nil
}

func UpdPg(tx gmdb.Transaction, id string, fv map[string]interface{}, fvw map[string]interface{}) (error) {
	do := gmdb.DbOpera{
		Table:gmdb.D_2,
		FV:fv,
		FVW:fvw,
	}
	if _, err := tx.Update(do); err != nil {
		log.AddError(err, do)
		return err
	}
	return nil
}

func DupP(tx gmdb.Transaction, idPg string) ([]string, error) {
	fvw := make(map[string]interface{})
	fvw["Paper_Grp_Id"] = idPg
	var ids []string
	var values [][]interface{}
	values = make([][]interface{}, 0)
	do := gmdb.DbOpera{ Table:gmdb.D_3, Name:[]string{"name","Paper_Grp_Id","type","ver","Create_Time","Author","Composed_Time","Remark","Status"}, FVW:fvw }
	if rows, err := tx.Query(do); err != nil {
		log.AddError(err,do)
		return ids, err
	} else {
		db := gmdb.GetDb()
		var count int
		for rows.Next() {
			var p Paper
			if err = rows.Scan(&p.Name,&p.Paper_Grp_Id,&p.Type,&p.Ver,&p.Create_Time,&p.Author,&p.Composed_Time,&p.Remark,&p.Status); err != nil {
				log.AddError(err, rows, values)
				return ids, err
			}
			if p.Id, err = UGuid(db, gmdb.D_3); err != nil {
				log.AddError(err, values)
				return ids, err
			}
			ids = append(ids, p.Id)
			p.Paper_Grp_Id = idPg
			temp := make([]interface{}, 0)
			temp = append(temp, p.Name,p.Paper_Grp_Id,p.Type,p.Ver,p.Create_Time,p.Author,p.Composed_Time,p.Remark,p.Status, p.Id)
			values = append(values, temp)
			count ++
		}
	}
	do.Name = append(do.Name, "id")
	dbom := gmdb.DbOM{ Table:gmdb.D_3, Name:do.Name, Value:values }
	fmt.Printf("%v\n", dbom)
	if _, err := tx.InsertMulti(dbom); err != nil {
		log.AddError(err, dbom, values)
		return ids, err
	}
	return ids, nil
}

func UpdP(tx gmdb.Transaction, fv map[string]interface{}, fvw []map[string]interface{}) (error) {
	do := gmdb.DbOpera{
		Table:gmdb.D_3,
		FV:fv,
	}
	fmt.Println(fvw)
	for _, v := range fvw {
		do.FVW = v
		fmt.Println(v)
		if _, err := tx.Update(do); err != nil {
			return err
		}
	}
	return nil
}

func DupQg(tx gmdb.Transaction, idPs []string) ([]string, error) {
	idQgs := []string{}
	for _, idP := range idPs {
		fvw := make(map[string]interface{})
		fvw["Paper_Id"] = idP
		var values [][]interface{}
		var count int
		do := gmdb.DbOpera{Table:gmdb.D_4, Name:[]string{"type","name","paper_id","desc","score","position","remark","status"}, FVW:fvw}
		if rows, err := tx.Query(do); err != nil {
			log.AddError(err, do)
			return idQgs, err
		} else {
			db := gmdb.GetDb()
			for rows.Next() {
				var qg Question_Grp
				if err = rows.Scan(&qg.Type,&qg.Name,&qg.Paper_Id,&qg.Desc,&qg.Score,&qg.Position,&qg.Remark,&qg.Status); err != nil {
					log.AddError(err)
					return idQgs, err
				}
				if qg.Id, err = UGuid(db, gmdb.D_4); err != nil {
					log.AddError(err, values)
					return idQgs, err
				}
				idQgs = append(idQgs, qg.Id)
				qg.Paper_Id = idP
				values[count] = append(values[count], qg.Type,qg.Name,qg.Paper_Id,qg.Desc,qg.Score,qg.Position,qg.Remark,qg.Status, qg.Id)
				count ++
			}
		}
		do.Name = append(do.Name, "id")
		dbom := gmdb.DbOM{ Table:gmdb.D_4, Name:do.Name, Value:values }
		if _, err := tx.InsertMulti(dbom); err != nil {
			log.AddError(err, dbom, values, idQgs)
			return idQgs, err
		}
	}
	return idQgs, nil
}

func UpdQg(tx gmdb.Transaction, fv map[string]interface{}, fvw []map[string]interface{}) (error) {
	do := gmdb.DbOpera{
		Table:gmdb.D_4,
		FV:fv,
	}
	for _, v := range fvw {
		do.FVW = v
		if _, err := tx.Update(do); err != nil {
			log.AddError(err, do, fv, fvw)
			return err
		}
	}
	return nil
}

//func DupQ(tx gmdb.Transaction, idQgs []string) ([]string, error) {
//	idQs := []string{}
//	for _, idQg := range idQgs {
//		fvw := make(map[string]interface{})
//		fvw[""]
//	}
//}