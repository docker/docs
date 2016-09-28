'use strict';

var $ = require('jquery');

BaseController.$inject = ['AuthService', 'MessageService', '$scope'];
function BaseController(AuthService, MessageService, $scope) {
  var vm = this;

  $scope.$on('ngRepeatFinished', function() {
    $('a.reload').click(function() {
      window.location.reload(true);
    });
  });

  AuthService.getMyAccount()
    .then(function(data) {
      vm.account = data.data;
    }, function(error) {
      MessageService.addErrorMessage('Unable to retrieve account details', error.data);
    });
}

module.exports = BaseController;
