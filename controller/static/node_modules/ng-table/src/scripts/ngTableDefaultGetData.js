/**
 * ngTable: Table + Angular JS
 *
 * @author Vitalii Savchuk <esvit666@gmail.com>
 * @url https://github.com/esvit/ng-table/
 * @license New BSD License <http://creativecommons.org/licenses/BSD/>
 */

(function(){
    'use strict';


    angular.module('ngTable')
        .provider('ngTableDefaultGetData', ngTableDefaultGetDataProvider);

    ngTableDefaultGetDataProvider.$inject = [];

    /**
     * @ngdoc provider
     * @name ngTableDefaultGetDataProvider
     * @description Allows for the configuration of the ngTableDefaultGetData service.
     *
     * Set filterFilterName to the name of a angular filter that knows how to apply the values returned by
     * `NgTableParams.filter()` to restrict an array of data.
     *
     * Set sortingFilterName to the name of a angular filter that knows how to apply the values returned by
     * `NgTableParams.orderBy()` to sort an array of data.
     *
     * Out of the box the `ngTableDefaultGetData` service will be configured to use the angular `filter` and `orderBy`
     * filters respectively
     */
    function ngTableDefaultGetDataProvider(){
        var provider = this;
        provider.$get = ngTableDefaultGetData;
        provider.filterFilterName = 'filter';
        provider.sortingFilterName = 'orderBy';

        ///////////

        ngTableDefaultGetData.$inject = ['$filter'];

        /**
         * @ngdoc service
         * @name ngTableDefaultGetData
         * @description A default implementation of the getData function that will apply the `filter`, `orderBy` and
         * paging values from the `NgTableParams` instance supplied to the data array supplied.
         *
         * The outcome will be to return the resulting array and to assign the total item count after filtering
         * to the `total` of the `NgTableParams` instance supplied
         */
        function ngTableDefaultGetData($filter) {

            var defaultDataOptions = {applyFilter: true, applySort: true, applyPaging: true};

            getData.applyPaging = applyPaging;
            getData.getFilterFn = getFilterFn;
            getData.getOrderByFn = getOrderByFn;

            return getData;

            function getFilterFn(params) {
                var filterOptions = params.settings().filterOptions;
                if (angular.isFunction(filterOptions.filterFn)){
                    return filterOptions.filterFn;
                } else {
                    return $filter(filterOptions.filterFilterName || provider.filterFilterName);
                }
            }

            function getOrderByFn (/*params*/){
                return $filter(provider.sortingFilterName);
            }

            function applyFilter(data, params) {
                if (!params.hasFilter()) {
                    return data;
                }

                var filter = params.filter(true);
                var filterKeys = Object.keys(filter);
                var parsedFilter = filterKeys.reduce(function(result, key){
                    result = setPath(result, filter[key], key);
                    return result;
                }, {});
                var filterFn = getFilterFn(params);
                return filterFn.call(params, data, parsedFilter, params.settings().filterOptions.filterComparator);
            }

            function applyPaging(data, params) {
                var pagedData = data.slice((params.page() - 1) * params.count(), params.page() * params.count());
                params.total(data.length); // set total for recalc pagination
                return pagedData;
            }

            function applySort(data, params) {
                var orderBy = params.orderBy();
                var orderByFn = getOrderByFn(params);
                return orderBy.length ? orderByFn(data, orderBy) : data;
            }

            function getData(data, params) {
                if (data == null){
                    return [];
                }

                var options = angular.extend({}, defaultDataOptions, params.settings().dataOptions);

                var fData = options.applyFilter ? applyFilter(data, params) : data;
                var orderedData = options.applySort ? applySort(fData, params) : fData;
                return options.applyPaging ? applyPaging(orderedData, params) : orderedData;
            }

            // Sets the value at any depth in a nested object based on the path
            // note: adapted from: underscore-contrib#setPath
            function setPath(obj, value, path) {
                var keys     = path.split('.');
                var ret      = obj;
                var lastKey  = keys[keys.length -1];
                var target   = ret;

                var parentPathKeys = keys.slice(0, keys.length -1);
                parentPathKeys.forEach(function(key) {
                    if (!target.hasOwnProperty(key)) {
                        target[key] = {};
                    }
                    target = target[key];
                });

                target[lastKey] = value;
                return ret;
            }
        }
    }
})();
