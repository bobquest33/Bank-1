package logic

//todo ID json should be comment
type Exam_Bank struct {
	Id					string          `json:"Id"`				//题库编码
	Name				string       	`json:"Name"`			//题库名称
	Type				string       	`json:"Type"`			//共享类型
	Class				string       	`json:"Class"`			//学科分类
	Create_Time			string       	`json:"Create_Time"`	//创建时间
	Remark				string       	`json:"Remark"`			//备注
	Status				string       	`json:"Status"`			//状态
}

// examBank
type EBI struct {
	Id					string
	Name				string
	Type				string
	Class				string
	Remark				string
	Status				string
}

//todo. how to analyze the close package json
type EBQ struct {
	EBI
	Create_Time			string
}

type Paper_Grp struct {
	Id					string			`json:"Id"`				//卷组编码
	Name				string        	`json:"Name"`			//卷组名称
	Type				string        	`json:"Type"`			//卷组类型
	Exam_Bank_Id		string 			`json:"Exam_Bank_Id"`	//隶属题库编码
	Remark				string        	`json:"Remark"`			//备注
	Status				string        	`json:"Status"`			//状态
}

type Paper struct {
	Id					string               `json:"Id"`//试卷编码
	Name				string            `json:"Name"`//试卷名称
	Paper_Grp_Id		string               //隶属卷组编码
	Type				string            `json:"Type"`//试卷类型
	Ver					string            `json:"Ver"`//试卷版本
	Create_Time			string            `json:"Create_Time"`//创建时间
	Author				int               `json:"Author"`//作者
	Composed_Time		string            `json:"Composed_Time"`//组卷完成时间
	Remark				string            `json:"Remark"`//备注
	Status				string            `json:"Status"`//状态，见数据库模型
}

type PaperI struct {
	Id					string               `json:"Id"`//试卷编码
	Name				string            `json:"Name"`//试卷名称
	Paper_Grp_Id		string               //隶属卷组编码
	Type				string            `json:"Type"`//试卷类型
	Ver					string            `json:"Ver"`//试卷版本
	Author				int               `json:"Author"`//作者
	Composed_Time		string            `json:"Composed_Time"`//组卷完成时间
	Remark				string            `json:"Remark"`//备注
	Status				string            `json:"Status"`//状态，见数据库模型
}

type Question_Grp struct {
	Id					string               `json:"Id"`//题组编码
	Type				string            `json:"Type"`//题组题型
	Name				string            `json:"Name"`//题组名称
	Old_Name			string            `json:"Old_Name"`
	Paper_Id			string               `json:"Paper_Id"`//隶属的试卷
	Paper_Name			string            `json:"Paper_Name"`//隶属的试卷名字
	Desc				string            `json:"Desc"`//题组说明
	Score				float32           `json:"Score"`//题目分数
	Position			int               `json:"Position"`//题组在试卷中的位置
	Remark				string            `json:"Remark"`//备注
	Status				string            `json:"Status"`//状态
}

type Paper_Question struct {
	Id					string             `json:"Id"`//试卷中试题编码
	Name				string          `json:"Name"`//试卷中的试题编码
	Old_Name			string          `json:"Old_Name"`
	Question_Id			string             `json:"Question_Id"`//引用试题编码
	Question_Name		string          `json:"Question_Name"`//引用试题名字
	Question_Grp_Id		string             `json:"Question_Grp_Id"`//卷组编码
	Question_Grp_Name	string          `json:"Question_Grp_Name"`//卷组名字
	Score				float32         `json:"Score"`//分数
	Position			int             `json:"Position"`//试卷中/卷组中的位置
	Required			bool           	`json:"Required"`//必做：0，选做：1
	Remark				string          `json:"Remark"`//备注
	Status				string          `json:"Status"`//状态
}

