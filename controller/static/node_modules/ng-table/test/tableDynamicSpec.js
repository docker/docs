describe('ng-table-dynamic', function() {
    var dataset = [
        { id: 1, name: "Moroni", age: 50, money: -10 },
        { id: 2, name: "Tiancum", age: 43, money: 120 },
        { id: 3, name: "Jacob", age: 27, money: 5.5 },
        { id: 4, name: "Nephi", age: 29, money: -54 },
        { id: 5, name: "Enos", age: 34, money: 110 },
        { id: 6, name: "Tiancum", age: 43, money: 1000 },
        { id: 7, name: "Jacob", age: 27, money: -201 },
        { id: 8, name: "Nephi", age: 29, money: 100 },
        { id: 9, name: "Enos", age: 34, money: -52.5 },
        { id: 10, name: "Tiancum", age: 43, money: 52.1 },
        { id: 11, name: "Jacob", age: 27, money: 110 },
        { id: 12, name: "Nephi", age: 29, money: -55 },
        { id: 13, name: "Enos", age: 34, money: 551 },
        { id: 14, name: "Tiancum", age: 43, money: -1410 },
        { id: 15, name: "Jacob", age: 27, money: 410 },
        { id: 16, name: "Nephi", age: 29, money: 100 },
        { id: 17, name: "Enos", age: 34, money: -100 }
    ];

    beforeEach(module('ngTable'));

    var scope;
    beforeEach(inject(function($rootScope) {
        scope = $rootScope.$new(true);
    }));

    describe('basics', function(){
        var elm;
        beforeEach(inject(function($compile, $q, NgTableParams) {
            elm = angular.element(
                    '<div>' +
                    '<table ng-table-dynamic="tableParams with cols">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td ng-repeat="col in $columns">{{user[col.field]}}</td>' +
                    '</tr>' +
                    '</table>' +
                    '</div>');

            function getCustomClass(context){
                if (context.$column.title().indexOf('Money') !== -1){
                    return 'moneyHeaderClass';
                } else{
                    return 'customClass';
                }
            }

            function money(/*$column*/) {
                var def = $q.defer();
                def.resolve([{
                    'id': 10,
                    'title': '10'
                }]);
                return def;
            }

            scope.tableParams = new NgTableParams({}, {});
            scope.cols = [
                {
                    'class': getCustomClass,
                    field: 'name',
                    filter: { name: 'text' },
                    headerTitle: 'Sort by Name',
                    sortable: 'name',
                    show: true,
                    title: 'Name of person'
                },
                {
                    'class': getCustomClass,
                    field: 'age',
                    headerTitle: 'Sort by Age',
                    sortable: 'age',
                    show: true,
                    title: 'Age'
                },
                {
                    'class': getCustomClass,
                    field: 'money',
                    filter: { action: 'select' },
                    headerTitle: 'Sort by Money',
                    filterData: money,
                    show: true,
                    title: 'Money'
                }
            ];

            $compile(elm)(scope);
            scope.$digest();
        }));

        it('should create table header', function() {
            var thead = elm.find('thead');
            expect(thead.length).toBe(1);

            var rows = thead.find('tr');
            expect(rows.length).toBe(2);

            var titles = angular.element(rows[0]).find('th');

            expect(titles.length).toBe(3);
            expect(angular.element(titles[0]).text().trim()).toBe('Name of person');
            expect(angular.element(titles[1]).text().trim()).toBe('Age');
            expect(angular.element(titles[2]).text().trim()).toBe('Money');

            expect(angular.element(rows[1]).hasClass('ng-table-filters')).toBeTruthy();
            var filters = angular.element(rows[1]).find('th');
            expect(filters.length).toBe(3);
            expect(angular.element(filters[0]).hasClass('filter')).toBeTruthy();
            expect(angular.element(filters[1]).hasClass('filter')).toBeTruthy();
            expect(angular.element(filters[2]).hasClass('filter')).toBeTruthy();
        });

        it('should create table header classes', inject(function($compile, $rootScope) {

            var thead = elm.find('thead');
            var rows = thead.find('tr');
            var titles = angular.element(rows[0]).find('th');

            expect(angular.element(titles[0]).hasClass('header')).toBeTruthy();
            expect(angular.element(titles[1]).hasClass('header')).toBeTruthy();
            expect(angular.element(titles[2]).hasClass('header')).toBeTruthy();

            expect(angular.element(titles[0]).hasClass('sortable')).toBeTruthy();
            expect(angular.element(titles[1]).hasClass('sortable')).toBeTruthy();
            expect(angular.element(titles[2]).hasClass('sortable')).toBeFalsy();

            expect(angular.element(titles[0]).hasClass('customClass')).toBeTruthy();
            expect(angular.element(titles[1]).hasClass('customClass')).toBeTruthy();
            expect(angular.element(titles[2]).hasClass('moneyHeaderClass')).toBeTruthy();
        }));

        it('should create table header titles', function() {

            var thead = elm.find('thead');
            var rows = thead.find('tr');
            var titles = angular.element(rows[0]).find('th');

            expect(angular.element(titles[0]).attr('title').trim()).toBe('Sort by Name');
            expect(angular.element(titles[1]).attr('title').trim()).toBe('Sort by Age');
            expect(angular.element(titles[2]).attr('title').trim()).toBe('Sort by Money');
        });

        it('should show data-title-text', inject(function(NgTableParams) {
            var tbody = elm.find('tbody');

            scope.tableParams = new NgTableParams({
                page: 1, // show first page
                count: 10 // count per page
            }, {
                dataset: dataset
            });
            scope.$digest();

            var filterRow = angular.element(elm.find('thead').find('tr')[1]);
            var filterCells = filterRow.find('th');
            expect(angular.element(filterCells[0]).attr('data-title-text').trim()).toBe('Name of person');
            expect(angular.element(filterCells[1]).attr('data-title-text').trim()).toBe('Age');
            expect(angular.element(filterCells[2]).attr('data-title-text').trim()).toBe('Money');

            var dataRows = elm.find('tbody').find('tr');
            var dataCells = angular.element(dataRows[0]).find('td');
            expect(angular.element(dataCells[0]).attr('data-title-text').trim()).toBe('Name of person');
            expect(angular.element(dataCells[1]).attr('data-title-text').trim()).toBe('Age');
            expect(angular.element(dataCells[2]).attr('data-title-text').trim()).toBe('Money');
        }));

        it('should show/hide columns', inject(function(NgTableParams) {
            var tbody = elm.find('tbody');

            scope.tableParams = new NgTableParams({
                page: 1, // show first page
                count: 10 // count per page
            }, {
                dataset: dataset
            });
            scope.$digest();

            var headerRow = angular.element(elm.find('thead').find('tr')[0]);
            expect(headerRow.find('th').length).toBe(3);

            var filterRow = angular.element(elm.find('thead').find('tr')[1]);
            expect(filterRow.find('th').length).toBe(3);

            var dataRow = angular.element(elm.find('tbody').find('tr')[0]);
            expect(dataRow.find('td').length).toBe(3);

            scope.cols[0].show = false;
            scope.$digest();
            expect(headerRow.find('th').length).toBe(2);
            expect(filterRow.find('th').length).toBe(2);
            expect(dataRow.find('td').length).toBe(2);
            expect(angular.element(headerRow.find('th')[0]).text().trim()).toBe('Age');
            expect(angular.element(headerRow.find('th')[1]).text().trim()).toBe('Money');
            expect(angular.element(filterRow.find('th')[0]).find('input').length).toBe(0);
            expect(angular.element(filterRow.find('th')[1]).find('select').length).toBe(1);
        }));
    });
   describe('changing column list', function(){
        var elm;
        beforeEach(inject(function($compile, $q, NgTableParams) {
            elm = angular.element(
                    '<div>' +
                    '<table ng-table-dynamic="tableParams with cols">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td ng-repeat="col in $columns">{{user[col.field]}}</td>' +
                    '</tr>' +
                    '</table>' +
                    '</div>');

            function getCustomClass(parmasScope){
                if (parmasScope.$column.title().indexOf('Money') !== -1){
                    return 'moneyHeaderClass';
                } else{
                    return 'customClass';
                }
            }

            function money(/*$column*/) {
                var def = $q.defer();
                def.resolve([{
                    'id': 10,
                    'title': '10'
                }]);
                return def;
            }

            scope.tableParams = new NgTableParams({}, {});
            scope.cols = [
                {
                    'class': getCustomClass,
                    field: 'name',
                    filter: { name: 'text' },
                    headerTitle: 'Sort by Name',
                    sortable: 'name',
                    show: true,
                    title: 'Name of person'
                },
                {
                    'class': getCustomClass,
                    field: 'age',
                    headerTitle: 'Sort by Age',
                    sortable: 'age',
                    show: true,
                    title: 'Age'
                }
            ];

            $compile(elm)(scope);
            scope.$digest();
        }));

        it('adding new column should update table header', function() {
            var newCol =  {
                 'class': 'moneyadd',
                 field: 'money',
                 filter: { action: 'select' },
                 headerTitle: 'Sort by Money',
                 show: true,
                 title: 'Money'
            };
            scope.cols.push(newCol);
            scope.$digest();
            var thead = elm.find('thead');
            expect(thead.length).toBe(1);

            var rows = thead.find('tr');
            expect(rows.length).toBe(2);

            var titles = angular.element(rows[0]).find('th');

            expect(titles.length).toBe(3);
            expect(angular.element(titles[0]).text().trim()).toBe('Name of person');
            expect(angular.element(titles[1]).text().trim()).toBe('Age');
            expect(angular.element(titles[2]).text().trim()).toBe('Money');

            var filterRow = angular.element(rows[1]);
            expect(filterRow.hasClass('ng-table-filters')).toBeTruthy();
            expect(filterRow.hasClass("ng-hide")).toBe(false);

            var filters = filterRow.find('th');
            expect(filters.length).toBe(3);
            expect(angular.element(filters[0]).hasClass('filter')).toBeTruthy();
            expect(angular.element(filters[1]).hasClass('filter')).toBeTruthy();
            expect(angular.element(filters[2]).hasClass('filter')).toBeTruthy();

        });

        it('removing new column should update table header', function() {
            scope.cols.splice(0, 1);
            scope.$digest();
            var thead = elm.find('thead');
            expect(thead.length).toBe(1);

            var rows = thead.find('tr');
            var titles = angular.element(rows[0]).find('th');
            expect(titles.length).toBe(1);
            expect(angular.element(titles[0]).text().trim()).toBe('Age');

            var filterRow = angular.element(rows[1]);
            expect(filterRow.hasClass("ng-hide")).toBe(true);
        });

        it('setting columns to null should remove all table columns from header', function() {
            scope.cols = null;
            scope.$digest();
            var thead = elm.find('thead');
            expect(thead.length).toBe(1);

            var rows = thead.find('tr');
            var titles = angular.element(rows[0]).find('th');
            expect(titles.length).toBe(0);

            var filterRow = angular.element(rows[1]);
            expect(filterRow.hasClass("ng-hide")).toBe(true);

            expect(filterRow.find('th').length).toBe(0);
        });

    });
    describe('title-alt', function() {

        var elm;
        beforeEach(inject(function($compile, NgTableParams) {
            elm = angular.element(
                    '<table ng-table-dynamic="tableParams with cols">' +
                    '<tr ng-repeat="user in $data">' +
                        '<td ng-repeat="col in $columns">{{user[col.field]}}</td>' +
                    '</tr>' +
                    '</table>');

            scope.cols = [
                { field: 'name',  title: 'Name of person', titleAlt: 'Name' },
                { field: 'age',   title: 'Age',            titleAlt: 'Age' },
                { field: 'money', title: 'Money',          titleAlt: '£' }
            ];
            scope.tableParams = new NgTableParams({
                page: 1, // show first page
                count: 10 // count per page
            }, {
                dataset: dataset
            });

            $compile(elm)(scope);
            scope.$digest();
        }));

        it('should show as data-title-text', inject(function($compile) {
            var filterRow = angular.element(elm.find('thead').find('tr')[1]);
            var filterCells = filterRow.find('th');

            expect(angular.element(filterCells[0]).attr('data-title-text').trim()).toBe('Name');
            expect(angular.element(filterCells[1]).attr('data-title-text').trim()).toBe('Age');
            expect(angular.element(filterCells[2]).attr('data-title-text').trim()).toBe('£');

            var dataRows = elm.find('tbody').find('tr');
            var dataCells = angular.element(dataRows[0]).find('td');
            expect(angular.element(dataCells[0]).attr('data-title-text').trim()).toBe('Name');
            expect(angular.element(dataCells[1]).attr('data-title-text').trim()).toBe('Age');
            expect(angular.element(dataCells[2]).attr('data-title-text').trim()).toBe('£');
        }));
    });

    describe('filters', function(){

        var elm;
        beforeEach(inject(function($compile, NgTableParams) {
            elm = angular.element(
                    '<table ng-table-dynamic="tableParams with cols">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td ng-repeat="col in $columns">{{user[col.field]}}</td>' +
                    '</tr>' +
                    '</table>');
        }));

        describe('filter specified as alias', function(){

            beforeEach(inject(function($compile, NgTableParams) {
                scope.cols = [
                    { field: 'name', filter: {username: 'text'} }
                ];
                scope.tableParams = new NgTableParams({}, {});
                $compile(elm)(scope);
                scope.$digest();
            }));

            it('should render named filter template', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(inputs.length).toBe(1);
                expect(inputs.eq(0).attr('type')).toBe('text');
                expect(inputs.eq(0).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(0).attr('name')).toBe('username');
            });

            it('should render named filter template - select template', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(inputs.length).toBe(1);
                expect(inputs.eq(0).attr('type')).toBe('text');
                expect(inputs.eq(0).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(0).attr('name')).toBe('username');
            });

            it('should databind ngTableParams.filter to filter input', function () {
                scope.tableParams.filter()['username'] = 'my name is...';
                scope.$digest();

                var input = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(input.val()).toBe('my name is...');
            });
        });

        describe('select filter', function(){

            beforeEach(inject(function ($compile, $q, NgTableParams) {
                scope.cols = [{
                    field: 'name',
                    filter: {username: 'select'},
                    filterData: getNamesAsDefer
                }, {
                    field: 'names2',
                    filter: {username2: 'select'},
                    filterData: getNamesAsPromise
                }, {
                    field: 'names3',
                    filter: {username3: 'select'},
                    filterData: getNamesAsArray
                }];
                scope.tableParams = new NgTableParams({}, {});
                $compile(elm)(scope);
                scope.$digest();

                function getNamesAsDefer(/*$column*/) {
                    var def = $q.defer();
                    def.resolve([{
                        'id': 10,
                        'title': 'Christian'
                    }, {
                        'id': 11,
                        'title': 'Simon'
                    }]);
                    return def;
                }

                function getNamesAsPromise(/*$column*/) {
                    return $q.when([{
                        'id': 20,
                        'title': 'Christian'
                    }, {
                        'id': 21,
                        'title': 'Simon'
                    }]);
                }

                function getNamesAsArray(/*$column*/) {
                    return [{
                        'id': 20,
                        'title': 'Christian'
                    }, {
                        'id': 21,
                        'title': 'Simon'
                    }];
                }

            }));

            it('should render select lists', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('select');
                expect(inputs.length).toBe(3);
                expect(inputs.eq(0).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(0).attr('name')).toBe('username');
                expect(inputs.eq(1).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(1).attr('name')).toBe('username2');
                expect(inputs.eq(2).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(2).attr('name')).toBe('username3');
            });

            it('should render list data return as a deferred', function() {
                /* WARNING: support for returning a $defer is depreciated */

                var inputs = elm.find('thead').find('tr').eq(1).find('th').eq(0).find('select');
                expect(inputs[0].options.length).toBeGreaterThan(0);
                var $column = inputs.eq(0).scope().$column;
                var plucker = _.partialRight(_.pick, ['id', 'title']);
                var actual = _.map($column.data, plucker);
                expect(actual).toEqual([{
                    'id': '',
                    'title': ''
                },{
                    'id': 10,
                    'title': 'Christian'
                }, {
                    'id': 11,
                    'title': 'Simon'
                }]);
            });

            it('should render select list return as a promise', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').eq(1).find('select');
                expect(inputs[0].options.length).toBeGreaterThan(0);
                var $column = inputs.eq(0).scope().$column;
                var plucker = _.partialRight(_.pick, ['id', 'title']);
                var actual = _.map($column.data, plucker);
                expect(actual).toEqual([{
                    'id': '',
                    'title': ''
                },{
                    'id': 20,
                    'title': 'Christian'
                }, {
                    'id': 21,
                    'title': 'Simon'
                }]);
            });

            it('should render select list return as an array', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').eq(2).find('select');
                expect(inputs[0].options.length).toBeGreaterThan(0);
                var $column = inputs.eq(0).scope().$column;
                var plucker = _.partialRight(_.pick, ['id', 'title']);
                var actual = _.map($column.data, plucker);
                expect(actual).toEqual([{
                    'id': '',
                    'title': ''
                },{
                    'id': 20,
                    'title': 'Christian'
                }, {
                    'id': 21,
                    'title': 'Simon'
                }]);
            });
        });

        describe('multiple filter inputs', function(){

            beforeEach(inject(function($compile, NgTableParams) {
                scope.cols = [
                    { field: 'name', filter: {name: 'text', age: 'text'} }
                ];
                scope.tableParams = new NgTableParams({}, {});
                $compile(elm)(scope);
                scope.$digest();
            }));

            it('should render filter template for each key/value pair ordered by key', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(inputs.length).toBe(2);
                expect(inputs.eq(0).attr('type')).toBe('text');
                expect(inputs.eq(0).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(1).attr('type')).toBe('text');
                expect(inputs.eq(1).attr('ng-model')).not.toBeUndefined();
            });

            it('should databind ngTableParams.filter to filter inputs', function () {
                scope.tableParams.filter()['name'] = 'my name is...';
                scope.tableParams.filter()['age'] = '10';
                scope.$digest();

                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(inputs.eq(0).val()).toBe('my name is...');
                expect(inputs.eq(1).val()).toBe('10');
            });
        });

        describe('dynamic filter', function(){

            var ageFilter;
            beforeEach(inject(function($compile, NgTableParams) {

                ageFilter = { age: 'text'};
                function getFilter(paramsScope){
                    if (paramsScope.$column.title() === 'Name of user') {
                        return {username: 'text'};
                    } else if (paramsScope.$column.title() === 'Age') {
                        return ageFilter;
                    }
                }

                scope.cols = [
                    { field: 'name', title: 'Name of user', filter: getFilter },
                    { field: 'age', title: 'Age', filter: getFilter }
                ];
                scope.tableParams = new NgTableParams({}, {});

                $compile(elm)(scope);
                scope.$digest();

            }));

            it('should render named filter template', function() {
                var usernameInput = elm.find('thead').find('tr').eq(1).find('th').eq(0).find('input');
                expect(usernameInput.attr('type')).toBe('text');
                expect(usernameInput.attr('name')).toBe('username');

                var ageInput = elm.find('thead').find('tr').eq(1).find('th').eq(1).find('input');
                expect(ageInput.attr('type')).toBe('text');
                expect(ageInput.attr('name')).toBe('age');
            });

            it('should databind ngTableParams.filter to filter input', function () {
                scope.tableParams.filter()['username'] = 'my name is...';
                scope.tableParams.filter()['age'] = '10';
                scope.$digest();

                var usernameInput = elm.find('thead').find('tr').eq(1).find('th').eq(0).find('input');
                expect(usernameInput.val()).toBe('my name is...');
                var ageInput = elm.find('thead').find('tr').eq(1).find('th').eq(1).find('input');
                expect(ageInput.val()).toBe('10');
            });

            it('should render new template as filter changes', inject(function($compile) {

                var scriptTemplate = angular.element(
                    '<script type="text/ng-template" id="ng-table/filters/number.html"><input type="number" name="{{name}}"/></script>');
                $compile(scriptTemplate)(scope);

                ageFilter.age = 'number';
                scope.$digest();

                var ageInput = elm.find('thead').find('tr').eq(1).find('th').eq(1).find('input');
                expect(ageInput.attr('type')).toBe('number');
                expect(ageInput.attr('name')).toBe('age');
            }));
        });
    });

    describe('reorder columns', function() {
        var elm;
        var getTitles = function () {
            var thead = elm.find('thead');
            var rows = thead.find('tr');
            var titles = angular.element(rows[0]).find('th');

            return angular.element(titles).text().trim().split(/\s+/g)
        };

        beforeEach(inject(function ($compile, $q, NgTableParams) {
            elm = angular.element(
                '<div>' +
                '<table ng-table-dynamic="tableParams with cols">' +
                '<tr ng-repeat="user in $data">' +
                "<td ng-repeat=\"col in $columns\">{{user[col.field]}}</td>" +
                '</tr>' +
                '</table>' +
                '</div>');

            scope.tableParams = new NgTableParams({}, {});
            scope.cols = [
                {
                    field: 'name',
                    title: 'Name'
                },
                {
                    field: 'age',
                    title: 'Age'
                },
                {
                    field: 'money',
                    title: 'Money'
                }
            ];

            $compile(elm)(scope);
            scope.$digest();
        }));

        it('"in place" switch of columns within array should reorder html table columns', function () {
            expect(getTitles()).toEqual([ 'Name', 'Age', 'Money' ]);

            var colToSwap = scope.cols[2];
            scope.cols[2] = scope.cols[1];
            scope.cols[1] = colToSwap;
            scope.$digest();

            expect(getTitles()).toEqual([ 'Name', 'Money', 'Age' ]);
        });

        it('"in place" reverse of column array should reorder html table columns', function () {
            expect(getTitles()).toEqual([ 'Name', 'Age', 'Money' ]);

            scope.cols.reverse();
            scope.$digest();

            expect(getTitles()).toEqual([ 'Money', 'Age', 'Name' ]);
        });

        it('html table columns should reflect order of columns in replacement array', function () {
            expect(getTitles()).toEqual([ 'Name', 'Age', 'Money' ]);

            var newArray = scope.cols.map(angular.identity);
            newArray.reverse();
            scope.cols = newArray;
            scope.$digest();

            expect(getTitles()).toEqual([ 'Money', 'Age', 'Name' ]);
        });
    });
});
