package router

import (
	"net/http"
	"logic"
)

func Init() {
	http.HandleFunc("/createExamBank",		logic.CebHandle)
	http.HandleFunc("/csaveExamBank",		logic.CSebHanlde)
	http.HandleFunc("/usaveExamBank",		logic.USebHandle)
	http.HandleFunc("/deleteExamBank",		logic.DebHandle)
	http.HandleFunc("/listExamBank",		logic.ListExamBank)

	http.HandleFunc("/createPaperGrp",		logic.CpgHandle)
	http.HandleFunc("/csavePaperGrp",		logic.CSpgHandle)
	http.HandleFunc("/usavePaperGrp",		logic.USpgHandle)
	http.HandleFunc("/DeletePaperGrp",		logic.DpgHandle)
	http.HandleFunc("/listPaperGrp",		logic.ListPaperGrp)

}