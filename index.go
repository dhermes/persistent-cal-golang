package main

// https://cloud.google.com/appengine/docs/go/tools/remoteapi
import _ "appengine/remote_api"

import (
	"encoding/json"
	"html/template"
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
	frequencyResponses = map[int]string{
		1:  `["once a week", "week"]`,
		4:  `["every two days", "two-day"]`,
		7:  `["once a day", "day"]`,
		14: `["twice a day", "half-day"]`,
		28: `["every six hours", "six-hrs"]`,
		56: `["every three hours", "three-hrs"]`,
	}
)

type UserCal struct {
	Email           string
	Calendars       []string
	UpdateIntervals []int
	Upcoming        []string
	Frequency       string `datastore:"-"`
	CalendarsJSON   string `datastore:"-"`
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

func getUserCal(c appengine.Context, u *user.User) (*UserCal, error) {
	key := datastore.NewKey(c, "UserCal", u.ID, 0, nil)
	userCal := &UserCal{}
	err := datastore.Get(c, key, userCal)
	if err != nil {
		c.Infof("Nothing found for user: %v", u)
		c.Infof("Got an error: %v", err)
		baseInterval := 0 // TODO: Add logic.
		userCal = &UserCal{
			Email:           u.Email,
			Calendars:       []string{},
			UpdateIntervals: []int{baseInterval},
			Upcoming:        []string{},
			Frequency:       `["once a week", "week"]`,
			CalendarsJSON:   "[]",
		}
		_, err = datastore.Put(c, key, &userCal)
	} else {
		c.Infof("User was found: %v", u)
		// TODO: Implement PropertyLoadSaver interface.
		if userCal.Calendars == nil {
			userCal.Calendars = []string{}
		}
		if userCal.UpdateIntervals == nil {
			userCal.UpdateIntervals = []int{}
		}
		if userCal.Upcoming == nil {
			userCal.Upcoming = []string{}
		}

		var b []byte
		b, err = json.Marshal(userCal.Calendars)
		userCal.CalendarsJSON = string(b[:])
		userCal.Frequency = frequencyResponses[len(userCal.UpdateIntervals)]
		c.Infof("ent.Frequency: %#v, ent.CalendarsJSON: %#v",
			userCal.Frequency, userCal.CalendarsJSON)
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
