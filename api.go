package main

import (
	"fmt"
	"net/http"
)

var notAllowed405 = `<html>
  <head>
    <title>405 Method Not Allowed</title>
  </head>
</html>
`

func init() {
	http.HandleFunc("/add", addSubscription)
	http.HandleFunc("/freq", changeFrequency)
	http.HandleFunc("/getinfo", getInfo)
}

func notAllowed(writer http.ResponseWriter) {
	fmt.Fprint(writer, notAllowed405)
	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func addSubscription(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		notAllowed(writer)
		return
	}
	fmt.Fprint(writer, `"no_user:fail"`)
}

func changeFrequency(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "PUT" {
		notAllowed(writer)
		return
	}
	fmt.Fprint(writer, `"no_user:fail"`)
}

func getInfo(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		notAllowed(writer)
		return
	}
	fmt.Fprint(writer, `"no_user:fail"`)
}
