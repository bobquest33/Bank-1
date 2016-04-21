package logic

import (
	"net/http"
	"util/log"
	"util/gmdb"
	"encoding/json"
	"fmt"
	"database/sql"
	"io/ioutil"
	"github.com/lib/pq/oid"
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
	QgMap		map[string]CI
	QMap		map[string]CI
	PqMap		map[string]CI
}

var ump UMP

func Init() {
	ump = UMP{
		EbMap:make(map[string]CI),
		PgMap:make(map[string]CI),
		PMap :make(map[string]CI),
		QgMap:make(map[string]CI),
		QMap :make(map[string]CI),
		PqMap:make(map[string]CI),
	}
	mapMapping(&ump)
}

//define function name's prefix
// C 	Create
// CS	Create Save
// US	Update Save
// D	Delete

func CebHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
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
	do := gmdb.DbOpera{Table:gmdb.D_1, FV:fv }
	_, err = db.Insert(do)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v\n%+v", do, eb))
		OutPut(w, 205, err.Error(), nil)
		return
	}
	ump.EbMap[eb.Id] = CI{C:false, I:true }
	log.AddLog("Create exambank succeed", fmt.Sprintf("%+v", eb))
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
	err := json.Unmarshal([]byte(data), &eb)
	if err != nil {
		log.AddError(err, fmt.Sprintf("%+v", data))
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
	eb := Exam_Bank{}
	if err := json.Unmarshal([]byte(data), &eb); err != nil {
		log.AddError(err, data)
		OutPut(w, 202, err.Error(), nil)
	}
	if !ump.EbMap[eb.Id].I {
		log.AddWarning("Delete exambank but the id isn't correct", eb)
		OutPut(w, 203, "Delete exambank but the is isn't correct", nil)
	} else {
		db := gmdb.GetDb()
		fvw := make(map[string]interface{})
		fvw["id"] = eb.Id
		do := gmdb.DbOpera{
			Table:gmdb.D_1,
			FVW:fvw,
		}
		if _, err := db.Delete(do, false); err != nil {
			log.AddError(err, do)
			OutPut(w, 204, err.Error(), nil)
		} else {
			log.AddLog("Delete exambank succeed", eb)
			OutPut(w, 200, "Delete exambank succeed", nil)
			return
		}
	}
	//ebm := []Exam_Bank{}
	//err := json.Unmarshal([]byte(data), &ebm)
	//if err != nil {
	//	log.AddError(err, fmt.Sprintf("%+v", data))
	//	OutPut(w, 202, err.Error(), nil)
	//	return
	//}
	//db := gmdb.GetDb()
	//res := MDEB(db, ebm, 203)
	//OutPut(w, res.Status, res.Msg, res.Data)
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
			log.AddError(err, eb, ebs, rows)
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
			ump.PgMap[pg.Id] = CI{ C:true }
			log.AddLog("Create papergrp id succeed", fmt.Sprintf("%+v", pg), uuid)
			OutPut(w, 200, "Create papergrp id succeed", nil)
		}
	}
}

func CSpgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if data := r.FormValue("data"); data == "" {
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
		if !ump.PgMap[pg.Exam_Bank_Id].I || !ump.EbMap[pg.Id].C {
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
				_, err = db.Insert(do)
			}
			if err != nil {
				log.AddError(err, pg)
				OutPut(w, 205, err.Error(), nil)
			} else {
				ump.PgMap[pg.Id] = CI{ C:false, I:true }
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
				do := gmdb.DbOpera{ Table:gmdb.D_2, FV:pgm, FVW:fvw}
				_, err = db.Update(do)
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
		var err error
		pg := Paper_Grp{}
		if err = json.Unmarshal([]byte(data), &pg); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
		}
		if !ump.PgMap[pg.Exam_Bank_Id].I || !ump.PgMap[pg.Id].I {
			log.AddWarning("Delete papergrp but the id or exam_bank_id isn't correct", pg)
			OutPut(w, 203, "Delete papergrp but the id or exam_bank id isn't correct", nil)
		} else {
			db := gmdb.GetDb()
			fvw := make(map[string]interface{})
			fvw["id"] = pg.Id
			fvw["exam_bank_id"] = pg.Exam_Bank_Id
			do := gmdb.DbOpera{ Table:gmdb.D_2, FVW:fvw }
			if _, err := db.Delete(do, false); err != nil {
				log.AddError(err, do)
				OutPut(w, 204, err.Error(), nil)
			} else {
				ump.PgMap[pg.Id] = CI{ I:true }
				log.AddLog("Delete papergrp succeed", pg)
				OutPut(w, 200, "Delete papergrp succeed", nil)
			}
		}
	}
}
func ListPaperGrp(w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadAll(r.Body); err != nil {
	} else {
		OutPut(w, 1, string(data), nil)
	}
	r.ParseForm()
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("List papergrp but the papergrp is null")
		OutPut(w, 201, "List papergrp but the papergrp is null", nil)
	} else {
		var err error
		pg := Paper_Grp{}
		if err = json.Unmarshal([]byte(data), &pg); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
		}
		if !ump.EbMap[pg.Exam_Bank_Id].I {
			log.AddWarning("List papergrp but the exam_bank_id isn't correct", pg)
			OutPut(w, 203, "List papergrp but the exam_bank_id isn't correct", nil)
		}
		fv := make(map[string]interface{})
		fv["exam_bank_id"] = pg.Exam_Bank_Id
		pgm := []Paper_Grp{}
		db := gmdb.GetDb()
		do := gmdb.DbOpera{
			Table:gmdb.D_2,
			Name:[]string{"id", "name", "type", "exam_bank_id", "remark", "status"},
			FV:fv,
		}
		var rows *sql.Rows
		if rows, err = db.Query(do); err != nil {
			log.AddError(err, do)
			OutPut(w, 204, err.Error(), nil)
		} else {
			for rows.Next() {
				pg := Paper_Grp{}
				if err = rows.Scan(&pg.Id, &pg.Name, &pg.Type, &pg.Exam_Bank_Id, &pg.Remark, &pg.Status); err != nil {
					log.AddError(err, pg, pgm, rows)
					OutPut(w, 205, err.Error(), nil)
				} else {
					pgm = append(pgm, pg)
				}
			}
			log.AddLog("List papergrp succeed", data, pgm)
			OutPut(w, 200, "List papergrp succeed", pgm)
		}
	}
}

