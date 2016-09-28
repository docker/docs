'use strict';

var primaries = {
  primary1: '#1aaaf8', // main blue
  primary2: '#4ec6ef',
  primary3: '#ffc400', // yellow
  primary4: '#00cbca',
  primary5: '#8f9ea8',
  primary6: '#EF4A53'
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

  infoColor: '#c4cdda',
  linkColor: primaries.primary1,
  linkHoverColor: primaries.primary1,

  coolGrey: '#8f9ea7',
  ghost: '#fafafa',
  iron: '#cccccc',
  paleOrange: '#ffbb5a',
  silver: '#c0c9ce',
  slate: '#445d6e',
  waterBlue: '#1d9fcb'
};

module.exports = {
  colors: {
    ...colors
  },
  reboot: {
    bodyColor: '#373a3c',
    bodyBackgroundColor: 'white',
    bodyFontFamily: '"Helvetica Neue", Helvetica, Arial, sans-serif',
    bodyFontBaseSize: '1rem',
    bodyLineHeight: '1.5',
    minFontSize: '16px',
    maxFontSize: '21px',
    linkColor: primaries.primary1,
    captionTableCellPadding: '.75rem',
    textMutedColor: '#e5e5e5'
  },
  typography: {
    headings: '"Helvetica Neue", Helvetica, Arial, sans-serif',
    monospace: 'Consolas, "Liberation Mono", Menlo, Courier, monospace'
  },
  media: {
    small: '40em', // ~640px
    medium: '64em', // ~1024px
    large: '90em', // ~1440px
    xlarge: '120em' // ~1920px
  },
  fonts: {
    headings: '"Helvetica Neue", Helvetica, Arial, sans-serif',
    base: '"Open Sans", "Helvetica Neue", helvetica, arial, sans-serif',
    fixed: '"Courier New", monospace',
    monospace: 'Consolas, "Liberation Mono", Menlo, Courier, monospace',
    rootSize: '16px',
    baseSize: '1rem'
  },
  layout: {
    topNavHeight: '64px',
    leftNavWidth: '98px', // Width of the left-hand menu
    leftNavWidthExpanded: '229px'
  },
  grid: {
    width: '1024px'
  }
};
