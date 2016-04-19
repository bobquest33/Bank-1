package logic

//
//func CreateEB(db gmdb.DbController, table string, exam Exam_Bank) (error) {
//	ebi := EBI{}
//	err := JS2S(exam, &ebi)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	do := gmdb.DbOpera{
//		Table:table,
//	}
//	do.FV, err = Struct2Map(ebi)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	fmt.Println(do)
//	_, err = db.Insert(do)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	return nil
//}
//
//func UpdateEB(db gmdb.DbController, table string, exam Exam_Bank) error {
//	ebi := EBI{}
//	err := JS2S(exam, &ebi)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	do := gmdb.DbOpera{
//		Table:table,
//	}
//	do.FV, err = Struct2Map(ebi)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	do.FVW = make(map[string]interface{})
//	do.FVW["id"] = exam.Id
//	fmt.Printf("%+v\n", do)
//	_, err = db.Update(do)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	return nil
//}
//func DeleteEB(db gmdb.DbController, table string, exam Exam_Bank) error {
//	ebi := EBI{}
//	err := JS2S(exam, &ebi)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	do := gmdb.DbOpera{
//		Table:table,
//	}
//	do.FVW = make(map[string]interface{})
//	do.FVW["id"] = exam.Id
//	_, err = db.Delete(do, false)
//	if err != nil {
//		log.AddError(err)
//		return err
//	}
//	return nil
//}
//
//func ListEB(db gmdb.DbController, table string, exam Exam_Bank) error {
//	//ebi := EBI{}
//	//names, err := Struct2Map(ebi)
//	//if err != nil {
//	//	return err
//	//	do := gmdb.DbOpera{
//	//		Table:table,
//	//		Name:names,
//	//	}
//	//}
//	return nil
//}