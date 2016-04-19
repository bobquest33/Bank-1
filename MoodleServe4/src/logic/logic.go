package logic

import (
	"net/http"
	"util/log"
	"util/gmdb"
	"encoding/json"
	"fmt"
	"w.gdy.io/dyf/dms/db"
)

type Result struct {
	Status			int
	Msg				string
	Data			interface{}
}

// Uuid map
type UMP struct {
	EbMap		map[string]bool
	PgMap		map[string]bool
	PMap		map[string]bool
}

var ump UMP

func Init() {
	ump.EbMap = make(map[string]bool)
	ump.PgMap = make(map[string]bool)
	ump.PMap  = make(map[string]bool)
}

//define function name's prefix
// C 	Create
// CS	Create Save
// US	Update Save
// D	Delete

func CebHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	res := Result { Status:200 }
	data := r.FormValue("data")
	if data == "" {
		//do something
	}
	db := gmdb.GetDb()
	uuid, err := UGuid(db, gmdb.D_1)
	if err != nil {
		log.AddError("Id create error", err)
		res.Status = 201
		res.Msg = "Id create error"
		OutPut(w, res)
		return
	}
	ump.EbMap[uuid] = true
	log.AddLog("Exambank id create succeed", uuid)
	res.Msg = "Exambank id create succeed"
	res.Data = uuid
	OutPut(w, res)
	return
}

func CSebHanlde(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{
		Status:200,
	}
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Create save exambank but the exambank is null")
		res.Status = 201
		res.Msg = "Create save exambank but the exambank is null"
		OutPut(w, res)
		return
	}
	eb := Exam_Bank{}
	err := json.Unmarshal([]byte(data), &eb)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", eb))
		res.Status = 202
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	if !ump.EbMap[eb.Id] {
		log.AddWarning("Create save exambank but Id isn't correct", fmt.Sprintf("%+v", eb))
		res.Status = 203
		res.Msg = "Create save exambank but Id isn't correct"
		OutPut(w, res)
		return
	}
	db := gmdb.GetDb()
	ebi := EBI{}
	fv, err := JS2M(eb, ebi)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", eb, ebi))
		res.Status = 204
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	do := gmdb.DbOpera{
		Table:gmdb.D_1,
		FV:fv,
	}
	_, err = db.Insert(do)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", do, eb))
		res.Status = 206
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	log.AddLog("Create exambank succeed",fmt.Sprintf("%+v", eb))
	res.Msg = "Create exambank succeed"
	OutPut(w, res)
	return
}

func USebHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{
		Status:200,
	}
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Update save exambank but the exambank is null")
		res.Status = 201
		res.Msg = "Update save exambank but the exambank is null"
		OutPut(w, res)
		return
	}
	eb := Exam_Bank{}
	err := json.Unmarshal([]byte(data), &eb)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", eb))
		res.Status = 202
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	db := gmdb.GetDb()
	idE, err := FindId(db.Mdb, gmdb.D_1, eb.Id)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n", eb))
		res.Status = 203
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	if idE == "" {
		log.AddWarning("Update save exambank but the exambank is not exist", fmt.Sprintf("%+v\n  Find id %s", eb, idE))
		res.Status = 204
		res.Msg = "Update save exambank but the exambank is not exist"
		OutPut(w, res)
		return
	}
	ebi := EBI{}
	fv, err := JS2M(eb, ebi)
	if err != nil {
		log.AddError("Unexpected system error!", err.Error(), fmt.Sprintf("%+v\n%+v", eb, ebi))
		res.Status = 205
		res.Msg = "Unexpected system error"
		OutPut(w, res)
		return
	}
	do := gmdb.DbOpera{
		Table:gmdb.D_1,
		FV:fv,
	}
	do.FVW = make(map[string]interface{})
	do.FVW["id"] = eb.Id
	_, err = db.Update(do)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", do, eb))
		res.Status = 206
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	log.AddLog("Update exambank succeed",fmt.Sprintf("%+v", eb))
	res.Msg = "Update exambank succeed"
	OutPut(w, res)
	return
}

func DebHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{ Status:200 }
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Delete exambank but the exambank is null")
		res.Status = 201
		res.Msg = "Delete exambank but the exambank is null"
		OutPut(w, res)
		return
	}
	ebm := []Exam_Bank{}
	err := json.Unmarshal([]byte(data), &ebm)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", ebm))
		res.Status = 202
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	db := gmdb.GetDb()
	res = MDEB(db, ebm)
	OutPut(w, res)
	return
}

func ListExamBank(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{
		Status:200,
	}
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
		res.Status = 201
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	ebs := []Exam_Bank{}
	for rows.Next() {
		var eb Exam_Bank
		err = rows.Scan(&eb.Id, &eb.Name, &eb.Type, &eb.Class, &eb.Create_Time, &eb.Remark, &eb.Status)
		if err != nil {
			log.AddError(err, eb)
			res.Status = 202
			res.Msg = err.Error()
			OutPut(w, res)
			return
		}
		ebs = append(ebs, eb)
	}
	log.AddLog("List exambank sueeccd", ebs)
	res.Msg = "List exambank succeed"
	res.Data = ebs
	OutPut(w, res)
	return
}

