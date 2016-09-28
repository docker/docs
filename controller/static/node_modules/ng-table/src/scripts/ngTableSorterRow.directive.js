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
        .directive('ngTableSorterRow', ngTableSorterRow);

    ngTableSorterRow.$inject = [];

    function ngTableSorterRow(){
        var directive = {
            restrict: 'E',
            replace: true,
            templateUrl: 'ng-table/sorterRow.html',
            scope: true,
            controller: 'ngTableSorterRowController'
        };
        return directive;
    }
})();
