package logic

import (
	"net/http"
	"util/log"
	"util/gmdb"
	"encoding/json"
	"fmt"
	"database/sql"
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
		if res, err := db.Delete(do, false); err != nil {
			log.AddError(err, do)
			OutPut(w, 204, err.Error(), nil)
		} else {
			delete(ump.EbMap, eb.Id)
			log.AddLog("Delete exambank succeed", eb)
			if affected, err := res.RowsAffected	(); err == nil {
				OutPut(w, 200, "Delete exambank succeed", affected)
			} else {
				log.AddError(err, res)
				OutPut(w, 200, "Delete exambank succeed", nil)
			}
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
			ump.PgMap[uuid] = CI{ C:true }
			log.AddLog("Create papergrp id succeed", fmt.Sprintf("%+v", pg), uuid)
			OutPut(w, 200, "Create papergrp id succeed", uuid)
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
		if !ump.PgMap[pg.Id].C || !ump.EbMap[pg.Exam_Bank_Id].I {
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
		if !ump.PgMap[pg.Id].I {
			log.AddWarning("Delete papergrp but the id or exam_bank_id isn't correct", pg)
			OutPut(w, 203, "Delete papergrp but the id or exam_bank id isn't correct", nil)
		} else {
			db := gmdb.GetDb()
			fvw := make(map[string]interface{})
			fvw["id"] = pg.Id
			do := gmdb.DbOpera{ Table:gmdb.D_2, FVW:fvw }
			if res, err := db.Delete(do, false); err != nil {
				log.AddError(err, do)
				OutPut(w, 204, err.Error(), nil)
			} else {
				delete(ump.PgMap, pg.Id)
				log.AddLog("Delete papergrp succeed", pg)
				if affected, err := res.RowsAffected(); err == nil {
					OutPut(w, 200, "Delete papergrp succeed", affected)
				} else {
					log.AddError(err, res)
					OutPut(w, 200, "Delete papergrp succeed", affected)
				}
			}
		}
	}
}
func ListPaperGrp(w http.ResponseWriter, r *http.Request) {
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
	//r.ParseForm()
	p := PaperI{}
	if rBody, err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err,string(rBody))
		OutPut(w, 201, err.Error(), nil)
		return
	}
	if !ump.PgMap[p.Paper_Grp_Id].I {
		log.AddWarning("Create paper but the paper_grp_id isn't correct", p)
		OutPut(w, 202, "Create paper but the paper_grp_id isn't correct", nil)
	} else {
		if uuid, err := Guid(); err == nil {
			ump.PMap[uuid] = CI{ C:true, I:false }
			log.AddLog("Paper create id succeed", uuid)
			OutPut(w, 200, "Paper create id succeed", uuid)
		} else {
			log.AddError(err, p)
			OutPut(w, 203, err.Error(), nil)
		}
	}
}

func CSpHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	data := r.FormValue("data")
	p := PaperI{}
	if data == "" {
		log.AddWarning("Create save paper but the paper is null")
		OutPut(w, 201, "Create save paper but the paper is null", nil)
	} else {
		if err := json.Unmarshal([]byte(data), &p); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		}
	}
	if !ump.PgMap[p.Paper_Grp_Id].I || !ump.PMap[p.Id].C {
		log.AddWarning("Create save paper but the id or paper_id isn't correct", p)
		OutPut(w, 203, "Create save paper but the id or paper_id isn't correct", nil)
	} else {
		if pm, err := JS2M(p, p); err != nil {
			log.AddError(err, p)
			OutPut(w, 204, err.Error(), nil)
		} else {
			db := gmdb.GetDb()
			do := gmdb.DbOpera{ Table:gmdb.D_3, FV:pm }
			if _, err := db.Insert(do); err != nil {
				log.AddError(err, do)
				OutPut(w, 204, err.Error(), nil)
			} else {
				ump.PMap[p.Id] = CI{ C:false, I:true }
				log.AddLog("Create save paper succeed", p)
				OutPut(w, 200 ,"Create save paper succeed", nil)
			}
		}
	}
}

