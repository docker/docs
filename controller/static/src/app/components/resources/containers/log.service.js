'use strict';

var $ = require('jquery');

LogService.$inject = ['$http', 'LogConstants'];
function LogService($http, LogConstants) {

  var defaultParams = {
    stdout: 1,
    stdin: 1,
    stderr: 1,
    follow: 0,
    since: 0,
    timestamps: 0,
    tail: LogConstants.DISABLE_TAIL
  };

  function get(id, params) {
    if(!params) {
      params = defaultParams;
    }
    var host = window.location.host;
    var scheme = window.location.protocol.replace('http', 'ws');
    var wsUrl = scheme + '//' + host + '/containerlogs?' + $.param(params) + '&id=' + id + '&token=';

    return $http.post('/api/containerlogs/' + id)
      .then(function(data) {
        return new WebSocket(wsUrl + data.data.token);
      });
  }

  return {
    get: get,
    defaultParams: defaultParams
  };
}

module.exports = LogService;
