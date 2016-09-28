import webpack from 'webpack';
import webpackDevMiddleware from 'webpack-dev-middleware';
import webpackHotMiddleware from 'webpack-hot-middleware';
import webpackConfig from '../webpack.config';
import path from 'path';
import Express from 'express';
import url from 'url';
import compression from 'compression';
import proxy from './proxy';

const server = new Express();
const port = 3000;

server.disable('x-powered-by');
server.use(compression());

server.use('/_health', (req, res) => {
  res.sendStatus(200);
});

// In development: enable hot reloading
// In production: serve pre-built static assets
if (process.env.NODE_ENV !== 'production') {
  const compiler = webpack(webpackConfig);
  server.use(webpackDevMiddleware(compiler, {
    noInfo: true,
    publicPath: webpackConfig.output.publicPath,
  }));

  // Proxy requests to avoid CORS errors
  // For /v2 endpoint, rewrite cookie Domains to localhost
  // This makes the login & logout work in development
  const proxyOptions = url.parse('https://store-stage.docker.com/v2');
  proxyOptions.cookieRewrite = 'localhost';
  server.use('/v2', proxy(proxyOptions));
  server.use('/api', proxy('https://store-stage.docker.com/api'));

  // Client-side hot reload
  server.use(webpackHotMiddleware(compiler));

  // Server-side hot reload
  compiler.plugin('done', () => {
    Object.keys(require.cache).forEach(id => {
      if (id.indexOf(path.resolve(__dirname, '..', 'src')) !== -1 ||
          id.indexOf(path.resolve(__dirname)) !== -1) {
        delete require.cache[id];
      }
    });
  });
} else {
  server.use('/dist', Express.static('dist'));
}

server.use((req, res) => {
  res.send(`
    <!doctype html>
    <html>
      <head>
        <title>Docker Store</title>
        <!-- Google Webmaster -->
        <meta
          name="google-site-verification"
          content="u4812of_thlIvAZUrmDNK4dCM30Us49hReSqGAlttNM"
        />
        <!-- Segment -->
        <script type="text/javascript">
          !function(){var analytics=window.analytics=window.analytics||[];if(!analytics.initialize)if(analytics.invoked)window.console&&console.error&&console.error("Segment snippet included twice.");else{analytics.invoked=!0;analytics.methods=["trackSubmit","trackClick","trackLink","trackForm","pageview","identify","reset","group","track","ready","alias","page","once","off","on"];analytics.factory=function(t){return function(){var e=Array.prototype.slice.call(arguments);e.unshift(t);analytics.push(e);return analytics}};for(var t=0;t<analytics.methods.length;t++){var e=analytics.methods[t];analytics[e]=analytics.factory(e)}analytics.load=function(t){var e=document.createElement("script");e.type="text/javascript";e.async=!0;e.src=("https:"===document.location.protocol?"https://":"http://")+"cdn.segment.com/analytics.js/v1/"+t+"/analytics.min.js";var n=document.getElementsByTagName("script")[0];n.parentNode.insertBefore(e,n)};analytics.SNIPPET_VERSION="3.1.0";
          analytics.load("PkiQ99OVaGVevM33khgOK18hXwwFSoPT");
          }}();
        </script>
        <script src="//d2wy8f7a9ursnm.cloudfront.net/bugsnag-2.min.js"
          data-apikey=${process.env.BUGSNAG_KEY}></script>
        ${process.env.NODE_ENV === 'production' ?
        '<link rel="stylesheet" href="/dist/main.css">' : ''}
      </head>
      <body>
        <div id="app"></div>
        <script>
          window.__INITIAL_STATE__ = ${'{}'}
        </script>
        <script src="/dist/main.js"></script>
      </body>
    </html>
  `);
});

server.listen(port, error => {
  if (error) {
    console.error(error);
    return;
  }

  console.info(`Server is listening at http://localhost:${port}.`);
});
