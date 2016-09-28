'use strict';

require('./vendor')();

var angular = require('angular');
var $ = require('jquery');

var routes = require('./config.routes');

angular
  .module('ducp', [
    'ui.router',
    'ui.codemirror',
    'angular-jwt',
    'angular-storage',
    'ncy-angular-breadcrumb',
    'ngSanitize',
    'ngTable',
    'ngclipboard',
    'ngFileUpload',
    'luegg.directives',

    // TODO: Collate these into a single common module
    require('./components/ng-table').name,
    require('./components/core').name,
    require('./components/directives').name,
    require('./components/services').name,
    require('./components/filters').name,

    // TODO: Deprecate this
    require('./components/core/angular-semantic-ui').name,

    require('./components/basetemplate').name,
    require('./components/login').name,
    require('./components/error').name,
    require('./components/dashboard').name,
    require('./components/usermanagement').name,
    require('./components/resources').name,
    require('./components/adminsettings').name,
    require('./components/userprofile').name
  ])
  .service('versionInterceptor', require('./components/interceptors/version.js'))
  .service('loaderInterceptor', require('./components/interceptors/loader.js'))
  .service('errorInterceptor', require('./components/interceptors/error.js'))
  .config(['$httpProvider', 'jwtInterceptorProvider', function ($httpProvider, jwtInterceptorProvider) {
    jwtInterceptorProvider.tokenGetter = ['store', function(store) {
      return store.get('sessionToken');
    }];

    $httpProvider.interceptors.push('jwtInterceptor');
    $httpProvider.interceptors.push('versionInterceptor');
    $httpProvider.interceptors.push('loaderInterceptor');
    $httpProvider.interceptors.push('errorInterceptor');
  }])
  .config(routes)
  .run(['$rootScope', '$state', '$stateParams', '$window', 'AuthService', '$location', 'LoaderService', 'MessageService', '$timeout', 'store',
     function ($rootScope, $state, $stateParams, $window, AuthService, $location, LoaderService, MessageService, $timeout, store) {
       $rootScope.$state = $state;
       $rootScope.$stateParams = $stateParams;

       // Disable angular-storage caching
       store.useCache = false;

       // If the user already has a token in their browser's local storage, or if they are authenticated via a TLS certificate,
       // then they should be able to grab their username and account details straightaway.  If they are not logged in, this request
       // should 401 and they will be redirected to the login page.
       AuthService.updateRootScopeAccountInfo();

       $rootScope.$on('$stateChangeStart', function(event, toState, toStateParams) {
         // If route has redirect, then use it
         if (toState.redirectTo) {
           event.preventDefault();
           $state.go(toState.redirectTo, toStateParams.params);
         }

         // Remove any created semantic modals from the DOM
         angular.element('.ui.modals').remove();

         LoaderService.setMessage('');
         MessageService.clearMessages();

         if(toState && toState.authenticate) {
           if (!store.get('sessionToken')) {
             event.preventDefault();
             $rootScope.returnState = toState.name;
             $rootScope.returnStateParams = toStateParams;
             $state.go('login');
           }
         }
       });

       $rootScope.$on('$stateChangeError', function(event, toState, toParams, fromState, fromParams) {
         event.preventDefault();
         $state.go('error');
       });

       $rootScope.$on('$stateChangeSuccess', function(event, toState, toParams, fromState, fromParams) {
         $rootScope.toState = toState;
         $rootScope.toParams = toParams;
         $rootScope.fromState = fromState;
         $rootScope.fromParams = fromParams;
       });
     }
  ]);
