'use strict';

function SplitFilter() {
  return function(input, splitChars) {
    if(!input) {
      return input;
    }
    var re = new RegExp('[' + splitChars + ']+');
    return input.split(re);
  };
}

module.exports = SplitFilter;