func USpHandle(w http.ResponseWriter, r *http.Request) {
	p := PaperI{}
	var data string
	if data = r.FormValue("data"); data == "" {
		log.AddWarning("Updata save paper but the paper is null")
		OutPut(w, 201, "Update save paper but the paper is null", nil)
		return
	}
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		log.AddError(err, data)
		OutPut(w, 202, err.Error(), nil)
		return
	}
	if !ump.PMap[p.Id].I || !ump.PgMap[p.Paper_Grp_Id].I {
		log.AddWarning("Update save paper but the id or paper_grp_id isn't correct", p)
		OutPut(w, 202, "Update save paper but the id or paper_grp_id isn't correct", nil)
		fmt.Println(ump.PMap,  ump.PgMap)
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
	if rBody, err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err, string(rBody))
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
		if res, err := db.Delete(do, false); err != nil {
			log.AddError(err, p, do)
			OutPut(w, 203, err.Error(), nil)
		} else{
			ump.PMap[p.Id] = CI{ C:false, I:false }
			if affected, err := res.RowsAffected(); err == nil {
				log.AddLog("Delete paper succeed", p, affected)
				OutPut(w, 200, "Delete paper succeed", affected)
			} else {
				log.AddError(err, res)
				OutPut(w, 200, "Delete paper succeed", nil)
			}
		}
	}
}

