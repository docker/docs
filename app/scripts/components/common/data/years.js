'use strict';
var years = function() {
  var D = new Date();
  var currentYear = D.getFullYear();
  var yearsList = [];
  for (var i = currentYear; i < currentYear + 30; i++) {
    yearsList.push({
     abbr: i,
     name: i
    });
  }
  return yearsList;
};

module.exports = years();
