package main

// https://cloud.google.com/appengine/docs/go/tools/remoteapi
import _ "appengine/remote_api"

import (
	"html/template"
	"log"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

var (
	templateIndex = template.Must(template.ParseFiles(
		"templates/index.html",
	))
	template404 = template.Must(template.ParseFiles(
		"templates/404.html",
	))
)

type UserCal struct {
	Id        string
	Email     string
	Calendars string
	Frequency string
}

func init() {
	// Handles "/" or any other route not matched.
	http.HandleFunc("/", indexOr404)
}

func loginRedirect(w http.ResponseWriter, c appengine.Context, r *http.Request) {
	url, err := user.LoginURL(c, r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

func getUserCal(c appengine.Context, u *user.User) (UserCal, error) {
	key := datastore.NewKey(c, "UserCal", u.ID, 0, nil)
	userCal := UserCal{}
	err := datastore.Get(c, key, userCal)
	if err != nil {
		log.Print("Nothing found for user: ", u)
		userCal = UserCal{
			Id:        u.ID,
			Email:     u.Email,
			Calendars: "[]",
			Frequency: "[]",
		}
		_, err = datastore.Put(c, key, &userCal)
	}

	return userCal, err
}

func index(writer http.ResponseWriter, request *http.Request) {
	c := appengine.NewContext(request)
	u := user.Current(c)
	if u == nil {
		loginRedirect(writer, c, request)
		return
	}

	userCal, err := getUserCal(c, u)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	err = templateIndex.Execute(writer, userCal)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func missing(writer http.ResponseWriter, unusedReq *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	err := template404.Execute(writer, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func indexOr404(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		index(writer, request)
	} else {
		missing(writer, request)
	}
}
