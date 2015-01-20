package main

import (
	"html/template"
	"net/http"
)

var templateAbout = template.Must(template.ParseFiles(
	"templates/about.html",
))

func init() {
	http.HandleFunc("/about", about)
}

func about(writer http.ResponseWriter, unusedReq *http.Request) {
	err := templateAbout.Execute(writer, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
