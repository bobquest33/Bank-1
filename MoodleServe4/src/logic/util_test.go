package logic

import (
	"testing"
	"fmt"
)

func TestUtil(t *testing.T) {
	eb := Exam_Bank{Name:"tangs", Old_Name:"lily", Type:"private", Class:"math", Remark:"nothing", Status:"draft"}
	ebi := EBI{}
	err := JS2S(eb, &ebi)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", ebi)

	ebm, err := Struct2Map(ebi)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", ebm)

	ebt := EBI{}
	err = Map2Struct(ebm, &ebt)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", ebt)
}
