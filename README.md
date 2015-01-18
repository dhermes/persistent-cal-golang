### Description

This repository contains server and client code for
the `persistent-cal` [App Engine][1] application.

The application enables persistent import of an [iCalendar][2] feed into a
Google Calendar.

**For example**, events from a [TripIt][3] feed can be
periodically added to a user's default Google Calendar (perhaps
one shared with coworkers or collaborators).

See the [about page][4] of the deployed version of this application
for more information.

### Pre-requisites

Make sure to turn on the Calendar API service in the APIs console.

### Dependencies

In order to parse ICAL feeds, we use:

```
$ go get github.com/laurent22/ical-go
```

For interacting with the Google Calendar API:

```
$ go get golang.org/x/oauth2
$ go get google.golang.org/api/calendar/v3
```

[1]: https://cloud.google.com/products/app-engine/
[2]: http://en.wikipedia.org/wiki/ICalendar
[3]: https://www.tripit.com/
[4]: http://persistent-cal.appspot.com/about
