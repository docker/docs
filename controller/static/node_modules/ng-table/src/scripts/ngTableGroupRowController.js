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
        .controller('ngTableGroupRowController', ngTableGroupRowController);

    ngTableGroupRowController.$inject = ['$scope'];

    function ngTableGroupRowController($scope){

        var groupFns = [];

        init();

        function init(){
            $scope.getGroupables = getGroupables;
            $scope.getGroupTitle = getGroupTitle;
            $scope.getVisibleColumns = getVisibleColumns;
            $scope.groupBy = groupBy;
            $scope.isSelectedGroup = isSelectedGroup;
            $scope.toggleDetail = toggleDetail;

            $scope.$watch('params.group()', setGroup, true);
        }

        function changeSortDirection(){
            var newDirection;
            if ($scope.params.hasGroup($scope.$selGroup, 'asc')) {
                newDirection = 'desc';
            } else if ($scope.params.hasGroup($scope.$selGroup, 'desc')){
                newDirection = '';
            } else {
                newDirection = 'asc';
            }
            $scope.params.group($scope.$selGroup, newDirection);
        }

        function findGroupColumn(groupKey) {
            return $scope.$columns.filter(function ($column) {
                return $column.groupable($scope) === groupKey;
            })[0];
        }

        function getGroupTitle(group){
            return angular.isFunction(group) ? group.title : group.title($scope);
        }

        function getGroupables(){
            var groupableCols = $scope.$columns.filter(function ($column) {
                return $column.groupable($scope);
            });
            return groupFns.concat(groupableCols);
        }

        function getVisibleColumns(){
            return $scope.$columns.filter(function($column){
                return $column.show($scope);
            })
        }

        function groupBy(group){
            if (isSelectedGroup(group)){
                changeSortDirection();
            } else {
                if (group.groupable){
                    $scope.params.group(group.groupable($scope));
                } else{
                    $scope.params.group(group);
                }
            }
        }

        function isSelectedGroup(group){
            if (group.groupable){
                return group.groupable($scope) === $scope.$selGroup;
            } else {
                return group === $scope.$selGroup;
            }
        }

        function setGroup(group){
            var existingGroupCol = findGroupColumn($scope.$selGroup);
            if (existingGroupCol && existingGroupCol.show.assign){
                existingGroupCol.show.assign($scope, true);
            }
            if (angular.isFunction(group)) {
                groupFns = [group];
                $scope.$selGroup = group;
                $scope.$selGroupTitle = group.title;
            } else {
                // note: currently only one group is implemented
                var groupKey = Object.keys(group || {})[0];
                var groupedColumn = findGroupColumn(groupKey);
                if (groupedColumn) {
                    $scope.$selGroupTitle = groupedColumn.title($scope);
                    $scope.$selGroup = groupKey;
                    if (groupedColumn.show.assign) {
                        groupedColumn.show.assign($scope, false);
                    }
                }
            }
        }

        function toggleDetail(){
            $scope.params.settings().groupOptions.isExpanded = !$scope.params.settings().groupOptions.isExpanded;
            return $scope.params.reload();
        }
    }
})();
