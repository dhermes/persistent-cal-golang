package main

// https://cloud.google.com/appengine/docs/go/tools/remoteapi
import _ "appengine/remote_api"

import (
	"encoding/json"
	"html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

var (
	templateIndex = template.Must(template.ParseFiles(
		"templates/index.html",
	))
	template404 = template.Must(template.ParseFiles(
		"templates/404.html",
	))
	FrequencyResponses = map[int]string{
		1:  `["once a week", "week"]`,
		4:  `["every two days", "two-day"]`,
		7:  `["once a day", "day"]`,
		14: `["twice a day", "half-day"]`,
		28: `["every six hours", "six-hrs"]`,
		56: `["every three hours", "three-hrs"]`,
	}
)

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

func getUserCal(c appengine.Context, u *user.User) (*UserCal, error) {
	userCal, err := GetUserCal(c, u)
	if err != nil {
		c.Infof("Nothing found for user: %v", u)
		c.Infof("Got an error: %v", err)
		baseInterval := 0 // TODO: Add logic.
		userCal = &UserCal{
			Email:           u.Email,
			Id:              &u.ID,
			Calendars:       []string{},
			UpdateIntervals: []int{baseInterval},
			Upcoming:        []string{},
		}
		err = userCal.Put(c)
	} else {
		c.Infof("User was found: %v", u)
	}

	return userCal, err
}

func renderIndex(writer http.ResponseWriter, userCal *UserCal) error {
	b, err := json.Marshal(userCal.Calendars)
	if err != nil {
		return err
	}

	userCal.CalendarsJSON = string(b[:])
	userCal.Frequency = FrequencyResponses[len(userCal.UpdateIntervals)]

	return templateIndex.Execute(writer, userCal)

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
		return
	}

	err = renderIndex(writer, userCal)
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
