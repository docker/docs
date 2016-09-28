/**
 * ngTable: Table + Angular JS
 *
 * @author Vitalii Savchuk <esvit666@gmail.com>
 * @url https://github.com/esvit/ng-table/
 * @license New BSD License <http://creativecommons.org/licenses/BSD/>
 */

(function(){
    /**
     * @ngdoc object
     * @name ngTableController
     *
     * @description
     * Each {@link ngTable ngTable} directive creates an instance of `ngTableController`
     */
    angular.module('ngTable').controller('ngTableController', ['$scope', 'NgTableParams', '$timeout', '$parse', '$compile', '$attrs', '$element',
        'ngTableColumn', 'ngTableEventsChannel',
        function($scope, NgTableParams, $timeout, $parse, $compile, $attrs, $element, ngTableColumn, ngTableEventsChannel) {
            var isFirstTimeLoad = true;
            $scope.$filterRow = {};
            $scope.$loading = false;

            // until such times as the directive uses an isolated scope, we need to ensure that the check for
            // the params field only consults the "own properties" of the $scope. This is to avoid seeing the params
            // field on a $scope higher up in the prototype chain
            if (!$scope.hasOwnProperty("params")) {
                $scope.params = new NgTableParams(true);
            }
            $scope.params.settings().$scope = $scope;

            var delayFilter = (function() {
                var timer = 0;
                return function(callback, ms) {
                    $timeout.cancel(timer);
                    timer = $timeout(callback, ms);
                };
            })();

            function onDataReloadStatusChange (newStatus/*, oldStatus*/) {
                if (!newStatus || $scope.params.hasErrorState()) {
                    return;
                }

                $scope.params.settings().$scope = $scope;

                var currentParams = $scope.params;
                var filterOptions = currentParams.settings().filterOptions;

                if (currentParams.hasFilterChanges()) {
                    var applyFilter = function () {
                        currentParams.page(1);
                        currentParams.reload();
                    };
                    if (filterOptions.filterDelay) {
                        delayFilter(applyFilter, filterOptions.filterDelay);
                    } else {
                        applyFilter();
                    }
                } else {
                    currentParams.reload();
                }
            }

            // watch for when a new NgTableParams is bound to the scope
            // CRITICAL: the watch must be for reference and NOT value equality; this is because NgTableParams maintains
            // the current data page as a field. Checking this for value equality would be terrible for performance
            // and potentially cause an error if the items in that array has circular references
            $scope.$watch('params', function(newParams, oldParams){
                if (newParams === oldParams || !newParams) {
                    return;
                }

                newParams.reload();
            }, false);

            $scope.$watch('params.isDataReloadRequired()', onDataReloadStatusChange);

            this.compileDirectiveTemplates = function () {
                if (!$element.hasClass('ng-table')) {
                    $scope.templates = {
                        header: ($attrs.templateHeader ? $attrs.templateHeader : 'ng-table/header.html'),
                        pagination: ($attrs.templatePagination ? $attrs.templatePagination : 'ng-table/pager.html')
                    };
                    $element.addClass('ng-table');
                    var headerTemplate = null;

                    // $element.find('> thead').length === 0 doesn't work on jqlite
                    var theadFound = false;
                    angular.forEach($element.children(), function(e) {
                        if (e.tagName === 'THEAD') {
                            theadFound = true;
                        }
                    });
                    if (!theadFound) {
                        headerTemplate = angular.element(document.createElement('thead')).attr('ng-include', 'templates.header');
                        $element.prepend(headerTemplate);
                    }
                    var paginationTemplate = angular.element(document.createElement('div')).attr({
                        'ng-table-pagination': 'params',
                        'template-url': 'templates.pagination'
                    });
                    $element.after(paginationTemplate);
                    if (headerTemplate) {
                        $compile(headerTemplate)($scope);
                    }
                    $compile(paginationTemplate)($scope);
                }
            };

            this.loadFilterData = function ($columns) {
                angular.forEach($columns, function ($column) {
                    var result;
                    result = $column.filterData($scope);
                    if (!result) {
                        delete $column.filterData;
                        return;
                    }

                    // if we're working with a deferred object or a promise, let's wait for the promise
                    /* WARNING: support for returning a $defer is depreciated */
                    if ((angular.isObject(result) && (angular.isObject(result.promise) || angular.isFunction(result.then)))) {
                        var pData = angular.isFunction(result.then) ? result : result.promise;
                        delete $column.filterData;
                        return pData.then(function(data) {
                            // our deferred can eventually return arrays, functions and objects
                            if (!angular.isArray(data) && !angular.isFunction(data) && !angular.isObject(data)) {
                                // if none of the above was found - we just want an empty array
                                data = [];
                            }
                            $column.data = data;
                        });
                    }
                    // otherwise, we just return what the user gave us. It could be a function, array, object, whatever
                    else {
                        return $column.data = result;
                    }
                });
            };

            this.buildColumns = function (columns) {
                var result = [];
                (columns || []).forEach(function(col){
                    result.push(ngTableColumn.buildColumn(col, $scope, result));
                });
                return result
            };

            this.parseNgTableDynamicExpr = function (attr) {
                if (!attr || attr.indexOf(" with ") > -1) {
                    var parts = attr.split(/\s+with\s+/);
                    return {
                        tableParams: parts[0],
                        columns: parts[1]
                    };
                } else {
                    throw new Error('Parse error (expected example: ng-table-dynamic=\'tableParams with cols\')');
                }
            };

            this.setupBindingsToInternalScope = function(tableParamsExpr){

                // note: this we're setting up watches to simulate angular's isolated scope bindings

                // note: is REALLY important to watch for a change to the ngTableParams *reference* rather than
                // $watch for value equivalence. This is because ngTableParams references the current page of data as
                // a field and it's important not to watch this
                var tableParamsGetter = $parse(tableParamsExpr);
                $scope.$watch(tableParamsGetter, (function (params) {
                    if (angular.isUndefined(params)) {
                        return;
                    }
                    $scope.paramsModel = tableParamsGetter;
                    $scope.params = params;
                }), false);

                setupFilterRowBindingsToInternalScope();
                setupGroupRowBindingsToInternalScope();
            };

            function setupFilterRowBindingsToInternalScope(){
                if ($attrs.showFilter) {
                    $scope.$parent.$watch($attrs.showFilter, function(value) {
                        $scope.show_filter = value;
                    });
                } else {
                    $scope.$watch(hasVisibleFilterColumn, function(value){
                        $scope.show_filter = value;
                    })
                }

                if ($attrs.disableFilter) {
                    $scope.$parent.$watch($attrs.disableFilter, function(value) {
                        $scope.$filterRow.disabled = value;
                    });
                }
            }

            function setupGroupRowBindingsToInternalScope(){
                $scope.$groupRow = {};
                if ($attrs.showGroup) {
                    var showGroupGetter = $parse($attrs.showGroup);
                    $scope.$parent.$watch(showGroupGetter, function(value) {
                        $scope.$groupRow.show = value;
                    });
                    if (showGroupGetter.assign){
                        // setup two-way databinding thus allowing ngTableGrowRow to assign to the showGroup expression
                        $scope.$watch('$groupRow.show', function(value) {
                            showGroupGetter.assign($scope.$parent, value);
                        });
                    }
                } else{
                    $scope.$watch('params.hasGroup()', function(newValue) {
                        $scope.$groupRow.show = newValue;
                    });
                }
            }

            function getVisibleColumns(){
                return ($scope.$columns || []).filter(function(c){
                    return c.show($scope);
                });
            }

            function hasVisibleFilterColumn(){
                if (!$scope.$columns) return false;

                return some($scope.$columns, function($column){
                    return $column.show($scope) && $column.filter($scope);
                });
            }

            function some(array, predicate){
                var found = false;
                for (var i = 0; i < array.length; i++) {
                    var obj = array[i];
                    if (predicate(obj)){
                        found = true;
                        break;
                    }
                }
                return found;
            }

            function commonInit(){
                ngTableEventsChannel.onAfterReloadData(bindDataToScope, $scope, isMyPublisher);
                ngTableEventsChannel.onPagesChanged(bindPagesToScope, $scope, isMyPublisher);

                function bindDataToScope(params, newDatapage){
                    var visibleColumns = getVisibleColumns();
                    if (params.hasGroup()) {
                        $scope.$groups = newDatapage || [];
                        $scope.$groups.visibleColumnCount = visibleColumns.length;
                    } else {
                        $scope.$data = newDatapage || [];
                        $scope.$data.visibleColumnCount = visibleColumns.length;
                    }
                }

                function bindPagesToScope(params, newPages){
                    $scope.pages = newPages
                }

                function isMyPublisher(publisher){
                    return $scope.params === publisher;
                }
            }

            commonInit();
        }]);
})();
