package logic

import (
	"util/gmdb"
	"util/log"
	"fmt"
)

func MDEB(db gmdb.DbController, ebm []Exam_Bank, status int) Result {
	res := Result{
		Status:200,
	}
	length := 0
	for i, v := range ebm {
		if !ump.EbMap[v.Id].I {
			log.AddWarning("Delete exambank but Id isn't correct", fmt.Sprintf("%+v", v))
			res.Status = status
			res.Msg = "Delete exambank but Id isn't correct"
			res.Data = i
			return res
		}
		do := gmdb.DbOpera{
			Table:gmdb.D_1,
		}
		do.FVW = make(map[string]interface{})
		do.FVW["id"] = v.Id
		_, err := db.Delete(do, false)
		if err != nil {
			status ++
			log.AddError(err, fmt.Sprintf("%+v", do))
			res.Status = status
			res.Msg = err.Error()
			res.Data = i
			return res
		}
		ump.EbMap[v.Id] = CI{}
		length = i
	}
	log.AddLog("Delete exambank succed",fmt.Sprintf("%+v", ebm))
	res.Msg = "Delete exambank succeed"
	res.Data = length + 1
	return  res
}

func MulDl(table string, ids []string, status int) Result {
	//res := Result{}
	//
	//for i, v := range ids {
	//	if  {
	//
	//	}
	//}
	return Result{}
}