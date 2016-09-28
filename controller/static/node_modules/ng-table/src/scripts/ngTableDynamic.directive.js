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
     * @name ngTableDynamic
     * @module ngTable
     * @restrict A
     *
     * @description
     * A dynamic version of the {@link ngTable ngTable} directive that accepts a dynamic list of columns
     * definitions to render
     */
    angular.module('ngTable').directive('ngTableDynamic', [function (){

        return {
            restrict: 'A',
            priority: 1001,
            scope: true,
            controller: 'ngTableController',
            compile: function(tElement) {
                var row;

                // IE 8 fix :not(.ng-table-group) selector
                angular.forEach(tElement.find('tr'), function(tr) {
                    tr = angular.element(tr);
                    if (!tr.hasClass('ng-table-group') && !row) {
                        row = tr;
                    }
                });
                if (!row) {
                    return;
                }

                angular.forEach(row.find('td'), function(item) {
                    var el = angular.element(item);
                    var getAttrValue = function(attr){
                        return el.attr('x-data-' + attr) || el.attr('data-' + attr) || el.attr(attr);
                    };

                    // this used in responsive table
                    var titleExpr = getAttrValue('title');
                    if (!titleExpr){
                        el.attr('data-title-text', '{{$columns[$index].titleAlt(this) || $columns[$index].title(this)}}');
                    }
                    var showExpr = el.attr('ng-if');
                    if (!showExpr){
                        el.attr('ng-if', '$columns[$index].show(this)');
                    }
                });
                return function (scope, element, attrs, controller) {
                    var expr = controller.parseNgTableDynamicExpr(attrs.ngTableDynamic);

                    controller.setupBindingsToInternalScope(expr.tableParams);
                    controller.compileDirectiveTemplates();

                    scope.$watchCollection(expr.columns, function (newCols/*, oldCols*/) {
                        scope.$columns = controller.buildColumns(newCols);
                        controller.loadFilterData(scope.$columns);
                    });
                };
            }
        };
    }]);
})();
