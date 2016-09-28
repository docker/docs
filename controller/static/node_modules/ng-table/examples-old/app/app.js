define([
    'angular',
    'ngTable'
], function (angular) {
    'use strict';

    var app = angular.module('main', ['ngTable']).
    controller('DemoCtrl', function($scope, $filter, $q, NgTableParams) {
        var data = [{name: "Moroni", age: 50, money: -10},
                    {name: "Tiancum", age: 43,money: 120},
                    {name: "Jacob", age: 27, money: 5.5},
                    {name: "Nephi", age: 29,money: -54},
                    {name: "Enos", age: 34,money: 110},
                    {name: "Tiancum", age: 43, money: 1000},
                    {name: "Jacob", age: 27,money: -201},
                    {name: "Nephi", age: 29, money: 100},
                    {name: "Enos", age: 34, money: -52.5},
                    {name: "Tiancum", age: 43, money: 52.1},
                    {name: "Jacob", age: 27, money: 110},
                    {name: "Nephi", age: 29, money: -55},
                    {name: "Enos", age: 34, money: 551},
                    {name: "Tiancum", age: 43, money: -1410},
                    {name: "Jacob", age: 27, money: 410},
                    {name: "Nephi", age: 29, money: 100},
                    {name: "Enos", age: 34, money: -100}];

        $scope.tableParams = new NgTableParams({
            $liveFiltering: true,
            page: 1,            // show first page
            total: data.length, // length of data
            count: 10           // count per page
        });

        $scope.names = function(column) {
            var def = $q.defer(),
                arr = [],
                names = [];
            angular.forEach(data, function(item){
                if ($.inArray(item.name, arr) === -1) {
                    arr.push(item.name);
                    names.push({
                        'id': item.name,
                        'title': item.name
                    });
                }
            });
            def.resolve(names);
            return def.promise;
        };

        $scope.$watch('tableParams', function(params) {
            // use built-in angular filter
            var orderedData = params.sorting ?
                                $filter('orderBy')(data, params.orderBy()) :
                                data;
            orderedData = params.filter ?
                                $filter('filter')(orderedData, params.filter) :
                                orderedData;

            params.total = orderedData.length; // set total for recalc pagination
            $scope.users = orderedData.slice((params.page - 1) * params.count, params.page * params.count);
        }, true);
    });
    return app;
});
