'use strict';

describe('Main Module', function() {

  var module;
  var dependencies;
  dependencies = [];

  var hasModule = function(module) {
    return dependencies.indexOf(module) >= 0;
  };

  beforeEach(function() {
    // Get module
    module = angular.module('angular-jwt');
    dependencies = module.requires;
  });

  it('should load interceptor module', function() {
    expect(hasModule('angular-jwt.interceptor')).to.be.ok;
  });

  it('should load jwt module', function() {
    expect(hasModule('angular-jwt.jwt')).to.be.ok;
  });


});
