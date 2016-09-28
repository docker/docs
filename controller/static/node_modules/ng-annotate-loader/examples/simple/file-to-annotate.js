'use strict';

angular.module('test', [])
	.controller('testCtrl', function($scope) {

	})
	.factory('testFactory', function($cacheFactory) {
		return {};
	})
	.service('testNotAnnotated', function() {
		return {};
	})
	.directive('testDirective', function ($timeout) {
		return {
			restrict: 'E',
			controller: function($scope) {

			}
		};
	})
	.service('namedFunction', namedFunction);

	function namedFunction($dep) {
		$dep.do();
	}