package router_test

import (
	"testing"
	"fmt"
	"util/gmdb"
	"util/log"
	"net/http"
	"io/ioutil"
	"net/url"
	"logic"
	"encoding/json"
	"strings"
)

func newQuestion() {
	client := new(http.Client)
	reg, err := http.NewRequest("get", "http://127.0.0.1:8186/bank", nil)
	if err != nil {
		fmt.Println("Error1 ", err)
		return
	}
	te := logic.Te{ Name:"name", Bank:"bank" }
	var data []byte
	if data, err = json.Marshal(te); err != nil {
		return
	}
	form := url.Values{"data":{string(data)}}
	formd, err := json.Marshal(form)

	reg.Body = ioutil.NopCloser(strings.NewReader(string(formd)))
	resp, err := client.Do(reg)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error2 ", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	} else {
		fmt.Println(string(body))
	}
	return
}

func NewQuestion(url string, method string, v interface{}) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, "http://127.0.0.1:8186" + url, nil)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	req.Body = ioutil.NopCloser(strings.NewReader(string(data)))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))
	return string(body), nil
}

//request url
func NetConn(url string, v url.Values) (string, error) {

	resp, err := http.PostForm("http://127.0.0.1:8186" + url, v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body), "NetConn end")
	return string(body), nil
}

