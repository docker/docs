'use strict';

var angular = require('angular');
var $ = require('jquery');

EnziController.$inject = ['EnziService', 'MessageService', 'AuthService', '$state', '$scope', '$interval', 'NgTableParams'];
function EnziController(EnziService, MessageService, AuthService, $state, $scope, $interval, NgTableParams) {
  var vm = this;
  vm.error = '';

  vm.username = '';
  vm.configuredAuthBackend = '';
  vm.authBackend = '';
  vm.ldapConfig = {};
  vm.legacyAuthSettings = {};
  vm.jobs = [];
  vm.advancedLdapVisible = false;

  vm.tableParams = new NgTableParams({
    count: 3,
    sorting: {
      lastUpdated: 'desc'
    },
    filter: {
      action: 'ldap-sync'
    }
  });

  vm.addBlankUserSearchConfig = addBlankUserSearchConfig;
  vm.updateAuthConfig = updateAuthConfig;
  vm.testLogin = testLogin;
  vm.triggerSync = triggerSync;
  vm.showJobLog = showJobLog;
  vm.refreshJobs = refreshJobs;

  EnziService.getAuthBackend()
    .then(function(backend) {
      vm.configuredAuthBackend = backend;
      vm.authBackend = backend;
    },
    function(error) {
      MessageService.addErrorMessage('Error getting auth backend', error.data);
    });

  EnziService.getLegacyAuthSettings()
    .then(function(response) {
      vm.legacyAuthSettings = response.data;
    },
    function(error) {
      MessageService.addErrorMessage('Error getting legacy auth settings', error.data);
    });

  function addBlankUserSearchConfig() {
    vm.ldapConfig.userSearchConfigs.push({
      usernameAttr: 'sAMAccountName',
      fullNameAttr: 'cn',
      scopeSubtree: false
    });
  }

  AuthService.getMyAccount()
    .then(function(data) {
      vm.username = data.data.username;
    });

  EnziService.getLdapConfig()
    .then(function(response) {
      vm.ldapConfig = response.data;
      vm.ldapSyncInterval = cronToHours(vm.ldapConfig.syncSchedule);

      // If no recovery admin username is configured, add the current user's username as a recommendation
      if(!vm.ldapConfig.recoveryAdminUsername || vm.ldapConfig.recoveryAdminUsername === '') {
        vm.ldapConfig.recoveryAdminUsername = vm.username;
      }

      // Initialise user search configs if it's empty
      if(!vm.ldapConfig.userSearchConfigs || vm.ldapConfig.userSearchConfigs.length === 0) {
        addBlankUserSearchConfig();
      }

      // If advanced settings are set, then show advanced section
      if(response.data.noSimplePagination || response.data.adminSyncOpts.enableSync) {
        vm.advancedLdapVisible = true;
      }
    },
    function(error) {
      MessageService.addErrorMessage('Error retrieving data from auth service', error.data);
    });

  $scope.$on('ngRepeatFinished', function() {
    $('.help.circle.icon').popup({
      inline: true,
      hoverable: true,
      delay: {
        show: 150,
        hide: 400
      }
    });
  });

  refreshJobs();

  function refreshJobs() {
    EnziService.getJobs()
      .then(function(response) {
        vm.tableParams.settings({
          dataset: response.data.jobs
        });
      }, function(error) {
          MessageService.addErrorMessage('Error retrieving jobs auth service', error.data);
      });
  }

  function refresh() {
    $state.go('dashboard.settings.enzi', {}, { reload: true });
  }

  function isFormValid() {
    if(vm.authBackend === 'managed') {
      return true;
    }
    return $('#enzi .ui.form')
      .form({
        inline: true,
        fields: {
          serverURL: {
            identifier: 'serverURL',
            rules: [{
              type: 'empty',
              prompt: 'Please enter a server URL'
            }]
          },
          recoveryAdminUsername: {
            identifier: 'recoveryAdminUsername',
            rules: [{
              type: 'empty',
              prompt: 'Please enter a Recovery Admin Username'
            }]
          },
          readerPassword: {
            identifier: 'readerPassword',
            depends: 'readerDN',
            rules: [{
              type: 'empty',
              prompt: 'Please enter a Reader Password'
            }]
          },
          userSearchBaseDN0: {
            identifier: 'userSearchBaseDN0',
            rules: [{
              type: 'empty',
              prompt: 'Please enter a Base DN for LDAP members'
            }]
          },
          userSearchUsernameAttr0: {
            identifier: 'userSearchUsernameAttr0',
            rules: [{
              type: 'empty',
              prompt: 'Please enter the LDAP attribute for member usernames'
            }]
          }
        }
      })
      .form('validate form');
  }

  function testLogin() {
    vm.testLoginError = '';

    if(!isFormValid()) {
      MessageService.addErrorMessage('LDAP Test Failed', 'Error validating LDAP configuration');
      return;
    }

    EnziService.testLogin(vm.ldapTestUsername, vm.ldapTestPassword, vm.ldapConfig)
      .then(function(data) {
        MessageService.addSuccessMessage('LDAP Test Succeeded', 'Success!');
      }, function(error) {
        for(var i = 0; i < error.data.errors.length; i++) {
          if(error.data.errors[i].code === 'INVALID_FORM_FIELD') {
            vm.testLoginError += error.data.errors[i].detail.field + ': ' + error.data.errors[i].detail.reason + '\n';
          } else {
            vm.testLoginError += error.data.errors[i].detail.reason + '\n';
          }
          MessageService.addErrorMessage('LDAP Test Failed', '', vm.testLoginError);
        }
      });
  }

  function cronToHours(cron) {
    // By default run hourly
    if(!cron || cron === '0 0 * * *' || cron === '@hourly') {
      return 1;
    }

    if(cron === '0 0 0 1 *') {
      return 24 * 28;
    }

    if(cron === '0 0 0 * *') {
      return 24;
    }

    var cronsplit = cron.split(' ');
    // If we can't parse CRON then assume hourly
    if(cronsplit.length < 5) {
      return 1;
    }

    var hours = 0;

    // Parse days
    var daySpec = cronsplit[3];
    if(daySpec !== '*') {
      hours += parseInt(daySpec.replace('*/', '')) * 24;
    }

    // Parse hours
    var hourSpec = cronsplit[2];
    if(hourSpec !== '*') {
      hours += parseInt(hourSpec.replace('*/', ''));
    }

    return hours;
  }

  function hoursToCron(hours) {
    // If less than or equal to 1 hour, then run hourly
    if (!hours || hours <= 1) {
      return '0 0 * * *';
    }

    if (hours < 24) {
      // Run every `hours` hours.
      return '0 0 */' + hours + ' * *';
    }

    // Truncate to the largest multiple of a day.
    var days = Math.floor(hours / 24);

    if (days === 1) {
      // Run daily.
      return '0 0 0 * *';
    }

    if (days < 28) {
      // Run every `days` days.
      return '0 0 0 */' + days + ' *';
    }

    // Since it's greater than a month, run at midnight the first day every
    // month
    return '0 0 0 1 *';
  }

  function updateAuthConfig() {
    if(!isFormValid()) {
      return;
    }
    EnziService.setLegacyAuthSettings(vm.legacyAuthSettings)
      .then(null, function(error) {
        reportErrors('Error updating default permissions', error.data);
      });

    if(vm.authBackend === 'ldap') {
      vm.ldapConfig.syncSchedule = hoursToCron(vm.ldapSyncInterval);

      EnziService.updateLdapConfig(vm.ldapConfig)
        .then(function(data) {
          changeAuthBackend('ldap');
        }, function(error) {
          reportErrors('Error updating LDAP config', error.data.errors);
        });
    } else {
      changeAuthBackend(vm.authBackend);
    }
  }

  function changeAuthBackend(backend) {
    // If the configured backend is already set as this value, do nothing
    if(vm.configuredAuthBackend === backend) {
      return;
    }

    // If we're changing to LDAP, show warning modal
    if(backend === 'ldap') {
      $('#ldap-warning-modal')
        .modal({
          onApprove: function() {
            EnziService.setAuthBackend('ldap')
              .then(function(backendData) {
                refresh();
                MessageService.addSuccessMessage('Success', 'Updated to use ldap auth');
              }, function(error) {
                reportErrors('Error updating backend to LDAP', error.data.errors);
              });
          }
        })
        .modal('show');
    } else {
      EnziService.setAuthBackend(backend)
        .then(function(backendData) {
          refresh();
          MessageService.addSuccessMessage('Success', 'Updated to use ' + backend + ' auth');
        }, function(error) {
          reportErrors('Error updating backend to ' + backend, error.data.errors);
        });
    }
  }

  function triggerSync() {
    EnziService.triggerSync()
      .then(function(data) {
        MessageService.addSuccessMessage('Triggered Auth Sync');
        refreshJobs();
      }, function(error) {
        reportErrors('Error triggering sync', error.data.errors);
      });
  }

  function reportErrors(msg, errors) {
    if(!errors) {
      MessageService.addErrorMessage(msg);
      return;
    }

    for(var i = 0; i < errors.length; i++) {
      if(errors[i].code === 'INVALID_FORM_FIELD') {
        MessageService.addErrorMessage(msg, errors[i].detail.field + ': ' + errors[i].detail.reason);
      } else {
        MessageService.addErrorMessage(msg, errors[i].message);
      }
    }
  }

  function showJobLog(id) {
    vm.jobLog = '';
    EnziService.getJobLog(id)
      .then(function(data) {
        vm.jobLog = data.data;

        $('#job-log-modal')
          .modal({
            onShow: function () {
              setTimeout(function () {
                $('#job-log-modal').modal('refresh');
              }, 100);
            }
          })
          .modal('show');
      }, function(error) {
        reportErrors('Error getting job logs');
      });
  }
}

module.exports = EnziController;
