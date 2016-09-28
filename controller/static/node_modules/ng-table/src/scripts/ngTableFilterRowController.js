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
        .controller('ngTableFilterRowController', ngTableFilterRowController);

    ngTableFilterRowController.$inject = ['$scope', 'ngTableFilterConfig'];

    function ngTableFilterRowController($scope, ngTableFilterConfig){

        $scope.config = ngTableFilterConfig;

        $scope.getFilterCellCss = function (filter, layout){
            if (layout !== 'horizontal') {
                return 's12';
            }

            var size = Object.keys(filter).length;
            var width = parseInt(12 / size, 10);
            return 's' + width;
        };

        $scope.getFilterPlaceholderValue = function(filterValue/*, filterName*/){
            if (angular.isObject(filterValue)) {
                return filterValue.placeholder;
            } else {
                return '';
            }
        };
    }
})();
