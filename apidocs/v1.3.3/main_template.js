

$(function () {
  var url = window.location.search.match(/url=([^&]+)/);
  if (url && url.length > 1) {
    url = decodeURIComponent(url[1]);
  } else {
    url = "../docs.json";
  }

  // Pre load translate...
  if(window.SwaggerTranslator) {
    window.SwaggerTranslator.translate();
  }
  window.swaggerUi = new SwaggerUi({
    spec:
INSERT SPEC HERE
,
    dom_id: "swagger-ui-container",
    supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
    validatorUrl: null,
    docExpansion: 'list',
    supportedSubmitMethods: [],
    onComplete: function(swaggerApi, swaggerUi){
      if(typeof initOAuth == "function") {
        initOAuth({
          clientId: "your-client-id",
          clientSecret: "your-client-secret",
          realm: "your-realms",
          appName: "your-app-name", 
          scopeSeparator: ","
        });
      }

      if(window.SwaggerTranslator) {
        window.SwaggerTranslator.translate();
      }

      $('pre code').each(function(i, e) {
        hljs.highlightBlock(e)
      });
    },
    onFailure: function(data) {
      log("Unable to Load SwaggerUI");
    },
    apisSorter: "alpha",
    showRequestHeaders: false
  });

  // if you have an apiKey you would like to pre-populate on the page for demonstration purposes...
  /*
    var apiKey = "myApiKeyXXXX123456789";
    $('#input_apiKey').val(apiKey);
  */

  window.swaggerUi.load();

  function log() {
    if ('console' in window) {
      console.log.apply(console, arguments);
    }
  }
});
