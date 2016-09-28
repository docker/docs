describe('NgTableParams', function () {
    var scope,
        NgTableParams,
        $rootScope;

    beforeEach(module("ngTable"));
    beforeEach(function(){
        module(function($provide){
            $provide.decorator('ngTableDefaultGetData', createSpy);


            createSpy.$inject = ['$delegate'];
            function createSpy(ngTableDefaultGetData){
                return jasmine.createSpy('ngTableDefaultGetDataSpy',ngTableDefaultGetData).and.callThrough();
            }
        });
    });

    beforeEach(inject(function ($controller, _$rootScope_, _NgTableParams_) {
        $rootScope = _$rootScope_;
        scope = $rootScope.$new();
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

    it('NgTableParams should be defined', function () {
        expect(NgTableParams).toBeDefined();
    });

    describe('generatePagesArray', function(){
        it('should generate pages array using arguments supplied', function(){
            var params = new NgTableParams();
            expect(params.generatePagesArray(1, 30, 10)).toEqual([
                { type: 'prev', number: 1, active: false },
                { type: 'first', number: 1, active: false, current: true },
                { type: 'page', number: 2, active: true, current: false },
                { type: 'last', number: 3, active: true, current: false },
                { type: 'next', number: 2, active: true }
            ]);
            expect(params.generatePagesArray(2, 30, 10)).toEqual([
                { type: 'prev', number: 1, active: true },
                { type: 'first', number: 1, active: true, current: false },
                { type: 'page', number: 2, active: false, current: true },
                { type: 'last', number: 3, active: true, current: false },
                { type: 'next', number: 3, active: true }
            ]);
            expect(params.generatePagesArray(2, 100, 10)).toEqual([
                { type: 'prev', number: 1, active: true },
                { type: 'first', number: 1, active: true, current: false },
                { type: 'page', number: 2, active: false, current: true },
                { type: 'page', number: 3, active: true, current: false },
                { type: 'page', number: 4, active: true, current: false },
                { type: 'page', number: 5, active: true, current: false },
                { type: 'page', number: 6, active: true, current: false },
                { type: 'page', number: 7, active: true, current: false },
                { type: 'more', active: false },
                { type: 'last', number: 10, active: true, current: false },
                { type: 'next', number: 3, active: true }
            ]);
        });

        it('should use own parameter values to generate pages when no arguments supplied', function(){
            var params = new NgTableParams({ page: 2, count: 10 }, { total: 30 });

            expect(params.generatePagesArray()).toEqual([
                { type: 'prev', number: 1, active: true },
                { type: 'first', number: 1, active: true, current: false },
                { type: 'page', number: 2, active: false, current: true },
                { type: 'last', number: 3, active: true, current: false },
                { type: 'next', number: 3, active: true }
            ]);

        });
    });

    it('NgTableParams `page` parameter', function () {
        var params = new NgTableParams();

        expect(params.page()).toBe(1);
        expect(params.page(2)).toEqual(params);
        expect(params.page()).toBe(2);

        params = new NgTableParams({
            page: 3
        });
        expect(params.page()).toBe(3);
    });

    it('NgTableParams parse url parameters', function () {
        var params = new NgTableParams({
            'sorting[name]': 'asc',
            'sorting[age]': 'desc',
            'filter[name]': 'test',
            'filter[age]': '20',
            'group[name]' : 'asc',
            'group[age]' : 'desc',
            'group[surname]' : undefined
        });

        expect(params.filter()).toEqual({ 'name': 'test', 'age': 20 });
        expect(params.sorting()).toEqual({ 'age': 'desc' }); // sorting only by one column - todo: remove restriction
        expect(params.group()).toEqual({ 'name' : 'asc', 'age': 'desc', 'surname': undefined });
    });

    it('NgTableParams return url parameters', function () {
        var params = new NgTableParams({
            'sorting[name]': 'asc',
            'sorting[age]': 'desc',
            'filter[name]': 'test',
            'filter[age]': '20',
            'group[name]' : 'asc',
            'group[age]' : 'desc',
            'group[surname]' : ''
        });
        expect(params.url()).toEqual({
            'page': '1',
            'count': '10',
            'filter[name]': 'test',
            'filter[age]': '20',
            'sorting[age]': 'desc',  // sorting only by one column - todo: remove restriction
            'group[name]' : 'asc',
            'group[age]' : 'desc',
            'group[surname]' : ''
        });
        expect(params.url(true)).toEqual([
            'page=1',
            'count=10',
            'filter[name]=test',
            'filter[age]=20',
            'sorting[age]=desc',  // sorting only by one column - todo: remove restriction
            'group[name]=asc',
            'group[age]=desc',
            'group[surname]='
        ]);
    });

    it('orderBy', function () {
        var params = new NgTableParams({
            'sorting[name]': 'asc'
        });

        expect(params.orderBy()).toEqual([ '+name' ]); // for angular sorting function

        params.sorting({ name: 'desc', age: 'asc' });

        expect(params.orderBy()).toEqual([ '-name', '+age' ]);
    });

    describe('group', function(){
        it('one group', function () {
            var params = new NgTableParams({
                group: 'role'
            }, {
                groupOptions: { defaultSort: 'desc' }
            });

            expect(params.hasGroup()).toBe(true);
            expect(params.hasGroup('role')).toBe(true);
            expect(params.hasGroup('role', 'desc')).toBe(true);
            expect(params.group()).toEqual({ role: 'desc' });

            params.group('age');
            expect(params.hasGroup('age', 'desc')).toBe(true);
            expect(params.group()).toEqual({ age: 'desc' });

            params.group('age', 'asc');
            expect(params.hasGroup('age', 'asc')).toBe(true);
            expect(params.group()).toEqual({ age: 'asc' });
        });

        it('one group (nested property)', function () {
            var params = new NgTableParams({
                group: 'details.personal.age'
            }, {
                groupOptions: { defaultSort: 'desc' }
            });

            expect(params.hasGroup()).toBe(true);
            expect(params.hasGroup('details.personal.age')).toBe(true);
            expect(params.hasGroup('details.personal.age', 'desc')).toBe(true);
            expect(params.group()).toEqual({ 'details.personal.age': 'desc' });

            params.group('details.country');
            expect(params.hasGroup('details.country', 'desc')).toBe(true);
            expect(params.group()).toEqual({ 'details.country': 'desc' });

            params.group('details.country', 'asc');
            expect(params.hasGroup('details.country', 'asc')).toBe(true);
            expect(params.group()).toEqual({ 'details.country': 'asc' });
        });

        it('one group function', function () {
            var params = new NgTableParams({
                group: angular.identity
            });

            expect(params.hasGroup()).toBe(true);
            expect(params.hasGroup(angular.identity)).toBe(true);
            expect(params.hasGroup(angular.identity, params.settings().groupOptions.defaultSort)).toBe(true);
            expect(params.group()).toEqual(angular.identity);


            var fn = function(){ };
            fn.sortDirection = 'desc';
            params.group(fn);
            expect(params.hasGroup(fn)).toBe(true);
            expect(params.hasGroup(fn, 'desc')).toBe(true);
            expect(params.group()).toEqual(fn);
        });

        it('can clear group', function () {
            var params = new NgTableParams({
                group: 'role'
            });

            expect(params.hasGroup()).toBe(true); // checking assumptions

            params.group({});
            expect(params.hasGroup()).toBe(false);
            expect(params.group()).toEqual({});
        });

        it('multiple groups', function () {
            var params = new NgTableParams({
                group: 'role'
            }, {
                groupOptions: { defaultSort: 'desc' }
            });

            var newGroups = _.extend({}, params.group(), { age: 'desc'});
            params.group(newGroups);
            expect(params.hasGroup()).toBe(true);
            expect(params.hasGroup('role')).toBe(true);
            expect(params.hasGroup('role', 'desc')).toBe(true);
            expect(params.hasGroup('age')).toBe(true);
            expect(params.hasGroup('age', 'desc')).toBe(true);
            expect(params.group()).toEqual({ role: 'desc', age: 'desc' });

            params.group({ role: 'asc', age: 'asc'});
            expect(params.hasGroup('age', 'desc')).toBe(false);
            expect(params.hasGroup('age', 'asc')).toBe(true);
            expect(params.group()).toEqual({ role: 'asc', age: 'asc' });
        });

        it('should apply current defaultSort from groupOptions', function(){
            var params = new NgTableParams({
                group: 'role'
            }, {
                groupOptions: { defaultSort: 'desc' }
            });

            params.settings({ groupOptions: { defaultSort: 'asc' }});

            params.group('age');
            expect(params.group()).toEqual({ age: 'asc' });
        });

        it('should not apply defaultSort from groupOptions when explicitly set to empty string', function(){
            var params = new NgTableParams({
                group: 'role'
            }, {
                groupOptions: { defaultSort: 'desc' }
            });
            expect(params.group()).toEqual({ role: 'desc' }); // checking assumptions

            params.group({ role: ''});
            expect(params.group()).toEqual({ role: '' });

            params.group({ age: undefined});
            expect(params.group()).toEqual({ age: 'desc' });

            params.group({ role: null});
            expect(params.group()).toEqual({ role: 'desc' });
        });
    });

    describe('settings', function(){

        it('defaults', function () {
            var params = new NgTableParams();

            var expectedSettings = {
                $loading: false,
                dataset: null,
                total: 0,
                defaultSort: 'desc',
                counts: [10, 25, 50, 100],
                interceptors: [],
                paginationMaxBlocks: 11,
                paginationMinBlocks: 5,
                sortingIndicator: 'span',
                filterOptions: {
                    filterComparator: undefined,
                    filterDelay: 500,
                    filterDelayThreshold: 10000,
                    filterFilterName: undefined,
                    filterFn: undefined,
                    filterLayout: 'stack'
                },
                groupOptions: { defaultSort: 'asc', isExpanded: true }
            };
            expect(params.settings()).toEqual(jasmine.objectContaining(expectedSettings));
            expect(params.settings().getData).toEqual(jasmine.any(Function));
            expect(params.settings().getGroups).toEqual(jasmine.any(Function));

            params = new NgTableParams({}, {
                total: 100,
                counts: [1,2],
                groupOptions: { isExpanded: false } });

            expectedSettings.total = 100;
            expectedSettings.counts = [1,2];
            expectedSettings.groupOptions = { defaultSort: 'asc', isExpanded: false };
            expect(params.settings()).toEqual(jasmine.objectContaining(expectedSettings));
        });

        it('changing settings().dataset should reset page to 1', function(){
            // given
            var tableParams = createNgTableParams({ count: 1, page: 2 }, { dataset: [1,2,3]});
            tableParams.reload();
            scope.$digest();
            expect(tableParams.page()).toBe(2); // checking assumptions

            // when
            tableParams.settings({ dataset: [1,2,3, 4]});

            // then
            expect(tableParams.page()).toBe(1);
        });

        it('should not set filterDelay when working with synchronous dataset', function(){
            // given
            var tableParams = new NgTableParams({}, { dataset: [1,2,3]});
            expect(tableParams.settings().filterOptions.filterDelay).toBe(0);
        });

        it('should not set filterDelay when working with synchronous dataset (empty dataset)', function(){
            // given
            var tableParams = new NgTableParams({}, { dataset: []});
            expect(tableParams.settings().filterOptions.filterDelay).toBe(0);
        });

        it('should set filterDelay when not certain working with synchronous dataset', function(){
            // given
            var tableParams = new NgTableParams({}, { dataset: [1,2], getData: function(){
                // am I sync or async?
            }});
            expect(tableParams.settings().filterOptions.filterDelay).toBe(500);
        });

        it('should set filterDelay when dataset exceeds filterDelayThreshold', function(){
            // given
            var tableParams = new NgTableParams({}, { filterOptions: { filterDelayThreshold: 5 }, dataset: [,2,3,4,5,6] });
            expect(tableParams.settings().filterOptions.filterDelay).toBe(500);
        });

        it('should allow filterDelay to be set explicitly', function(){
            // given
            var tableParams = new NgTableParams({}, { filterOptions: { filterDelay: 100}, dataset: [1,2] });
            expect(tableParams.settings().filterOptions.filterDelay).toBe(100);
        });
    });

    describe('reload', function(){

        it('should call getData to retrieve data', function(){
            var tp = createNgTableParams();
            tp.reload();
            scope.$digest();
            expect(tp.settings().getData.calls.count()).toBe(1);
        });

        it('should add the results returned by getData to the data field', function(){
            var tp = createNgTableParams({ getData: function(){
                return [1,2,3];
            }});
            tp.reload();
            scope.$digest();
            expect(tp.data).toEqual([1,2,3]);
        });

        it('should use ngTableDefaultGetData function when NgTableParams not supplied a getData function', inject(function(ngTableDefaultGetData){
            // given
            var tp = createNgTableParams();

            // when
            tp.reload();
            scope.$digest();

            // then
            expect(ngTableDefaultGetData).toHaveBeenCalled();
        }));

        it('should propagate rejection reason from getData', inject(function($q){
            // given
            var tp = createNgTableParams({ getData: function(){
                return $q.reject('bad response')
            }});

            // when
            var actualRejection;
            tp.reload().catch(function(reason){
                actualRejection = reason;
            });
            scope.$digest();

            // then
            expect(actualRejection).toBe('bad response');
        }));
    });

    describe('getGroups', function(){

        var dataset;
        beforeEach(function(){
            dataset = [
                { 'name': 'Hanson', 'role': 'Accounting' },
                { 'name': 'Eaton', 'role': 'Customer Service' },
                { 'name': 'Perry', 'role': 'Customer Service' },
                { 'name': 'George', 'role': 'Accounting' },
                { 'name': 'Jennings', 'role': 'Asset Management' },
                { 'name': 'Whitney', 'role': 'Accounting' },
                { 'name': 'Weaver', 'role': 'Payroll' },
                { 'name': 'Gibson', 'role': 'Payroll' },
                { 'name': 'Wells', 'role': 'Media Relations' },
                { 'name': 'Willis', 'role': 'Finances' },
                { 'name': 'Donovan', 'role': 'Customer Relations' },
                { 'name': 'Mcdonald', 'role': 'Finances' },
                { 'name': 'Young', 'role': 'Asset Management' }
            ];
        });


        it('should group data then apply paging to groups', function () {
            var tp = createNgTableParams({ count: 2, group: { role: '' } }, { dataset: dataset });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Accounting',
                    data: [
                        { "name": 'Hanson', "role": 'Accounting' },
                        { "name": 'George', "role": 'Accounting' },
                        { "name": 'Whitney', "role": 'Accounting' }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Customer Service',
                    data: [
                        { "name": 'Eaton', "role": 'Customer Service' },
                        { "name": 'Perry', "role": 'Customer Service' }
                    ]
                }
            ]);
        });

        it('should sort data, then group, then page groups', function () {
            var tp = createNgTableParams({
                count: 2,
                sorting: { name: 'desc'},
                group: { role: '' }
            }, { dataset: dataset });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Asset Management',
                    data: [
                        { "name": 'Young', "role": 'Asset Management' },
                        { "name": 'Jennings', "role": 'Asset Management' }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Finances',
                    data: [
                        { "name": 'Willis', "role": 'Finances' },
                        { "name": 'Mcdonald', "role": 'Finances' }
                    ]
                }
            ]);
        });

        it('should use group function to group data', function () {
            var grouper = function (item) {
                return item.name[0];
            };
            grouper.sortDirection = '';
            var tp = createNgTableParams({
                count: 2,
                sorting: {name: 'desc'},
                group: grouper
            }, {dataset: dataset});

            var actualRoleGroups;
            tp.reload().then(function (groups) {
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Y',
                    data: [
                        { "name": "Young", "role": "Asset Management" }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'W',
                    data: [
                        { "name": "Willis", "role": "Finances" },
                        { "name": "Whitney", "role": "Accounting" },
                        { "name": "Wells", "role": "Media Relations" },
                        { "name": "Weaver", "role": "Payroll" }
                    ]
                }
            ]);
        });

        it('should filter and sort data, then group, then page groups', function () {
            var tp = createNgTableParams({
                count: 2,
                sorting: { name: 'desc' },
                filter: { name: 'e' },
                group: { role: '' }
            }, {
                dataset: dataset
            });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Accounting',
                    data: [
                        { 'name': 'Whitney', 'role': 'Accounting' },
                        { 'name': 'George', 'role': 'Accounting' }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Media Relations',
                    data: [
                        { 'name': 'Wells', 'role': 'Media Relations' }
                    ]
                }
            ]);
        });

        it('should filter and sort data, then group, then apply group sortDirection, and finally page groups', function () {
            var tp = createNgTableParams({
                count: 3,
                sorting: { name: 'desc' },
                filter: { name: 'e' },
                group: { role: 'desc' }
            }, {
                dataset: dataset,
                groupOptions: {
                    // this value will be overridden by group: { role: 'desc' }
                    defaultSort: undefined
                }
            });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Payroll',
                    data: [
                        { 'name': 'Weaver', 'role': 'Payroll' }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Media Relations',
                    data: [
                        { 'name': 'Wells', 'role': 'Media Relations' }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Customer Service',
                    data: [
                        { 'name': 'Perry', 'role': 'Customer Service' },
                        { 'name': 'Eaton', 'role': 'Customer Service' }
                    ]
                }
            ]);
        });

        it('should use sortDirection defined on group function to sort groups', function () {
            var groupFn = function(item){
                return item.role;
            };
            groupFn.sortDirection = 'desc';
            var tp = createNgTableParams({
                count: 3,
                sorting: { name: 'desc' },
                filter: { name: 'e' },
                group: groupFn
            }, {
                dataset: dataset,
                groupOptions: {
                    // this value will be overridden by groupFn.sortDirection
                    defaultSort: undefined
                }
            });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Payroll',
                    data: [
                        { 'name': 'Weaver', 'role': 'Payroll' }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Media Relations',
                    data: [
                        { 'name': 'Wells', 'role': 'Media Relations' }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Customer Service',
                    data: [
                        { 'name': 'Perry', 'role': 'Customer Service' },
                        { 'name': 'Eaton', 'role': 'Customer Service' }
                    ]
                }
            ]);
        });
    });

    describe('getGroups with nested property', function(){

        var dataset;
        beforeEach(function(){
            dataset = [
                { 'details': { 'name': 'Hanson', 'role': 'Accounting' } },
                { 'details': { 'name': 'Eaton', 'role': 'Customer Service' } },
                { 'details': { 'name': 'Perry', 'role': 'Customer Service' } },
                { 'details': { 'name': 'George', 'role': 'Accounting' } },
                { 'details': { 'name': 'Jennings', 'role': 'Asset Management' } },
                { 'details': { 'name': 'Whitney', 'role': 'Accounting' } },
                { 'details': { 'name': 'Weaver', 'role': 'Payroll' } },
                { 'details': { 'name': 'Gibson', 'role': 'Payroll' } },
                { 'details': { 'name': 'Wells', 'role': 'Media Relations' } },
                { 'details': { 'name': 'Willis', 'role': 'Finances' } },
                { 'details': { 'name': 'Donovan', 'role': 'Customer Relations' } },
                { 'details': { 'name': 'Mcdonald', 'role': 'Finances' } },
                { 'details': { 'name': 'Young', 'role': 'Asset Management' } }
            ];
        });


        it('should group data then apply paging to groups', function () {
            var tp = createNgTableParams({ count: 2, group: { 'details.role': '' } }, { dataset: dataset });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Accounting',
                    data: [
                        { 'details': { 'name': 'Hanson', 'role': 'Accounting' } },
                        { 'details': { 'name': 'George', 'role': 'Accounting' } },
                        { 'details': { 'name': 'Whitney', 'role': 'Accounting' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Customer Service',
                    data: [
                        { 'details': { 'name': 'Eaton', 'role': 'Customer Service' } },
                        { 'details': { 'name': 'Perry', 'role': 'Customer Service' } }
                    ]
                }
            ]);
        });

        it('should sort data, then group, then page groups', function () {
            var tp = createNgTableParams({
                count: 2,
                sorting: { 'details.name': 'desc'},
                group: { 'details.role': '' }
            }, { dataset: dataset });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Asset Management',
                    data: [
                        { 'details': { 'name': 'Young', 'role': 'Asset Management' } },
                        { 'details': { 'name': 'Jennings', 'role': 'Asset Management' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Finances',
                    data: [
                        { 'details': { 'name': 'Willis', 'role': 'Finances' } },
                        { 'details': { 'name': 'Mcdonald', 'role': 'Finances' } }
                    ]
                }
            ]);
        });

        it('should use group function to group data', function () {
            var grouper = function (item) {
                return item.details.name[0];
            };
            grouper.sortDirection = '';
            var tp = createNgTableParams({
                count: 2,
                sorting: {'details.name': 'desc'},
                group: grouper
            }, {dataset: dataset});

            var actualRoleGroups;
            tp.reload().then(function (groups) {
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Y',
                    data: [
                        { 'details': { 'name': 'Young', 'role': 'Asset Management' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'W',
                    data: [
                        { 'details': { 'name': 'Willis', 'role': 'Finances' } },
                        { 'details': { 'name': 'Whitney', 'role': 'Accounting' } },
                        { 'details': { 'name': 'Wells', 'role': 'Media Relations' } },
                        { 'details': { 'name': 'Weaver', 'role': 'Payroll' } }
                    ]
                }
            ]);
        });

        it('should filter and sort data, then group, then page groups', function () {
            var tp = createNgTableParams({
                count: 2,
                sorting: { 'details.name': 'desc' },
                filter: { 'details.name': 'e' },
                group: { 'details.role': '' }
            }, {
                dataset: dataset
            });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Accounting',
                    data: [
                        { 'details': { 'name': 'Whitney', 'role': 'Accounting' } },
                        { 'details': { 'name': 'George', 'role': 'Accounting' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Media Relations',
                    data: [
                        { 'details': { 'name': 'Wells', 'role': 'Media Relations' } }
                    ]
                }
            ]);
        });

        it('should filter and sort data, then group, then apply group sortDirection, and finally page groups', function () {
            var tp = createNgTableParams({
                count: 3,
                sorting: { 'details.name': 'desc' },
                filter: { 'details.name': 'e' },
                group: { 'details.role': 'desc' }
            }, {
                dataset: dataset,
                groupOptions: {
                    // this value will be overridden by group: { role: 'desc' }
                    defaultSort: undefined
                }
            });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Payroll',
                    data: [
                        { 'details': { 'name': 'Weaver', 'role': 'Payroll' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Media Relations',
                    data: [
                        { 'details': { 'name': 'Wells', 'role': 'Media Relations' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Customer Service',
                    data: [
                        { 'details': { 'name': 'Perry', 'role': 'Customer Service' } },
                        { 'details': { 'name': 'Eaton', 'role': 'Customer Service' } }
                    ]
                }
            ]);
        });

        it('should use sortDirection defined on group function to sort groups', function () {
            var groupFn = function(item){
                return item.details.role;
            };
            groupFn.sortDirection = 'desc';
            var tp = createNgTableParams({
                count: 3,
                sorting: { 'details.name': 'desc' },
                filter: { 'details.name': 'e' },
                group: groupFn
            }, {
                dataset: dataset,
                groupOptions: {
                    // this value will be overridden by groupFn.sortDirection
                    defaultSort: undefined
                }
            });

            var actualRoleGroups;
            tp.reload().then(function(groups){
                actualRoleGroups = groups;
            });
            $rootScope.$digest();

            expect(actualRoleGroups).toEqual([
                {
                    $hideRows: false,
                    value: 'Payroll',
                    data: [
                        { 'details': { 'name': 'Weaver', 'role': 'Payroll' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Media Relations',
                    data: [
                        { 'details': { 'name': 'Wells', 'role': 'Media Relations' } }
                    ]
                },
                {
                    $hideRows: false,
                    value: 'Customer Service',
                    data: [
                        { 'details': { 'name': 'Perry', 'role': 'Customer Service' } },
                        { 'details': { 'name': 'Eaton', 'role': 'Customer Service' } }
                    ]
                }
            ]);
        });
    });

    it('ngTableParams test defaults', inject(function ($q, ngTableDefaults) {
        ngTableDefaults.params.count = 2;
        ngTableDefaults.settings.counts = [];
        var tp = new NgTableParams();

        expect(tp.count()).toEqual(2);
        expect(tp.page()).toEqual(1);

        var settings = tp.settings();
        expect(settings.counts.length).toEqual(0);
        expect(settings.interceptors.length).toEqual(0);
        expect(settings.filterOptions.filterDelay).toEqual(500);

        ngTableDefaults.settings.interceptors = [ { response: angular.identity }];
        tp = new NgTableParams();
        expect(tp.settings().interceptors.length).toEqual(1);
    }));

    describe('hasFilter', function(){
        var tp;

        beforeEach(inject(function(){
            tp = new NgTableParams({}, {});
        }));

        it('should return false for an empty filter object', function(){
            tp.filter({});
            expect(tp.hasFilter()).toBe(false);
        });

        it('should return true for when filter has a field with a significant value', function(){
            tp.filter({ a: 'b' });
            expect(tp.hasFilter()).toBe(true);

            tp.filter({ a: 0 });
            expect(tp.hasFilter()).toBe(true);
        });

        it('should return false when filter only has insignificant field values', function(){
            tp.filter({ a: '' });
            expect(tp.hasFilter()).toBe(false);

            tp.filter({ a: null });
            expect(tp.hasFilter()).toBe(false);

            tp.filter({ a: undefined });
            expect(tp.hasFilter()).toBe(false);

            tp.filter({ a: undefined, b: '', c: undefined });
            expect(tp.hasFilter()).toBe(false);

            //tp.filter({ a: NaN });
            //expect(tp.hasFilter()).toBe(false);
        });
    });


    describe('isDataReloadRequired', function(){
        var tp;

        beforeEach(inject(function(){
            tp = createNgTableParams();
        }));

        it('should return true after construction', function(){
            expect(tp.isDataReloadRequired()).toBe(true);
        });

        it('should return false once reload called after construction', function(){
            tp.reload();
            // note: we don't have to wait for the getData promise to be resolved before considering reload
            // to be unnecessary - that's why we're not having to run a $digest
            expect(tp.isDataReloadRequired()).toBe(false);
        });

        it('should return false when getData fails', inject(function($q){
            tp.settings({ getData: function(){
                return $q.reject('bad response');
            }});
            tp.reload();
            scope.$digest();
            expect(tp.isDataReloadRequired()).toBe(false);
        }));

        it('should detect direct changes to parameters', inject(function($q){
            // given
            tp.reload();
            scope.$digest();
            expect(tp.isDataReloadRequired()).toBe(false); // checking assumptions

            // when
            tp.filter().newField = 99;
            expect(tp.isDataReloadRequired()).toBe(true);
        }));

        it('should return true until changed parameters have been reloaded', function(){
            // given
            tp.reload();
            scope.$digest();

            // when, then...
            verifyIsDataReloadRequired(function(){
                tp.filter({ age: 1});
            });
            verifyIsDataReloadRequired(function(){
                tp.page(55);
            });
            verifyIsDataReloadRequired(function(){
                tp.sorting({ age: 'desc'});
            });
            verifyIsDataReloadRequired(function(){
                tp.count(100);
            });
        });

        it('should return true when `$` field on `filter` changes', function(){
            // given
            tp.reload();
            scope.$digest();
            expect(tp.isDataReloadRequired()).toBe(false); // checking assumptions

            // when
            tp.filter().$ = 'cc';
            expect(tp.isDataReloadRequired()).toBe(true);
        });

        it('should return true when group function sort direction changes', function(){
            // given
            tp.reload();
            scope.$digest();
            expect(tp.isDataReloadRequired()).toBe(false); // checking assumptions

            // when, then...

            var grouper = function(){ };
            tp.group(grouper);
            expect(tp.isDataReloadRequired()).toBe(true);

            tp.reload();
            scope.$digest();
            expect(tp.isDataReloadRequired()).toBe(false);

            grouper.sortDirection = 'desc';
            expect(tp.isDataReloadRequired()).toBe(true);
        });

        it('should return false when parameters "touched" but not modified', function(){
            // given
            tp.filter({ age: 1});
            tp.reload();
            scope.$digest();

            // when, then...
            tp.filter({ age: 1});
            expect(tp.isDataReloadRequired()).toBe(false);
        });

        it('should return true when new settings dataset array supplied', function(){
            // given
            tp.reload();
            scope.$digest();

            verifyIsDataReloadRequired(function(){
                tp.settings({ dataset: [11,22,33]});
            });
        });

        it('should return true when existing settings dataset array is unset', function(){
            // given
            tp = createNgTableParams({ dataset: [1,2,3]});
            tp.reload();
            scope.$digest();

            verifyIsDataReloadRequired(function(){
                tp.settings({ dataset: null});
            });
        });

        it('status should not change when settings called without a dataset array', function(){
            // given
            tp = createNgTableParams({ dataset: [1,2,3]});
            tp.reload();
            scope.$digest();

            // when, then...
            tp.settings({ total : 100 });
            expect(tp.isDataReloadRequired()).toBe(false);
        });


        function verifyIsDataReloadRequired(modifer){
            modifer();
            expect(tp.isDataReloadRequired()).toBe(true);
            tp.reload();
            scope.$digest();
            expect(tp.isDataReloadRequired()).toBe(false);
        }
    });

    describe('hasFilterChanges', function(){
        var tp;

        beforeEach(inject(function(){
            tp = createNgTableParams();
        }));

        it('should return true after construction', function(){
            expect(tp.hasFilterChanges()).toBe(true);
        });

        it('should return false once reload called after construction', function(){
            tp.reload();
            // note: we don't have to wait for the getData promise to be resolved before considering reload
            // to be unnecessary - that's why we're not having to run a $digest
            expect(tp.hasFilterChanges()).toBe(false);
        });

        it('should return true when getData fails', inject(function($q){
            tp.settings({ getData: function(){
                return $q.reject('bad response');
            }});
            tp.reload();
            scope.$digest();
            expect(tp.hasFilterChanges()).toBe(false);
        }));

        it('should detect direct changes to filters', inject(function($q){
            // given
            tp.reload();
            scope.$digest();
            expect(tp.hasFilterChanges()).toBe(false); // checking assumptions

            // when
            tp.filter().newField = 99;
            expect(tp.hasFilterChanges()).toBe(true);
        }));

        it('should return true when `$` field on `filter` changes', function(){
            // given
            tp.reload();
            scope.$digest();
            expect(tp.hasFilterChanges()).toBe(false); // checking assumptions

            // when
            tp.filter().$ = 'cc';
            expect(tp.hasFilterChanges()).toBe(true);
        });

        it('should return true until changed filters have been reloaded', function(){
            // given
            tp.reload();
            scope.$digest();

            // when, then...
            tp.filter({ age: 1});
            expect(tp.hasFilterChanges()).toBe(true);
            tp.reload();
            scope.$digest();
            expect(tp.hasFilterChanges()).toBe(false);
        });

        it('should return false when filters "touched" but not modified', function(){
            // given
            tp.filter({ age: 1});
            tp.reload();
            scope.$digest();

            // when, then...
            tp.filter({ age: 1});
            expect(tp.hasFilterChanges()).toBe(false);
        });

        it('status should not change just because new settings dataset array supplied', function(){
            // given
            tp.reload();
            scope.$digest();

            // when, then...
            tp.settings({ dataset: [11,22,33]});
            expect(tp.hasFilterChanges()).toBe(false);
        });

    });

    describe('hasErrorState', function(){
        var tp;

        it('should return false until reload fails', inject(function($q){
            // given
            tp = createNgTableParams({ getData: function(){
                if (tp.settings().getData.calls.count() > 2){
                    return $q.reject('bad response');
                }
                return [1,2];
            }});

            // when, then
            tp.reload();
            scope.$digest();
            expect(tp.hasErrorState()).toBe(false);
            tp.reload();
            scope.$digest();
            expect(tp.hasErrorState()).toBe(false);
            tp.reload();
            scope.$digest();
            expect(tp.hasErrorState()).toBe(true);
        }));

        it('should return false once parameter values change', inject(function($q){
            // given
            tp = createNgTableParams({ getData: function(){
                return $q.reject('bad response');
            }});

            // when, then
            tp.reload();
            scope.$digest();
            expect(tp.hasErrorState()).toBe(true);
            tp.filter({ age: 598});
            expect(tp.hasErrorState()).toBe(false);
        }));
    });

    describe('backwards compatibility shim', function(){

        it('shim should supply getData original arguments', inject(function(ngTableGetDataBcShim){
            // given
            var callCount = 0;
            var adaptedFn = ngTableGetDataBcShim(originalGetDataFn);

            // when
            var tableParams = new NgTableParams({}, {});
            adaptedFn(tableParams);

            // then
            expect(callCount).toBe(1);

            function originalGetDataFn($defer, params){
                callCount++;
                expect($defer).toBeDefined();
                expect($defer.resolve).toBeDefined();
                expect(params).toBeDefined();
                expect(params).toEqual(jasmine.any(NgTableParams));
            }

        }));

        it('shim should return the getData "$defer" promise', inject(function(ngTableGetDataBcShim){
            // given
            var callCount = 0;
            var adaptedFn = ngTableGetDataBcShim(originalGetDataFn);

            // when
            var tableParams = new NgTableParams({}, {});
            var results = adaptedFn(tableParams);

            // then
            expect(results).toBeDefined();
            expect(results.then).toBeDefined();
            results.then(function(data){
                expect(data.length).toBe(3);
            });
            scope.$digest();

            function originalGetDataFn($defer/*, params*/){
                $defer.resolve([1,2,3]);
            }
        }));

        it('shim should return the data', inject(function(ngTableGetDataBcShim){
            // given
            var callCount = 0;
            var adaptedFn = ngTableGetDataBcShim(originalGetDataFn);

            // when
            var tableParams = new NgTableParams({}, {});
            var results = adaptedFn(tableParams);

            // then
            expect(results).toBeDefined();
            expect(results).toEqual([1,2,3]);

            function originalGetDataFn(/*$defer, params*/){
                return [1,2,3];
            }
        }));

        it('shim should be applied when getData function supplied has more than one parameter', function(){
            // given
            var tp = createNgTableParams({ getData: originalGetDataFn});

            // when
            var dataFetched = tp.reload();

            // then
            var actualData;
            dataFetched.then(function(data){
                actualData = data;
            });
            scope.$digest();
            expect(actualData).toEqual([1,2,3]);

            function originalGetDataFn($defer, params){
                params.total(3);
                $defer.resolve([1,2,3]);
            }
        });

        it('shim should NOT be applied when getData has new signature', function(){
            // given
            var tp = createNgTableParams({ getData: newGetDataFn});

            // when
            var dataFetched = tp.reload();

            // then
            var actualData;
            dataFetched.then(function(data){
                actualData = data;
            });
            scope.$digest();
            expect(actualData).toEqual([1,2,3]);

            function newGetDataFn(params){
                params.total(3);
                return [1,2,3];
            }
        });
    });

    describe('interceptors', function(){

        it('can register interceptor', function(){
            var interceptor = { response: angular.identity };
            var tp = createNgTableParams({ interceptors: [interceptor]});
            expect(tp.settings().interceptors).toEqual([interceptor]);
        });

        it('can register multiple interceptor', function(){
            var interceptors = [{ response: angular.identity }, { response: angular.identity }];
            var tp = createNgTableParams({ interceptors: interceptors});
            expect(tp.settings().interceptors).toEqual(interceptors);
        });

        it('can register interceptors after NgTableParams created', function(){
            var interceptor = { response: angular.identity };
            var interceptor2 = { response: angular.identity };
            var tp = createNgTableParams({ interceptors: [interceptor]});
            var interceptors = tp.settings().interceptors.concat([interceptor2]);
            tp.settings({ interceptors: interceptors});
            expect(tp.settings().interceptors).toEqual(interceptors);
        });

        describe('one response interceptor', function(){

            it('should receive response from call to getData', function(){
                // given
                var interceptor = {
                    response: function(/*data, params*/){
                        this.hasRun = true;
                    }
                };
                var tp = createNgTableParams({ interceptors: [interceptor]});

                // when
                tp.reload();
                scope.$digest();

                // then
                expect(interceptor.hasRun).toBeTruthy();
            });

            it('should be able to modify data returned by getData', function(){
                // given
                var interceptor = {
                    response: function(data/*, params*/){
                        data.forEach(function(item){
                            item.modified = true;
                        });
                        return data;
                    }
                };
                var tp = createNgTableParams({ interceptors: [interceptor], getData: function(){
                    return [{}, {}];
                }});

                // when
                var actualData;
                tp.reload().then(function(data){
                    actualData = data;
                });
                scope.$digest();

                // then
                expect(actualData).toEqual([{ modified: true }, { modified: true }]);
            });

            it('should be able to replace data returned by getData', function(){
                // given
                var interceptor = {
                    response: function(data/*, params*/){
                        return data.map(function(item){
                            return item * 2;
                        });
                    }
                };
                var tp = createNgTableParams({ interceptors: [interceptor], getData: function(){
                    return [2, 3];
                }});

                // when
                var actualData;
                tp.reload().then(function(data){
                    actualData = data;
                });
                scope.$digest();

                // then
                expect(actualData).toEqual([4, 6]);
            });

            it('should be able to access tableParams supplied to getData', function(){
                // given
                var interceptor = {
                    response: function(data, params){
                        params.total(data.total);
                        return data.results;
                    }
                };
                var tp = createNgTableParams({ interceptors: [interceptor], getData: function(){
                    return { results: [1,2,3], total: 3};
                }});

                // when
                var actualData;
                tp.reload().then(function(data){
                    actualData = data;
                });
                scope.$digest();

                // then
                expect(actualData).toEqual([1,2,3]);
                expect(tp.total()).toEqual(3);
            });
        });

        describe('one responseError interceptor', function(){

            it('should receive rejections from getData', inject(function($q){
                // given
                var interceptor = {
                    responseError: function(reason/*, params*/){
                        this.actualReason = reason;
                    }
                };
                var tp = createNgTableParams({
                    interceptors: [interceptor],
                    getData: function(/*params*/){
                        return $q.reject('BANG!');
                    }
                });

                // when
                tp.reload();
                scope.$digest();

                // then
                expect(interceptor.actualReason).toBe('BANG!');
            }));

            it('should NOT receive errors from getData', function(){
                // given
                var interceptor = {
                    responseError: function(/*reason, params*/){
                        this.hasRun = true;
                    }
                };
                var tp = createNgTableParams({
                    interceptors: [interceptor],
                    getData: function(/*params*/){
                        throw new Error('BANG!');
                    }
                });

                // when, then
                expect(tp.reload).toThrow();
                expect(interceptor.hasRun).toBeFalsy();
            });

            it('should be able to modify reason returned by getData', inject(function($q){
                // given
                var interceptor = {
                    responseError: function(reason/*, params*/){
                        reason.code = 400;
                        return $q.reject(reason);
                    }
                };
                var tp = createNgTableParams({
                    interceptors: [interceptor],
                    getData: function(/*params*/){
                        return $q.reject({ description: 'crappy data'});
                    }
                });

                // when
                var actualReason;
                tp.reload().catch(function(reason){
                    actualReason = reason;
                });
                scope.$digest();

                // then
                expect(actualReason.code).toBe(400);
                expect(actualReason.description).toBe('crappy data');
            }));

            it('should be able to replace reason returned by getData', inject(function($q){
                // given
                var interceptor = {
                    responseError: function(/*reason, params*/){
                        return $q.reject('Cancelled by user');
                    }
                };
                var tp = createNgTableParams({
                    interceptors: [interceptor],
                    getData: function(/*params*/){
                        return $q.reject('BANG!');
                    }
                });

                // when
                var actualReason;
                tp.reload().catch(function(reason){
                    actualReason = reason;
                });
                scope.$digest();

                // then
                expect(actualReason).toBe('Cancelled by user');
            }));

            it('should be able to access tableParams supplied to getData', function(){
                // given
                var interceptor = {
                    response: function(reason, params){
                        this.actualParams = params;
                    }
                };
                var tp = createNgTableParams({ interceptors: [interceptor]});

                // when
                var actualData;
                tp.reload();
                scope.$digest();

                // then
                expect(interceptor.actualParams).toBe(tp);
            });
        });

        describe('one response plus responseError interceptor', function(){

            it('should NOT call responseError on same interceptor whose response method fails', inject(function($q){
                // given
                var interceptor = {
                    response: function(/*data, params*/){
                        return $q.reject('BANG!');
                    },
                    responseError: function(reason/*, params*/){
                        this.hasRun = true;
                    }
                };
                var tp = createNgTableParams({
                    interceptors: [interceptor],
                    getData: function(/*params*/){
                        return [1,2,3];
                    }
                });

                // when
                tp.reload();
                scope.$digest();

                // then
                expect(interceptor.hasRun).toBeFalsy();
            }));
        });

        describe('multiple response interceptors', function(){

            it('should run interceptors in the order they were registered', function(){
                // given
                var callCount = 0;
                var interceptor = {
                    response: function(/*data, params*/){
                        this.sequence = callCount;
                        callCount++;
                    }
                };
                var interceptors = [interceptor, angular.copy(interceptor)];
                var tp = createNgTableParams({ interceptors: interceptors});

                // when
                tp.reload();
                scope.$digest();

                // then
                expect(interceptors[0].sequence).toBe(0);
                expect(interceptors[1].sequence).toBe(1);
            });

            it('results of one interceptor should be piped to the next validator', function(){
                // given
                var interceptor = {
                    response: function(data/*, params*/){
                        return data.map(function(item){
                            return item * 2;
                        });
                    }
                };
                var interceptor2 = {
                    response: function(data/*, params*/){
                        return data.map(function(item){
                            return item.toString() + '0';
                        });
                    }
                };
                var tp = createNgTableParams({ interceptors: [interceptor, interceptor2], getData: function(){
                    return [2, 3];
                }});

                // when
                var actualData;
                tp.reload().then(function(data){
                    actualData = data;
                });
                scope.$digest();

                // then
                expect(actualData).toEqual(['40', '60']);
            });
        });

        describe('multiple response and responseError interceptors', function(){

            it('responseError of next interceptor should receive failures from previous interceptor', inject(function($q){
                // given
                var badInterceptor = {
                    response: function(/*data, params*/){
                        return $q.reject('BANG!');
                    }
                };
                var nextInterceptor = {
                    responseError: function(/*reason, params*/){
                        this.hasRun = true;
                    }
                };
                var tp = createNgTableParams({ interceptors: [badInterceptor, nextInterceptor]});

                // when
                tp.reload();
                scope.$digest();

                // then
                expect(nextInterceptor.hasRun).toBe(true);
            }));

            it('should call next response interceptor when previous interceptor recovers from failure', inject(function($q){
                // given
                var badInterceptor = {
                    response: function(/*data, params*/){
                        return $q.reject('BANG!');
                    }
                };
                var recoveringInterceptor = {
                    responseError: function(/*reason, params*/){
                        return [8894,58];
                    }
                };
                var recoveredData;
                var nextInterceptor = {
                    response: function(data/*, params*/){
                        this.hasRun = true;
                        recoveredData = data;
                    }
                };
                var tp = createNgTableParams({ interceptors: [badInterceptor, recoveringInterceptor, nextInterceptor]});

                // when
                tp.reload();
                scope.$digest();

                // then
                expect(nextInterceptor.hasRun).toBe(true);
                expect(recoveredData).toEqual([8894,58]);
            }));
        });
    });

    describe('events', function(){

        var actualEventArgs,
            actualPublisher,
            fakeTableParams,
            ngTableEventsChannel;

        beforeEach(inject(function(_ngTableEventsChannel_){
            ngTableEventsChannel = _ngTableEventsChannel_;
            fakeTableParams = {};
            actualPublisher = undefined;
            actualEventArgs = undefined;
        }));

        function getSubscriberCount(){
            var allEventNames = Object.keys($rootScope.$$listenerCount);
            var ngTableEvents = allEventNames.filter(function(event){
                return event.indexOf('ngTable:') === 0;
            });
            return ngTableEvents.reduce(function(result, event){
                result += $rootScope.$$listenerCount[event];
                return result;
            }, 0);
        }

        describe('general pub/sub mechanics', function(){

            var supportedEvents = ['DatasetChanged', 'AfterReloadData', 'PagesChanged', 'AfterCreated'];

            it('should be safe to publish event when no subscribers', function () {

                function test(event) {
                    ngTableEventsChannel['publish' + event](fakeTableParams);
                }

                supportedEvents.forEach(test);
            });

            it('publishing event should notify registered subscribers (one)', function(){

                function test(event){
                    // given
                    var cbCount = 0;
                    ngTableEventsChannel['on' + event](function(){
                        cbCount++;
                    });

                    // when
                    ngTableEventsChannel['publish' + event](fakeTableParams);

                    // then
                    expect(cbCount).toBe(1);
                }

                supportedEvents.forEach(test);
            });

            it('publishing event should notify registered subscribers (multiple)', function(){

                function test(event){
                    // given
                    var cb1Count = 0;
                    ngTableEventsChannel['on' + event](function(){
                        cb1Count++;
                    });
                    var cb2Count = 0;
                    ngTableEventsChannel['on' + event](function(){
                        cb2Count++;
                    });

                    // when
                    ngTableEventsChannel['publish' + event](fakeTableParams);

                    // then
                    expect(cb1Count).toBe(1);
                    expect(cb2Count).toBe(1);
                }

                supportedEvents.forEach(test);
            });

            it('subscriber should be able to unregister their callback', function(){

                function test(event){
                    // given
                    var cbCount = 0;
                    var subscription = ngTableEventsChannel['on' + event](function(){
                        cbCount++;
                    });
                    ngTableEventsChannel['publish' + event](fakeTableParams);
                    expect(cbCount).toBe(1); // checking assumptions
                    expect(getSubscriberCount()).toBe(1); // checking assumptions
                    cbCount = 0; // reset

                    // when
                    subscription(); // unsubscribe
                    ngTableEventsChannel['publish' + event](fakeTableParams);

                    // then
                    expect(cbCount).toBe(0);
                    expect(getSubscriberCount()).toBe(0);
                }

                supportedEvents.forEach(test);
            });

            it('subscriber should be able specify the scope to receive events', function(){

                // this is useful as it allows all subscriptions to be removed by calling $destroy on that scope

                function test(event){
                    // given
                    var childScope = $rootScope.$new();
                    var cbCount = 0;
                    ngTableEventsChannel['on' + event](function(){
                        cbCount++;
                    }, childScope);

                    // when, then
                    ngTableEventsChannel['publish' + event](fakeTableParams);
                    expect(cbCount).toBe(1);

                    cbCount = 0; // reset

                    // when, then
                    childScope.$destroy();
                    ngTableEventsChannel['publish' + event](fakeTableParams);

                    // then
                    expect(cbCount).toBe(0);
                }

                supportedEvents.forEach(test);
            });

            it('should not notify subscribers who have filter out the publisher', function(){

                function test(event){
                    // given
                    var cbCount = 0;
                    ngTableEventsChannel['on' + event](function(){
                        cbCount++;
                    }, function(publisher){
                        return publisher === fakeTableParams;
                    });

                    // when
                    ngTableEventsChannel['publish' + event](fakeTableParams);
                    var anoParams = {};
                    ngTableEventsChannel['publish' + event](anoParams);

                    // then
                    expect(cbCount).toBe(1);
                }

                supportedEvents.forEach(test);
            });

            it('should not notify subscribers who have filter not to receive event based on arg values', function(){

                function test(event){
                    // given
                    var cbCount = 0;
                    ngTableEventsChannel['on' + event](function(){
                        cbCount++;
                    }, function(publisher, arg1){
                        return arg1 === 1;
                    });

                    // when
                    ngTableEventsChannel['publish' + event](fakeTableParams, 'cc');
                    ngTableEventsChannel['publish' + event](fakeTableParams, 1);

                    // then
                    expect(cbCount).toBe(1);
                }

                supportedEvents.forEach(test);
            });

            it('should support a shorthand for subscribers to receive events from specific NgTableParams instance', function () {

                function test(event) {
                    // given
                    var cbCount = 0;
                    ngTableEventsChannel['on' + event](function () {
                        cbCount++;
                    }, fakeTableParams);

                    // when
                    ngTableEventsChannel['publish' + event](fakeTableParams);
                    var anoParams = {};
                    ngTableEventsChannel['publish' + event](anoParams);

                    // then
                    expect(cbCount).toBe(1);
                }

                supportedEvents.forEach(test);
            });

            it('publisher should be supplied to subscriber callback', function(){

                function test(event){
                    // given
                    var cbCount = 0;
                    ngTableEventsChannel['on' + event](function(params){
                        actualPublisher = params;
                    });

                    // when
                    ngTableEventsChannel['publish' + event](fakeTableParams);

                    // then
                    expect(actualPublisher).toBe(fakeTableParams);
                }

                supportedEvents.forEach(test);
            });

            it('remaining event args should be supplied to subscriber callback', function(){

                function test(event){
                    // given
                    ngTableEventsChannel['on' + event](function(params/*, ...args*/){
                        actualPublisher = params;
                        actualEventArgs = _.rest(arguments);
                    });

                    // when
                    var arg1 = [1, 2];
                    var arg2 = [1];
                    ngTableEventsChannel['publish' + event](fakeTableParams, arg1, arg2);

                    // then
                    expect(actualEventArgs).toEqual([arg1, arg2]);
                }

                supportedEvents.forEach(test);
            });

            it('subscribers should never receive events from null instance of tableParams', function(){

                function test(event){
                    // given
                    var cbCount = 0;
                    ngTableEventsChannel['on' + event](function(params/*, ...args*/){
                        cbCount++;
                    });

                    // when
                    fakeTableParams.isNullInstance = true;
                    ngTableEventsChannel['publish' + event](fakeTableParams);

                    // then
                    expect(cbCount).toEqual(0);
                }

                supportedEvents.forEach(test);
            });

        });

        describe('afterCreated', function(){

            it('should fire when a new NgTableParams has been constructed', function(){
                // given
                ngTableEventsChannel.onAfterCreated(function(params){
                    actualPublisher = params;
                });

                // when
                var params = createNgTableParams();

                // then
                expect(actualPublisher).toBe(params);
            });
        });

        describe('afterReloadData', function(){

            it('should fire when a reload completes', function(){
                // given
                ngTableEventsChannel.onAfterReloadData(function(params, newVal, oldVal){
                    actualPublisher = params;
                    actualEventArgs = [newVal, oldVal];
                });
                var newDatapage = [1, 5, 6];
                var params = createNgTableParams({ getData: function(){
                    return newDatapage;
                }});

                // when
                params.reload();
                scope.$digest();

                // then
                expect(actualPublisher).toBe(params);
                expect(actualEventArgs).toEqual([newDatapage, []]);
            });

            it('should fire on reload even if datapage remains the same array', function(){
                // given
                var callCount = 0;
                ngTableEventsChannel.onAfterReloadData(function(/*params, newVal, oldVal*/){
                    callCount++;
                });
                var dataPage = [1,2,3];
                var params = createNgTableParams({ getData: function(){
                    return dataPage;
                }});

                // when
                params.reload();
                scope.$digest();
                params.reload();
                scope.$digest();

                // then
                expect(callCount).toBe(2);
            });

            it('should fire after afterCreated event', function(){
                // given
                var events = [];
                ngTableEventsChannel.onAfterReloadData(function(/*params, newVal, oldVal*/){
                    events.push('afterReloadData');
                });
                ngTableEventsChannel.onAfterCreated(function(/*params*/){
                    events.push('afterCreated');
                });

                // when
                var params = createNgTableParams({}, {dataset: [1,2,3,4,5,6]});
                params.reload();
                scope.$digest();

                // then
                expect(events[0]).toEqual('afterCreated');
                expect(events[1]).toEqual('afterReloadData');
            });

        });

        describe('pagesChanged', function(){

            it('should fire when a reload completes', function(){
                // given
                ngTableEventsChannel.onPagesChanged(function(params, newVal, oldVal){
                    actualPublisher = params;
                    actualEventArgs = [newVal, oldVal];
                });
                var params = createNgTableParams({ count: 5 }, { counts: [5,10], dataset: [1,2,3,4,5,6]});

                // when
                params.reload();
                scope.$digest();

                // then
                var expectedPages = params.generatePagesArray(params.page(), params.total(), params.count());
                expect(expectedPages.length).toBeGreaterThan(0); // checking assumptions
                expect(actualEventArgs).toEqual([expectedPages, undefined]);
            });

            it('should fire when a reload completes - no data', function(){
                // given
                ngTableEventsChannel.onPagesChanged(function(params, newVal, oldVal){
                    actualPublisher = params;
                    actualEventArgs = [newVal, oldVal];
                });
                var params = createNgTableParams({ count: 5 }, { counts: [5,10], dataset: []});

                // when
                params.reload();
                scope.$digest();

                // then
                expect(actualEventArgs).toEqual([[], undefined]);
            });

            it('should fire when a reload completes (multiple)', function(){
                // given
                var callCount = 0;
                ngTableEventsChannel.onPagesChanged(function(/*params, newVal, oldVal*/){
                    callCount++;
                });
                var params = createNgTableParams({ count: 5 }, { counts: [5,10], dataset: [1,2,3,4,5,6]});

                // when
                params.reload();
                scope.$digest();
                params.page(2); // trigger a change to pages data structure
                params.reload();
                scope.$digest();

                // then
                expect(callCount).toBe(2);
            });

            it('should not fire on reload when pages remain the same', function(){
                // given
                var callCount = 0;
                ngTableEventsChannel.onPagesChanged(function(/*params, newVal, oldVal*/){
                    callCount++;
                });
                var params = createNgTableParams({ count: 5 }, { counts: [5,10], dataset: [1,2,3,4,5,6]});
                params.reload();
                scope.$digest();

                // when
                params.reload();
                scope.$digest();

                // then
                expect(callCount).toBe(1);
            });

            it('should fire after afterCreated event', function(){
                // given
                var events = [];
                ngTableEventsChannel.onPagesChanged(function(/*params, newVal, oldVal*/){
                    events.push('pagesChanged');
                });
                ngTableEventsChannel.onAfterCreated(function(/*params*/){
                    events.push('afterCreated');
                });

                // when
                var params = createNgTableParams({ count: 5 }, { counts: [5,10], dataset: [1,2,3,4,5,6]});
                params.reload();
                scope.$digest();

                // then
                expect(events[0]).toEqual('afterCreated');
                expect(events[1]).toEqual('pagesChanged');
            });
        });

        describe('datasetChanged', function(){

            it('should fire when a initial dataset is supplied as a settings value', function(){
                // given
                ngTableEventsChannel.onDatasetChanged(function(params, newVal, oldVal){
                    actualPublisher = params;
                    actualEventArgs = [newVal, oldVal];
                });

                // when
                var initialDs = [5, 10];
                var params = createNgTableParams({ dataset: initialDs});

                // then
                expect(actualPublisher).toBe(params);
                expect(actualEventArgs).toEqual([initialDs, null]);
            });

            it('should fire when a new dataset is supplied as a settings value', function(){
                // given
                var initialDs = [1, 2];
                ngTableEventsChannel.onDatasetChanged(function(params, newVal, oldVal){
                    actualPublisher = params;
                    actualEventArgs = [newVal, oldVal];
                });
                var params = createNgTableParams({ dataset: initialDs});

                // when
                var newDs = [5, 10];
                params.settings({ dataset: newDs});

                // then
                expect(actualPublisher).toBe(params);
                expect(actualEventArgs).toEqual([newDs, initialDs]);
            });

            it('should fire when a dataset is removed from settings value', function(){
                // given
                var initialDs = [1, 2];
                ngTableEventsChannel.onDatasetChanged(function(params, newVal, oldVal){
                    actualPublisher = params;
                    actualEventArgs = [newVal, oldVal];
                });
                var params = createNgTableParams({ dataset: initialDs});

                // when
                var newDs = null;
                params.settings({ dataset: newDs});

                // then
                expect(actualPublisher).toBe(params);
                expect(actualEventArgs).toEqual([newDs, initialDs]);
            });

            it('should NOT fire when the same dataset array is supplied as a new settings value', function(){
                // given
                var callCount = 0;
                var initialDs = [1, 2];
                ngTableEventsChannel.onDatasetChanged(function(/*params, newVal, oldVal*/){
                    callCount++;
                });
                var params = createNgTableParams({ dataset: initialDs});

                // when
                params.settings({ dataset: initialDs});

                // then
                expect(callCount).toBe(1);
            });

            it('settings().dataset on publisher should reference the new dataset', function(){
                var initialDs = [1, 2];
                var newDs = [1, 2, 3];
                var params = createNgTableParams({ dataset: initialDs});
                ngTableEventsChannel.onDatasetChanged(function(params, newVal/*, oldVal*/){
                    expect(params.settings().dataset).toBe(newVal);
                });
                params.settings({ dataset: newDs});
            });

            it('should fire after afterCreated event', function(){
                // given
                var events = [];
                var initialDs = [1, 2];
                ngTableEventsChannel.onDatasetChanged(function(/*params, newVal, oldVal*/){
                    events.push('datasetChanged');
                });
                ngTableEventsChannel.onAfterCreated(function(/*params*/){
                    events.push('afterCreated');
                });

                // when
                var params = createNgTableParams({ dataset: initialDs});

                // then
                expect(events[0]).toEqual('afterCreated');
                expect(events[1]).toEqual('datasetChanged');
            });
        });
    })
});
