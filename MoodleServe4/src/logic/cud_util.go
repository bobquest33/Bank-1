package logic

import (
	"util/gmdb"
	"util/log"
	"fmt"
)

func MDEB(db gmdb.DbController, ebm []Exam_Bank) Result {
	res := Result{
		Status:200,
	}
	length := 0
	for i, v := range ebm {
		idE, err := FindId(db.Mdb, gmdb.D_1, v.Id)
		if err != nil {
			log.AddError(err, fmt.Sprintf("%+v\n  Find id %s", v, idE))
			res.Status = 203
			res.Msg = err.Error()
			res.Data = i
			return res
		}
		if idE == "" {
			log.AddWarning("Delete exambank but the exambank is not exist", fmt.Sprintf("%+v", v))
			res.Status = 204
			res.Msg = "Delete exambank but the exambank is not exist"
			res.Data = i
			return res
		}
		do := gmdb.DbOpera{
			Table:gmdb.D_1,
		}
		do.FVW = make(map[string]interface{})
		do.FVW["id"] = v.Id
		_, err = db.Delete(do, false)
		if err != nil {
			log.AddError(err, fmt.Sprintf("%+v", do))
			res.Status = 205
			res.Msg = err.Error()
			res.Data = i
			return res
		}
		length = i
	}
	log.AddLog("Delete exambank succed",fmt.Sprintf("%+v", ebm))
	res.Msg = "Delete exambank succeed"
	res.Data = length + 1
	return  res
}

//func MultiDelEBM(table string, ebm []Exam_Bank) (items, status int, err error) {
//	db := gmdb.GetDb()
//	var i int
//	for i, eb := range ebm {
//		idE, err := FindId(db.Mdb, gmdb.D_1, eb.Name)
//		if err != nil {
//			return i, 202, err
//		}
//		if idE < 1 {
//			return i, 203, errors.New("Delete exambank but exambank is not exist\n")
//		}
//		eb.Id = idE
//		do := gmdb.DbOpera{
//			Table:table,
//		}
//		do.FVW = make(map[string]interface{})
//		do.FVW["id"] = idE
//		_, err = db.Delete(do, false)
//		if err != nil {
//			return i, 204, err
//		}
//	}
//	return i, 200, nil
//}