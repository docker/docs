describe('ngTableDefaultGetData', function () {

    describe('provider', function () {
        var ngTableDefaultGetDataProvider;

        beforeEach(module("ngTable"));
        beforeEach(function(){
            module(function(_ngTableDefaultGetDataProvider_){
                ngTableDefaultGetDataProvider=_ngTableDefaultGetDataProvider_;
            });
        });
        beforeEach(inject());

        it('should be configured to use built-in angular filters', function () {
            expect(ngTableDefaultGetDataProvider.filterFilterName).toBe('filter');
            expect(ngTableDefaultGetDataProvider.sortingFilterName).toBe('orderBy');
        });
    });

    describe('service', function () {
        var ngTableDefaultGetData,
            tableParams;

        beforeEach(module('ngTable'));

        beforeEach(inject(function (_ngTableDefaultGetData_, NgTableParams) {
            ngTableDefaultGetData = _ngTableDefaultGetData_;
            tableParams = new NgTableParams({count: 10}, {counts: [10]});
        }));

        describe('sorting', function () {
            it('empty sorting', function () {
                // given
                tableParams.sorting({});
                // when
                var actualResults = ngTableDefaultGetData([{age: 1}, {age: 2}, {age: 3}], tableParams);
                // then
                expect(actualResults).toEqual([{age: 1}, {age: 2}, {age: 3}]);
            });

            it('single property sort ascending', function () {
                // given
                tableParams.sorting({age: 'asc'});
                // when
                var actualResults = ngTableDefaultGetData([{age: 1}, {age: 3}, {age: 2}], tableParams);
                // then
                expect(actualResults).toEqual([{age: 1}, {age: 2}, {age: 3}]);
            });

            it('single property sort descending', function () {
                // given
                tableParams.sorting({age: 'desc'});
                // when
                var actualResults = ngTableDefaultGetData([{age: 1}, {age: 3}, {age: 2}], tableParams);
                // then
                expect(actualResults).toEqual([{age: 3}, {age: 2}, {age: 1}]);
            });

            it('multiple property sort ascending', function () {
                // given
                tableParams.sorting({age: 'asc', name: 'asc'});
                // when
                var data = [
                    {age: 1, name: 'b'}, {age: 3, name: 'a'}, {age: 2, name: 'a'}, {age: 1, name: 'a'}
                ];
                var actualResults = ngTableDefaultGetData(data, tableParams);
                // then
                var expectedData = [
                    {age: 1, name: 'a'}, {age: 1, name: 'b'}, {age: 2, name: 'a'}, {age: 3, name: 'a'}
                ];
                expect(actualResults).toEqual(expectedData);
            });

        });

        describe('filters', function () {
            it('empty filter', function () {
                // given
                tableParams.filter({});
                // when
                var actualResults = ngTableDefaultGetData([{age: 1}, {age: 2}, {age: 3}], tableParams);
                // then
                expect(actualResults).toEqual([{age: 1}, {age: 2}, {age: 3}]);
            });

            it('empty filter - simple values', function () {
                // given
                tableParams.filter({});
                // when
                var actualResults = ngTableDefaultGetData([1, 2, 3], tableParams);
                // then
                expect(actualResults).toEqual([1, 2, 3]);
            });

            it('single property filter', function () {
                // given
                tableParams.filter({age: 1});
                // when
                var actualResults = ngTableDefaultGetData([{age: 1}, {age: 2}, {age: 3}], tableParams);
                // then
                expect(actualResults).toEqual([{age: 1}]);
            });

            it('multiple property filter', function () {
                // given
                var data = [{age: 1, name: 'A'}, {age: 2, name: 'B'}, {age: 3, name: 'B'}];
                tableParams.filter({age: 2, name: 'B'});
                // when
                var actualResults = ngTableDefaultGetData(data, tableParams);
                // then
                expect(actualResults).toEqual([{age: 2, name: 'B'}]);
            });

            it('should remove null and undefined values before applying', function () {
                // given
                var data = [{age: 1, name: 'A'}, {age: 2, name: 'B'}, {age: 3, name: 'B'}];
                tableParams.filter({age: null, name: 'B'});
                // when
                var actualResults = ngTableDefaultGetData(data, tableParams);
                // then
                expect(actualResults).toEqual([{age: 2, name: 'B'}, {age: 3, name: 'B'}]);
            });

            it('should remove empty string value before applying', function () {
                // given
                var data = [{age: 1, name: 'A'}, {age: 2, name: 'B'}, {age: 3, name: 'B'}];
                tableParams.filter({age: 2, name: ''});
                // when
                var actualResults = ngTableDefaultGetData(data, tableParams);
                // then
                expect(actualResults).toEqual([{age: 2, name: 'B'}]);
            });

            it('single nested property, one level deep', function () {
                // given
                tableParams.filter({'details.age': 1});
                // when
                var data = [{
                    details: {age: 1}
                }, {
                    details: {age: 2}
                }, {
                    details: {age: 3}
                }, {
                    age: 1
                }];
                var actualResults = ngTableDefaultGetData(data, tableParams);
                // then
                expect(actualResults).toEqual([{details: {age: 1}}]);
            });

            it('single nested property, two levels deep', function () {
                // given
                tableParams.filter({'details.personal.age': 1});
                // when
                var data = [{
                    details: {personal: {age: 1}}
                }, {
                    details: {personal: {age: 2}}
                }, {
                    details: {personal: {age: 3}}
                }, {
                    age: 1
                }];
                var actualResults = ngTableDefaultGetData(data, tableParams);
                // then
                expect(actualResults).toEqual([{details: {personal: {age: 1}}}]);
            });

            it('multiple nested property, two levels deep', function () {
                // given
                tableParams.filter({'details.personal.age': 1, 'job.money': 100});
                // when
                var data = [{
                    details: {personal: {age: 1}},
                    job: {money: 200}
                }, {
                    details: {personal: {age: 1}},
                    job: {money: 100}
                }, {
                    details: {personal: {age: 3}},
                    job: {money: 100}
                }];
                var actualResults = ngTableDefaultGetData(data, tableParams);
                // then
                var expected = [{
                    details: {personal: {age: 1}},
                    job: {money: 100}
                }];
                expect(actualResults).toEqual(expected);
            });

            describe('filterComparator', function(){

                it('function', function () {
                    // given
                    var comparer = function (actual, expected) {
                        return angular.equals(actual, expected);
                    };
                    tableParams.settings({ filterOptions: { filterComparator: comparer }});
                    tableParams.filter({age: 1});
                    // when
                    var actualResults = ngTableDefaultGetData([{age: 10}, {age: 1}, {age: 101}, {age: 2}], tableParams);
                    // then
                    expect(actualResults).toEqual([{age: 1}]);
                });

                it('"true"', function () {
                    // given
                    tableParams.settings({ filterOptions: { filterComparator: true }});
                    tableParams.filter({age: 1});
                    // when
                    var actualResults = ngTableDefaultGetData([{age: 10}, {age: 1}, {age: 101}, {age: 2}], tableParams);
                    // then
                    expect(actualResults).toEqual([{age: 1}]);
                });

                it('"false"', function () {
                    // given
                    tableParams.settings({ filterOptions: { filterComparator: false }});
                    tableParams.filter({age: 1});
                    // when
                    var actualResults = ngTableDefaultGetData([{age: 10}, {age: 1}, {age: 101}, {age: 2}], tableParams);
                    // then
                    expect(actualResults).toEqual([{age: 10}, {age: 1}, {age: 101}]);
                });

                it('"undefined" (the default)', function () {
                    // given
                    tableParams.settings({ filterOptions: { filterComparator: undefined }});
                    tableParams.filter({age: 1});
                    // when
                    var actualResults = ngTableDefaultGetData([{age: 10}, {age: 1}, {age: 101}, {age: 2}], tableParams);
                    // then
                    expect(actualResults).toEqual([{age: 10}, {age: 1}, {age: 101}]);
                });
            })
        });
    });

    describe('service, custom filters', function () {
        var ngTableDefaultGetData, tableParams;

        beforeEach(function () {
            // add a custom filter available to our tests
            module('ngTable', function($provide){
                $provide.factory('myCustomFilterFilter', myCustomFilter)
            });

            myCustomFilter.$inject = [];
            function myCustomFilter() {

                return jasmine.createSpy('myCustomFilterSpy',filter).and.callThrough();

                function filter(people, criteriaObj/*, comparator*/) {
                    return people.filter(function (p) {
                        return criteriaObj.ages.indexOf(p.age) !== -1;
                    });
                }
            }
        });

        beforeEach(inject(function (_ngTableDefaultGetData_, NgTableParams) {
            ngTableDefaultGetData = _ngTableDefaultGetData_;
            tableParams = new NgTableParams({count: 10}, {counts: [10]});
        }));

        it('filterFilterName override', function () {
            // given
            tableParams.settings({ filterOptions: {filterFilterName: 'myCustomFilter'}});
            tableParams.filter({ages: [1, 2]});
            // when
            var actualResults = ngTableDefaultGetData([{age: 1}, {age: 2}, {age: 3}], tableParams);
            // then
            expect(actualResults).toEqual([{age: 1}, {age: 2}]);
        });

        it('`this` context of custom filter should be set to the NgTableParams instance', inject(function (myCustomFilterFilter) {
            // given
            tableParams.settings({ filterOptions: {filterFilterName: 'myCustomFilter'}});
            tableParams.filter({ages: [1, 2]});
            // when
            ngTableDefaultGetData([{age: 1}, {age: 2}, {age: 3}], tableParams);
            // then
            expect(myCustomFilterFilter.calls.mostRecent().object).toBe(tableParams);
        }));

        it('custom filter function', function () {
            // given
            var filterFn = function (data, criteriaObj/*, comparator*/) {
                return data.filter(function (p) {
                    return criteriaObj.ages.indexOf(p.age) !== -1;
                });
            };
            tableParams.settings({ filterOptions: { filterFn: filterFn }});
            tableParams.filter({ages: [1, 2]});
            // when
            var actualResults = ngTableDefaultGetData([{age: 1}, {age: 2}, {age: 3}], tableParams);
            // then
            expect(actualResults).toEqual([{age: 1}, {age: 2}]);
        });

        it('`this` context of custom filter function should be set to the NgTableParams instance', function () {
            // given
            var filterFnSpy = jasmine.createSpy('filterFn', angular.identity).and.callThrough();
            tableParams.settings({ filterOptions: { filterFn: filterFnSpy }});
            tableParams.filter({age: 1});
            // when
            var actualResults = ngTableDefaultGetData([{age: 1}], tableParams);
            // then
            expect(filterFnSpy.calls.mostRecent().object).toBe(tableParams);
        });
    });
});
