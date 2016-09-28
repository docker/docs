'use strict';

import babelTestMsg from './to-import';
console.log(babelTestMsg);

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
	.controller('someCtrl', someCtrl);
  
function toAnnotate($scope) {
  'ngInject';
}

class someCtrl {
	constructor($scope) {
		this.doSomething();
	}

	doSomething() {

	}
}

console.log('after annotated function');