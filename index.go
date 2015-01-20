package main

// https://cloud.google.com/appengine/docs/go/tools/remoteapi
import _ "appengine/remote_api"

import (
	"html/template"
	"net/http"
)

var (
	templateIndex = template.Must(template.ParseFiles(
		"templates/index.html",
	))
	template404 = template.Must(template.ParseFiles(
		"templates/404.html",
	))
)

type IndexData struct {
	Id        string
	Calendars string
	Frequency string
}

func init() {
	http.HandleFunc("/", index)
}

func index(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		data := IndexData{Id: "Foo", Calendars: "Bar", Frequency: "Baz"}
		err := templateIndex.Execute(writer, data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := template404.Execute(writer, nil)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}
