package main

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/user"
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

func currentUser(request *http.Request) *user.User {
	c := appengine.NewContext(request)
	return user.Current(c)
}

func addSubscription(writer http.ResponseWriter, request *http.Request) {
	// Check for the correct verb.
	if request.Method != "POST" {
		notAllowed(writer)
		return
	}
	// Check for a signed-in user.
	u := currentUser(request)
	if u == nil {
		fmt.Fprint(writer, `"no_user:fail"`)
	}

	fmt.Fprint(writer, `"limit:fail"`)
}

func changeFrequency(writer http.ResponseWriter, request *http.Request) {
	// Check for the correct verb.
	if request.Method != "PUT" {
		notAllowed(writer)
		return
	}
	// Check for a signed-in user.
	u := currentUser(request)
	if u == nil {
		fmt.Fprint(writer, `"no_user:fail"`)
	}

	fmt.Fprint(writer, `"no_cal:fail"`)
}

func getInfo(writer http.ResponseWriter, request *http.Request) {
	// Check for the correct verb.
	if request.Method != "GET" {
		notAllowed(writer)
		return
	}
	// Check for a signed-in user.
	u := currentUser(request)
	if u == nil {
		fmt.Fprint(writer, `"no_user:fail"`)
	}

	fmt.Fprint(writer, `"no_cal:fail"`)
}
