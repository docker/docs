'use strict';

var angular = require('angular');

// TODO: Remove the need for this by removing mentions of <dropdown> and <dropdown-group>
module.exports = angular.module('angularify.semantic.dropdown', [])
  .controller('DropDownController', ['$scope',
    function ($scope) {
      $scope.items = [];
      this.add_item = function (scope) {
        $scope.items.push(scope);
        return $scope.items;
      };
      this.remove_item = function (scope) {
        var index = $scope.items.indexOf(scope);
        if (index !== -1) {
          $scope.items.splice(index, 1);
        }
      };
      this.update_title = function (title) {
        for (var i in $scope.items) {
          $scope.items[i].title = title;
        }
      };
    }
  ])
  .directive('dropdown', function () {
    return {
      restrict: 'E',
      replace: true,
      transclude: true,
      controller: 'DropDownController',
      scope: {
        title: '@',
        open: '@',
        model: '=ngModel'
      },
      template: '<div class="{{ dropdown_class }}">' +
        '<div class="default text">{{ title }}</div>' +
        '<i class="dropdown icon"></i>' +
        '<div class="{{ menu_class }}"  ng-transclude>' +
        '</div>' +
        '</div>',
      link: function (scope, element, attrs, DropDownController) {
        scope.dropdown_class = 'ui selection dropdown';
        scope.menu_class = 'menu transition hidden';
        scope.original_title = scope.title;

        if (scope.open === 'true') {
          scope.is_open = true;
          scope.dropdown_class = scope.dropdown_class + ' active visible';
          scope.menu_class = scope.menu_class + ' visible';
        } else {
          scope.is_open = false;
        }
        DropDownController.add_item(scope);

        /*
         * Watch for title changing
         */
        scope.$watch('title', function (val, oldVal) {
            if (val === undefined || val === oldVal || val === scope.original_title) {
            return;

      }
      scope.model = val;
    });

    /*
     * Watch for ng-model changing
     */
        scope.$watch('model', function (val) {
            // update title or reset the original title if its empty
            scope.model = val;
            DropDownController.update_title(val || scope.original_title);
            });

    /*
     * Click handler
     */
        element.bind('click', function () {
            if (scope.is_open === false) {
            scope.$apply(function () {
                scope.dropdown_class = 'ui selection dropdown active visible';
                scope.menu_class = 'menu transition visible';
                });
            } else {
            if (scope.title !== scope.original_title) {
            scope.model = scope.title;
            }
            scope.$apply(function () {
                scope.dropdown_class = 'ui selection dropdown';
                scope.menu_class = 'menu transition hidden';
                });
            }
            scope.is_open = !scope.is_open;
            });
      }
    };
  })
  .directive('dropdownGroup', function () {
    return {
      restrict: 'AE',
      replace: true,
      transclude: true,
      require: '^dropdown',
      scope: {
        title: '=title'
      },
      template: '<div class="item" ng-transclude >{{ item_title }}</div>',
      link: function (scope, element, attrs, DropDownController) {

        // Check if title= was set... if not take the contents of the dropdown-group tag
        // title= is for dynamic variables from something like ng-repeat {{variable}}
        if (scope.title === undefined) {
          scope.item_title = element.children()[0].innerHTML;
        } else {
          scope.item_title = scope.title;
        }

        //
        // Menu item click handler
        //
        element.bind('click', function () {
          DropDownController.update_title(scope.item_title);
        });
      }
    };
  });
