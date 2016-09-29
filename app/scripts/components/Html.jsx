'use strict';
var React = require('react');
var debug = require('debug')('html');
/**
 * React class to handle the rendering of the HTML head section
 *
 * @class Head
 * @constructor
 */
var Html = React.createClass({
    /**
     * Refer to React documentation render
     *
     * @method render
     * @return {Object} HTML head section
     */
  render() {

    let bugsnag = null;
    if(process.env.BUGSNAG_API_KEY) {
      debug('BUGSNAGGG', process.env.ENV, process.env.BUGSNAG_API_KEY);
      bugsnag = (
        <script src="https://d2wy8f7a9ursnm.cloudfront.net/bugsnag-2.min.js"
                data-apikey={process.env.BUGSNAG_API_KEY}></script>
      );
    }

    let gtm;
    if (process.env.GOOGLE_TAG_MANAGER === 'gtmActive') {
      const GTMScript = `(
        function(w,d,s,l,i){
          w[l]=w[l]||[];
          w[l].push({'gtm.start': new Date().getTime(),event:'gtm.js'});
          var f=d.getElementsByTagName(s)[0],
              j=d.createElement(s),
              dl=l!='dataLayer'?'&l='+l:'';
          j.async=true;
          j.src='//www.googletagmanager.com/gtm.js?id='+i+dl;
          f.parentNode.insertBefore(j,f);
        }
      )(window,document,'script','dataLayer','GTM-KB4JTX');`;
      gtm = (<script dangerouslySetInnerHTML={{__html: GTMScript}}></script>);
    }

    const botTrackScript = `(
      function(){
        window._pxAppId ='${process.env.BOT_TRACKING_ID}';
        window._pxPubHost = 'collector.a';
        var p = document.getElementsByTagName('script')[0],
        s = document.createElement('script');
        s.async = 1;
        s.src = '//client.a.pxi.pub/${process.env.BOT_TRACKING_ID}/main.min.js';
        p.parentNode.insertBefore(s,p);
      }());`;
    const botTracking = (<script type="text/javascript" dangerouslySetInnerHTML={{__html: botTrackScript}}></script>);
    const botNoScript = (
      <noscript>
        <div style={{position: 'fixed', top: '0', left: '0', width: '1', height: '1'}}><img src={`//collector.a.pxi.pub/api/v1/collector/noScript.gif?appId=${process.env.BOT_TRACKING_ID}`}/></div>
      </noscript>
    );

    return (
        <html>
        <head>
            <meta charSet="utf-8" />
            <title>{this.props.title}</title>
            <meta name="viewport" content="width=device-width, initial-scale=1" />
            <link rel="stylesheet" href="/public/styles/main.css" />
            <link rel="stylesheet" href="/public/styles/font-awesome.min.css" />
            <link rel="stylesheet" href={`/public/styles/${process.env.CSS_FILENAME}`} />
        </head>
        <body>
            <div id="app" dangerouslySetInnerHTML={{__html: this.props.markup}}></div>
        </body>
        <script dangerouslySetInnerHTML={{__html: this.props.state}}></script>
        {bugsnag}
        <script src="https://js.recurly.com/v3/recurly.js"></script>
        <script src={`/public/js/${process.env.CLIENT_JS_FILENAME}`} defer></script>
        <script src="https://app-sj05.marketo.com/js/forms2/js/forms2.min.js"></script>
        {gtm}
        {botTracking}
        {botNoScript}
        </html>
    );
  }
});

module.exports = Html;
