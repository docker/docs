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
        .directive('ngTableGroupRow', ngTableGroupRow);

    ngTableGroupRow.$inject = [];

    function ngTableGroupRow(){
        var directive = {
            restrict: 'E',
            replace: true,
            templateUrl: 'ng-table/groupRow.html',
            scope: true,
            controller: 'ngTableGroupRowController',
            controllerAs: 'dctrl'
        };
        return directive;
    }
})();
