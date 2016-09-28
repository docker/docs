const debug = require('debug')('webpack-debug');
var ENV_CONFIG = require('./_webpack/_envConfig.js');
var fs = require('fs');
var path = require('path');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var _ = require('lodash');
var webpack = require('webpack');

var loaders = require('./_webpack/_commonLoaders');

/**
 * blacklist this array from being included in `externals`.
 *
 * This has the effect of making any modules in this list be
 * resolved at build time instead of runtime. This affects the
 * server bundle
 */
var blacklist = ['.bin', 'hub-js-sdk', 'dux'];
var node_modules = fs.readdirSync('node_modules').filter(function(x) {
  return !_.includes(blacklist, x);
});

/* Dux Button Config */
var elementButton = require('@dux/element-button/defaults');
var buttons = elementButton.mkButtons([{
  name: 'primary',
  color: '#FFF',
  bg: '#22B8EB'
},{
  name: 'secondary',
  color: '#FFF',
  bg: '#232C37'
},{
  name: 'coral',
  color: '#FFF',
  bg: '#FF85AF'
},{
  name: 'success',
  color: '#FFF',
  bg: '#0FD85A'
},{
  name: 'warning',
  color: '#FFF',
  bg: '#FF8546'
},{
  name: 'yellow',
  color: '#FFF',
  bg: '#FFDE50'
},{
  name: 'alert',
  color: '#FFF',
  bg: '#EB3E46'
}]);
debug('modules that will be runtime require dependencies of the server if the server requires them: ', node_modules);
var commonConfig = {
  resolve: {
    extensions: ['', '.js', '.jsx', '.json'],
    root: [
      path.resolve(__dirname, './app/scripts/'),
      path.resolve(__dirname, './app/scripts/components/')
    ],
    modulesDirectories: ['node_modules', 'app/scripts']
  },
  module: {
    preLoaders: loaders.preLoaders,
    loaders: loaders.commonLoaders
  },
  plugins: [
    ENV_CONFIG,
    new webpack.optimize.DedupePlugin(),
    new ExtractTextPlugin('public/styles/style.css', { allChunks: true })
  ],
  postcss: [
    require('postcss-import')(),
    require('postcss-constants')({
        defaults: _.merge(require('@dux/element-card/defaults')({
          capBackground: '#f1f6fb',
          borderColor: '#c4cdda'
        }),
        {
          duxElementButton: {
            radius: '.25rem',
            buttons: buttons
          }
        })
    }),
    require('postcss-each'),
    require('postcss-cssnext')({
      browsers: 'last 2 versions',
      features: {
        // https://github.com/robwierzbowski/node-pixrem/issues/40
        rem: false
      }
    }),
    require('postcss-nested'),
    require('lost')({
      gutter: '1.25rem',
      flexbox: 'flex'
    }),
    require('postcss-cssstats')(function(stats) {
      /**
       * this is in test-phase because it runs on all
       * files individually. We should either figure out
       * that that is useful or get it to run on the full postcss
       * AST or extracted CSS file.
       */
      debug(stats);
    }),
    require('postcss-url')(),
    require('cssnano')(),
    require('postcss-browser-reporter')
  ],
  eslint: {
    failOnError: true
  },
  profile: true
}

var clientBundle = _.assign({},
                            commonConfig,
                            {
                              // client.js
                              entry: './app/scripts/client.js',
                              devtool: 'eval-source-map',
                              output: {
                                path: 'app/.build/public/',
                                filename: 'js/client.js'
                              }
                            });

var serverBundle = _.assign({},
                            commonConfig,
                            {
                              // server.js
                              entry: './app/scripts/server.js',
                              output: {
                                path: 'app/.build/',
                                filename: 'server.js',
                                libraryTarget: 'commonjs2'
                              },
                              target: 'node',
                              externals: node_modules,
                              node: {
                                __dirname: '/opt/hub/'
                              }
                            });

module.exports = [
  clientBundle,
  serverBundle
];
