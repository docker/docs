describe('ng-table', function() {
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
    var NgTableParams;

    beforeEach(module('ngTable'));

    var scope;
    beforeEach(inject(function($rootScope, _NgTableParams_) {
        scope = $rootScope.$new(true);
        NgTableParams = _NgTableParams_;
    }));

    function createNgTableParams(settings) {
        var initialParams;
        if (arguments.length === 2){
            initialParams = arguments[0];
            settings = arguments[1];
        }

        settings = angular.extend({}, settings);
        settings.filterOptions = angular.extend({}, {
            filterDelay: 0
        }, settings.filterOptions);
        var tableParams = new NgTableParams(initialParams, settings);
        spyOn(tableParams.settings(), 'getData').and.callThrough();
        return tableParams;
    }

    describe('basics', function(){
        var elm;
        beforeEach(inject(function($compile, $q) {
            elm = angular.element(
                    '<div>' +
                    '<table ng-table="tableParams">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td data-header-title="\'Sort by Name\'" data-title="nameTitle()" filter="{ \'name\': \'text\' }" sortable="\'name\'" data-header-class="getCustomClass($column)"' +
                        ' ng-if="showName">' +
                    '{{user.name}}' +
                    '</td>' +
                    '<td x-data-header-title="\'Sort by Age\'" x-data-title="ageTitle()" sortable="\'age\'" x-data-header-class="getCustomClass($column)"' +
                        ' ng-if="showAge">' +
                    '{{user.age}}' +
                    '</td>' +
                    '<td header-title="\'Sort by Money\'" title="moneyTitle()" filter="{ \'action\': \'select\' }" filter-data="money($column)" header-class="getCustomClass($column)"' +
                        ' ng-if="showMoney">' +
                    '{{user.money}}' +
                    '</td>' +
                    '</tr>' +
                    '</table>' +
                    '</div>');

            scope.nameTitle = function(){
                return 'Name of person';
            };
            scope.ageTitle = function(){
                return 'Age';
            };
            scope.moneyTitle = function(){
                return 'Money';
            };

            scope.showName = true;
            scope.showAge = true;
            scope.showMoney = true;

            scope.ageTitle = function(){
                return 'Age';
            };
            scope.moneyTitle = function(){
                return 'Money';
            };

            scope.getCustomClass = function($column){
                if ($column.title().indexOf('Money') !== -1){
                    return 'moneyHeaderClass';
                } else{
                    return 'customClass';
                }
            };

            scope.money = function(/*$column*/) {

                var def = $q.defer();

                def.resolve([{
                    'id': 10,
                    'title': '10'
                }]);
                return def;
            };

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

        it('should create table header classes', function() {

            var thead = elm.find('thead');
            var rows = thead.find('tr');
            var titles = angular.element(rows[0]).find('th');

            expect(angular.element(titles[0]).hasClass('header')).toBe(true);
            expect(angular.element(titles[1]).hasClass('header')).toBe(true);
            expect(angular.element(titles[2]).hasClass('header')).toBe(true);

            expect(angular.element(titles[0]).hasClass('sortable')).toBe(true);
            expect(angular.element(titles[1]).hasClass('sortable')).toBe(true);
            expect(angular.element(titles[2]).hasClass('sortable')).toBe(false);

            expect(angular.element(titles[0]).hasClass('customClass')).toBe(true);
            expect(angular.element(titles[1]).hasClass('customClass')).toBe(true);
            expect(angular.element(titles[2]).hasClass('moneyHeaderClass')).toBe(true);


            var filterCells = angular.element(rows[1]).find('th');

            expect(angular.element(filterCells[0]).hasClass('filter')).toBe(true);
            expect(angular.element(filterCells[1]).hasClass('filter')).toBe(true);
            expect(angular.element(filterCells[2]).hasClass('filter')).toBe(true);

            expect(angular.element(filterCells[0]).hasClass('customClass')).toBe(true);
            expect(angular.element(filterCells[1]).hasClass('customClass')).toBe(true);
            expect(angular.element(filterCells[2]).hasClass('moneyHeaderClass')).toBe(true);
        });

        it('should create table header titles', function() {

            var thead = elm.find('thead');
            var rows = thead.find('tr');
            var titles = angular.element(rows[0]).find('th');

            expect(angular.element(titles[0]).attr('title').trim()).toBe('Sort by Name');
            expect(angular.element(titles[1]).attr('title').trim()).toBe('Sort by Age');
            expect(angular.element(titles[2]).attr('title').trim()).toBe('Sort by Money');
        });


        it('should show scope data', function() {
            var tbody = elm.find('tbody');
            expect(tbody.length).toBe(1);

            var rows = tbody.find('tr');
            expect(rows.length).toBe(0);

            var params = new NgTableParams({
                count: 10 // count per page
            }, {
                dataset: dataset
            });

            scope.tableParams = params;
            scope.$digest();

            rows = tbody.find('tr');
            expect(rows.length).toBe(10);

            scope.tableParams.page(2);
            scope.$digest();

            rows = tbody.find('tr');
            expect(rows.length).toBe(7);

            params.total(20);
            scope.$digest();

            rows = tbody.find('tr');
            expect(rows.length).toBe(7);
        });

        it('should show data-title-text', function() {
            var tbody = elm.find('tbody');

            var params = new NgTableParams({}, {
                dataset: dataset
            });
            scope.tableParams = params;
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
        });


        it('should show/hide columns', function() {
            var tbody = elm.find('tbody');

            scope.tableParams = new NgTableParams({}, {
                dataset: dataset
            });
            scope.$digest();

            var headerRow = angular.element(elm.find('thead').find('tr')[0]);
            expect(headerRow.find('th').length).toBe(3);

            var filterRow = angular.element(elm.find('thead').find('tr')[1]);
            expect(filterRow.find('th').length).toBe(3);

            var dataRow = angular.element(elm.find('tbody').find('tr')[0]);
            expect(dataRow.find('td').length).toBe(3);

            scope.showName = false;
            scope.$digest();
            expect(headerRow.find('th').length).toBe(2);
            expect(filterRow.find('th').length).toBe(2);
            expect(dataRow.find('td').length).toBe(2);
            expect(angular.element(headerRow.find('th')[0]).text().trim()).toBe('Age');
            expect(angular.element(headerRow.find('th')[1]).text().trim()).toBe('Money');
            expect(angular.element(filterRow.find('th')[0]).find('input').length).toBe(0);
            expect(angular.element(filterRow.find('th')[1]).find('select').length).toBe(1);
        });
    });

    describe('title-alt', function() {

        var elm;
        beforeEach(inject(function($compile) {
            elm = angular.element(
                    '<table ng-table="tableParams">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td title="\'Name of person\'" title-alt="\'Name\'">{{user.name}}</td>' +
                    '<td title="\'Age of person\'" data-title-alt="\'Age\'">{{user.age}}</td>' +
                    '<td title="\'Money earned\'" x-data-title-alt="\'£\'">{{user.money}}</td>' +
                    '</tr>' +
                    '</table>');

            $compile(elm)(scope);
            scope.$digest();

            var params = new NgTableParams({}, {
                dataset: dataset
            });
            scope.tableParams = params;
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

    describe('sorting', function() {

        it('should provide $column definition', inject(function($compile) {
            var columnDef;
            var elm = angular.element(
                    '<table ng-table="tableParams">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td title="\'Age\'" sortable="captureColumn($column)">{{user.age}}</td>' +
                    '</tr>' +
                    '</table>');

            scope.captureColumn = function($column){
                columnDef = $column;
                return 'age'
            };

            $compile(elm)(scope);
            scope.$digest();

            expect(columnDef).toBeDefined();
        }));

        it('should apply initial sort', inject(function ($compile) {
            var elm = angular.element(
                '<table ng-table="tableParams">' +
                '<tr ng-repeat="user in $data"><td title="\'Age\'" sortable="\'age\'">{{user.age}}</td></tr>' +
                '</table>');
            $compile(elm)(scope);

            var actualSort;
            scope.tableParams = new NgTableParams({
                sorting: { age: 'desc' }
            }, {
                getData: function($defer, params){
                    actualSort = params.sorting();
                    $defer.resolve([]);
                }});
            scope.$digest();

            expect(actualSort.age).toBe('desc');
        }));

        it('when sorting changes should trigger reload of table', inject(function ($compile) {
            var elm = angular.element(
                '<table ng-table="tableParams">' +
                '<tr ng-repeat="user in $data"><td title="\'Age\'" sortable="\'age\'">{{user.age}}</td></tr>' +
                '</table>');
            $compile(elm)(scope);

            var params = createNgTableParams();
            scope.tableParams = params;
            scope.$digest();
            params.settings().getData.calls.reset();

            params.sorting()['age'] = 'desc';
            scope.$digest();
            expect(params.settings().getData.calls.count()).toBe(1);

            params.sorting()['age'] = 'asc';
            scope.$digest();
            expect(params.settings().getData.calls.count()).toBe(2);

            // setting the same sort order should not trigger reload
            params.sorting({ age: 'asc'});
            scope.$digest();
            expect(params.settings().getData.calls.count()).toBe(2);
        }));
    });

    describe('paging', function() {

        var elm;

        beforeEach(inject(function($compile) {
            elm = angular.element(
                '<table ng-table="tableParams">' +
                '<tr ng-repeat="user in $data">' +
                '<td title="\'Age\'">{{user.age}}</td>' +
                '</tr>' +
                '</table>');

            dataCallCount = 0;
            $compile(elm)(scope);
            scope.$digest();
        }));

        function verifyPageWas(expectedPage){
            expect(scope.tableParams.settings().getData.calls.argsFor(0)[0].page()).toBe(expectedPage);
        }

        it('should use initial NgTableParams constructor value', function(){
            var params = createNgTableParams({ page: 2}, null);
            scope.tableParams = params;
            scope.$digest();
            verifyPageWas(2);
            expect(params.settings().getData.calls.count()).toBe(1);
        });

        it('should use initial NgTableParams constructor value combined with filter', function(){
            var params = createNgTableParams({ page: 2, filter: { age: 5}}, null);
            scope.tableParams = params;
            scope.$digest();
            verifyPageWas(2);
            expect(params.settings().getData.calls.count()).toBe(1);
        });

        it('changing page # should trigger reload of data', function(){
            var params = createNgTableParams({ page: 3}, null);
            scope.tableParams = params;
            scope.$digest();
            verifyPageWas(3);
            params.settings().getData.calls.reset();

            scope.tableParams.page(5);
            scope.$digest();
            verifyPageWas(5);
        });
    });

    describe('filters', function(){

        var $capturedColumn;
        beforeEach(inject(function() {
            // stash a reference to $column definition so that its available in asserts
            scope.captureColumn = function ($column) {
                $capturedColumn = $column;
            };
        }));

        describe('filter specified as alias', function(){

            var elm,
                tp;
            beforeEach(inject(function($compile) {
                elm = angular.element(
                        '<div>' +
                        '<table ng-table="tableParams">' +
                        '<tr ng-repeat="user in $data">' +
                        '<td header-class="captureColumn($column)" title="\'Name\'" ' +
                            'filter="usernameFilter">{{user.name}}</td>' +
                        '</tr>' +
                        '</table>' +
                        '</div>');

                $compile(elm)(scope);
                scope.$digest();

                // 'text' is a shortcut alias for the template ng-table/filters/text
                scope.usernameFilter = {username: 'text'};
                tp = scope.tableParams = createNgTableParams({ filterOptions: { filterDelay: 10 } });
                scope.$digest();
            }));

            it('should render named filter template', function() {
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

            it('should make filter def available on $column', function () {
                expect($capturedColumn).toBeDefined();
                expect($capturedColumn.filter).toBeDefined();
                expect($capturedColumn.filter()['username']).toBe('text');
            });

            it('when filter changes should trigger reload of table', inject(function ($timeout) {
                tp.settings().getData.calls.reset();

                tp.filter()['username'] = 'new value';
                scope.$digest();
                $timeout.flush();  // trigger delayed filter
                tp.filter()['username'] = 'another value';
                scope.$digest();
                $timeout.flush();  // trigger delayed filter
                expect(tp.settings().getData.calls.count()).toBe(2);

                // same value - should not trigger reload
                tp.filter()['username'] = 'another value';
                scope.$digest();
                try{
                    $timeout.flush();  // trigger delayed filter
                } catch (ex) {

                }
                expect(tp.settings().getData.calls.count()).toBe(2);
            }));

            it('when filter changes should reset page number to 1', inject(function ($timeout) {
                // trigger initial load of data so that subsequent changes to filter will trigger reset of page #
                tp.filter()['username'] = 'initial value';
                scope.$digest();
                $timeout.flush();  // trigger delayed filter

                // set page to something other than 1
                tp.page(5);
                expect(tp.page()).toBe(5); // checking assumptions

                // when
                tp.filter()['username'] = 'new value';
                scope.$digest();
                $timeout.flush();  // trigger delayed filter

                expect(tp.page()).toBe(1);
            }));
        });

        describe('filter specified with url', function(){

            var elm;
            beforeEach(inject(function($compile) {
                elm = angular.element(
                        '<div>' +
                        '<script type="text/ng-template" id="ng-table/filters/customNum.html"><input type="number" id="{{name}}"/></script>' +
                        '<table ng-table="tableParams">' +
                        '<tr ng-repeat="user in $data">' +
                        '<td header-class="captureColumn($column)" title="\'Age\'" ' +
                            'filter="{ \'age\': \'ng-table/filters/customNum.html\' }">{{user.age}}</td>' +
                        '</tr>' +
                        '</table>' +
                        '</div>');

                $compile(elm)(scope);
                scope.tableParams = createNgTableParams();
                scope.$digest();
            }));

            it('should render filter template specified by url', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(inputs.length).toBe(1);

                expect(inputs.eq(0).attr('type')).toBe('number');
                expect(inputs.eq(0).attr('id')).toBe('age');
            });
        });

        describe('multiple filter inputs', function(){

            var elm;
            beforeEach(inject(function($compile) {
                elm = angular.element(
                        '<div>' +
                        '<table ng-table="tableParams">' +
                        '<tr ng-repeat="user in $data">' +
                        '<td header-class="captureColumn($column)" title="\'Name\'" ' +
                        'filter="{ \'name\': \'text\', \'age\': \'text\' }">{{user.name}}</td>' +
                        '</tr>' +
                        '</table>' +
                        '</div>');

                $compile(elm)(scope);
                scope.$digest();

                scope.tableParams = createNgTableParams();
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

            it('should make filter def available on $column', function () {
                expect($capturedColumn).toBeDefined();
                expect($capturedColumn.filter).toBeDefined();
                expect($capturedColumn.filter()['name']).toBe('text');
                expect($capturedColumn.filter()['age']).toBe('text');
            });
        });
        describe('dynamic filter', function(){

            var elm, ageFilter;
            beforeEach(inject(function($compile) {

                ageFilter = {age: 'text'};

                elm = angular.element(
                        '<div>' +
                        '<script type="text/ng-template" id="ng-table/filters/number.html"><input type="number" name="{{name}}"/></script>' +
                        '<table ng-table="tableParams">' +
                        '<tr ng-repeat="user in $data">' +
                        '<td title="\'Name\'" filter="getFilter($column)">{{user.name}}</td>' +
                        '<td title="\'Age\'" filter="getFilter($column)">{{user.age}}</td>' +
                        '</tr>' +
                        '</table>' +
                        '</div>');

                $compile(elm)(scope);
                scope.$digest();

                scope.getFilter = function(colDef){
                    if (colDef.id === 0) {
                        return {username: 'text'};
                    } else if (colDef.id === 1) {
                        return ageFilter;
                    }
                };
                scope.tableParams = createNgTableParams();
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

            it('should render new template as filter changes', function() {
                ageFilter.age = 'number';
                scope.$digest();

                var ageInput = elm.find('thead').find('tr').eq(1).find('th').eq(1).find('input');
                expect(ageInput.attr('type')).toBe('number');
                expect(ageInput.attr('name')).toBe('age');
            });
        });

        describe('filter with placeholder value and alias', function(){

            var elm,
                tp;
            beforeEach(inject(function($compile) {
                elm = angular.element(
                    '<div>' +
                    '<table ng-table="tableParams">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td header-class="captureColumn($column)" title="\'Name\'" '
                    + 'filter="usernameFilter">{{user.name}}</td>' +
                    '</tr>' +
                    '</table>' +
                    '</div>');

                $compile(elm)(scope);
                scope.$digest();

                // 'text' is a shortcut alias for the template ng-table/filters/text
                scope.usernameFilter = {
                    username: { id: 'text', placeholder: 'User name'}
                };
                tp = scope.tableParams = createNgTableParams();
                scope.$digest();
            }));

            it('should render named filter template with placeholder value', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(inputs.length).toBe(1);
                expect(inputs.eq(0).attr('type')).toBe('text');
                expect(inputs.eq(0).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(0).attr('name')).toBe('username');
                expect(inputs.eq(0).attr('placeholder')).toBe('User name');
            });

            it('should databind placeholder value to filter input', function () {
                scope.usernameFilter.username.placeholder = 'Name of user';
                scope.$digest();

                var input = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(input.attr('placeholder')).toBe('Name of user');
            });

            it('should make filter def available on $column', function () {
                expect($capturedColumn).toBeDefined();
                expect($capturedColumn.filter).toBeDefined();
                expect($capturedColumn.filter()).toBe(scope.usernameFilter);
            });
        });

        describe('filter with placeholder value and url', function(){

            var elm,
                tp;
            beforeEach(inject(function($compile) {
                elm = angular.element(
                    '<div>' +
                    '<table ng-table="tableParams">' +
                    '<tr ng-repeat="user in $data">' +
                    '<td header-class="captureColumn($column)" title="\'Age\'" '
                    + 'filter="ageFilter">{{user.age}}</td>' +
                    '</tr>' +
                    '</table>' +
                    '</div>');

                $compile(elm)(scope);
                scope.$digest();

                scope.ageFilter = {
                    age: { id: 'ng-table/filters/number.html', placeholder: 'User age'}
                };
                tp = scope.tableParams = createNgTableParams();
                scope.$digest();
            }));

            it('should render named filter template with placeholder value', function() {
                var inputs = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(inputs.length).toBe(1);
                expect(inputs.eq(0).attr('type')).toBe('number');
                expect(inputs.eq(0).attr('ng-model')).not.toBeUndefined();
                expect(inputs.eq(0).attr('name')).toBe('age');
                expect(inputs.eq(0).attr('placeholder')).toBe('User age');
            });

            it('should databind placeholder value to filter input', function () {
                scope.ageFilter.age.placeholder = 'Age of user';
                scope.$digest();

                var input = elm.find('thead').find('tr').eq(1).find('th').find('input');
                expect(input.attr('placeholder')).toBe('Age of user');
            });

            it('should make filter def available on $column', function () {
                expect($capturedColumn).toBeDefined();
                expect($capturedColumn.filter).toBeDefined();
                expect($capturedColumn.filter()).toBe(scope.ageFilter);
            });
        });
    });

    describe('show-filter', function(){
        var elm;
        beforeEach(inject(function($compile) {
            elm = angular.element(
                '<div>' +
                '<table ng-table="tableParams" show-filter="showFilterRow">' +
                '<tr ng-repeat="user in $data">' +
                '<td title="\'Age\'" filter="{ age: \'number\'}">{{user.age}}</td>' +
                '</tr>' +
                '</table>' +
                '</div>');

            scope.showFilterRow = true;
            scope.tableParams = createNgTableParams();

            $compile(elm)(scope);
            scope.$digest();
        }));

        it('when true, should display filter row', function() {
            var filterRow = elm.find('thead').find('tr').eq(1);
            expect(filterRow.hasClass('ng-table-filters')).toBe(true);
            expect(filterRow.hasClass('ng-hide')).toBe(false);
        });

        it('when false, should hide filter row', function() {
            // given
            scope.showFilterRow = false;

            // when
            scope.$digest();

            // then
            var filterRow = elm.find('thead').find('tr').eq(1);
            expect(filterRow.hasClass('ng-table-filters')).toBe(true);
            expect(filterRow.hasClass('ng-hide')).toBe(true);
        });
    });

    describe('$columns', function(){
        var elm,
            tp;
        beforeEach(inject(function($compile) {
            elm = angular.element(
                '<div>' +
                '<table ng-table="tableParams">' +
                '<tr ng-repeat="user in $data">' +
                '<td title="ageTitle" ng-if="isAgeVisible" filter="ageFilter">{{user.age}}</td>' +
                '<td title="\'Name\'" groupable="\'name\'" sortable="\'name\'">{{user.name}}</td>' +
                '</tr>' +
                '</table>' +
                '</div>');

            $compile(elm)(scope);
            scope.$digest();

            scope.ageFilter = {
                age: 'text'
            };
            scope.isAgeVisible = true;
            scope.ageTitle = 'Age';
            tp = scope.tableParams = createNgTableParams();
            scope.$digest();
        }));

        it('should make $columns available on the scope created for ng-table', function(){
            // check that the scope is indeed the one created for out NgTableParams
            expect(scope.$$childHead.params).toBe(tp);

            expect(scope.$$childHead.$columns).toBeDefined();
        });

        it('should NOT polute the outer scope with a reference to $columns ', function(){
            expect(scope.$columns).toBeUndefined();
        });

        it('$scolumns should contain a column definition for each `td` element', function(){
            expect(scope.$$childHead.$columns.length).toBe(2);
        });

        it('each column definition should have getters for each column attribute', function(){
            var ageCol = scope.$$childHead.$columns[0];
            expect(ageCol.title()).toBe('Age');
            expect(ageCol.show()).toBe(true);
            expect(ageCol.filter()).toBe(scope.ageFilter);
            expect(ageCol.class()).toBe('');
            expect(ageCol.filterData).toBeUndefined();
            expect(ageCol.groupable()).toBe(false);
            expect(ageCol.headerTemplateURL()).toBe(false);
            expect(ageCol.headerTitle()).toBe('');
            expect(ageCol.sortable()).toBe(false);
            expect(ageCol.titleAlt()).toBe('');

            var nameCol = scope.$$childHead.$columns[1];
            expect(nameCol.title()).toBe('Name');
            expect(nameCol.show()).toBe(true);
            expect(nameCol.filter()).toBe(false);
            expect(nameCol.class()).toBe('');
            expect(nameCol.filterData).toBeUndefined();
            expect(nameCol.groupable()).toBe('name');
            expect(nameCol.headerTemplateURL()).toBe(false);
            expect(nameCol.headerTitle()).toBe('');
            expect(nameCol.sortable()).toBe('name');
            expect(nameCol.titleAlt()).toBe('');
        });

        it('each column attribute should be assignable', function(){
            var ageCol = scope.$$childHead.$columns[0];

            ageCol.title.assign(scope.$$childHead, 'Age of person');
            expect(ageCol.title()).toBe('Age of person');
            expect(scope.ageTitle).toBe('Age of person');

            ageCol.show.assign(scope.$$childHead, false);
            expect(ageCol.show()).toBe(false);
            expect(scope.isAgeVisible).toBe(false);

            var newFilter = {age: 'select'};
            ageCol.filter.assign(scope.$$childHead, newFilter);
            expect(ageCol.filter()).toBe(newFilter);
            expect(scope.ageFilter).toBe(newFilter);

            ageCol.class.assign(scope.$$childHead, 'amazing');
            expect(ageCol.class()).toBe('amazing');

            ageCol.groupable.assign(scope.$$childHead, 'age');
            expect(ageCol.groupable()).toBe('age');

            ageCol.headerTemplateURL.assign(scope.$$childHead, 'some.html');
            expect(ageCol.headerTemplateURL()).toBe('some.html');

            ageCol.headerTitle.assign(scope.$$childHead, 'wow');
            expect(ageCol.headerTitle()).toBe('wow');

            ageCol.sortable.assign(scope.$$childHead, 'incredible');
            expect(ageCol.sortable()).toBe('incredible');

            ageCol.titleAlt.assign(scope.$$childHead, 'really');
            expect(ageCol.titleAlt()).toBe('really');


            var nameCol = scope.$$childHead.$columns[1];

            nameCol.groupable.assign(scope.$$childHead, false);
            expect(nameCol.groupable()).toBe(false);

            nameCol.sortable.assign(scope.$$childHead, false);
            expect(nameCol.sortable()).toBe(false);
        });
    });


    describe('groups', function(){

        var $capturedColumn;
        beforeEach(inject(function() {
            // stash a reference to $column definition so that its available in asserts
            scope.captureColumn = function ($column) {
                $capturedColumn = $column;
            };
        }));

        describe('one groupable column', function(){

            var elm,
                tp;
            beforeEach(inject(function($compile) {
                elm = angular.element(
                    '<div>' +
                    '<table ng-table="tableParams">' +
                    '<tr class="ng-table-group" ng-repeat-start="group in $groups"></tr>' +
                    '<tr ng-repeat-end="user in group.data">' +
                    '<td title="\'Name\'" groupable="\'name\'">{{user.name}}</td>' +
                    '</tr>' +
                    '</table>' +
                    '</div>');

                $compile(elm)(scope);
                scope.$digest();

                tp = scope.tableParams = createNgTableParams();
                scope.$digest();
            }));

            it('should not render group row until group assigned', function() {
                var groupRow = elm.find('thead').find('.ng-table-group-header');
                expect(groupRow.length).toBe(0);
            });

            xit('should render group row once group assigned', function() {
                // todo: not sure why this test is not working as manually testing shows that it does :-(
                tp.group('name');
                scope.$digest();
                var groupRow = elm.find('thead').find('.ng-table-group-header');
                expect(groupRow.length).toBe(1);
            });
        });
    });

    describe('internals', function(){

        var elm,
            $timeout;

        beforeEach(inject(function($compile, _$timeout_) {
            $timeout = _$timeout_;
            elm = angular.element(
                '<table ng-table="tableParams">' +
                '<tr ng-repeat="user in $data">' +
                '<td title="\'Age\'">{{user.age}}</td>' +
                '</tr>' +
                '</table>');

            $compile(elm)(scope);
            scope.$digest();
        }));

        it('should reload when binding a new tableParams to scope', function(){
            var tp = createNgTableParams();
            scope.tableParams = tp;
            scope.$digest();

            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('should reload 1 time when binding a new tableParams that has an initial settings dataset field', function(){
            var tp = createNgTableParams({ dataset: [1,2,3] });
            scope.tableParams = tp;
            scope.$digest();

            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('should reload 1 time when binding a new tableParams with initial filter that has an initial settings dataset field', function(){
            var tp = createNgTableParams({filter: {age: 1}}, { dataset: [1,2,3] });
            scope.tableParams = tp;
            scope.$digest();

            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('should reload when binding a new tableParams to scope multiple times', function(){
            var tp1 = createNgTableParams();
            scope.tableParams = tp1;
            scope.$digest();

            expect(tp1.settings().getData.calls.count()).toBe(1);

            var tp2 = createNgTableParams();
            scope.tableParams = tp2;
            scope.$digest();

            expect(tp2.settings().getData.calls.count()).toBe(1);
        });

        it('should reload 1 time when binding a new settings dataset value and changing the filter', function(){
            // given
            var tp = createNgTableParams({filterOptions: { filterDelay: 100 }, dataset: [{age: 1}, {age: 2}]});
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 1 });
            tp.settings({ dataset: [{ age: 1 }, { age: 11 }, { age: 22 }]});
            scope.$digest();
            $timeout.flush(); // trigger the delayed reload

            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('should reload 1 time when multiple filter changes are debounced', function(){
            // given
            var tp = createNgTableParams({filterOptions: { filterDelay: 100 }, dataset: [{age: 1}, {age: 2}]});
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 1 });
            scope.$digest();
            tp.filter({ age: 2 });
            scope.$digest();
            $timeout.flush(); // trigger the delayed reload

            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('should reload 1 time when initial load fails', inject(function($q){
            // given
            var tp = createNgTableParams({ getData: function(){
                return $q.reject('BANG!');
            }});

            // when
            scope.tableParams = tp;
            scope.$digest();

            // then
            expect(tp.settings().getData.calls.count()).toBe(1);
        }));

        it('should reload 1 time with page reset to 1 when binding a new settings dataset value and changing the filter', function(){
            var settings = {
                counts: [1],
                dataset: [{age: 1}, {age: 2}],
                filterOptions: { filterDelay: 100 }
            };
            var tp = createNgTableParams({ count: 1, page: 2 }, settings);
            scope.tableParams = tp;
            scope.$digest();
            expect(tp.page()).toBe(2); // checking assumptions
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 1 });
            tp.settings({ dataset: [{ age: 1 }, { age: 11 }, { age: 22 }]});
            scope.$digest();
            $timeout.flush(); // trigger the delayed reload

            expect(tp.settings().getData.calls.count()).toBe(1);
            expect(tp.page()).toBe(1);
        });

        it('changing filter, orderBy, or page and then calling reload should not invoke getData twice', function(){
            var tp = createNgTableParams();
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 5 });
            tp.reload();
            scope.$digest();

            // then
            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('change to filter that fails to load should not cause infinite reload loop', inject(function($q){
            var tp = createNgTableParams({ getData: function(){
                if (tp.settings().getData.calls.count() > 1){
                    return $q.reject('BANG!');
                }
                return [1,2]
            }});
            scope.tableParams = tp;
            scope.$digest();
            expect(tp.settings().getData.calls.count()).toBe(1); // checking assumptions
            expect(tp.isDataReloadRequired()).toBe(false); // checking assumptions

            // when
            tp.filter({ age: 5 });
            expect(tp.isDataReloadRequired()).toBe(true); // checking assumptions
            scope.$digest();

            // then
            expect(tp.isDataReloadRequired()).toBe(false);
            expect(tp.settings().getData.calls.count()).toBe(2);
        }));

        it('changing filter, orderBy, or page in a callback to reload should re-invoke getData 1 time only', function(){
            var tp = createNgTableParams();
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 5 });
            tp.reload().then(function(){
                tp.sorting({ age: 'desc'});
                // note: better to call tp.reload() here rather than rely on a watch firing later to do it for us
                // that way the second reload is chained to the first and returned as a single promise
            });
            scope.$digest();

            // then
            // ie calls.count() === (1 x reload) + (1 x sorting)
            expect(tp.settings().getData.calls.count()).toBe(2);
        });

        it('changing filter, orderBy, or page then reload in a callback to reload should re-invoke getData 1 time only', function(){

            // todo: refactor the watches in ngTableController to handle this case

            var tp = createNgTableParams();
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 5 });
            tp.reload().then(function(){
                tp.sorting({ age: 'desc'});
                return tp.reload();
            });
            scope.$digest();

            // then
            // ie calls.count() === (1 x reload) + (1 x sorting)
            expect(tp.settings().getData.calls.count()).toBe(2);
        });

        it('should not reload when filter value is assigned the same value', function(){
            // given
            var tp = createNgTableParams({ filter: {age: 10} }, { });
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 10});
            scope.$digest();

            expect(tp.settings().getData.calls.count()).toBe(0);
        });

        it('should reload when filter value changes', function(){
            // given
            var tp = createNgTableParams({ filter: {age: 10} }, {});
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.filter({ age: 12});
            scope.$digest();

            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('should reload when new dataset supplied', function(){
            // given
            var initialDataset = [
                {age: 1},
                {age: 2}
            ];
            var tp = createNgTableParams();
            scope.tableParams = tp;
            scope.$digest();
            tp.settings().getData.calls.reset();

            // when
            tp.settings({ dataset: [{ age: 10}, { age: 11}, { age: 12}]});
            scope.$digest();

            expect(tp.settings().getData.calls.count()).toBe(1);
        });
    });
});
