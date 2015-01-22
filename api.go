package main

import (
	"encoding/json"
	"errors"
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
var (
	getInfoResponses = map[int]string{
		1:  "once a week",
		4:  "every two days",
		7:  "once a day",
		14: "twice a day",
		28: "every six hours",
		56: "every three hours",
	}
	frequencies = map[string]int{
		"three-hrs": 56,
		"six-hrs":   28,
		"half-day":  14,
		"day":       7,
		"two-day":   4,
		"week":      1,
	}
)

func init() {
	http.HandleFunc("/add", addSubscription)
	http.HandleFunc("/freq", changeFrequency)
	http.HandleFunc("/getinfo", getInfo)
}

func notAllowed(c appengine.Context, writer http.ResponseWriter, method string) {
	c.Infof("%v method not allowed", method)
	fmt.Fprint(writer, notAllowed405)
	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func prepareMethod(writer http.ResponseWriter, request *http.Request, desiredMethod string) (appengine.Context, *UserCal, error) {
	c := appengine.NewContext(request)

	// Check for the correct verb.
	if request.Method != desiredMethod {
		notAllowed(c, writer, request.Method)
		return c, nil, errors.New("Wrong HTTP method")
	}

	// Check for a signed-in user.
	u := user.Current(c)
	if u == nil {
		fmt.Fprint(writer, `"no_user:fail"`)
		return c, nil, errors.New("No signed-in User")
	}

	// Check that the user actually has a calendar.
	userCal, err := GetUserCal(c, u)
	if userCal == nil || err != nil {
		c.Infof("no_cal:fail")
		fmt.Fprint(writer, `"no_cal:fail"`)
		if err == nil {
			err = errors.New("No UserCal found")
		}
		return c, userCal, err
	}

	return c, userCal, nil
}

func addSubscription(writer http.ResponseWriter, request *http.Request) {
	_, _, err := prepareMethod(writer, request, "POST")
	if err != nil {
		return
	}
	fmt.Fprint(writer, `"limit:fail"`)
}

func getFrequency(request *http.Request) (int, error) {
	err := request.ParseForm()
	freqVals := request.PostForm["frequency"]

	if err != nil {
		return 0, err
	}
	if len(freqVals) != 1 {
		return 0, errors.New(`"frequency" not found in request`)
	}

	numFreq := frequencies[freqVals[0]]
	if numFreq == 0 {
		return 0, errors.New(`"frequency" not an accepted interval`)
	}

	return numFreq, nil
}

func changeFrequency(writer http.ResponseWriter, request *http.Request) {
	c, userCal, err := prepareMethod(writer, request, "PUT")
	if err != nil {
		return
	}

	// Get the frequency from the PUT body.
	var numFreq int
	numFreq, err = getFrequency(request)
	if err != nil {
		c.Infof("wrong_freq:fail")
		fmt.Fprint(writer, `"wrong_freq:fail"`)
		return
	}

	// Use valid frequency to update `UpdateIntervals`.
	userCal.UpdateFrequency(numFreq)

	// Attempt to store the newly updated
	err = userCal.Put(c)
	if err == nil {
		c.Infof("Updating frequency succeeded")
		fmt.Fprint(writer, FrequencyResponses[len(userCal.UpdateIntervals)])
	} else {
		c.Infof("invalid_put:fail")
		fmt.Fprint(writer, `"invalid_put:fail"`)
	}
}

func getInfo(writer http.ResponseWriter, request *http.Request) {
	c, userCal, err := prepareMethod(writer, request, "GET")
	if err != nil {
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
