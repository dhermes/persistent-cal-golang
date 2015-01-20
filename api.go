package main

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/add", addSubscription)
	http.HandleFunc("/freq", changeFrequency)
	http.HandleFunc("/getinfo", getInfo)
}

func addSubscription(writer http.ResponseWriter, unusedReq *http.Request) {
	fmt.Fprint(writer, `"no_user:fail"`)
}

func changeFrequency(writer http.ResponseWriter, unusedReq *http.Request) {
	fmt.Fprint(writer, `"no_user:fail"`)
}

func getInfo(writer http.ResponseWriter, unusedReq *http.Request) {
	fmt.Fprint(writer, `"no_user:fail"`)
}
