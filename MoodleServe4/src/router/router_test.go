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
	reg, err := http.NewRequest("get", "http://127.0.0.1:8186/listPaperGrp", nil)
	if err != nil {
		fmt.Println("Error1 ", err)
		return
	}
	reg.Write("id=dldjlkjdakljfdlafjlfd")
	form := url.Values{
		"id":{"fdjljalfjdlajdlfjeojfjdaifdl"},
	}
	reg.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	resp, err := client.Do(reg)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error2 ", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
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
	if res.Status == 200 {
		idE, ok = res.Data.(string)
		if !ok {
			t.Error("data type miss matched")
		}
		fmt.Println("exambank id get 1: ", idE)
	}
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

	//---------------papergrp-------------------------------------------------------
	//case 1
	
}