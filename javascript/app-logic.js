// H/T to Pete Lepage.
// petelepage.com/blog/2011/07/showing-hiding-panels-with-html-and-css/

function togglePanel(prefix) {
  var elem = document.getElementById(prefix + "-panel");
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

function draw_cal(data) {
  var calendars = JSON.parse(data);

  if (calendars.length) {
    $('#subscriptions-row').show()
  }

  for (var i = 0; i < calendars.length; i++) {
    $('#cal-' + i).html(calendars[i]);
  }

  if (calendars.length > 3) {
    $('#add-panel').remove()
  } else {
    $('#cal-button').show();
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

  $('#calendar-link').val('');
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
    $('#freq-val').html(frequency_verbose);
    $('#frequency').val(frequency_val);
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

  function removeElt() {
    $('#alert-panel').remove();
  }
  // TODO(dhermes), make the transition to 0 px rather
  // than to -145px
  // alert_div.css('left', (-width) + 'px');
  setTimeout("$('#alert-panel').remove()", 500);
}

function spawnAlert(text) {
  // first check if one exists, and remove it if it has not been
  var alert = $('#alert-panel');
  if (alert.length) {
    removeAlert();
  }

  var container = $('#alerts'),
      alert_div = $(document.createElement('div')),
      alert_text = $(document.createElement('span')),
      alert_anchor = $(document.createElement('a'));
  alert_text.text(text);
  alert_text.css('position', 'relative');
  alert_text.css('top', '12px');
  alert_div.append(alert_text);
  alert_div.attr('id', 'alert-panel');
  alert_div.addClass('panel');

  alert_anchor.attr('href', '#');
  alert_anchor.attr('onclick', 'removeAlert();');
  alert_anchor.addClass('controller');
  alert_anchor.text('X')
  alert_div.append(alert_anchor);

  container.append(alert_div);
  // max(text_length, 170) since 170 is the standard
  var width = Math.max(140, alert_text.get(0).offsetWidth) + 30;
  alert_div.width(width);

  togglePanel('alert');
}

$(window).load(function () {
  draw_cal(persistentCal.calendars);
  freq_set(persistentCal.frequency);

  $('#cal-button').click(function () {
    $('#cal-data').show();
    $(this).hide();
  });

  $('#add').submit(function () {
    $.post('/add',
           {'calendar-link': $('#calendar-link').val()},
           reset);
    return false;
  });

  $('#freq').submit(function () {
    $.ajax({
       type: 'PUT',
       url: '/freq',
       data: {'frequency': $('#frequency').val()},
       success: freq_reset
    });
    return false;
  });
});