func TestRouter(t *testing.T) {
	log.Init()
	gmdb.Init()
	res := logic.Result{}

	// case 1
	// get uuid
	str, err := NetConn("/createExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	var idE string
	var ok bool
	idE, ok = res.Data.(string)
	if !ok {
		t.Error("data type miss matched")
	}
	fmt.Println("exambank id get 1: ", idE, res)
	//create exambank
	fmt.Println("ide is ", idE)
	eb := logic.Exam_Bank{Id:idE, Name:"examTangs", Type:"common", Class:"Math", Remark:"Nothing", Status:"draft"}
	fmt.Println(eb)
	data, err := json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/csaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error()
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	//update exambank
	eb = logic.Exam_Bank{Id:idE, Name:"examLily", Type:"private", Class:"Math", Remark:"Remark", Status:"auditing"}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/usaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	// Del exambank
	eb = logic.Exam_Bank{Id:idE}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/deleteExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	// list exambank
	str, err = NetConn("/listExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v\n", res)


	// case 2
	str, err = NetConn("/createExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	if res.Status == 200 {
		idE, ok = res.Data.(string)
		if !ok {
			t.Error("data type miss matched")
		}
		fmt.Println("exambank id get 2: ", idE)
	}
	//create exambank
	fmt.Println("ide is ", idE)
	eb = logic.Exam_Bank{Id:"", Name:"examTangs", Type:"common", Class:"Math", Remark:"Nothing", Status:"draft"}
	fmt.Println(eb)
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/csaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error()
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	//update exambank
	eb = logic.Exam_Bank{Id:idE, Name:"examLily", Type:"private", Class:"Math", Remark:"Remark", Status:"auditing"}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/usaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err, fmt.Sprintf("%+v   %+v", str, res))
	}
	fmt.Sprintf("%+v", res)
	// Del exambank
	eb = logic.Exam_Bank{Id:idE}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/deleteExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	// list exambank
	str, err = NetConn("/listExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v\n", res)

	// case 3
	str, err = NetConn("/createExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	if res.Status == 200 {
		idE, ok = res.Data.(string)
		if !ok {
			t.Error("data type miss matched")
		}
		fmt.Println("exambank id get 3: ", idE)
	}
	//create exambank
	fmt.Println("ide is ", idE)
	eb = logic.Exam_Bank{Id:idE, Name:"examTangs", Type:"common", Class:"Math", Remark:"Nothing", Status:"draft"}
	fmt.Println(eb)
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/csaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error()
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	//update exambank
	eb = logic.Exam_Bank{Id:"", Name:"examLily", Type:"private", Class:"Math", Remark:"Remark", Status:"auditing"}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/usaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	// Del exambank
	eb = logic.Exam_Bank{Id:idE}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/deleteExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	// list exambank
	str, err = NetConn("/listExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v\n", res)

	// case 4
	str, err = NetConn("/createExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	if res.Status == 200 {
		idE, ok = res.Data.(string)
		if !ok {
			t.Error("data type miss matched")
		}
		fmt.Println("exambank id get 4: ", idE)
	}
	//create exambank
	fmt.Println("ide is ", idE)
	eb = logic.Exam_Bank{Id:idE, Name:"examTangs", Type:"common", Class:"Math", Remark:"Nothing", Status:"draft"}
	fmt.Println(eb)
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/csaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error()
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	//update exambank
	eb = logic.Exam_Bank{Id:idE, Name:"examLily", Type:"private", Class:"Math", Remark:"Remark", Status:"auditing"}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/usaveExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	// Del exambank
	eb = logic.Exam_Bank{Id:""}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/deleteExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)

	// Del exambank
	eb = logic.Exam_Bank{Id:idE}
	data, err = json.Marshal(eb)
	if err != nil {
		t.Error(err)
	}
	str, err = NetConn("/deleteExamBank", url.Values{"data":{string(data)}})
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v", res)
	// list exambank
	str, err = NetConn("/listExamBank", nil)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintf("%+v\n", res)
	////---------------papergrp-------------------------------------------------------
	////case 1
	//create papergrp
	pg := logic.Paper_Grp{}
	var idPg string
	pg.Exam_Bank_Id = "9f3b7417c70e044659f72faaf04604bd"
	if data, err = json.Marshal(pg); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/createPaperGrp", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	}
	if err := json.Unmarshal([]byte(str), &res); err != nil {
		t.Error(err)
	}
	idPg, ok = res.Data.(string)
	if !ok {
		fmt.Sprintf("%+v", res)
		t.Error("data format error")
	} else {
		fmt.Println(idPg)
	}
	//create save papergrp
	pg = logic.Paper_Grp{Id:idPg, Name:"tangs", Type:"tl", Exam_Bank_Id:"13943fb24ca335cf96dc5624f1ccc64c", Remark:"tRemark", Status:"draft"}
	if data, err = json.Marshal(pg); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/csavePaperGrp", url.Values{"data":{string(data)}}); err != nil {
		t.Error()
	} else {
		fmt.Sprintf("%+v", str)
	}
	//update save papergrp
	pg = logic.Paper_Grp{ Id:idPg, Name:"lily", Type:"common", Exam_Bank_Id: "13943fb24ca335cf96dc5624f1ccc64c", Remark:"nothing", Status:"draft" }
	if data, err = json.Marshal(pg); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/usavePaperGrp", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	} else {
		fmt.Sprintf("%v", str)
	}
	//delete papergrp
	pg = logic.Paper_Grp{ Id:idPg }
	if data, err = json.Marshal(pg); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/deletePaperGrp", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	} else {
		fmt.Sprintf("%v", str)
	}
	//list papergrp
	pg = logic.Paper_Grp{Exam_Bank_Id:"13943fb24ca335cf96dc5624f1ccc64c" }
	if data, err = json.Marshal(pg); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/listPaperGrp", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	} else {
		fmt.Sprintf("%v", str)
	}

	//---------------paper-------------------------------------------------------
	//case 1
	//create paper
	p := logic.PaperI{ Paper_Grp_Id:"df4c6b88a1c71b5b506b4be734eb75f1" }
	if str, err = NewQuestion("/createPaper", "get", p); err != nil {
		t.Error(err)
	}
	fmt.Println(str, p)
	if err := json.Unmarshal([]byte(str), &res); err != nil {
		t.Error(err)
	}
	idP, ok := res.Data.(string)
	if !ok {
		t.Error("data format error")
	} else {
		fmt.Println(idP)
	}
	//create save paper
	p = logic.PaperI{ Id:idP, Name:"tngs", Paper_Grp_Id:"df4c6b88a1c71b5b506b4be734eb75f1", Type:"math", Ver:"cook", Author:2, Remark:"none", Status:"draft" }
	if data, err = json.Marshal(p); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/csavePaper", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	} else {
		fmt.Sprintf("%v", str)
	}
	//update save paper
	p = logic.PaperI{ Id:idP, Name:"lily", Paper_Grp_Id:"df4c6b88a1c71b5b506b4be734eb75f1", Type:"cook", Composed_Time:"",Remark:"ok", Status:"draft" }
	if data, err = json.Marshal(p); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/usavePaper", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	} else {
		fmt.Sprintf("%v", str)
	}
	//delete paper
	p = logic.PaperI{ Id:idP }
	if str, err = NewQuestion("/deletePaper", "get", p); err != nil {
		t.Error(err)
	}
	fmt.Println(str)
	//------------------ question_grp ---------------------------------------
	//case 1
	// create question_grp
	qg := logic.Question_Grp{Paper_Id:"bc441fb5a60d7ac3976bab50aa2d152a"}
	var idQg string
	if str, err = NewQuestion("/createQuestionGrp", "get", qg); err != nil {
		t.Error(err)
	} else {
		if err = json.Unmarshal([]byte(str), &res); err != nil {
			t.Error(err)
		} else {
			idQg = res.Data.(string)
			fmt.Println("idQg is ", idQg)
		}
	}
	// create save question_grp
	qg = logic.Question_Grp{ Id:idQg, Type:"math", Name:"mathQG", Paper_Id:"bc441fb5a60d7ac3976bab50aa2d152a", Desc:"none", Score:98.9, Position:1 }
	if data, err = json.Marshal(qg); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/csaveQuestionGrp", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	}
	// update save quesion_grp
	qg = logic.Question_Grp{ Id:idQg, Type:"math", Name:"mathQgrp",Paper_Id:"bc441fb5a60d7ac3976bab50aa2d152a", Desc:"none", Score:99.9, Position:2, Remark:"none", Status:"draft" }
	if data, err = json.Marshal(qg); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/usaveQuestionGrp", url.Values{"data":{string(data)}}); err != nil {
		t.Error()
	}
	// delete save question_grp
	qg = logic.Question_Grp{ Id:idQg }
	if str, err = NewQuestion("/deleteQuestionGrp", "get", qg); err != nil {
		t.Error(err)
	}
	// list question_grp
	qg = logic.Question_Grp{ Paper_Id:"bc441fb5a60d7ac3976bab50aa2d152a" }
	if str, err = NewQuestion("/listQuestionGrp", "get", qg); err != nil {
		t.Error(err)
	}
	//------------------ question ---------------------------------------
	// create question
	var idQ string
	q := logic.Question{ Exam_Bank_Id:"4e8e085e991b554a767ae69fdf7a46eb" }
	if str, err = NewQuestion("/createQuestion", "get", q); err != nil {
		t.Error(err)
	} else {
		if err = json.Unmarshal([]byte(str), &res); err != nil {
			t.Error(err)
		}
		idQ = res.Data.(string)
		fmt.Println("idQ is ", idQ)
	}
	// create save question
	q = logic.Question{ Id:idQ,Name:"q1",Type:"math",Base_Type:"",Spec:"none",Ver:"1", Exam_Bank_Id:"4e8e085e991b554a767ae69fdf7a46eb",Stem:"this stem", Choice_Answer:0,Analyze:"none", Remark:"none", Status:"null" }
	if data, err = json.Marshal(q); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/csaveQuestion", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	}
	// update save question
	q = logic.Question{ Id:idQ,Name:"q1",Type:"Chinese",Base_Type:"",Spec:"none",Ver:"2", Exam_Bank_Id:"4e8e085e991b554a767ae69fdf7a46eb",Stem:"this stem", Choice_1:"ok", Choice_Answer:1,Analyze:"have none", Remark:"none", Status:"null" }
	if data, err = json.Marshal(q); err != nil {
		t.Error(err)
	}
	if str, err = NetConn("/usaveQuestion", url.Values{"data":{string(data)}}); err != nil {
		t.Error(err)
	}
	// delete question
	q = logic.Question{ Id:idQ }
	if str, err = NewQuestion("/deleteQuestion", "get", q); err != nil {
		t.Error(err)
	}
	// list question
	q = logic.Question{Exam_Bank_Id:"4e8e085e991b554a767ae69fdf7a46eb"}
	if str, err = NewQuestion("/listQuestion", "get", q); err != nil {
		t.Error(err)
	}
	//------------------ paper_question ---------------------------------------
	// case 1
	// create paper_question
	var idPq string
	pq := logic.Paper_Question{Question_Grp_Id:"352d98c947128f842277e65fc0940891", Question_Id:"5ee2721ea8c3f381eeec203555a8799b"}
	if str, err = NewQuestion("/createPaperQuestion", "get", pq); err != nil {
		t.Error(err)
	} else {
		if err = json.Unmarshal([]byte(str), &res); err != nil {
			t.Error(err)
		}
		idPq = res.Data.(string)
		fmt.Println("idPq is ", idPq)
	}
	// create save paper_question
	pq = logic.Paper_Question{ Id:idPq, Question_Id:"5ee2721ea8c3f381eeec203555a8799b", Question_Grp_Id:"352d98c947128f842277e65fc0940891", Score:99,Position:1, Required:true, Remark:"none", Status:"draft" }
	if data, err = json.Marshal(pq); err != nil {
		t.Error(err)
	} else {
		if str, err = NetConn("/csavePaperQuestion", url.Values{"data":{string(data)}}); err != nil {
			t.Error(err)
		}
	}
	// update save paper_question
	pq = logic.Paper_Question{ Id:idPq, Question_Id:"5ee2721ea8c3f381eeec203555a8799b", Question_Grp_Id:"352d98c947128f842277e65fc0940891", Score:96,Position:2, Required:true, Remark:"null", Status:"draft" }
	if data, err = json.Marshal(pq); err != nil {
		t.Error(err)
	} else {
		if str, err = NetConn("/usavePaperQuestion", url.Values{"data":{string(data)}}); err != nil {
			t.Error(err)
		}
	}
	// delete paper_question
	pq = logic.Paper_Question{Id:idPq}
	if str, err = NewQuestion("/deletePaperQuestion", "get", pq); err != nil {
		t.Error(err)
	}
	// list paper_question
	pq = logic.Paper_Question{Question_Grp_Id:"352d98c947128f842277e65fc0940891"}
	if str, err = NewQuestion("/listPaperQuestion", "get", pq); err != nil {
		t.Error(err)
	}
	//------------------ test ---------------------------------------

	//te := logic.Te{Name:"tangs", Bank:"bank"}
	//if data, err = json.Marshal(te); err != nil {
	//	t.Error(err)
	//} else {
	//	str, err = NewQuestion("/bank", "get", te)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	if err = json.Unmarshal([]byte(str), &res); err != nil {
	//		t.Error(err)
	//	} else {
	//		fmt.Println(res)
	//	}
	//	//ttt := logic.Te{}
	//	//if err = json.Unmarshal([]byte(res.Data.(string)), &ttt); err != nil {
	//	//	t.Error(err)
	//	//} else {
	//	//	fmt.Println(ttt, ttt.Bank)
	//	//}
	//}
}