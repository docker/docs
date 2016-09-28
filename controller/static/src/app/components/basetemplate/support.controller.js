'use strict';

var $ = require('jquery');

SupportController.$inject = ['MessageService', '$http'];
function SupportController(MessageService, $http) {
  var vm = this;

  vm.downloadSupportDump = downloadSupportDump;

  function downloadSupportDump() {
    $http
      .post('/api/support', {}, { responseType: 'arraybuffer' })
      .then(function(response) {
        var file = new Blob([ response.data ], {
          type: 'application/zip'
        });
        var filename = 'docker-support.zip';
        var rx = response.headers('Content-Disposition').match(/inline; filename='(.*?)'/);
        if(rx.length > 0) {
          filename = rx[1];
        }
        //trick to download store a file having its URL
        var fileURL = URL.createObjectURL(file);
        var a = document.createElement('a');
        a.href = fileURL;
        a.target = '_blank';
        a.download = filename;
        document.body.appendChild(a);
        a.click();
      },
      function(error) {
        MessageService.addErrorMessage('Error creating support dump', error.statusText);
      });
  }
}

module.exports = SupportController;
