package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

var notAllowed405 = `<html>
  <head>
    <title>405 Method Not Allowed</title>
  </head>
</html>
`
var getInfoResponses = map[int]string{
	1:  "once a week",
	4:  "every two days",
	7:  "once a day",
	14: "twice a day",
	28: "every six hours",
	56: "every three hours",
}

type UserCal struct {
	Email           string
	Calendars       []string
	UpdateIntervals []int
	Upcoming        []string
	Frequency       string `datastore:"-"`
	CalendarsJSON   string `datastore:"-"`
}

func GetUserCal(c appengine.Context, u *user.User) (*UserCal, error) {
	key := datastore.NewKey(c, "UserCal", u.ID, 0, nil)
	userCal := &UserCal{}
	err := datastore.Get(c, key, userCal)

	if err == nil {
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

		return userCal, nil
	} else {
		return nil, err
	}
}

func init() {
	http.HandleFunc("/add", addSubscription)
	http.HandleFunc("/freq", changeFrequency)
	http.HandleFunc("/getinfo", getInfo)
}

func notAllowed(writer http.ResponseWriter, request *http.Request) {
	c := appengine.NewContext(request)
	c.Infof("%v method not allowed", request.Method)
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
		notAllowed(writer, request)
		return
	}
	// Check for a signed-in user.
	u := currentUser(request)
	if u == nil {
		fmt.Fprint(writer, `"no_user:fail"`)
		return
	}

	fmt.Fprint(writer, `"limit:fail"`)
}

func changeFrequency(writer http.ResponseWriter, request *http.Request) {
	// Check for the correct verb.
	if request.Method != "PUT" {
		notAllowed(writer, request)
		return
	}
	// Check for a signed-in user.
	u := currentUser(request)
	if u == nil {
		fmt.Fprint(writer, `"no_user:fail"`)
		return
	}

	fmt.Fprint(writer, `"no_cal:fail"`)
}

func getInfo(writer http.ResponseWriter, request *http.Request) {
	c := appengine.NewContext(request)

	// Check for the correct verb.
	if request.Method != "GET" {
		notAllowed(writer, request)
		return
	}

	// Check for a signed-in user.
	u := currentUser(request)
	if u == nil {
		c.Infof("no_user:fail")
		fmt.Fprint(writer, `"no_user:fail"`)
		return
	}

	userCal, err := GetUserCal(c, u)
	if userCal == nil || err != nil {
		c.Infof("no_cal:fail")
		fmt.Fprint(writer, `"no_cal:fail"`)
		return
	}

	userInfo := userCal.Calendars
	freq := getInfoResponses[len(userCal.UpdateIntervals)]
	b, err := json.Marshal([]interface{}{userInfo, freq})
	if err != nil {
		c.Infof("encoding_error:fail")
		fmt.Fprint(writer, `"encoding_error:fail"`)
	} else {
		fmt.Fprint(writer, string(b[:]))
	}
}
