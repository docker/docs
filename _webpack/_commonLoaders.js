var ExtractTextPlugin = require("extract-text-webpack-plugin");

var babelcfg = 'babel?optional[]=runtime&stage=0';

var preLoaders = [
  { test: /\.jsx?$/, exclude: /node_modules/, loader: 'eslint'}
]

var commonLoaders = [
  // This loader matches .js and .jsx files
  { test: /\.json$/, loader: 'json' },
  { test: /\.jsx?$/, exclude: /node_modules/, loader: babelcfg},
  { test: /dux.*\.jsx?$/, loader: babelcfg},
  { test: /hub-js-sdk.*\.jsx?$/, exclude: /hub-js-sdk.*node_modules.*\.jsx?$/, loader: babelcfg },
  { test: /\.css$/, loader: ExtractTextPlugin.extract('style-loader', 'css-loader?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss-loader') },
  { test: /\.svg$/, loader: 'svg-inline' }
]

module.exports = {
  preLoaders: preLoaders,
  commonLoaders: commonLoaders
}