//func CpgHandle(w http.ResponseWriter, r *http.Request) {
//	r.ParseForm()
//
//	res := Result{ Status:200 }
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
//		log.AddError("Create papergrp!", err)
//		res.Status = 202
//		res.Msg = "Create papergrp!" + err.Error()
//		OutPut(w, res)
//		return
//	}
//	db := gmdb.GetDb()
//	fvw := make(map[string]interface{})
//	fvw["Id"] = pg.Id
//	mpi := Mapping{ equal:true, name:[]string{"Id"} }
//}

func CpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{
		Status:200,
	}
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Create papergrp but the papergrp is null")
		res.Status = 201
		res.Msg = "Create papergrp but the papergrp is null"
		OutPut(w, res)
		return
	}
	pg := Paper_Grp{}
	err := json.Unmarshal([]byte(data), &pg)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", pg))
		res.Status = 202
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	db := gmdb.GetDb()
	idE, err := FindId(db.Mdb, gmdb.D_2, pg.Exam_Bank_Id)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", pg))
		res.Status = 203
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	if idE == "" {
		log.AddWarning("Create paper but paper's exam_bank_id is not exist")
		res.Status = 204
		res.Msg = "Create paper but paper's exam_bank_id is not exist"
		OutPut(w, res)
		return
	}
	uuid, err := UGuid(db, gmdb.D_2)
	if err != nil {
		log.AddError(err)
		res.Status = 205
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	log.AddLog("Papergrp id create succeed", uuid)
	res.Msg = "Papergrp id create succeed"
	res.Data = uuid
	OutPut(w, res)
	return
}

func CSpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{ Status:200 }
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Create save papergrp but the papergrp is null")
		res.Status = 201
		res.Msg = "Create save papergrp but the papergrp is null"
		OutPut(w, res)
		return
	}
	pg := Paper_Grp{}
	err := json.Unmarshal([]byte(data), &pg)
	if err != nil {
		log.AddError(err)
		res.Status = 202
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	db := gmdb.GetDb()
	idE, err := FindId(db.Mdb, gmdb.D_1, pg.Exam_Bank_Id)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", pg))
		res.Status = 203
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	if idE == "" {
		log.AddWarning("Create save papergrp but papergrp's exam_bank_id is not exist", fmt.Sprintf("%+v\n  Find id %s", pg, idE))
		res.Status = 204
		res.Msg = "Create save papergrp but papergrp's exam_bank_id is not exist"
		OutPut(w, res)
		return
	}
	idPg, err := FindId(db.Mdb, gmdb.D_2, pg.Id)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", pg))
		res.Status = 205
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	if idPg != "" {
		log.AddWarning("Create save papergrp but the papergrp's id is already exist", fmt.Sprintf("%+v\n Find id %s", pg, idPg))
		res.Status = 206
		res.Msg = "Create save papergrp but the papergrp's id is already exist"
		OutPut(w, res)
		return
	}
	fv, err := JS2M(pg, pg)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, fv))
		res.Status = 207
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	do := gmdb.DbOpera{ Table:gmdb.D_2, FV:fv }
	_, err = db.Insert(do)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, do))
		res.Status = 208
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	log.AddLog("Create papergrp succeed", fmt.Sprintf("%+v", do))
	res.Msg = "Create papergrp succeed"
	OutPut(w, res)
	return
}

func USpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{ Status:200 }
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Update save papergrp but the papergrp is null")
		res.Status = 201
		res.Msg = "Update save papergrp but the papergrp is null"
		OutPut(w, res)
		return
	}
	pg := Paper_Grp{}
	err := json.Unmarshal([]byte(data), &pg)
	if err != nil {
		log.AddError(err)
		res.Status = 202
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	db := gmdb.GetDb()
	idE, err := FindId(db.Mdb, gmdb.D_1, pg.Exam_Bank_Id)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", pg))
		res.Status = 203
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	if idE == "" {
		log.AddWarning("Update save papergrp but papergrp's exam_bank_id is not exist", fmt.Sprintf("%+v\n  Find id %s", pg, idE))
		res.Status = 204
		res.Msg = "Update save papergrp but papergrp's exam_bank_id is not exist"
		OutPut(w, res)
		return
	}
	idPg, err := FindId(db.Mdb, gmdb.D_2, pg.Id)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", pg))
		res.Status = 205
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	if idPg == "" {
		log.AddWarning("Update save papergrp but the papergrp's id is not exist", fmt.Sprintf("%+v\n Find id %s", pg, idPg))
		res.Status = 206
		res.Msg = "Update save papergrp but the papergrp's id is not exist"
		OutPut(w, res)
		return
	}
	fv, err := JS2M(pg, pg)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, fv))
		res.Status = 207
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	fvw := make(map[string]interface{})
	fvw["Id"] = pg.Id
	do := gmdb.DbOpera{ Table:gmdb.D_2, FV:fv, FVW:fvw }
	_, err = db.Update(do)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", pg, do))
		res.Status = 208
		res.Msg = err.Error()
		OutPut(w, res)
		return
	}
	log.AddLog("Update papergrp succeed", fmt.Sprintf("%+v", do))
	res.Msg = "Update papergrp succeed"
	OutPut(w, res)
	return
}

func DpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := Result{ Status:200 }
	data := r.FormValue("data")
	if data == "" {
		log.AddWarning("Delete papergrp but the papergrp is null")
		res.Status = 201
		res.Msg = "Delete papergrp but the papergrp is null"
		OutPut(w, res)
		return
	}
}
func ListPaperGrp(w http.ResponseWriter, r *http.Request) {

}