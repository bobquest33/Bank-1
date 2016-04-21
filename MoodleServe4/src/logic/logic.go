package logic

import (
	"net/http"
	"util/log"
	"util/gmdb"
	"encoding/json"
	"fmt"
	"mydb"
)

type Result struct {
	Status			int
	Msg				string
	Data			interface{}
}

// if create ,if Insert
type CI struct {
	C		bool
	I		bool
}

// Uuid map
// mapping Ram with database for reduce times to visit database
type UMP struct {
	EbMap		map[string]CI
	PgMap		map[string]CI
	PMap		map[string]CI
}

var ump UMP

func Init() {
	ump.EbMap = make(map[string]CI)
	ump.PgMap = make(map[string]CI)
	ump.PMap  = make(map[string]CI)
}

//define function name's prefix
// C 	Create
// CS	Create Save
// US	Update Save
// D	Delete

func CebHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue("data")
	if data == "" {
		//do something
	}
	db := gmdb.GetDb()
	uuid, err := UGuid(db, gmdb.D_1)
	if err != nil {
		log.AddError("Exambank id create error", err)
		OutPut(w, 201, "Exambank id create error", nil)
		return
	}
	ump.EbMap[uuid] = CI{ C:true }
	log.AddLog("Exambank id create succeed", uuid)
	OutPut(w, 200, "Exambank id create succeed", uuid)
	return
}

func CSebHanlde(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Create save exambank but the exambank is null")
		OutPut(w, 201, "Create save exambank but the exambank is null", nil)
		return
	}
	eb := Exam_Bank{}
	err := json.Unmarshal([]byte(data), &eb)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", data))
		OutPut(w, 202, err.Error(), nil)
		return
	}
	if !ump.EbMap[eb.Id].C {
		log.AddWarning("Create save exambank but Id isn't correct", fmt.Sprintf("%+v", eb))
		OutPut(w, 203, "Create save exambank but Id isn't correct", nil)
		return
	}
	db := gmdb.GetDb()
	ebi := EBI{}
	fv, err := JS2M(eb, ebi)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", eb, ebi))
		OutPut(w, 204, err.Error(), nil)
		return
	}
	do := gmdb.DbOpera{ Table:gmdb.D_1, FV:fv }
	_, err = db.Insert(do)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", do, eb))
		OutPut(w, 205, err.Error(), nil)
		return
	}
	ump.EbMap[eb.Id] = CI{ C:false, I:true }
	log.AddLog("Create exambank succeed",fmt.Sprintf("%+v", eb))
	OutPut(w, 200, "Create exambank succeed", nil)
	return
}

func USebHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Update save exambank but the exambank is null", nil)
		OutPut(w, 201, "Update save exambank but the exambank is null", nil)
		return
	}
	eb := Exam_Bank{}
	err := json.Unmarshal([]byte(data), &data)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", eb))
		OutPut(w, 202, err.Error(), nil)
		return
	}
	if !ump.EbMap[eb.Id].I {
		log.AddWarning("Update save exambank but Id isn't correct", fmt.Sprintf("%+v", eb))
		OutPut(w, 203, "Update save exambank but Id isn't correct", nil)
		return
	}
	db := gmdb.GetDb()
	ebi := EBI{}
	fv, err := JS2M(eb, ebi)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", eb, ebi))
		OutPut(w, 204, err.Error(), nil)
		return
	}
	do := gmdb.DbOpera{ Table:gmdb.D_1, FV:fv }
	do.FVW = make(map[string]interface{})
	do.FVW["id"] = eb.Id
	_, err = db.Update(do)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", do, eb))
		OutPut(w, 205, err.Error(), nil)
		return
	}
	log.AddLog("Update exambank succeed",fmt.Sprintf("%+v", eb))
	OutPut(w, 200, "Update exambank succeed", nil)
	return
}

func DebHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Delete exambank but the exambank is null")
		OutPut(w, 201, "Delete exambank but the exambank is null", nil)
		return
	}
	ebm := []Exam_Bank{}
	err := json.Unmarshal([]byte(data), &data)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", ebm))
		OutPut(w, 202, err.Error(), nil)
		return
	}
	db := gmdb.GetDb()
	res := MDEB(db, ebm, 203)
	OutPut(w, res.Status, res.Msg, res.Data)
	return
}

func ListExamBank(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue("data")
	if data == "" {
		// todo verify uer's message.do what?
	}
	do := gmdb.DbOpera{
		Table:gmdb.D_1,
		Name:[]string{"Id", "Name", "Type", "Class", "Create_Time", "Remark", "Status"},
	}
	db := gmdb.GetDb()
	rows, err := db.Query(do)
	if err != nil {
		log.AddError(err)
		OutPut(w, 201, err.Error(), nil)
		return
	}
	ebs := []Exam_Bank{}
	for rows.Next() {
		var eb Exam_Bank
		err = rows.Scan(&eb.Id, &eb.Name, &eb.Type, &eb.Class, &eb.Create_Time, &eb.Remark, &eb.Status)
		if err != nil {
			log.AddError(err, eb)
			OutPut(w, 202, err.Error(), nil)
			return
		}
		ebs = append(ebs, eb)
	}
	log.AddLog("List exambank sueeccd", ebs)
	OutPut(w, 200, "List exambank succeed", ebs)
	return
}

func CpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Create papergrp but the papergrp is null")
		OutPut(w, 201, "Create papergrp but the papergrp is null", nil)
	} else {
		var uuid string
		var err error
		pg := Paper_Grp{}
		if err := json.Unmarshal([]byte(data), &pg); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		}
		if !ump.EbMap[pg.Exam_Bank_Id].I {
			log.AddWarning("Create papergrp but exam_bank_id isn't correct", fmt.Sprintf("%+v", pg))
			OutPut(w, 203, "Create papergrp but exam_bank_id isn't correct", nil)
			return
		}
		db := gmdb.GetDb()
		if uuid, err = UGuid(db, gmdb.D_2); err != nil {
			log.AddWarning("Create papergrp id error", err.Error(), fmt.Sprintf("%+v", pg))
			OutPut(w, 204, "Create papergrp id error." + err.Error(), nil)
		} else {
			ump.PgMap[pg.Id].I = CI{ C:true }
			log.AddLog("Create papergrp id succeed", fmt.Sprintf("%+v", pg), uuid)
			OutPut(w, 200, "Create papergrp id succeed", nil)
		}
	}
}

func CSpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if data := r.FormValue("data"); data != nil {
		log.AddWarning("Create save papergrp but the papergrp is null")
		OutPut(w, 201, "Create save papergrp but the papergrp is null", nil)
	} else {
		var err error
		pg := Paper_Grp{}
		if err = json.Unmarshal([]byte(data), &pg); err != nil {
			log.AddError(err.Error(), data)
			OutPut(w, 202, err.Error(), nil)
			return
		}
		if !ump.PgMap[pg.Exam_Bank_Id].C || !ump.EbMap[pg.Id].I {
			log.AddWarning("Create save papergrp but the id or exam_bank_id isn't correct", pg)
			OutPut(w, 203, "Create save papergrp but the id or exam_bank_id isn't correct", nil)
		} else {
			if pgm, err := JS2M(pg, pg); err != nil {
				log.AddError(err, pg)
				OutPut(w, 204, err.Error(), nil)
				return
			} else {
				db := gmdb.GetDb()
				do := gmdb.DbOpera{ Table:gmdb.D_2, FV:pgm }
				err = db.Insert(do)
			}
			if err != nil {
				log.AddError(err, pg)
				OutPut(w, 205, err.Error(), nil)
			} else {
				log.AddLog("Create save papergrp succeed", pg)
				OutPut(w, 200, "Create save papergrp succeed", nil)
			}
		}
	}
}

func USpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Update save papergrp but the papergrp is null")
		OutPut(w, 201, "Update save papergrp but the papergrp is null", nil)
	} else {
		var err error
		pg := Paper_Grp{}
		if err = json.Unmarshal([]byte(data), &pg); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		}
		if !ump.EbMap[pg.Exam_Bank_Id].I || !ump.PgMap[pg.Id].I {
			log.AddWarning("Update save papergrp but the id or exam_bank_id isn't correct", pg)
			OutPut(w, 203, "Update save papergrp but the id or exam_bank_id isn't correct", nil)
		} else {
			if pgm, err := JS2M(pg, pg); err != nil {
				log.AddError(err, pg)
				OutPut(w, 204, err.Error(), nil)
			} else {
				db := gmdb.GetDb()
				fvw := make(map[string]interface{})
				fvw["id"] = pg.Id
				do := mydb.DbOpera{ Table:gmdb.D_2, FV:pgm, FVW:fvw}
				err = db.Update(do)
			}
			if err != nil {
				log.AddError(err, pg)
				OutPut(w, 205, err.Error(), nil)
			} else {
				log.AddLog("Update save papergrp succeed", pg)
				OutPut(w, 200, "Updata save papergrp succedd", nil)
			}
		}
	}
}

func DpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Delete papergrp but the papergrp is null")
		OutPut(w, 201, "Delete papergrp but the papergrp is null", nil)
	} else {
		
	}
}
//
//func CpgHandle(w http.ResponseWriter, r *http.Request) {
//	r.ParseForm()
//
//	res := Result{
//		Status:200,
//	}
//	data := r.FormValue("data")
//	if data == "" {
//		log.AddWarning("Create papergrp but the papergrp is null")
//		res.Status = 201
//		res.Msg = "Create papergrp but the papergrp is null"
//		OutPut(w, res)
//		return
//	}
//	pg := Paper_Grp{}
//	err := json.Unmarshal([]byte(data), &pg)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v", pg))
//		res.Status = 202
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	db := gmdb.GetDb()
//	idE, err := FindId(db.Mdb, gmdb.D_2, pg.Exam_Bank_Id)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v", pg))
//		res.Status = 203
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	if idE == "" {
//		log.AddWarning("Create paper but paper's exam_bank_id is not exist")
//		res.Status = 204
//		res.Msg = "Create paper but paper's exam_bank_id is not exist"
//		OutPut(w, res)
//		return
//	}
//	uuid, err := UGuid(db, gmdb.D_2)
//	if err != nil {
//		log.AddError(err)
//		res.Status = 205
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	log.AddLog("Papergrp id create succeed", uuid)
//	res.Msg = "Papergrp id create succeed"
//	res.Data = uuid
//	OutPut(w, res)
//	return
//}
//
//func CSpgHandle(w http.ResponseWriter, r *http.Request) {
//	r.ParseForm()
//
//	res := Result{ Status:200 }
//	data := r.FormValue("data")
//	if data == "" {
//		log.AddWarning("Create save papergrp but the papergrp is null")
//		res.Status = 201
//		res.Msg = "Create save papergrp but the papergrp is null"
//		OutPut(w, res)
//		return
//	}
//	pg := Paper_Grp{}
//	err := json.Unmarshal([]byte(data), &pg)
//	if err != nil {
//		log.AddError(err)
//		res.Status = 202
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	db := gmdb.GetDb()
//	idE, err := FindId(db.Mdb, gmdb.D_1, pg.Exam_Bank_Id)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v", pg))
//		res.Status = 203
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	if idE == "" {
//		log.AddWarning("Create save papergrp but papergrp's exam_bank_id is not exist", fmt.Sprintf("%+v\n  Find id %s", pg, idE))
//		res.Status = 204
//		res.Msg = "Create save papergrp but papergrp's exam_bank_id is not exist"
//		OutPut(w, res)
//		return
//	}
//	idPg, err := FindId(db.Mdb, gmdb.D_2, pg.Id)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v", pg))
//		res.Status = 205
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	if idPg != "" {
//		log.AddWarning("Create save papergrp but the papergrp's id is already exist", fmt.Sprintf("%+v\n Find id %s", pg, idPg))
//		res.Status = 206
//		res.Msg = "Create save papergrp but the papergrp's id is already exist"
//		OutPut(w, res)
//		return
//	}
//	fv, err := JS2M(pg, pg)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, fv))
//		res.Status = 207
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	do := gmdb.DbOpera{ Table:gmdb.D_2, FV:fv }
//	_, err = db.Insert(do)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, do))
//		res.Status = 208
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	log.AddLog("Create papergrp succeed", fmt.Sprintf("%+v", do))
//	res.Msg = "Create papergrp succeed"
//	OutPut(w, res)
//	return
//}
//
//func USpgHandle(w http.ResponseWriter, r *http.Request) {
//	r.ParseForm()
//
//	res := Result{ Status:200 }
//	data := r.FormValue("data")
//	if data == "" {
//		log.AddWarning("Update save papergrp but the papergrp is null")
//		res.Status = 201
//		res.Msg = "Update save papergrp but the papergrp is null"
//		OutPut(w, res)
//		return
//	}
//	pg := Paper_Grp{}
//	err := json.Unmarshal([]byte(data), &pg)
//	if err != nil {
//		log.AddError(err)
//		res.Status = 202
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	db := gmdb.GetDb()
//	idE, err := FindId(db.Mdb, gmdb.D_1, pg.Exam_Bank_Id)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v", pg))
//		res.Status = 203
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	if idE == "" {
//		log.AddWarning("Update save papergrp but papergrp's exam_bank_id is not exist", fmt.Sprintf("%+v\n  Find id %s", pg, idE))
//		res.Status = 204
//		res.Msg = "Update save papergrp but papergrp's exam_bank_id is not exist"
//		OutPut(w, res)
//		return
//	}
//	idPg, err := FindId(db.Mdb, gmdb.D_2, pg.Id)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v", pg))
//		res.Status = 205
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	if idPg == "" {
//		log.AddWarning("Update save papergrp but the papergrp's id is not exist", fmt.Sprintf("%+v\n Find id %s", pg, idPg))
//		res.Status = 206
//		res.Msg = "Update save papergrp but the papergrp's id is not exist"
//		OutPut(w, res)
//		return
//	}
//	fv, err := JS2M(pg, pg)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, fv))
//		res.Status = 207
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	fvw := make(map[string]interface{})
//	fvw["Id"] = pg.Id
//	do := gmdb.DbOpera{ Table:gmdb.D_2, FV:fv, FVW:fvw }
//	_, err = db.Update(do)
//	if err != nil {
//		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, do))
//		res.Status = 208
//		res.Msg = err.Error()
//		OutPut(w, res)
//		return
//	}
//	log.AddLog("Update papergrp succeed", fmt.Sprintf("%+v", do))
//	res.Msg = "Update papergrp succeed"
//	OutPut(w, res)
//	return
//}
//
//func DpgHandle(w http.ResponseWriter, r *http.Request) {
//	r.ParseForm()
//
//	res := Result{ Status:200 }
//	data := r.FormValue("data")
//	if data == "" {
//		log.AddWarning("Delete papergrp but the papergrp is null")
//		res.Status = 201
//		res.Msg = "Delete papergrp but the papergrp is null"
//		OutPut(w, res)
//		return
//	}
//}
//func ListPaperGrp(w http.ResponseWriter, r *http.Request) {
//
//}