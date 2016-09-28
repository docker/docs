'use strict';

var moment = require('moment');

function fromCalendar() {
  return function(input) {
    return moment(input).calendar().toLowerCase();
  };
}

module.exports = fromCalendar;
