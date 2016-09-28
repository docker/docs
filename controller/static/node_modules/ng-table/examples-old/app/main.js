require.config({
    paths: {
        jquery: '../js/jquery-1.9.1.min',
        angular: '../js/angular.min',
        ngTable: '../../dist'
    },
    shim: {
        'angular': {'exports': 'angular'}
    }
});

require([
    'jquery',
    'angular',
    'app'
], function ($, angular, app) {
    'use strict';

    angular.bootstrap(document, ['main']);
});
