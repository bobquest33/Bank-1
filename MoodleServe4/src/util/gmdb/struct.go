package gmdb

const (
	D_1		=	"t_exam_bank"
	D_2		=	"t_paper_grp"
	D_3		= 	"t_paper"
	D_4		=	"t_question_grp"
	D_5		=	"t_question"
	D_6		=	"t_paper_question"
	D_7		= 	"t_exam"
	D_8		= 	"t_invigiation"
	D_9		=	"t_answer_paper"
	D_10	=	"t_answer_question"
	D_11	=	"t_audit"
	D_12	=	"t_file"
	D_13	=	"t_attr"
	D_14	=	"t_attr_value"
	D_15	=	"t_obj_owner"
	D_T		=	"test"
)

type DbConfigInfo struct {
	DbUser		string    `xml:"dbuser"`
	DbCert		string    `xml:"dbcert"`
	DbName		string    `xml:"dbname"`
	DbPort		string    `xml:"dbport"`
}

type DbOpera struct {
	Table 	string						//table name
	Name	[]string					//select fields
	FV		map[string]interface{}		//field and value
	FVW 	map[string]interface{} 		//where filed, use with select and update
	NEqual	map[string]string			//conditions field's operated type, default is '='
}