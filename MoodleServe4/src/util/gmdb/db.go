package gmdb

import (
	"os"
	"io/ioutil"
	"encoding/xml"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"util/log"
	"fmt"
)


//Mysql Database
type DbController struct {
	Mdb		*sql.DB
}

type Result struct {
	res		sql.Result
}

type Rows struct {
	rows	*sql.Rows
}

var db DbController

func Init() {
	_, err := DbConnect()
	if err != nil {
		panic(err)
		log.AddError("Connect database but failed.", err)
		//how to solve the question if project run but db can't connect
	}
	log.AddLog("Connect database succeed.")
	fmt.Println("Database Init")
}

func DbConfig() (confStr DbConfigInfo, err error) {
	path := "/Users/tangs/IdeaProjects/MoodleServe4/src/util/gmdb/"
	f, err := os.Open(path + "configdb.xml")	//use with test gmdb
	defer f.Close()
	if err != nil {
		return confStr, err
	}
	conf, err := ioutil.ReadAll(f)
	if err != nil {
		return confStr, err
	}
	err = xml.Unmarshal([]byte(conf), &confStr)
	if err != nil {
		return confStr, err
	}
	return confStr, nil
}

func DbConnect() (DbController, error) {
	conf, err := DbConfig()
	if err != nil {
		return db, err
	}
	//root:cert@/dbname?charset=utf8
	cmd := conf.DbUser + ":" + conf.DbCert + "@/" + conf.DbName + "?charset=utf8"
	db.Mdb, err = sql.Open("mysql", cmd)
	if err != nil {
		return db, err
	}

	db.Mdb.SetMaxOpenConns(2000)			//最大连接
	db.Mdb.SetMaxIdleConns(1000)			//最大空闲
	return db, nil
}

func GetDb() (DbController) {
	err := db.Mdb.Ping()
	if err != nil {
		log.AddError("GetDb error happened!", err)
		db, err = DbConnect()
		if err != nil {
			log.AddError("GetDb error happened again!", err)
			log.AddWarning("GetDb error happend again!", err)
		}
		return db
	}
	log.AddLog("GetDb succeed!")
	return db
}