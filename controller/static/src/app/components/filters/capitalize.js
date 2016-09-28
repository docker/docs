'use strict';

function CapitalizeFilter() {
  return function(input) {
    if(!input) {
      return '';
    }
    return (input.charAt(0).toUpperCase() + input.slice(1));
  };
}

module.exports = CapitalizeFilter;

