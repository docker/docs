var NAVSERVICE_CONSTANT = require('./nav-service');
var ngTableDoc = angular.module('ngTableDoc', ['ui.router', 'ngMessages'])
    .config([
        'NAVSERVICE',
        '$stateProvider',
        '$urlRouterProvider',
        function (NAVSERVICE, $stateProvider, $urlRouterProvider) {
            $urlRouterProvider.otherwise("/docs/api/");
            $stateProvider
                .state('api', {
                    url:'/docs/api/:doc',
                    views:{
                        'navigation':{
                            templateUrl:'partials/nav.html',
                            controller:'NavController'
                        },
                        'content':{
                            templateUrl: function($stateParams){
                                return $stateParams.doc ? $stateParams.doc : 'partials/api/ngTable/index.html';
                            }
                        }
                    }
                })
                // note: runnable examples are not yet implemented (instead we link to external demo site)
                .state('examples',{
                    url:'/docs/api/examples',
                    templateUrl:'/partials/examples/'
                });
        }
    ]).run(function($rootScope, NavStateService){
        $rootScope.$on('$stateChangeSuccess', function(event, toState){
            NavStateService.setView(toState.name);
            console.log('Switched state to', toState.name)
        });
    })
    .factory('NavStateService', function(){
        //default state
        var navState = "gettingStarted";
        return{
            getView:function(){
                return navState;
            },
            setView:function(viewName){
                navState = viewName
            }
        }
    })
    .controller('MainController', function ($scope) {
        $scope.vm = {};
    })
    .controller('NavController', function ($scope,$location, NAVSERVICE, NavStateService) {
        var model = {};
        var currentView = NavStateService.getView();
        model.navigation = NAVSERVICE.filter(function(navItem){
            return currentView === navItem.name
        });
        model.navigation = model.navigation[0];
        $scope.vm = model;
    });

ngTableDoc.constant('NAVSERVICE', NAVSERVICE_CONSTANT.data);
