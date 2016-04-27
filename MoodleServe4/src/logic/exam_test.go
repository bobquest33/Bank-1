package logic_test

import (
	"testing"
	"logic"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
	"net/url"
)

func TestAnalyse(t *testing.T) {
	fmt.Println(logic.Analyze(13))
}

func NewRequest(url string, method string, v interface{}) (string, error) {
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

func TestRelease(t *testing.T) {
	exam := logic.Exam{ Id:"852327da0c6111e6846a94911dddfd5c", Paper_Grp_Id:"df4c6b88a1c71b5b506b4be734eb75f1" }
	//if _, err := NewRequest("/releaseExam", "get", exam); err != nil {
	//	t.Error(err)
	//}
	if data, err := json.Marshal(exam); err != nil {
		t.Error(err)
	} else {
		if _, err = NetConn("/releaseExam", url.Values{"data":{string(data)}}); err != nil {
			t.Error(err)
		}
	}
}