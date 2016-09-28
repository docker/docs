'use strict';

VersionInterceptor.$inject = ['MessageService'];
function VersionInterceptor(MessageService) {
  var service = this;

  /*global TAG*/
  var ucpVersion = TAG;

  service.response = function(r) {
    var headerVersion = r.headers('ucp-version');
    if(!headerVersion) {
      return r;
    }

    if(ucpVersion !== headerVersion.split(' ')[0]) {
      MessageService.addWarningMessage('The UCP server has been updated', '<a class="reload" alt="Reload">Please click here to reload the UI now</a>');
    }
    return r;
  };
}

module.exports = VersionInterceptor;


