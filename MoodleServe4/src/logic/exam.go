package logic

import (
	"net/http"
	"util/gmdb"
	"util/log"
	"encoding/json"
	"database/sql"
)

func CeHandle(w http.ResponseWriter, r *http.Request) {
	e := Exam{}
	if data, err := UnmarshalJ(r, &e); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
	} else {
		if !ump.PgMap[e.Paper_Grp_Id].I {
			log.AddWarning("Create exam but the paper_grp_id isn't correct", e)
			OutPut(w, 202, "Create exam but the paper_grp_id isn't correct", nil)
			return
		}
	}
	db := gmdb.GetDb()
	if uuid, err := UGuid(db, gmdb.D_7); err != nil {
		log.AddError(err, e)
		OutPut(w, 203, err.Error(), nil)
	} else {
		ump.EMap[uuid] = CI{ C:true }
		log.AddLog("Exam create id succeed", uuid)
		OutPut(w, 200, "Exam create id succeed", uuid)
	}
}
func AddeHandle(w http.ResponseWriter, r *http.Request) {
	var exam Exam
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Add an exam but the exam is null")
		OutPut(w, 201, "Add an exam but the exam is null", nil)
		return
	} else {
		if err := json.Unmarshal([]byte(data), &exam); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		} else {
			if !ump.PgMap[exam.Paper_Grp_Id].I || !ump.EMap[exam.Id].C {
				log.AddWarning("Add an exam but the id or paper_grp_id isn't correct", exam)
				OutPut(w, 203, "Add an exam but the id or paper_grp_id isn't correct", nil)
				return
			}
		}
	}
	if err := realeseExam(exam); err != nil {
		log.AddError(err, exam)
		OutPut(w, 204, err.Error(), nil)
	} else {
		ump.EMap[exam.Id] = CI{ I:true }
		log.AddLog("Release exam succeed", exam)
		OutPut(w, 200, "Release exam succeed", nil)
	}
}
func realeseExam(exam Exam) (error) {
	auditType := Analyze(exam.Audit_Type)
	if tx, err := gmdb.GetTx(); err != nil {
		log.AddError(err, auditType, exam)
		return err
	} else {
		if res, err := insertExam(tx, exam); err != nil {
			log.AddError(err, res, exam)
			return err
		}
		//ReleaseE()
		if err = ReleasePaperG(tx, exam.Paper_Grp_Id, auditType); err != nil {
			log.AddError(err, auditType, exam)
			return err
		}
		tx.Tx.Commit()
	}
	return nil
}
func insertExam(tx gmdb.Transaction, exam Exam) (sql.Result, error) {
	if fv, err := JS2M(exam, exam); err != nil {
		log.AddError(err, exam)
		return nil, err
	} else {
		do := gmdb.DbOpera{
			Table:gmdb.D_7,
			FV:fv,
		}
		if res, err := tx.Insert(do); err != nil {
			log.AddError(err, do)
			return res, err
		} else {
			return res, nil
		}
	}
}
func ReleaseQuestion(tx gmdb.Transaction, idQgs []string, ifAudit []int) (error) {
	if idPqs, err := DupPq(tx, idQgs, ifAudit); err != nil {
		log.AddError(err, idQgs, ifAudit)
		return err
	} else {
		fv := make(map[string]interface{})
		var fvw []map[string]interface{}
		if IfRelease(Rel_1, ifAudit) {
			fv["status"] = "2"
			ifAudit = ifAudit[1:]
		} else {
			fv["status"] = "6"
		}
		for _, v := range idPqs {
			fvwt := make(map[string]interface{})
			fvwt["id"] = v
			fvw = append(fvw, fvwt)
		}
		if err = UpdPq(tx, fv, fvw); err != nil {
			log.AddError(err, fv, fvw)
			return err
		}
		return nil
	}
}
func ReleaseQuestionG(tx gmdb.Transaction, idPs []string, ifAudit []int) (error) {
	if idQgs, err := DupQg(tx, idPs); err != nil {
		log.AddError(err)
		return err
	} else {
		fv := make(map[string]interface{})
		var fvw []map[string]interface{}
		if IfRelease(Rel_2, ifAudit) {
			fv["status"] = "2"
			ifAudit = ifAudit[1:]
		} else {
			fv["status"] = "6"
		}
		for _, v := range idQgs {
			fvwt := make(map[string]interface{})
			fvwt["id"] = v
			fvw = append(fvw, fvwt)
		}
		if err = UpdQg(tx, fv, fvw); err != nil {
			log.AddError(err, tx, fv, fvw)
			return err
		}
		return ReleaseQuestion(tx, idQgs, ifAudit)
	}
}
func ReleasePaper(tx gmdb.Transaction, IdPg string, ifAudit []int) (error) {
	if idPs, err := DupP(tx, IdPg); err != nil {
		log.AddError(err)
		return err
	} else {
		fv := make(map[string]interface{})
		var fvw []map[string]interface{}
		if IfRelease(Rel_4, ifAudit) {
			fv["status"] = "2"
			ifAudit = ifAudit[1:]
		} else {
			fv["status"] = "6"
		}
		for _, v := range idPs {
			fvwt := make(map[string]interface{})
			fvwt["id"] = v
			fvw = append(fvw, fvwt)
		}
		if err = UpdP(tx, fv, fvw); err != nil {
			log.AddError(err, fv, fvw)
			return err
		}
		return ReleaseQuestionG(tx, idPs, ifAudit)
		//return nil
	}
}

