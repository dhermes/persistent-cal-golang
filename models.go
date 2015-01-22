package main

import (
	"errors"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type UserCal struct {
	Email           string
	Calendars       []string
	UpdateIntervals []int
	Upcoming        []string
	Id              *string `datastore:"-"`
	Frequency       string  `datastore:"-"`
	CalendarsJSON   string  `datastore:"-"`
}

func CurrentInterval() int {
	t := time.Now().UTC()
	return 8*int(t.Weekday()) + t.Hour()/3 + 1
}

func (userCal *UserCal) UpdateFrequency(numFreq int) {
	var baseInterval int
	if userCal.UpdateIntervals == nil || len(userCal.UpdateIntervals) == 0 {
		baseInterval = CurrentInterval()
	} else {
		baseInterval = userCal.UpdateIntervals[0]
	}

	updateIntervals := make([]int, numFreq)
	delta := 56 / numFreq
	updateIntervals[0] = baseInterval
	for i := 1; i < numFreq; i++ {
		updateIntervals[i] = (updateIntervals[i-1] + delta) % 56
	}
	userCal.UpdateIntervals = updateIntervals
}

func (userCal *UserCal) Put(c appengine.Context) error {
	if userCal.Id == nil {
		return errors.New("No ID set on UserCal")
	}

	key := datastore.NewKey(c, "UserCal", *userCal.Id, 0, nil)
	_, err := datastore.Put(c, key, userCal)
	return err
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
		userCal.Id = &u.ID

		return userCal, nil
	} else {
		return nil, err
	}
}
