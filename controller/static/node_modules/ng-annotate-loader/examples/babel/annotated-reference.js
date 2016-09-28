/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId])
/******/ 			return installedModules[moduleId].exports;
/******/
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			exports: {},
/******/ 			id: moduleId,
/******/ 			loaded: false
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.loaded = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(0);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ function(module, exports, __webpack_require__) {

	'use strict';
	
	toAnnotate.$inject = ["$scope"];
	var _createClass = (function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ('value' in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; })();
	
	function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }
	
	function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError('Cannot call a class as a function'); } }
	
	var _toImport = __webpack_require__(1);
	
	var _toImport2 = _interopRequireDefault(_toImport);
	
	console.log(_toImport2['default']);
	
	angular.module('test', []).controller('testCtrl', ["$scope", function ($scope) {}]).factory('testFactory', ["$cacheFactory", function ($cacheFactory) {
		return {};
	}]).service('testNotAnnotated', function () {
		return {};
	}).directive('testDirective', ["$timeout", function ($timeout) {
		return {
			restrict: 'E',
			controller: ["$scope", function controller($scope) {}]
		};
	}]).controller('someCtrl', someCtrl);
	
	function toAnnotate($scope) {
		'ngInject';
	}
	
	var someCtrl = (function () {
		someCtrl.$inject = ["$scope"];
		function someCtrl($scope) {
			_classCallCheck(this, someCtrl);
	
			this.doSomething();
		}
	
		_createClass(someCtrl, [{
			key: 'doSomething',
			value: function doSomething() {}
		}]);
	
		return someCtrl;
	})();
	
	console.log('after annotated function');

/***/ },
/* 1 */
/***/ function(module, exports) {

	'use strict';
	
	Object.defineProperty(exports, '__esModule', {
	  value: true
	});
	exports['default'] = 'babel-test';
	module.exports = exports['default'];

/***/ }
/******/ ]);
//# sourceMappingURL=build.js.map