func ReleasePaperG(tx gmdb.Transaction, idPg string, ifAudit []int) (error) {
	if id, err := DupPg(tx, idPg); err != nil {
		log.AddError(err, idPg)
		return err
	} else {
		fv := make(map[string]interface{})
		fvw := make(map[string]interface{})
		if IfRelPG(ifAudit) {
			fv["status"] = "2"
		} else {
			fv["status"] = "6"
		}
		fvw["id"] = idPg
		if err = UpdPg(tx, id, fv, fvw); err != nil {
			log.AddError(err, id, fv, fvw)
			return err
		}
	}
	return ReleasePaper(tx, idPg, ifAudit)
	//return nil
}

func CInviHandle(w http.ResponseWriter, r *http.Request) {
	invi := Invigilation{}
	if data, err := UnmarshalJ(r, &invi); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
		return
	} else {
		if !ump.EMap[invi.Exam_Id].I {
			log.AddWarning("Create invigilation but the exam_id isn't correct", invi)
			OutPut(w, 202, "Create invigilation but the exam_id isn't correct", nil)
			return
		}
	}
	db := gmdb.GetDb()
	if uuid ,err := UGuid(db, gmdb.D_8); err != nil {
		log.AddError(err, invi)
		OutPut(w, 203, err.Error(), nil)
	} else {
		ump.IMap[uuid] = CI{ C:true }
		log.AddLog("Invigilation craete id succeed", uuid)
		OutPut(w, 200, "Invigilation create id succeed", uuid)
	}
}

func CSInviHandle(w http.ResponseWriter, r *http.Request) {
	invi := Invigilation{}
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Create save invigilation but the invigilaton is null")
		OutPut(w, 201, "Create save invigilation but the invigilaton is null", nil)
		return
	} else  {
		if err := json.Unmarshal([]byte(data), &invi); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		} else {
			if !ump.EbMap[invi.Exam_Id].I || !ump.IMap[invi.Id].C {
				log.AddWarning("Create save invigilation but the id or exam_id isn't correct", invi)
				OutPut(w, 203, "Create save invigilation but the id or exam_id isn't correct", nil)
				return
			}
		}
	}
	if fv, err := JS2M(invi, invi); err != nil {
		log.AddError(err, invi)
		OutPut(w, 204, err.Error(), nil)
	} else {
		db := gmdb.GetDb()
		do := gmdb.DbOpera{ Table:gmdb.D_8, FV:fv }
		if _, err = db.Insert(do); err != nil {
			log.AddError(err, do, invi)
			OutPut(w, 205, err.Error(), nil)
		} else {
			log.AddLog("Create save succeed", invi)
			OutPut(w, 200, "Create save succeed", nil)
		}
	}
}

func USInviHandle(w http.ResponseWriter, r *http.Request) {
	invi := Invigilation{}
	if data := r.FormValue("data"); data == "" {
		log.AddWarning("Update save invigilation but the invigilaton is null")
		OutPut(w, 201, "Update save invigilation but the invigilaton is null", nil)
		return
	} else  {
		if err := json.Unmarshal([]byte(data), &invi); err != nil {
			log.AddError(err, data)
			OutPut(w, 202, err.Error(), nil)
			return
		} else {
			if !ump.EbMap[invi.Exam_Id].I || !ump.IMap[invi.Id].C {
				log.AddWarning("Update save invigilation but the id or exam_id isn't correct", invi)
				OutPut(w, 203, "Update save invigilation but the id or exam_id isn't correct", nil)
				return
			}
		}
	}
	if fv, err := JS2M(invi, invi); err != nil {
		log.AddError(err, invi)
		OutPut(w, 204, err.Error(), nil)
	} else {
		db := gmdb.GetDb()
		fvw := make(map[string]interface{})
		fvw["id"] = invi.Id
		do := gmdb.DbOpera{ Table:gmdb.D_8, FV:fv, FVW:fvw }
		if _, err = db.Update(do); err != nil {
			log.AddError(err, do, invi)
			OutPut(w, 205, err.Error(), nil)
		} else {
			log.AddLog("Update save invigilation succeed", invi)
			OutPut(w, 200, "Update save invigilation succeed", nil)
		}
	}
}

func DInviHandle(w http.ResponseWriter, r *http.Request) {
	invi := Invigilation{}
	if data, err := UnmarshalJ(r, &invi); err != nil {
		log.AddError(err, string(data))
		OutPut(w, 201, err.Error(), nil)
		return
	} else {
		if !ump.IMap[invi.Id].I {
			log.AddWarning("Delete invigilation but the id isn't correct", invi)
			OutPut(w, 202, "Delete invigilation but the id isn't correct", nil)
			return
		}
	}
	fvw := make(map[string]interface{})
	fvw["id"] = invi.Id
	db := gmdb.GetDb()
	do := gmdb.DbOpera{ Table:gmdb.D_8, FVW:fvw }
	if res, err := db.Delete(do, false); err != nil {
		log.AddError(err, do, invi)
		OutPut(w, 203, err.Error(), nil)
	} else {}
}
