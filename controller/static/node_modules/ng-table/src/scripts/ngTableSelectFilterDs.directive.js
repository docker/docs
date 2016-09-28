/**
 * ngTable: Table + Angular JS
 *
 * @author Vitalii Savchuk <esvit666@gmail.com>
 * @url https://github.com/esvit/ng-table/
 * @license New BSD License <http://creativecommons.org/licenses/BSD/>
 */

(function(){
    'use strict';

    /**
     * @ngdoc directive
     * @name ngTableSelectFilterDs
     * @module ngTable
     * @restrict A
     *
     * @description
     * Takes the array returned by $column.filterData and makes it available as `$selectData` on the `$scope`.
     *
     * The resulting `$selectData` array will contain an extra item that is suitable to represent the user
     * "deselecting" an item from a `<select>` tag
     *
     * This directive is is focused on providing a datasource to an `ngOptions` directive
     */
    angular.module('ngTable')
        .directive('ngTableSelectFilterDs', ngTableSelectFilterDs);

    ngTableSelectFilterDs.$inject = [];

    function ngTableSelectFilterDs(){
        // note: not using isolated or child scope "by design"
        // this is to allow this directive to be combined with other directives that do

        var directive = {
            restrict: 'A',
            controller: ngTableSelectFilterDsController
        };
        return directive;
    }

    ngTableSelectFilterDsController.$inject = ['$scope', '$parse', '$attrs', '$q'];
    function ngTableSelectFilterDsController($scope, $parse, $attrs, $q){

        var $column = {};
        init();

        function init(){
            $column = $parse($attrs.ngTableSelectFilterDs)($scope);
            $scope.$watch(function(){
                return $column.data;
            }, bindDataSource);
        }

        function bindDataSource(){
            getSelectListData($column).then(function(data){
                if (data && !hasEmptyOption(data)){
                    data.unshift({ id: '', title: ''});
                }
                data = data || [];
                $scope.$selectData = data;
            });
        }

        function hasEmptyOption(data) {
            var isMatch;
            for (var i = 0; i < data.length; i++) {
                var item = data[i];
                if (item && item.id === '') {
                    isMatch = true;
                    break;
                }
            }
            return isMatch;
        }

        function getSelectListData($column) {
            var data = angular.isFunction($column.data) ? $column.data() : $column.data;
            return $q.when(data);
        }
    }
})();
