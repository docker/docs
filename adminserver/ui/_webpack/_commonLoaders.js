var ExtractTextPlugin = require("extract-text-webpack-plugin");

var babelcfg = 'babel?optional[]=runtime&stage=0';

var preLoaders = [
  { test: /\.jsx?$/, exclude: /node_modules(?!\/dtr-js-sdk)/, loader: 'eslint'}
]

module.exports = {
  preLoaders: preLoaders,
  babelcfg: babelcfg
}
