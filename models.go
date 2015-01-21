package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
)

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
