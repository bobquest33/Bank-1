package logic

import (
	"net/http"
	"util/gmdb"
	"util/log"
	"encoding/json"
	"fmt"
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
			if !ump.PgMap[exam.Paper_Grp_Id].I || !ump.EMap[exam.Id].I {
				fmt.Println(ump.PgMap, ump.EMap)
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
		//ReleaseE()
		if err = ReleasePaperG(tx, exam.Paper_Grp_Id, auditType); err != nil {
			log.AddError(err, auditType, exam)
			return err
		}
		tx.Tx.Commit()
	}
	return nil
}

func ReleaseE(tx gmdb.Transaction, idE string, idPg int, ifAudit []int) (error) {
	return nil
}
func NoRelease(tx gmdb.Transaction, id int, ifAudit []int) {

}
func ReleaseQuestion(tx gmdb.Transaction, id int) {

}
func ReleaseQuestionG(tx gmdb.Transaction, ids []string) (error) {
	return nil
}
func ReleasePaper(tx gmdb.Transaction, IdPg string, ifAudit []int) (error) {
	if ids, err := DupP(tx, IdPg); err != nil {
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
		for _, v := range ids {
			fvwt := make(map[string]interface{})
			fvwt["id"] = v
			fvw = append(fvw, fvwt)
		}
		fmt.Println(fvw)
		if err = UpdP(tx, fv, fvw); err != nil {
			log.AddError(err, fv, fvw)
			return err
		}
		return ReleaseQuestionG(tx, ids)
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

func ReleaseExam(tx gmdb.Transaction, idE int) {
	
}
func ReleaseExaminationRoom(tx gmdb.Transaction, id int) {

}