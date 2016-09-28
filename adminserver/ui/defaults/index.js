'use strict';

var grid = {
  width: '1024px'
};

var primaries = {
  primary1: '#1aaaf8', // main blue
  primary2: '#4ec6ef',
  primary3: '#ffc400', // yellow
  primary4: '#00cbca',
  primary5: '#8f9ea8',
  primary6: '#f44336'
};

var colors = {
  ...primaries,

  warningColor: primaries.primary3,
  alertColor: primaries.primary6,
  successColor: primaries.primary4,

  // TODO: To be updated
  secondary1: '#233137',
  secondary2: '#3f5167',
  secondary3: '#556473',
  secondary4: '#7a8491',
  secondary5: '#c4cdda',
  secondary6: '#e6edf4',
  base: '#aaa',

  // TODO: These are copied from dux/dux and need to be updated
  textMuted: '#818a91',
  textColor: '#778085',
  headingLinkColor: '#445d6e',

  // Non-DUX colors (from Zeplin)
  coolGrey: '#8f9ea7',
  ghost: '#fafafa',
  iron: '#cccccc',
  paleOrange: '#ffbb5a',
  silver: '#c0c9ce',
  slate: '#445d6e',
  waterBlue: '#1d9fcb'
};


var elementButton = require('@dux/element-button/defaults');
var buttons = elementButton.mkButtons([{
  name: 'primary',
  color: '#FFF',
  bg: colors.primary1
}, {
  name: 'secondary',
  color: '#FFF',
  bg: colors.primary2
}, {
  name: 'coral',
  color: '#FFF',
  bg: '#FF85AF'
}, {
  name: 'success',
  color: '#FFF',
  bg: colors.successColor
}, {
  name: 'warning',
  color: '#FFF',
  bg: colors.warningColor
}, {
  name: 'yellow',
  color: '#FFF',
  bg: colors.warningColor
}, {
  name: 'alert',
  color: '#FFF',
  bg: colors.alertColor
}]);

module.exports = {
  layout: {
    topNavHeight: '64px',
    leftNavWidth: '98px', // Width of the left-hand menu
    leftNavWidthExpanded: '229px'
  },
  fonts: {
    base: '"Open Sans", "Helvetica Neue", helvetica, arial, sans-serif',
    fixed: '"Courier New", monospace'
  },
  reboot: {
    bodyBackgroundColor: 'white'
  },
  colors: colors,
  duxElementButton: {
    radius: '.25rem',
    buttons: buttons
  }
};
