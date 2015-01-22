/* jshint nocomma: false */
/* jshint quotmark: true */
/* jshint enforceall: true */
/* global console, document, escape, setTimeout, window, XMLHttpRequest */

// H/T to Pete Lepage.
// petelepage.com/blog/2011/07/showing-hiding-panels-with-html-and-css/

function togglePanel(prefix) {
  var elem = document.getElementById(prefix + '-panel');
  if (!elem) {
    console.log('No panel to toggle for prefix:', prefix);
  } else {
    if (elem.classList) {
      elem.classList.toggle('show');
    } else {
      var classes = elem.className;
      if (classes.indexOf('show') >= 0) {
        elem.className = classes.replace('show', '');
      } else {
        elem.className = classes + ' show';
      }
      console.log(elem.className);
    }
  }
}

function showById(idStr) {
  var elt = document.getElementById(idStr);
  if (elt) {
    elt.style.display = '';
  } else {
    console.log('Not found:', idStr);
  }
}

function removeElt(idStr) {
  var elem = document.getElementById(idStr);
  if (elem) {
    elem.remove();
  }
}

function drawCal(data) {
  var calendars = JSON.parse(data);

  if (calendars.length) {
    showById('subscriptions-row');
  }

  for (var i = 0; i < calendars.length; i++) {
    document.getElementById('cal-' + i).innerHTML = calendars[i];
  }

  if (calendars.length > 3) {
    removeElt('add-panel');
  }
}

function reset(data) {
  var parsedData = JSON.parse(data);

  if (parsedData === 'whitelist:fail') {
    spawnAlert('Feed is not on whitelist.');
  } else if (parsedData === 'limit:fail') {
    spawnAlert('You have reached the maximum number of feeds.');
  } else if (parsedData === 'contained:fail') {
    spawnAlert('You are already subscribed to this calendar feed.');
  } else if (parsedData === 'no_user:fail') {
    spawnAlert('No user was provided.');
  } else {
    drawCal(data);
  }

  document.getElementById('calendar-link').value = '';
  togglePanel('add');

  return false;
}

function freqSet(data) {
  var frequency = JSON.parse(data);

  if (frequency === 'no_cal:fail') {
    spawnAlert('You have no calendar to update.');
  } else if (frequency === 'wrong_freq:fail') {
    spawnAlert('That is not a valid frequency.');
  } else if (frequency === 'no_user:fail') {
    spawnAlert('No user was provided.');
  } else if (frequency === 'method_not_supported:fail') {
    spawnAlert('That method is not supported.');
  } else {
    var frequencyVerbose = frequency[0];
    var frequencyVal = frequency[1];
    document.getElementById('freq-val').innerHTML = frequencyVerbose;
    document.getElementById('frequency').value = frequencyVal;
  }

  return false;
}

function freqReset(data) {
  freqSet(data);
  togglePanel('freq');

  return false;
}

function removeAlert() {
  togglePanel('alert');
  // TODO(dhermes): Make the transition to 0 px rather
  //                than to -145px.
  setTimeout(function() { removeElt('alert-panel'); }, 500);
}

function spawnAlert(text) {
  // first check if one exists, and remove it if it has not been
  var elem = document.getElementById('alert-panel');
  if (elem) {
    removeAlert();
  }

  var alertText = document.createElement('span');
  alertText.style.position = 'relative';
  alertText.style.top = '12px';
  alertText.textContent = text;

  var alertAnchor = document.createElement('a');
  alertAnchor.href = '#';
  alertAnchor.setAttribute('onclick', 'removeAlert();');
  alertAnchor.classList.add('controller');
  alertAnchor.textContent = 'X';

  var alertDiv = document.createElement('div');
  alertDiv.id = 'alert-panel';
  alertDiv.classList.add('panel');
  alertDiv.appendChild(alertText);
  alertDiv.appendChild(alertAnchor);

  var container = document.getElementById('alerts');
  container.appendChild(alertDiv);
  // max(text_length, 170) since 170 is the standard
  var width = Math.max(140, alertText.offsetWidth) + 30;
  alertDiv.style.width = width + 'px';

  togglePanel('alert');
}

function makeHttp(httpVerb, uriPath, payload, callback) {
    var http = new XMLHttpRequest();
    http.open(httpVerb, uriPath, true);
    http.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    http.onreadystatechange = function() {
      if (http.readyState === 4 && http.status === 200) {
        callback(http.responseText);
      }
    };

    http.send(payload);
}

window.onload = function() {
  var appData = document.getElementById('persistentCalData');
  var calendars = appData.getAttribute('data-calendars');
  var frequency = appData.getAttribute('data-frequency');
  drawCal(calendars);
  freqSet(frequency);

  document.getElementById('add').onsubmit = function() {
    var calLink = document.getElementById('calendar-link').value;
    var params = 'calendar-link=' + escape(calLink);
    makeHttp('POST', '/add', params, reset);
    return false;
  };

  document.getElementById('freq').onsubmit = function() {
    var freq = document.getElementById('frequency').value;
    var params = 'frequency=' + escape(freq);
    makeHttp('PUT', '/freq', params, freqReset);
    return false;
  };
};
