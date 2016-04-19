package main

import (
	"net/http"
	"fmt"
	"util/log"
	"util/gmdb"
	"router"
)

func main() {
	fmt.Println("Start MS4 ... ")
	//---------- Log Init -----------------------------
	log.Init()

	//---------- Mysql Init----------------------------
	gmdb.Init()

	//---------- Router Init---------------------------
	router.Init()

	http.Handle("/", http.FileServer(http.Dir("../www")))
	fmt.Println(http.ListenAndServe(":8186", nil))
}
