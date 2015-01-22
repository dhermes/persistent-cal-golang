// H/T to Pete Lepage.
// petelepage.com/blog/2011/07/showing-hiding-panels-with-html-and-css/

function togglePanel(prefix) {
  var elem = document.getElementById(prefix + "-panel");
  if (!elem) {
    console.log('No panel to toggle for prefix:', prefix);
  } else {
    if (elem.classList) {
      elem.classList.toggle("show");
    } else {
      var classes = elem.className;
      if (classes.indexOf("show") >= 0) {
        elem.className = classes.replace("show", "");
      } else {
        elem.className = classes + " show";
      }
      console.log(elem.className);
    }
  }
}

function showById(idStr) {
  var elt = document.getElementById(idStr);
  if (elt) {
    elt.style.display = "";
  } else {
    console.log('Not found:', idStr);
  }
}

function hideById(idStr) {
  var elt = document.getElementById(idStr);
  if (elt) {
    elt.style.display = "none";
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

function draw_cal(data) {
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
  var parsed_data = JSON.parse(data);

  if (parsed_data === 'whitelist:fail') {
    spawnAlert('Feed is not on whitelist.');
  } else if (parsed_data === 'limit:fail') {
    spawnAlert('You have reached the maximum number of feeds.');
  } else if (parsed_data === 'contained:fail') {
    spawnAlert('You are already subscribed to this calendar feed.');
  } else if (parsed_data === 'no_user:fail') {
    spawnAlert('No user was provided.');
  } else {
    draw_cal(data);
  }

  document.getElementById('calendar-link').value = '';
  togglePanel('add');

  return false;
}

function freq_set(data) {
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
    var frequency_verbose = frequency[0],
        frequency_val = frequency[1];
    document.getElementById('freq-val').innerHTML = frequency_verbose;
    document.getElementById('frequency').value = frequency_val;
  }

  return false;
}

function freq_reset(data) {
  freq_set(data);
  togglePanel('freq');

  return false;
}

function removeAlert() {
  togglePanel('alert');
  // TODO(dhermes): Make the transition to 0 px rather
  //                than to -145px.
  setTimeout("removeElt('alert-panel');", 500);
}

function spawnAlert(text) {
  // first check if one exists, and remove it if it has not been
  var elem = document.getElementById('alert-panel');
  if (elem) {
    removeAlert();
  }

  var alert_text = document.createElement('span');
  alert_text.style.position = 'relative';
  alert_text.style.top = '12px';
  alert_text.textContent = text;

  var alert_anchor = document.createElement('a');
  alert_anchor.href = '#';
  alert_anchor.setAttribute('onclick', 'removeAlert();');
  alert_anchor.classList.add('controller');
  alert_anchor.textContent = 'X';

  var alert_div = document.createElement('div');
  alert_div.id = 'alert-panel';
  alert_div.classList.add('panel');
  alert_div.appendChild(alert_text);
  alert_div.appendChild(alert_anchor);

  var container = document.getElementById('alerts');
  container.appendChild(alert_div);
  // max(text_length, 170) since 170 is the standard
  var width = Math.max(140, alert_text.offsetWidth) + 30;
  alert_div.style.width = width + 'px';

  togglePanel('alert');
}

window.onload = function () {
  var appData = document.getElementById('persistentCalData');
  var calendars = appData.getAttribute('data-calendars');
  var frequency = appData.getAttribute('data-frequency');
  draw_cal(calendars);
  freq_set(frequency);

  document.getElementById('add').onsubmit = function() {
    $.post('/add',
           {'calendar-link': document.getElementById('calendar-link').value},
           reset);
    return false;
  };

  document.getElementById('freq').onsubmit = function() {
    $.ajax({
       type: 'PUT',
       url: '/freq',
       data: {'frequency': document.getElementById('frequency').value},
       success: freq_reset
    });
    return false;
  };
};