func CpHandle(w http.ResponseWriter, r *http.Request) {
	p := PaperI{}
	if err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err)
		OutPut(w, 201, err.Error(), nil)
		return
	}
	if !ump.PgMap[p.Paper_Grp_Id].I {
		log.AddWarning("Create paper but the paper_grp_id isn't correct", p)
		OutPut(w, 202, "Create paper but the paper_grp_id isn't correct", nil)
	} else {
		if uuid, err := Guid(); err == nil {
			log.AddLog("Paper create id succeed", uuid)
			OutPut(w, 200, "Paper create id succeed", uuid)
		} else {
			log.AddError(err, p)
			OutPut(w, 203, err.Error(), nil)
		}
	}
}

func CSpHandle(w http.ResponseWriter, r *http.Request) {
	p := PaperI{}
	if err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err)
		OutPut(w, 201, err.Error(), nil)
		return
	}
	if !ump.PgMap[p.Paper_Grp_Id].I || !ump.PMap[p.Id].C {
		log.AddWarning("Create save paper but the id or paper_id isn't correct", p)
		OutPut(w, 202, "Create save paper but the id or paper_id isn't correct", nil)
	} else {
		if pm, err := JS2M(p, p); err != nil {
			log.AddError(err, p)
			OutPut(w, 203, err.Error(), nil)
		} else {
			db := gmdb.GetDb()
			do := gmdb.DbOpera{ Table:gmdb.D_3, FV:pm }
			if _, err := db.Insert(do); err != nil {
				log.AddError(err, do)
				OutPut(w, 204, err.Error(), nil)
			} else {
				log.AddLog("Create save paper succeed", p)
				OutPut(w, 200 ,"Create save paper succeed", nil)
			}
		}
	}
}

func USpHandle(w http.ResponseWriter, r *http.Request) {
	p := PaperI{}
	if err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err)
		OutPut(w, 201, err.Error(), nil)
		return
	}
	if !ump.PMap[p.Paper_Grp_Id].I || !ump.PgMap[p.Id].I {
		log.AddWarning("Update save paper but the id or paper_grp_id isn't correct", p)
		OutPut(w, 202, "Update save paper but the id or paper_grp_id isn't correct", nil)
	} else {
		if fv, err := JS2M(p, p); err == nil {
			db := gmdb.GetDb()
			fvw := make(map[string]interface{})
			fvw["id"] = p.Id
			do := gmdb.DbOpera{ Table:gmdb.D_3, FV:fv, FVW:fvw }
			if _, err := db.Update(do); err != nil {
				log.AddError(err, do)
				OutPut(w, 203, err.Error(), nil)
			} else {
				log.AddLog("Update save paper succeed", p)
				OutPut(w, 200, "Update save paper succeed", nil)
			}
		}
	}
}

func DpHandle(w http.ResponseWriter, r *http.Request) {
	p := Paper{}
	if err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err)
		OutPut(w, 201, err.Error(), nil)
		return
	}
	if !ump.PMap[p.Id].I {
		log.AddWarning("Delete paper but the id isn't correct")
		OutPut(w, 202, "Delete paper but the id isn't correct", nil)
	} else {
		db := gmdb.GetDb()
		fvw := make(map[string]interface{})
		fvw["id"] = p.Id
		do := gmdb.DbOpera{ Table:gmdb.D_3, FVW:fvw }
		if _, err := db.Delete(do, false); err != nil {
			log.AddError(err, p, do)
			OutPut(w, 203, err.Error(), nil)
		} else{
			log.AddLog("Delete paper succeed", p)
			OutPut(w, 200, "Delete paper succeed", nil)
		}
	}
}

func ListPaper(w http.ResponseWriter, r *http.Request) {
	p := Paper{}
	if err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err)
		OutPut(w, 201, err.Error(), nil)
		return
	}
	if !ump.PgMap[p.Paper_Grp_Id].I {
		log.AddWarning("List paper but the id isn't correct", p)
		OutPut(w, 202, "List paper but the id isn't correct", nil)
	} else {
		db := gmdb.GetDb()
		pm := []Paper{}
		fvw := make(map[string]interface{})
		fvw["Paper_Grp_Id"] = p.Paper_Grp_Id
		do := gmdb.DbOpera{
			Table:gmdb.D_3,
			Name:[]string{"id", "name", "paper_grp_id", "type", "ver", "create_time", "author", "composed_time", "remark", "status"},
			FVW:fvw,
		}
		if rows, err := db.Query(do); err != nil {
			for rows.Next() {
				pt := Paper{}
				if err = rows.Scan(&pt.Id,&pt.Name,&pt.Paper_Grp_Id,&pt.Type,&pt.Ver,&pt.Create_Time,&pt.Author,&pt.Composed_Time,&pt.Remark,&pt.Status); err != nil {
					log.AddError(err, pt, pm, rows)
					OutPut(w, 203, err.Error(), nil)
				}
			}
			log.AddLog("List paper succeed", pm)
			OutPut(w, 200, "List paper succeed", nil)
		}
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