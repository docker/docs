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
     * @name ngTable
     * @module ngTable
     * @restrict A
     *
     * @description
     * Directive that instantiates {@link ngTableController ngTableController}.
     */
    angular.module('ngTable').directive('ngTable', ['$q', '$parse',
        function($q, $parse) {
            'use strict';

            return {
                restrict: 'A',
                priority: 1001,
                scope: true,
                controller: 'ngTableController',
                compile: function(element) {
                    var columns = [],
                        i = 0,
                        dataRow,
                        groupRow,
                        rows = [];

                    angular.forEach(element.find('tr'), function(tr) {
                        rows.push(angular.element(tr))
                    });
                    dataRow = rows.filter(function(tr){
                        return !tr.hasClass('ng-table-group');
                    })[0];
                    groupRow = rows.filter(function(tr){
                        return tr.hasClass('ng-table-group');
                    })[0];

                    if (!dataRow) {
                        return;
                    }
                    angular.forEach(dataRow.find('td'), function(item) {
                        var el = angular.element(item);
                        if (el.attr('ignore-cell') && 'true' === el.attr('ignore-cell')) {
                            return;
                        }

                        var getAttrValue = function(attr){
                            return el.attr('x-data-' + attr) || el.attr('data-' + attr) || el.attr(attr);
                        };
                        var setAttrValue = function(attr, value){
                            if (el.attr('x-data-' + attr)){
                                el.attr('x-data-' + attr, value)
                            } else if (el.attr('data' + attr)){
                                el.attr('data' + attr, value)
                            } else {
                                el.attr(attr, value)
                            }
                        };

                        var parsedAttribute = function(attr) {
                            var expr = getAttrValue(attr);
                            if (!expr){
                                return undefined;
                            }

                            var localValue;
                            var getter = function (context) {
                                if (localValue !== undefined){
                                    return localValue;
                                }
                                return $parse(expr)(context);
                            };
                            getter.assign = function($scope, value){
                                var parsedExpr = $parse(expr);
                                if (parsedExpr.assign) {
                                    // we should be writing back to the parent scope as this is where the expression
                                    // came from
                                    parsedExpr.assign($scope.$parent, value);
                                } else {
                                    localValue = value;
                                }
                            };
                            return getter;
                        };
                        var titleExpr = getAttrValue('title-alt') || getAttrValue('title');
                        if (titleExpr){
                            el.attr('data-title-text', '{{' + titleExpr + '}}'); // this used in responsive table
                        }
                        // NOTE TO MAINTAINERS: if you add extra fields to a $column be sure to extend ngTableColumn with
                        // a corresponding "safe" default
                        columns.push({
                            id: i++,
                            title: parsedAttribute('title'),
                            titleAlt: parsedAttribute('title-alt'),
                            headerTitle: parsedAttribute('header-title'),
                            sortable: parsedAttribute('sortable'),
                            'class': parsedAttribute('header-class'),
                            filter: parsedAttribute('filter'),
                            groupable: parsedAttribute('groupable'),
                            headerTemplateURL: parsedAttribute('header'),
                            filterData: parsedAttribute('filter-data'),
                            show: el.attr("ng-if") ? parsedAttribute('ng-if') : undefined
                        });

                        if (groupRow || el.attr("ng-if")){
                            // change ng-if to bind to our column definition which we know will be writable
                            // because this will potentially increase the $watch count, only do so if we already have an
                            // ng-if or when we definitely need to change visibility of the columns.
                            // currently only ngTableGroupRow directive needs to change visibility
                            setAttrValue('ng-if', '$columns[' + (columns.length - 1) + '].show(this)');
                        }
                    });
                    return function(scope, element, attrs, controller) {
                        scope.$columns = columns = controller.buildColumns(columns);

                        controller.setupBindingsToInternalScope(attrs.ngTable);
                        controller.loadFilterData(columns);
                        controller.compileDirectiveTemplates();
                    };
                }
            }
        }
    ]);
})();
