var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');

var PRODUCTION = process.env.NODE_ENV === 'production';

module.exports = {
  devtool: PRODUCTION ? 'none' : 'eval-source-map',
  entry: PRODUCTION ?
    ['babel-polyfill', './client/index.js'] :
    ['webpack-hot-middleware/client', 'babel-polyfill', './client/index.js'],
  resolve: {
    modulesDirectories: ['client', 'src', 'src/components', 'node_modules'],
    extensions: ['', '.js', '.json', '.md'],
  },
  output: {
    path: path.join(__dirname, 'dist'),
    filename: 'main.js',
    publicPath: '/dist/',
  },
  plugins: [
    new ExtractTextPlugin('main.css'),
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify(process.env.NODE_ENV),
        DOCKERSTORE_TOKEN: JSON.stringify(process.env.DOCKERSTORE_TOKEN),
        DOCKERSTORE_CSRF: JSON.stringify(process.env.DOCKERSTORE_CSRF),
        BUGSNAG_KEY: JSON.stringify('19b15287f55d44aad4df968d78856e3d'), // bugsnag key for 'store ui'
      },
    }),
  ].concat(PRODUCTION ? [
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false,
      },
    }),
    new webpack.optimize.DedupePlugin(),
  ] : [
    new webpack.HotModuleReplacementPlugin(),
  ]),
  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel',
        exclude: /node_modules/,
        query: {
          presets: PRODUCTION ? [] : ['react-hmre'],
          plugins: ['add-module-exports'],
        },
      },
      {
        test: /\.(txt|md)$/,
        loader: 'raw-loader',
      },
      {
        test: /\.css$/,
        exclude: /(github-markdown|normalize|react-select)\.css$/,
        loader: PRODUCTION ?
        /* eslint-disable max-len*/
          ExtractTextPlugin.extract('style', 'css?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss', { allChunks: true }) :
          'style!css?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss',
        /* eslint-enable max-len*/
      },
      {
        test: /(github-markdown|normalize|react-select)\.css$/,
        loader: 'style!css',
      },
      {
        test: /\.json$/,
        loader: 'json-loader',
      },
      {
        test: /\.(png|jpg|jpeg|gif|svg|woff|woff2)$/,
        loader: 'url-loader?limit=10000',
      },
      {
        test: /\.(eot|ttf|wav|mp3)$/,
        loader: 'file-loader',
      },
    ],
  },
  postcss: (bundler) => {
    return [
      require('postcss-import')({
        addDependencyTo: bundler,
        path: [path.resolve(__dirname, 'src', 'lib', 'css')],
      }),
      require('postcss-each')(),
      require('postcss-mixins')(),
      require('postcss-simple-vars')(),
      require('postcss-cssnext')({
        autoprefixer: [
          'Android 2.3',
          'Android >= 4',
          'Chrome >= 35',
          'Firefox >= 31',
          'Explorer >= 9',
          'iOS >= 7',
          'Opera >= 12',
          'Safari >= 7.1',
        ],
      }),
      require('lost'),
      require('postcss-nested')(),
      require('postcss-calc')(),
      require('postcss-color-function')(),
      require('postcss-custom-media'),
      require('postcss-media-minmax'),
    ];
  },
};
