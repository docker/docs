/**
 * ngTable: Table + Angular JS
 *
 * @author Vitalii Savchuk <esvit666@gmail.com>
 * @url https://github.com/esvit/ng-table/
 * @license New BSD License <http://creativecommons.org/licenses/BSD/>
 */

(function () {
    /**
     * @ngdoc service
     * @name ngTableColumn
     * @module ngTable
     * @description
     * Service to construct a $column definition used by {@link ngTable ngTable} directive
     */
    angular.module('ngTable').factory('ngTableColumn', [function () {

        return {
            buildColumn: buildColumn
        };

        //////////////

        /**
         * @ngdoc method
         * @name ngTableColumn#buildColumn
         * @description Creates a $column for use within a header template
         *
         * @param {Object} column an existing $column or simple column data object
         * @param {Scope} defaultScope the $scope to supply to the $column getter methods when not supplied by caller
         * @param {Array} columns a reference to the columns array to make available on the context supplied to the
         * $column getter methods
         * @returns {Object} a $column object
         */
        function buildColumn(column, defaultScope, columns){
            // note: we're not modifying the original column object. This helps to avoid unintended side affects
            var extendedCol = Object.create(column);
            var defaults = createDefaults();
            for (var prop in defaults) {
                if (extendedCol[prop] === undefined) {
                    extendedCol[prop] = defaults[prop];
                }
                if(!angular.isFunction(extendedCol[prop])){
                    // wrap raw field values with "getter" functions
                    // - this is to ensure consistency with how ngTable.compile builds columns
                    // - note that the original column object is being "proxied"; this is important
                    //   as it ensure that any changes to the original object will be returned by the "getter"
                    (function(prop1){
                        var getterSetter = function getterSetter(/*[value] || [$scope, locals]*/) {
                            if (arguments.length === 1 && !isScopeLike(arguments[0])) {
                                getterSetter.assign(null, arguments[0]);
                            } else {
                                return column[prop1];
                            }
                        };
                        getterSetter.assign = function($scope, value){
                            column[prop1] = value;
                        };
                        extendedCol[prop1] = getterSetter;
                    })(prop);
                }
                (function(prop1){
                    // satisfy the arguments expected by the function returned by parsedAttribute in the ngTable directive
                    var getterFn = extendedCol[prop1];
                    extendedCol[prop1] = function () {
                        if (arguments.length === 1 && !isScopeLike(arguments[0])){
                            getterFn.assign(null, arguments[0]);
                        } else {
                            var scope = arguments[0] || defaultScope;
                            var context = Object.create(scope);
                            angular.extend(context, {
                                $column: extendedCol,
                                $columns: columns
                            });
                            return getterFn.call(column, context);
                        }
                    };
                    if (getterFn.assign){
                        extendedCol[prop1].assign = getterFn.assign;
                    }
                })(prop);
            }
            return extendedCol;
        }

        function createDefaults(){
            return {
                'class': createGetterSetter(''),
                filter: createGetterSetter(false),
                groupable: createGetterSetter(false),
                filterData: angular.noop,
                headerTemplateURL: createGetterSetter(false),
                headerTitle: createGetterSetter(''),
                sortable: createGetterSetter(false),
                show: createGetterSetter(true),
                title: createGetterSetter(''),
                titleAlt: createGetterSetter('')
            };
        }

        function createGetterSetter(initialValue){
            var value = initialValue;
            var getterSetter = function getterSetter(/*[value] || [$scope, locals]*/){
                if (arguments.length === 1 && !isScopeLike(arguments[0])) {
                    getterSetter.assign(null, arguments[0]);
                } else {
                    return value;
                }
            };
            getterSetter.assign = function($scope, newValue){
                value = newValue;
            };
            return getterSetter;
        }

        function isScopeLike(object){
            return object != null && angular.isFunction(object.$new);
        }
    }]);
})();
