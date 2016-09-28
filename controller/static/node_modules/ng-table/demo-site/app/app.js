(function(){
    'use strict';

    angular.module('main', ['ngRoute', 'ngTable', 'ngSocial', 'embedCodepen', 'ngSanitize'])
        .config(['$routeProvider',
            function ($routeProvider) {
                $routeProvider.
                    when('/', {
                        templateUrl: 'views/intro/overview.html',
                        controller: 'introController',
                        controllerAs: 'vm'
                    }).
                    when('/demo/todo', {
                        templateUrl: 'views/todo.html',
                        controller: function () {

                        }
                    }).
                    when('/:section/:article', {
                        templateUrl: function ($routeParams) {
                            return 'views/' + $routeParams.section + '/' + $routeParams.article + '.html';
                        },
                        controller: function () {

                        }
                    });
            }])
        .controller('introController', function(NgTableParams){
            var self = this;
            var data = [{name: "Moroni", age: 50},
                {name: "Simon", age: 43},
                {name: "Jacob", age: 27},
                {name: "Nephi", age: 29},
                {name: "Christian", age: 34},
                {name: "Tiancum", age: 43},
                {name: "Jacob", age: 27}
            ];
            self.tableParams = new NgTableParams({ count: 5}, { counts: [5, 10, 25], dataset: data});
        })
        .directive('prettycode', function() {
            return {
                restrict: 'A',
                scope: {
                    'code': '=prettycode',
                    'prettyLang': '@prettyLang'
                },
                link: function postLink(scope, element, attrs) {
                    scope.$watch('code', function(code) {
                        if (angular.isUndefined(code)) {
                            return;
                        }
                        element.html(hljs.highlight(scope.prettyLang || 'html', code).value);
                    });
                    element.html(hljs.highlight(scope.prettyLang || 'html', element.text()).value);
                }
            };
        })
        .controller('menuController', function(){
            var self = this;
            self.sections = [{
                title: 'Intro',
                hasMore: false,
                url: '#/',
                items: [{
                    title: 'Real world example',
                    url: '#/intro/demo-real-world'
                }]
            }, {
                title: 'Loading data',
                hasMore: false,
                items: [{
                    title: 'Overview',
                    url: '#/loading/overview'
                }, {
                    title: 'Managed array',
                    url: '#/loading/demo-managed-array'
                }, {
                    title: 'Lazy loading managed array',
                    url: '#/loading/demo-lazy-loaded'
                }, {
                    title: 'External array (eg server-side)',
                    url: '#/loading/demo-external-array'
                }]
            }, {
                title: 'Pagination',
                hasMore: true,
                items: [{
                    title: 'Basic example',
                    url: '#/pagination/demo-pager-basic'
                }, {
                    title: 'Change page / page controls programmatically',
                    url: '#/pagination/demo-api'
                }, {
                    title: 'Custom template',
                    url: '#/demo/todo'
                }]
            }, {
                title: 'Sorting',
                hasMore: true,
                items: [{
                    title: 'Basic example',
                    url: '#/sorting/demo-sorting-basic'
                }, {
                    title: 'Change sort order programmatically',
                    url: '#/sorting/demo-api'
                }, {
                    title: 'Enable sorting programmatically',
                    url: '#/sorting/demo-enabling'
                }]
            }, {
                title: 'Filtering',
                hasMore: true,
                items: [{
                    title: 'Basic example',
                    url: '#/filtering/demo-filtering-basic'
                }, {
                    title: 'Nested property filters',
                    url: '#/filtering/demo-nested-property'
                }, {
                    title: 'Select filters',
                    url: '#/filtering/demo-select'
                }, {
                    title: 'Custom filter template',
                    url: '#/filtering/demo-custom-template'
                }, {
                    title: 'Multiple template filters',
                    url: '#/filtering/demo-multi-template'
                }, {
                    title: 'Change filter values programmatically',
                    url: '#/filtering/demo-api'
                }, {
                    title: 'Enable filters programmatically',
                    url: '#/filtering/demo-enabling'
                }, {
                    title: 'Customize filter algorithm',
                    url: '#/filtering/demo-customize-algorithm'
                }]
            }, {
                title: 'Grouping',
                hasMore: true,
                items: [{
                    title: 'Basic example',
                    url: '#/grouping/demo-grouping-basic'
                }, {
                    title: 'Custom grouping function',
                    url: '#/grouping/demo-grouping-fn'
                }, {
                    title: 'Summary row',
                    url: '#/grouping/demo-summary'
                }, {
                    title: 'Change grouping programmatically',
                    url: '#/grouping/demo-api'
                }, {
                    title: 'Enable grouping programmatically',
                    url: '#/grouping/demo-enabling'
                }]
            }, {
                title: 'Formatting table',
                hasMore: true,
                items: [{
                    title: '<em>ngTable</em>: Data cell template',
                    url: '#/formatting/demo-cell-values'
                }, {
                    title: '<em>ngTableDynamic</em>: Data cell template (via JS)',
                    url: '#/formatting/demo-dynamic-js-values'
                }, {
                    title: '<em>ngTableDynamic</em>: Data cell template (via HTML)',
                    url: '#/formatting/demo-dynamic-html-values'
                }, {
                    title: 'Data row template',
                    url: '#/formatting/demo-row'
                }, {
                    title: 'Header cell template - basic',
                    url: '#/formatting/demo-header-cell-basic'
                }, {
                    title: 'Header cell template - full',
                    url: '#/formatting/demo-header-cell-full'
                }, {
                    title: 'Custom header',
                    url: '#/formatting/demo-custom-header'
                }]
            }, {
                title: 'Working with columns',
                hasMore: false,
                items: [{
                    title: 'Show/hide columns',
                    url: '#/columns/demo-visibility'
                }, {
                    title: 'Reorder columns',
                    url: '#/columns/demo-reordering'
                }]
            }, {
                title: 'Editing data',
                hasMore: false,
                items: [{
                    title: 'Inline row edit',
                    url: '#/editing/demo-inline'
                }, {
                    title: 'Batch table edit',
                    url: '#/editing/demo-batch'
                }]
            }, {
                title: 'Events',
                hasMore: false,
                items: [{
                    title: 'Subscribe to events',
                    url: '#/events/demo-subscribe'
                }, {
                    title: 'Unsubscribe from events',
                    url: '#/events/demo-unsubscribe'
                }]
            }, {
                title: 'Global customizations',
                hasMore: false,
                items: [{
                    title: 'Response interceptors',
                    url: '#/global-customization/demo-response-interceptors'
                }, {
                    title: 'Change default parameters and settings',
                    url: '#/global-customization/demo-defaults'
                }, {
                    title: 'Replace default filter/sorting algorithm',
                    url: '#/demo/todo'
                }]
            }, {
                title: 'Miscellaneous',
                hasMore: false,
                items: [{
                    title: 'Saving params in url',
                    url: '#/demo/todo'
                }]
            }, {
                title: 'Plugins',
                hasMore: false,
                items: [{
                    title: 'Resize columns',
                    url: '#/demo/todo'
                }, {
                    title: 'Export to CSV',
                    url: '#/demo/todo'
                }]
            }];
        });
})();