type Question struct {
	Id					string               `json:"Id"`			//题目编码
	Name				string            `json:"Name"` 	 	//题目名称
	Old_Name			string            `json:"Old_Name"`
	Type				string            `json:"Type"`			//用户自定义类型.如情景题，阅读理解题
	Base_Type			string            `json:"Base_Type"`	//题目类型：单选，多选，判断，填空，问答
	Spec				string            `json:"Spec"`			//题目规则
	Ver					string            `json:"Ver"`			//版本
	Exam_Bank_Id		string										//隶属题库编码
	Exam_Bank_Name		string            `json:"Exam_Bank_Name"`	//隶属题库名字
	Stem				string            `json:"Stem"`			//题干
	Choice_1			string            `json:"Choice_1"`		//选项1
	Choice_2			string            `json:"Choice_2"`		//选项2
	Choice_3			string            `json:"Choice_3"`		//选项3
	Choice_4			string            `json:"Choice_4"`		//选项4
	Choice_5			string            `json:"Choice_5"`		//选项5
	Choice_6			string            `json:"Choice_6"`		//选项6
	Choice_7			string            `json:"Choice_7"`		//选项7
	Choice_8			string            `json:"Choice_8"`		//选项8
	Choice_Answer		int               `json:"Choice_Answer"`	//选择正确答案（单选或多选）
	//0：　未设置正确答案
	//1：　第1个选项正确
	//2：　第2个选项正确
	//4：　第3个选项正确
	//8：　第4个选项正确
	//16：  第5个选项正确
	//32：  第6个选项正确
	//64：  第7个选项正确
	//128：第8个选项正确
	//
	//判断题正确答案
	//0：　错误，　1：　正确
	//
	//填空题正确答案
	//整数值表示有多少个空，choice_1代表第1个空的参考答案，choice_2代表第2个空的参考答案，以此类推，每个空的多个参考答案用“｜“分隔
	Analyze				string            `json:"Analyze"`		//题目解析
	Tips				string            `json:"Tips"`			//题目提示
	Remark				string            `json:"Remark"`		//备注
	Status				string            `json:"Status"`		//状态
}

type Exam struct {
	Id					string                   `json:"Id"`//考试编码
	Exam_Name			string                `json:"Exam_Name"`//考试名称
	Old_Name			string                `json:"Old_Name"`
	Paper_Grp_Id		string                   `json:"Paper_Grp_Id"`//使用的卷组编码
	Paper_Grp_Name		string                `json:"Paper_Grp_Name"`
	Audit_Type			int                   `json:"Audit_Type"`//审核类型
	Exam_Create_Time	string                `json:"Exam_Create_Time"`//考试创建时间
	Exam_Author			int                   `json:"Exam_Author"`//考试创建者
	Exam_Target			string                `json:"Exam_Target"`//考试对象
	Time_Timit			int                   `json:"Time_Timit"`//考试限时
	Creator				int                   `json:"Creator"`//创建者
	Remark				string                `json:"Remark"`//备注
	Status				string                `json:"Status"`//状态
}

type Audit struct {
	Auditor				string            	 	  `json:"Auditor"`
	Audit_Opinion		string          	  `json:"Audit_Opinion"`
	Target_Obj_Id		string
	Target_Obj_Name		string                `json:"Target_Obj_Name"`
	Paper_Id			string
	Paper_Grp_Name		string            	  `json:"Paper_Grp_Name"`
	Question_Grp_Id		string
	Question_Grp_Name	string                `json:"Question_Grp_Name"`
	Target_Type			string        	      `json:"Target_Type"`
	Remark				string    		      `json:"Remark"`
	Status				string
}

type AuditI struct {
	Auditor				int
	Audit_Opinion		string
	Target_Obj_Id		string
	Target_Type			string
	Remark				string
	Status				string
}

type Invigilation struct {
	Id					string                   `json:"Id"`//监考编码
	Name				string                `json:"Name"` //监考编码名字e.g A1000
	Old_Name			string                `json:"Old_Name"`
	Exam_Id				string                   `json:"Exam_Id"`//考试编码
	Exam_Name			string                `json:"Exam_Name"`
	Expect_Examinee		int                   `json:"Expect_Examinee"`//应考人数
	Actural_Examinee	int                   `json:"Actural_Examinee"`//实考人数
	Admins				int                   `json:"Admins"`//监考人
	Cert				string                `json:"Cert"`//考试密码
	Rule				string                `json:"Rule"`//考试规则
	Location			string                `json:"Location"`//考试地点
	Start_Time			string                `json:"Start_Time"`//考试开始时间
	Stop_Time			string                `json:"Stop_Time"`//考试截止时间
	Remark				string                `json:"Remark"`//备注
	Status				string                `json:"Status"`//状态
}

