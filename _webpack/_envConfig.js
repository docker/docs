var webpack = require('webpack');
var debug = require('debug')('_envConfig');

debug(process.env.ENV);

function isProd() {
  return process.env.ENV === 'production';
}

/**
 * `development` is staging.
 * `staging` is also staging.
 *
 * The keys in these objects are possible ENV configurations
 *
 * Pointing production to staging during Alpha
 */

var HUB_URLS = {
  local: 'https://hub.dev.docker.com',
  development: 'https://hub-stage.docker.com',
  staging: 'https://hub-stage.docker.com',
  production: 'https://hub.docker.com'
}

var RECURLY_PUBLIC_KEY = {
  local: 'sjc-9XwqFDBZALFs9BP9dn3J8e',
  development: 'sjc-9XwqFDBZALFs9BP9dn3J8e',
  production: 'sjc-JIfmXVz2OVkg3xg10NhWm1'
}

// NO LONGER USED
// var MUNCHKIN_CODE = {
//   staging: '453-IHP-147',
//   production: '929-FJL-178'
// }

var BUGSNAG_API_KEY = {
  staging: 'ec43d0373895ee5eb76ec75301157a85',
  production: 'd639ea00dd6e493b739de27a7ee0f90c'
}

var TUTUM_SIGNIN_URLS = {
  development: 'https://dashboard-staging.tutum.co/login/docker/',
  production: 'https://dashboard.tutum.co/login/docker/'
}

var COOKIE_DOMAIN = {
  local: 'hub.dev.docker.com',
  development: 'bagels.docker.com',
  staging: 'hub-stage.docker.com',
  production: 'hub.docker.com'
};

// NODE_ENV is an express thing but is NOT being used by us
if (isProd() || process.env.ENV === 'staging') {
  process.env.NODE_ENV = 'production'
}

if ( !~['development', 'local', 'staging', 'production'].indexOf(process.env.ENV) ) {
  process.env.ENV = 'development';
}

// Override some ENV vars
process.env.HUB_API_BASE_URL = process.env.HUB_API_BASE_URL || HUB_URLS[process.env.ENV] || 'https://hub-stage.docker.com';
process.env.REGISTRY_API_BASE_URL = HUB_URLS[process.env.ENV] || 'https://hub-stage.docker.com';
process.env.RECURLY_PUBLIC_KEY = RECURLY_PUBLIC_KEY[process.env.ENV] || 'sjc-9XwqFDBZALFs9BP9dn3J8e';


process.env.CLIENT_JS_FILENAME = process.env.CLIENT_JS_FILENAME || 'client.js';
process.env.CSS_FILENAME = process.env.CSS_FILENAME || 'style.css';
process.env.COOKIE_DOMAIN = COOKIE_DOMAIN[process.env.ENV] || 'bagels.docker.com';

process.env.BOT_TRACKING_ID = 'PXbPb4C2uT';

if(isProd()) {
  process.env.BUGSNAG_API_KEY = BUGSNAG_API_KEY.production;
  process.env.GOOGLE_TAG_MANAGER = 'gtmActive';
  process.env.BOT_TRACKING_ID = 'PXPmP8ILuI';
} else if(process.env.ENV === 'staging') {
  process.env.BUGSNAG_API_KEY = BUGSNAG_API_KEY.staging;
  process.env.GOOGLE_TAG_MANAGER = 'gtmDisabled';
}

process.env.TUTUM_SIGNIN_URL = TUTUM_SIGNIN_URLS[process.env.ENV] || 'https://dashboard-staging.tutum.co/login/docker/';

process.env.NAUTILUS_API_BASE_URL = HUB_URLS[process.env.ENV] + '/api/nautilus/v1';

debug(process.env);

module.exports = new webpack.EnvironmentPlugin([
  'BUGSNAG_API_KEY',
  'CLIENT_JS_FILENAME',
  'CSS_FILENAME',
  'ENV',
  'NODE_ENV',
  'GOOGLE_TAG_MANAGER',
  'BOT_TRACKING_ID',
  'HUB_API_BASE_URL',
  'NAUTILUS_API_BASE_URL',
  'RECURLY_PUBLIC_KEY',
  'REGISTRY_API_BASE_URL',
  'TUTUM_SIGNIN_URL'
]);