func ListPaper(w http.ResponseWriter, r *http.Request) {
	p := Paper{}
	if rBody, err := UnmarshalJ(r, &p); err != nil {
		log.AddError(err, string(rBody))
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

type Te struct {
	Name		string
	Bank		string
}
func C(w http.ResponseWriter, r *http.Request) {
	var te Te
	if data, err := UnmarshalJ(r, &te); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		log.AddLog(data)
		if err = json.Unmarshal(data, &te); err == nil {
			fmt.Println(te)
		}
		fmt.Println(te, string(data))
		tt := Te{Name:"tang", Bank:"aa"}
		if ted, err := json.Marshal(tt); err == nil {
			OutPut(w, 200, "", string(ted))
		} else {
			log.AddError(err.Error())
			OutPut(w, 202, err.Error(), nil)
		}
	}
}
func CqgHandle(w http.ResponseWriter, r *http.Request) {
	qg := Question_Grp{}
	if data, err := UnmarshalJ(r, &qg); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
		return
	}
	if !ump.PMap[qg.Paper_Id].I {
		log.AddWarning("Create question_grp but the paper_id isn't correct", qg)
		OutPut(w, 202, "Create question_grp but the paper_id isn't correct", nil)
	} else {
		if uuid, err := Guid(); err == nil {
			ump.QgMap[uuid] = CI{ C:true }
			log.AddLog("Question_Grp create id succeed", uuid)
			OutPut(w, 200, "Question_Grp create id succeed", uuid)
		} else {
			log.AddError(err, qg)
			OutPut(w, 203, err.Error(), nil)
		}
	}
}
func CSqgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var data string
	qg := Question_Grp{}
	if data = r.FormValue("data"); data == "" {
		log.AddWarning("Create save question_grp but the question_grp is null")
		OutPut(w, 201, "Create save question_grp but the question_grp is null", nil)
	} else {
		if err := json.Unmarshal([]byte(data), &qg); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		}
	}
	if !ump.QgMap[qg.Id].C || !ump.PMap[qg.Paper_Id].I {
		log.AddError("Create save question_grp but the id or paper_id isn't correct", qg)
		OutPut(w, 203, "Create save question_grp but the id or paper_id isn't correct", nil)
	} else {
		db := gmdb.GetDb()
		if fv, err := JS2M(qg, qg); err != nil {
			log.AddError(err, qg)
			OutPut(w, 204, err.Error(), nil)
		} else {
			do := gmdb.DbOpera{ Table:gmdb.D_4, FV:fv }
			if _, err := db.Insert(do); err != nil {
				log.AddError(err, do, qg)
				OutPut(w, 205, err.Error(), nil)
			} else {
				ump.QgMap[qg.Id] = CI{ I:true }
				log.AddLog("Create save question_grp succeed", qg)
				OutPut(w, 200, "Create save question_grp succeed", nil)
			}
		}
	}
}
func USqgHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	qg := Question_Grp{}
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Update save question_grp but the question_grp is null")
		OutPut(w, 201, "Update save question_grp but the question_grp is null", nil)
	} else {
		if err := json.Unmarshal([]byte(data), &qg); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		}
	}
	if !ump.QgMap[qg.Id].I || !ump.PMap[qg.Paper_Id].I {
		log.AddWarning("Updata save question_grp but the id or paper_id isn't correct", qg)
		OutPut(w, 203, "Update save question_grp but the id or paper_id isn't correct", nil)
	} else {
		db := gmdb.GetDb()
		fvw := make(map[string]interface{})
		fvw["id"] = qg.Id
		if fv, err := JS2M(qg, qg); err != nil {
			log.AddError(err, qg)
			OutPut(w, 204, err.Error(), nil)
		} else {
			do := gmdb.DbOpera{ Table:gmdb.D_4, FV:fv, FVW:fvw }
			if _, err = db.Update(do); err != nil {
				log.AddError(err, do, qg)
				OutPut(w, 205, err.Error(), nil)
			} else {
				log.AddLog("Update save question_grp succeed", qg)
				OutPut(w, 200, "Update save question_grp succeed", nil)
			}
		}
	}
}
func DqgHandle(w http.ResponseWriter, r *http.Request) {
	qg := Question_Grp{}
	if data, err := UnmarshalJ(r, &qg); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.QgMap[qg.Id].I {
			log.AddWarning("Delete question_grp but the id isn't correct", qg)
			OutPut(w, 202, "Detele quesiton_grp but the id isn't correct", nil)
			return
		}
	}
	db := gmdb.GetDb()
	fvw := make(map[string]interface{})
	fvw["id"] = qg.Id
	do := gmdb.DbOpera{ Table:gmdb.D_4, FVW:fvw }
	if res, err := db.Delete(do, false); err != nil {
		log.AddError(err, do, qg)
		OutPut(w, 203, err.Error(), nil)
	} else {
		delete(ump.QgMap, qg.Id)
		items, err := res.RowsAffected()
		log.AddLog("Delete question_grp succeed", qg, items)
		if err != nil {
			log.AddError("Delete question_grp succeed but affected items error", err, qg, res)
		}
		OutPut(w, 200, "Delete question_grp succeed", items)
	}
}
func ListQuestionGrp(w http.ResponseWriter, r *http.Request) {
	qg := Question_Grp{}
	if data, err := UnmarshalJ(r, &qg); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.PMap[qg.Paper_Id].I {
			log.AddWarning("List question_grp but the paper_id isn't correct", data)
			OutPut(w, 202, "List question_grp but the paper_id isn't correct", nil)
			return
		}
	}
	db := gmdb.GetDb()
	fvw := make(map[string]interface{})
	fvw["paper_id"] = qg.Paper_Id
	do := gmdb.DbOpera{
		Table:gmdb.D_4,
		Name:[]string{"id", "type", "name", "paper_id", "desc", "score", "position", "remark", "status"},
		FVW:fvw,
	}
	qgm := []Question_Grp{}
	if rows, err := db.Query(do); err == nil {
		for rows.Next() {
			qgt := Question_Grp{}
			if err = rows.Scan(&qgt.Id, &qgt.Type, &qgt.Name, &qgt.Paper_Id, &qgt.Desc, &qgt.Score, &qgt.Position, &qgt.Remark, &qgt.Status); err != nil {
				log.AddError(err, rows, qg)
				OutPut(w, 203, err.Error(), nil)
				return
			}
			qgm = append(qgm, qgt)
		}
		log.AddLog("List question_grp succeed", qg, qgm)
		OutPut(w, 200, "List question_grp succeed", qgm)
	} else {
		log.AddError(err, do, qg)
		OutPut(w, 203, err.Error(), nil)
	}
}
func CpqHandle(w http.ResponseWriter, r *http.Request) {
	var pq Paper_Question
	if data, err := UnmarshalJ(r, &pq); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.QMap[pq.Question_Id].I || !ump.QgMap[pq.Question_Grp_Id].I {
			log.AddWarning("Create paper_question but the question_id or question_grp_id isn't correct", pq)
			OutPut(w, 202, "Create paper_question but the question_id or question_grp_id isn't correct", nil)
			return
		}
	}
	if uuid, err := Guid(); err != nil {
		log.AddError(err, pq)
		OutPut(w, 203, err.Error(), nil)
	} else {
		ump.PqMap[uuid] = CI{ C:true }
		log.AddLog("Paper_question create id succeed", uuid)
		OutPut(w, 200, "Paper_question create id succeed", uuid)
	}

}
func CSpqHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var pq Paper_Question
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Create save paper_question but the paper_question is null")
		OutPut(w, 201, "Create save paper_question but the paper_question is null", nil)
	} else {
		if err := json.Unmarshal([]byte(data), &pq); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
		} else {
			if !ump.PqMap[pq.Id].C || !ump.QMap[pq.Question_Id].I || !ump.QgMap[pq.Question_Grp_Id].I {
				log.AddWarning("Create save paper_question but the id or question_id or question_grp_id isn't correct", pq)
				OutPut(w, 203, "Create save paper_question but the id or question_id or questoin_grp_id isn't correct", nil)
				return
			}
		}
	}
	if fv, err := JS2M(pq, pq); err != nil {
		log.AddError(err, pq)
		OutPut(w, 204, err.Error(), nil)
	} else {
		db := gmdb.GetDb()
		do := gmdb.DbOpera{ Table:gmdb.D_6, FV:fv }
		if _,err = db.Insert(do); err != nil {
			log.AddError(err, do, pq)
			OutPut(w, 205, err.Error(), nil)
		} else {
			ump.PqMap[pq.Id] = CI{ I:true }
			log.AddLog("Create save paper_question succeed", pq)
			OutPut(w, 200, "Create save paper_question succeed", nil)
		}
	}

}
func USpqHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var pq Paper_Question
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Update save paper_question but the paper_question is null")
		OutPut(w, 201, "Update save paper_question but the paper_question is null", nil)
	} else {
		if err := json.Unmarshal([]byte(data), &pq); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
		} else {
			if !ump.PqMap[pq.Id].I || !ump.QMap[pq.Question_Id].I || !ump.QgMap[pq.Question_Grp_Id].I {
				log.AddWarning("Update save paper_question but the id or question_id or question_grp_id isn't correct", pq)
				OutPut(w, 203, "Update save paper_question but the id or question_id or questoin_grp_id isn't correct", nil)
				return
			}
		}
	}
	if fv, err := JS2M(pq, pq); err != nil {
		log.AddError(err, pq)
		OutPut(w, 204, err.Error(), nil)
	} else {
		db := gmdb.GetDb()
		fvw := make(map[string]interface{})
		fvw["id"] = pq.Id
		do := gmdb.DbOpera{ Table:gmdb.D_6, FV:fv, FVW:fvw }
		if _, err = db.Update(do); err != nil {
			log.AddError(err, do, pq)
			OutPut(w, 205, err.Error(), nil)
		} else {
			log.AddLog("Update save paper_question succeed", pq)
			OutPut(w, 200, "Update save paper_question succeed", nil)
		}
	}
}
func DpqHandle(w http.ResponseWriter, r *http.Request) {
	var pq Paper_Question
	if data, err := UnmarshalJ(r, &pq); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.PqMap[pq.Id].I {
			log.AddWarning("Delete paper_question but the id isn't correct", pq)
			OutPut(w, 202, "Delete paper_question but the id isn't correct", nil)
			return
		}
	}
	fvw := make(map[string]interface{})
	fvw["id"] = pq.Id
	do := gmdb.DbOpera{ Table:gmdb.D_6, FVW:fvw }
	db := gmdb.GetDb()
	if res, err := db.Delete(do, false); err != nil {
		log.AddError(err, res, pq)
		OutPut(w, 203, err.Error(), nil)
	} else {
		delete(ump.PqMap, pq.Id)
		items, err := res.RowsAffected()
		log.AddLog("Delete paper_question succeed", items, pq)
		if err != nil {
			log.AddError("Delete paper_question succeed but affected items error", res, pq, items)
		}
		OutPut(w, 200, "Delete paper_question succeed", items)
	}
}
func ListPaperQuestion(w http.ResponseWriter, r *http.Request) {
	var pq Paper_Question
	if data, err := UnmarshalJ(r, &pq); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.QgMap[pq.Question_Grp_Id].I {
			log.AddWarning("List paper_question but the question_grp_id isn't correct", pq)
			OutPut(w, 202, "List paper_question but the question_grp_id isn't correct", nil)
			return
		}
	}
	fvw := make(map[string]interface{})
	fvw["question_grp_id"] = pq.Question_Grp_Id
	do := gmdb.DbOpera{
		Table:gmdb.D_6,
		Name:[]string{"id", "question_id", "question_grp_id", "score", "position", "required", "remark", "status" },
		FVW:fvw,
	}
	db := gmdb.GetDb()
	qm := []Paper_Question{}
	if rows, err := db.Query(do); err != nil {
		log.AddError(err, do, pq)
		OutPut(w, 203, err.Error(), nil)
	} else {
		for rows.Next() {
			qt := Paper_Question{}
			if err = rows.Scan(&qt.Id,&qt.Question_Id,&qt.Question_Grp_Id,&qt.Score,&qt.Position,&qt.Required,&qt.Remark,&qt.Status); err != nil {
				log.AddError(err, rows, pq)
				OutPut(w, 204, err.Error(), nil)
			}
			qm = append(qm, qt)
		}
		log.AddLog("List paper_question succeed", pq, qm)
		OutPut(w, 200, "List paper_question succeed", qm)
	}
}
func CqHandle(w http.ResponseWriter, r *http.Request) {
	q := Question{}
	if data, err := UnmarshalJ(r, &q); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.EbMap[q.Exam_Bank_Id].I {
			log.AddWarning("Create question but the exam_bank_id isn't correct", data)
			OutPut(w, 202, "Create question but the exam_bank_id isn't correct", nil)
			return
		}
	}
	if uuid, err := Guid(); err != nil {
		log.AddError(err, uuid, q)
		OutPut(w, 203, err.Error(), nil)
	} else {
		ump.QMap[uuid] = CI{ C:true }
		log.AddLog("Question create id succeed", uuid)
		OutPut(w, 200, "Question create id succeed", uuid)
	}
}
func CSqHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := Question{}
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Create save question but the question is null")
		OutPut(w, 201, "Create save question but the question is null", nil)
	} else {
		if err := json.Unmarshal([]byte(data), &q); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		} else {
			if !ump.EbMap[q.Exam_Bank_Id].I || !ump.QMap[q.Id].C {
				log.AddWarning("Create save question but the id or exam_bank_id isn't correct", q)
				OutPut(w, 203, "Create save question but the id or exam_bank_id isn't correct", nil)
				return
			}
		}
	}
	if fv, err := JS2M(q, q); err != nil {
		log.AddError(err, q)
		OutPut(w, 204, err.Error(), nil)
	} else {
		db := gmdb.GetDb()
		do := gmdb.DbOpera{ Table:gmdb.D_5, FV:fv }
		if _, err = db.Insert(do); err != nil {
			log.AddError(err, do, q)
			OutPut(w, 205, err.Error(), nil)
		} else {
			ump.QMap[q.Id] = CI{ I:true }
			log.AddLog("Create save question succeed", q)
			OutPut(w, 200, "Create save question succeed", nil)
		}
	}
}
func USqHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := Question{}
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Update save question but the question is null")
		OutPut(w, 201, "Update save question but the question is null", nil)
	} else {
		if err := json.Unmarshal([]byte(data), &q); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		} else {
			if !ump.EbMap[q.Exam_Bank_Id].I || !ump.QMap[q.Id].I {
				log.AddWarning("Update save question but the id or exam_bank_id isn't correct", q)
				OutPut(w, 203, "Update save question but the id or exam_bank_id isn't correct", nil)
				return
			}
		}
	}
	fvw := make(map[string]interface{})
	fvw["id"] = q.Id
	if fv, err := JS2M(q, q); err != nil {
		log.AddError(err, q)
		OutPut(w, 204, err.Error(), nil)
	} else {
		db := gmdb.GetDb()
		do := gmdb.DbOpera{ Table:gmdb.D_5, FV:fv, FVW:fvw }
		if _, err = db.Update(do); err != nil {
			log.AddError(err, do, q)
			OutPut(w, 205, err.Error(), nil)
		} else {
			log.AddLog("Update save question succeed", q)
			OutPut(w, 200, "Update save question succeed", nil)
		}
	}
}
func DqHandle(w http.ResponseWriter, r *http.Request) {
	q := Question{}
	if data, err := UnmarshalJ(r, &q); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.QMap[q.Id].I {
			log.AddWarning("Delete question but the id or exam_bank_id isn't correct", q)
			OutPut(w, 202, "Delete question but the id or exam_bank_id isn't correct", nil)
			return
		}
	}
	fvw := make(map[string]interface{})
	fvw["id"] = q.Id
	db := gmdb.GetDb()
	do := gmdb.DbOpera{ Table:gmdb.D_5, FVW:fvw }
	if res, err := db.Delete(do, false); err != nil {
		log.AddError(err, res, q)
		OutPut(w, 203, err.Error(), nil)
	} else {
		delete(ump.QMap, q.Id)
		items, err := res.RowsAffected()
		log.AddLog("Delete question succeed", q, items)
		if err != nil {
			log.AddError("Delete question succeed but affected items error", err, q, res)
		}
		OutPut(w, 200, "Delete question succeed", items)
	}
}
func ListQuestion(w http.ResponseWriter, r *http.Request) {
	var q Question
	if data, err := UnmarshalJ(r, &q); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.EbMap[q.Exam_Bank_Id].I {
			log.AddWarning("List question but the exam_bank_id isn't correct", data)
			OutPut(w, 202, "List question but the exam_bank_id isn't correct", nil)
			return
		}
	}
	db := gmdb.GetDb()
	fvw := make(map[string]interface{})
	fvw["exam_bank_id"] = q.Exam_Bank_Id
	do := gmdb.DbOpera{
		Table:gmdb.D_5,
		Name:[]string{"id", "name", "type", "base_type", "spec", "ver", "exam_bank_id", "stem", "choice_1", "choice_2", "choice_3", "choice_4", "choice_5", "choice_6", "choice_7", "choice_8", "choice_answer", "analyze", "tips", "remark", "status" },
		FVW:fvw,
	}
	var qm []Question
	if rows, err := db.Query(do); err != nil {
		log.AddError(err, do, q)
		OutPut(w, 203, err.Error(), nil)
	} else {
		for rows.Next() {
			qt := Question{}
			if err = rows.Scan(&qt.Id,&qt.Name,&qt.Type,&qt.Base_Type,&qt.Spec,&qt.Ver,&qt.Exam_Bank_Id,&qt.Stem,&qt.Choice_1,&qt.Choice_2,&qt.Choice_3,&qt.Choice_4,&qt.Choice_5,&qt.Choice_6,&qt.Choice_7,&qt.Choice_8,&qt.Choice_Answer,&qt.Analyze,&qt.Tips,&qt.Remark,&qt.Status); err != nil {
				log.AddError(err, rows, q)
				OutPut(w, 204, err.Error(), nil)
			}
			qm = append(qm, qt)
		}
		log.AddLog("List question succeed", q, qm)
		OutPut(w, 200, "List question succeed", qm)
	}
}