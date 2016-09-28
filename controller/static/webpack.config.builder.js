'use strict';

var webpack = require('webpack'),
    path = require('path');

var HtmlWebpackPlugin = require('html-webpack-plugin');

var APP = __dirname;

if(!process.env.TAG) {
  console.log('TAG is not set');
  return;
}

module.exports = function(options) {
  // Defaults
  var entry = {
    app: ['./src/app/entrypoint.js']
  };

  var eslint = {
      configFile: './.eslintrc',
      formatter: require('eslint/lib/formatters/stylish'),
      failOnWarning: false,
      failOnError: false
  };

  var plugins = [
    // HtmlWebpackPlugin generates index.html from our blueimp template and then
    // copies it to index.html
    new HtmlWebpackPlugin({
      favicon: './src/favicon.ico',
      template: './src/index.html',
      inject: 'head'
    }),

    new webpack.DefinePlugin({
      ORG: JSON.stringify(process.env.ORG),
      TAG: JSON.stringify(process.env.TAG),
      REQUIRE_LICENSE: JSON.stringify(process.env.REQUIRE_LICENSE)
    })
  ];

  var loaders = [
    { test: /\.js$/, loader: 'eslint-loader', exclude: /node_modules|semantic\/dist/ },
		{ test: /angular-ui-codemirror.*\.js$/, loader: 'ng-annotate?map=false' },
    { test: /\.(woff|woff2|ttf|eot|png|svg)(\?]?.*)?$/, loader: 'file-loader' },
    { test: /\.html$/, loader: 'html-loader', exclude: /node_modules/ },
    { test: /\.css$/, loader: 'style-loader!css-loader' }
  ];

  if(options.failOnError) {
    // NoErrorsPlugin will fail the build if an error is detected while bundling
    plugins.unshift(new webpack.NoErrorsPlugin());
  }

  // Parse options
  if(options.minify) {
    plugins.unshift(new webpack.optimize.OccurrenceOrderPlugin());

    // Remove duplicate code
    plugins.unshift(new webpack.optimize.DedupePlugin());

    // Minification plugin
    plugins.unshift(new webpack.optimize.UglifyJsPlugin({
      sourceMap: false,
      compress: {
        warnings: false
      }
    }));
  }

  var devServer = {};
  if(options.devServer) {
    if(!process.env.CONTROLLER_URL) {
      console.log('CONTROLLER_URL not set');
      return {};
    }


    devServer = {
      watchOptions: {
        poll: process.env.POLL || false
      },
      contentBase: 'dist/',
      proxy: {
        '*': {
          target: process.env.CONTROLLER_URL,
          secure: false
        }
      },
      stats: { colors: true }
    };
  }

  return {
    context: APP,
    eslint: eslint,
    debug: options.debug,
    devtool: options.devtool,
    entry: entry,
    output: {
      path: APP + '/dist/',
      filename: './bundle.js'
    },
    module: {
      loaders: loaders
    },
    devServer: devServer,
    plugins: plugins
  };
};
