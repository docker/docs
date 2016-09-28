'use strict';
require('babel-register');

//const debug = require('debug')('webpack-debug');
var webpack = require('webpack');
var path = require('path');

var loaders = require('./_webpack/_commonLoaders');
var flags = require('./_webpack/_flags');

var defaults = require('./src/scripts/defaults');

var resolve = {
  extensions: ['', '.js', '.jsx', 'jsx', 'js', '.css'],
  fallback: path.join(__dirname, 'node_modules'),
  root: path.join(__dirname, 'src', 'scripts')
};

module.exports = {
  entry: [
    './src/scripts/index'
  ],
  devtool: 'eval',
  output: {
    path: path.join(__dirname, 'src'),
    filename: 'bundle.js',
    sourceMapFileName: 'bundle.map',
    publicPath: '/public/'
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
        'NODE_ENV': JSON.stringify('development')
      }
    }),
    new webpack.NoErrorsPlugin(),
    flags
  ],
  module: {
    preLoaders: loaders.preLoaders,
    loaders: [
      {
        test: /\.json$/,
        loaders: [ 'json' ]
      },
      {
        test: /\.jsx?$/,
        include: /(dtr-js-sdk|@dux|redux-ui)/,
        loader: 'babel-loader',
        query: {
            plugins: ['transform-decorators-legacy'],
            presets: ['es2015', 'stage-0', 'react']
        }
      },
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
        query: {
            plugins: ['transform-decorators-legacy'],
            presets: ['es2015', 'stage-0', 'react']
        }
      },
      {
        test: /\.css$/,
        include: /node_modules/,
        exclude: /node_modules\/@dux/,
        loader: 'style-loader!css-loader'
      },
      {
        test: /\.css$/,
        include: /node_modules\/@dux/,
        loader: 'style-loader!css-loader?modules&importLoaders=1&localIdentName=dux__[name]__[local]!postcss-loader'
      },
      {
        test: /\.css$/,
        exclude: /node_modules/,
        loader: 'style-loader!css-loader?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss-loader'
      }
    ]
  },
  resolve: resolve,
  resolveLoader: resolve,
  postcss: [
    require('postcss-mixins'),
    require('postcss-simple-vars'),
    require('postcss-constants')({
      defaults: defaults
    }),
    require('postcss-each'),
    require('postcss-cssnext')({
      browsers: 'last 2 versions',
      features: {
        // https://github.com/robwierzbowski/node-pixrem/issues/40
        rem: false
      },
      import: true,
      compress: false,
      messages: true
    }),
    require('postcss-nested'),
    require('lost')
  ],
  eslint: {
    failOnError: false
  },
  profile: true
};
