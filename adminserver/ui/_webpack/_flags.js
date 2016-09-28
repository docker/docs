var webpack = require('webpack');

// This script uses three environment variables:
//
// - local:  using webpack-dev-server for hot-reload and mocked API
// - development:  using DTR's admin server locally in a dev environment
// - production:   using DTR's admin server in a prod-like environment

process.env.DTR_API_BASE_URL = '/api/';

// We should only mock the API using webpack-dev-server manually
process.env.DTR_MOCK_API = process.env.DTR_MOCK_API || false;

module.exports = new webpack.EnvironmentPlugin([
  'DTR_MOCK_API',
  'DTR_API_BASE_URL'
]);