type Answer_Paper struct {
	Id						string                `json:"Id"`//答卷编码
	Name					string             `json:"Name"`//答卷编码名字
	Old_Name				string             `json:"Old_Name"`
	Invigilation_Id			string                `json:"Invigilation_Id"`//考试编码
	Invigilation_Name		string             `json:"Invigilation_Name"`//考试名字
	Examinee				int                `json:"Examinee"`//作答人
	Answer_Start_Time		string             `json:"Answer_Start_Time"`//作答开始时间
	Answer_Commit_Time		string             `json:"Answer_Commit_Time"`//作答提交时间
	Answer_Commit_Mode		string             `json:"Answer_Commit_Mode"`//作答提交方式
	Score_Type				string             `json:"Score_Type"`//计分方式
	Score					int                `json:"Score"`//分数
	Review_Opinion			string             `json:"Review_Opinion"`//评阅意见
	Reviewer				int                `json:"Reviewer"`//评阅者
	Review_Start_Time		string             `json:"Review_Start_Time"`//评阅开始时间
	Review_Completed_Time	string             `json:"Review_Completed_Timed"`//评阅完成时间
	Score_Publish_Time		string             `json:"Score_Publish_Time"`//成绩发布时间
	Score_Publisher			int                `json:"Score_Publisher"`//成绩发布人
	Remark					string             `json:"Remark"`//备注
	Status					string             `json:"Status"`//状态
}

type Answer_Question struct {
	Id						string                `json:"Id"`//作答编码
	Name					string             `json:"Name"`//试卷中的编码
	Old_Name				string             `json:"Old_Name"`//
	Answer_Paper_Id			string                `json:"Answer_Paper_Id"`//作答试卷编码
	Answer_Paper_Name		string             `json:"Answer_Paper_Name"`
	Invigilation_Id			string                `json:"Invigilation_Id"`//考场编码
	Invigilation_Name		string             `json:"Invigilation_Name"`
	Exam_Id					string                `json:"Exam_Id"`//考试编码
	Exam_Name				string             `json:"Exam_Name"`
	Paper_Grp_Id			string                `json:"Paper_Grp_Id"`//卷组编码
	Paper_Grp_Name			string             `json:"Paper_Grp_Name"`
	Paper_Id				string                `json:"Paper_Id"`//试卷编码
	Paper_Name				string             `json:"Paper"`
	Question_Grp_Id			string                `json:"Question_Grp_Id"`//题组编码
	Question_Grp_Name		string             `json:"Question_Grp_Name"`
	Paper_Question_Id		string                `json:"Paper_Question_Id"`//作答试题编码
	Paper_Question_Name		string             `json:"Paper_Question_Name"`
	Question_Id				string                `json:"Question_Id"`//题目编码
	Question_Name			string             `json:"Question_Name"`
	Start_Time				string             `json:"Start_Time"`//开始作答时间
	Stop_Time				string             `json:"Stop_Time"`//结束作答时间
	Answer_Text				string             `json:"Answer_Text"`//作答文本值
	Answer_Number			int                `json:"Answer_Number"`//作答整数值
	Score					float32            `json:"Score"`//成绩
	Reviewer				int                `json:"Reviewer"`//阅卷人
	Opinion					string             `json:"Opinion"`//教师评阅意见
	Remark					string             `json:"Remark"`//备注
	Status					string             `json:"Status"`//状态
}

type ListInfo_ExamBank struct {
	Name				string
	Type				string
	Class				string
	Create_Time			string
	Remark				string
	Status				string
}

type ListInfo_Question struct {
	Name					string
	Type					string
	Base_Type				string
	Spec					string
	Ver						string
	Remark					string
	Status					string
}