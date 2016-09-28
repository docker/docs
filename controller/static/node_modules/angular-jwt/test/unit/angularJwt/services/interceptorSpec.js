'use strict';

describe('interceptor', function() {

  beforeEach(function() {
    module('angular-jwt.interceptor');
  });

  afterEach(inject(function($httpBackend) {
    $httpBackend.verifyNoOutstandingExpectation();
    $httpBackend.verifyNoOutstandingRequest();
  }));


  it('should intercept requests when added to $httpProvider.interceptors and set token', function (done) {
    module( function ($httpProvider, jwtInterceptorProvider) {
      jwtInterceptorProvider.tokenGetter = function() {
        return 123;
      }
      $httpProvider.interceptors.push('jwtInterceptor');
    });

    inject(function ($http, $httpBackend) {
        $http({url: '/hello'}).success(function (data) {
          expect(data).to.be.equal('hello');
          done();
        });

        $httpBackend.expectGET('/hello', function (headers) {
          return headers.Authorization === 'Bearer 123';
        }).respond(200, 'hello');
        $httpBackend.flush();
    });

  });

  it('should work with promises', function (done) {
    module( function ($httpProvider, jwtInterceptorProvider) {
      jwtInterceptorProvider.tokenGetter = function($q) {
        return $q.when(345);
      }
      $httpProvider.interceptors.push('jwtInterceptor');
    });

    inject(function ($http, $httpBackend) {
        $http({url: '/hello'}).success(function (data) {
          expect(data).to.be.equal('hello');
          done();
        });

        $httpBackend.expectGET('/hello', function (headers) {
          return headers.Authorization === 'Bearer 345';
        }).respond(200, 'hello');
        $httpBackend.flush();
    });

  });

  it('should not send it if no tokenGetter', function (done) {
    module( function ($httpProvider, jwtInterceptorProvider) {
      $httpProvider.interceptors.push('jwtInterceptor');
    });

    inject(function ($http, $httpBackend) {
        $http({url: '/hello'}).success(function (data) {
          expect(data).to.be.equal('hello');
          done();
        });

        $httpBackend.expectGET('/hello', function (headers) {
          return !headers.Authorization;
        }).respond(200, 'hello');
        $httpBackend.flush();
    });

  });

  it('should add the token to the url params when the configuration option is set', function (done) {
    module( function ($httpProvider, jwtInterceptorProvider) {
      jwtInterceptorProvider.urlParam = 'access_token';
      jwtInterceptorProvider.tokenGetter = function() {
        return 123;
      }
      $httpProvider.interceptors.push('jwtInterceptor');
    });

    inject(function ($http, $httpBackend) {
        $http({url: '/hello'}).success(function (data) {
          expect(data).to.be.equal('hello');
          done();
        });

        $httpBackend.expectGET('/hello?access_token=123', function (headers) {
          return headers.Authorization === undefined;
        }).respond(200, 'hello');
        $httpBackend.flush();
    });

  });
});
