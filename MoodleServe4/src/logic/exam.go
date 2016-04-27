package logic

import (
	"net/http"
	"util/gmdb"
	"util/log"
	"encoding/json"
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

}
func realeseExam(exam Exam) (error) {
	auditType := Analyze(exam.Audit_Type)
	if tx, err := gmdb.GetTx(); err != nil {
		return err
	} else {
		ReleaseE(tx, exam.Paper_Grp_Id, auditType)
	}
	return nil
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

func ReleaseE(tx gmdb.Transaction, idPg string, ifAudit []int) {
	if IfRelease(Rel_8, ifAudit) {

	} else {

	}
}
func NoRelease(tx gmdb.Transaction, id int, ifAudit []int) {

}
func ReleaseQuestion(tx gmdb.Transaction, id int) {

}
func ReleaseQuestionG(tx gmdb.Transaction, id int) {

}
func ReleasePaper(tx gmdb.Transaction, IdPg int) {

}
func ReleaseExam(tx gmdb.Transaction, idE int) {
	
}
func ReleaseExaminationRoom(tx gmdb.Transaction, id int) {

}
