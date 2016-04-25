package router

import (
	"net/http"
	"logic"
)

func Init() {
	logic.Init()
	http.HandleFunc("/createExamBank",		logic.CebHandle)
	http.HandleFunc("/csaveExamBank",		logic.CSebHanlde)
	http.HandleFunc("/usaveExamBank",		logic.USebHandle)
	http.HandleFunc("/deleteExamBank",		logic.DebHandle)
	http.HandleFunc("/listExamBank",		logic.ListExamBank)

	http.HandleFunc("/createPaperGrp",		logic.CpgHandle)
	http.HandleFunc("/csavePaperGrp",		logic.CSpgHandle)
	http.HandleFunc("/usavePaperGrp",		logic.USpgHandle)
	http.HandleFunc("/deletePaperGrp",		logic.DpgHandle)
	http.HandleFunc("/listPaperGrp",		logic.ListPaperGrp)

	http.HandleFunc("/createPaper",			logic.CpHandle)
	http.HandleFunc("/csavePaper",			logic.CSpHandle)
	http.HandleFunc("/usavePaper",			logic.USpHandle)
	http.HandleFunc("/deletePaper",			logic.DpHandle)
	http.HandleFunc("/listPaper",			logic.ListPaper)

	http.HandleFunc("/createQuestionGrp",	logic.CqgHandle)
	http.HandleFunc("/csaveQuestionGrp",	logic.CSqgHandle)
	http.HandleFunc("/usaveQuestionGrp",	logic.USqgHandle)
	http.HandleFunc("/deleteQuestionGrp",	logic.DqgHandle)
	http.HandleFunc("/listQuestionGrp",		logic.ListQuestionGrp)

	http.HandleFunc("/createPaperQuestion",	logic.CpqHandle)
	http.HandleFunc("/csavePaperQuestion",	logic.CSpqHandle)
	http.HandleFunc("/usavePaperQuestion",	logic.USpqHandle)
	http.HandleFunc("/deletePaperQuestion",	logic.DpqHandle)
	http.HandleFunc("/listPaperQuestion",	logic.ListPaperQuestion)

	http.HandleFunc("/createQuestion",		logic.CqHandle)
	http.HandleFunc("/csaveQuestion",		logic.CSqHandle)
	http.HandleFunc("/usaveQuestion",		logic.USqHandle)
	http.HandleFunc("/deleteQuestion",		logic.DqHandle)
	http.HandleFunc("/listQuestion",		logic.ListQuestion)

	http.HandleFunc("/bank",				logic.C)
}