/**
 * ngTable: Table + Angular JS
 *
 * @author Vitalii Savchuk <esvit666@gmail.com>
 * @url https://github.com/esvit/ng-table/
 * @license New BSD License <http://creativecommons.org/licenses/BSD/>
 */

(function(){
    /**
     * @ngdoc directive
     * @name ngTablePagination
     * @module ngTable
     * @restrict A
     */
    angular.module('ngTable').directive('ngTablePagination', ['$compile', 'ngTableEventsChannel',
        function($compile, ngTableEventsChannel) {
            'use strict';

            return {
                restrict: 'A',
                scope: {
                    'params': '=ngTablePagination',
                    'templateUrl': '='
                },
                replace: false,
                link: function(scope, element/*, attrs*/) {

                    ngTableEventsChannel.onAfterReloadData(function(pubParams) {
                        scope.pages = pubParams.generatePagesArray();
                    }, scope, function(pubParams){
                        return pubParams === scope.params;
                    });

                    scope.$watch('templateUrl', function(templateUrl) {
                        if (angular.isUndefined(templateUrl)) {
                            return;
                        }
                        var template = angular.element(document.createElement('div'));
                        template.attr({
                            'ng-include': 'templateUrl'
                        });
                        element.append(template);
                        $compile(template)(scope);
                    });
                }
            };
        }
    ]);

})();
