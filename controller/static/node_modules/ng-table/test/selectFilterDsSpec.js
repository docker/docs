describe('ngTableSelectFilterDs directive', function(){

    var $scope,
        elem,
        $compile;
    beforeEach(module('ngTable'));

    beforeEach(inject(function($rootScope, _$compile_){
        $scope = $rootScope.$new();
        $compile = _$compile_;
        elem = '<select ng-table-select-filter-ds="$column"></select>';
    }));

    describe('array datasource', function(){

        it('should add array to current scope', function(){
            // given
            var data = [{id: 1, title: 'A'}];
            $scope.$column = {
                data: data
            };
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect($scope.$selectData).toBe(data);
        });

        it('should turn null/undefined array into empty array', function(){
            // given
            $scope.$column = {
                data: null
            };
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect($scope.$selectData).toEqual([]);
        });

        it('should keep the array on scope in sync with data array on $column', function(){
            // given
            $scope.$column = {
                data: null
            };
            $compile(elem)($scope);
            $scope.$digest();

            // when
            var newArray = [{id: 1, title: 'A'}];
            $scope.$column.data = newArray;
            $scope.$digest();

            // then
            expect($scope.$selectData).toBe(newArray);
        });

        it('should add empty option to array', function(){
            // note: modifying the array supplied is not great as this can cause unexpected side effects
            // however, it does mean that a consumer can update the array and have this reflected in the select list
            // and so therefore it increases the utility of this directive

            // given
            var data = [{id: 1, title: 'A'}];
            $scope.$column = {
                data: data
            };
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect(data).toEqual([{ id: '', title: ''}, {id: 1, title: 'A'}]);
        });

        it('should add empty option to empty array', function(){
            // this is useful as it allows for app to add items to array at a future date and still
            // allow for the user to select an empty select option thus removing column filter

            // given
            var data = [];
            $scope.$column = {
                data: data
            };
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect(data).toEqual([{ id: '', title: ''}]);
        });

        it('should not add empty option to array if already present', function(){

            // given
            var data = [{id: '', title: ''}, {id: 1, title: 'A'}];
            $scope.$column = {
                data: data
            };
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect(data).toEqual([{id: '', title: ''}, {id: 1, title: 'A'}]);
        });

        it('should add empty option to a new arriving array', function(){
            // given
            $scope.$column = {
                data: [{id: 1, title: 'A'}]
            };
            $compile(elem)($scope);
            $scope.$digest();

            // when
            $scope.$column.data = [{id: 1, title: 'B'}];
            $scope.$digest();

            // then
            expect($scope.$selectData).toEqual([{ id: '', title: ''}, {id: 1, title: 'B'}]);
        });

    });

    describe('function datasource', function(){

        var data;
        beforeEach(function(){
            $scope.$column = {
                data: function(){
                    return data;
                }
            };
        });

        it('should add array to current scope', function(){
            // given
            data = [{id: 1, title: 'A'}];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect($scope.$selectData).toBe(data);
        });

        it('should turn null/undefined array into empty array', function(){
            // given
            data = null;
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect($scope.$selectData).toEqual([]);
        });

        it('should keep the array on scope in sync with data array on $column', function(){
            // given
            data = [{id: 1, title: 'A'}];
            $compile(elem)($scope);
            $scope.$digest();

            // when
            var newArray = [{id: 1, title: 'A'}];
            $scope.$column.data = function(){
                return newArray;
            };
            $scope.$digest();

            // then
            expect($scope.$selectData).toBe(newArray);
        });


        it('should add empty option to array', function(){
            // given
            data = [{id: 1, title: 'A'}];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect(data).toEqual([{ id: '', title: ''}, {id: 1, title: 'A'}]);
        });

        it('should add empty option to empty array', function(){
            // this is useful as it allows for app to add items to array at a future date and still
            // allow for the user to select an empty select option thus removing column filter

            // given
            data = [];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect(data).toEqual([{ id: '', title: ''}]);
        });

        it('should not add empty option to array if already present', function(){
            // given
            data = [{id: '', title: ''}, {id: 1, title: 'A'}];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            // then
            expect(data).toEqual([{id: '', title: ''}, {id: 1, title: 'A'}]);
        });

        it('should add empty option to a new arriving array', function(){
            // given
            data = [{id: 1, title: 'A'}];
            $compile(elem)($scope);
            $scope.$digest();

            // when
            $scope.$column.data = function(){
                return [{id: 1, title: 'B'}];
            };
            $scope.$digest();

            // then
            expect($scope.$selectData).toEqual([{ id: '', title: ''}, {id: 1, title: 'B'}]);
        });
    });

    describe('asyn function datasource', function(){
        var data;
        var $timeout;
        beforeEach(inject(function(_$timeout_){
            $timeout = _$timeout_;
            $scope.$column = {
                data: function(){
                    return $timeout(function(){
                        return data;
                    }, 10);
                }
            };
        }));

        it('should add array to current scope', function(){
            // given
            data = [{id: 1, title: 'A'}];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            $timeout.flush();
            // then
            expect($scope.$selectData).toBe(data);
        });

        it('should turn null/undefined array into empty array', function(){
            // given
            data = null;
            // when
            $compile(elem)($scope);
            $scope.$digest();
            $timeout.flush();
            // then
            expect($scope.$selectData).toEqual([]);
        });

        it('should keep the array on scope in sync with data array on $column', function(){
            // given
            data = [{id: 1, title: 'A'}];
            $compile(elem)($scope);
            $scope.$digest();
            $timeout.flush();

            // when
            var newArray = [{id: 1, title: 'A'}];
            $scope.$column.data = function(){
                return $timeout(function(){
                    return newArray;
                }, 10);
            };
            $scope.$digest();
            $timeout.flush();

            // then
            expect($scope.$selectData).toBe(newArray);
        });


        it('should add empty option to array', function(){
            // given
            data = [{id: 1, title: 'A'}];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            $timeout.flush();
            // then
            expect(data).toEqual([{ id: '', title: ''}, {id: 1, title: 'A'}]);
        });

        it('should add empty option to empty array', function(){
            // this is useful as it allows for app to add items to array at a future date and still
            // allow for the user to select an empty select option thus removing column filter

            // given
            data = [];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            $timeout.flush();
            // then
            expect(data).toEqual([{ id: '', title: ''}]);
        });

        it('should not add empty option to array if already present', function(){
            // given
            data = [{id: '', title: ''}, {id: 1, title: 'A'}];
            // when
            $compile(elem)($scope);
            $scope.$digest();
            $timeout.flush();
            // then
            expect(data).toEqual([{id: '', title: ''}, {id: 1, title: 'A'}]);
        });

        it('should add empty option to a new arriving array', function(){
            // given
            data = [{id: 1, title: 'A'}];
            $compile(elem)($scope);
            $scope.$digest();
            $timeout.flush();

            // when
            $scope.$column.data = function(){
                return $timeout(function(){
                    return [{id: 1, title: 'B'}];
                }, 10);
            };
            $scope.$digest();
            $timeout.flush();

            // then
            expect($scope.$selectData).toEqual([{ id: '', title: ''}, {id: 1, title: 'B'}]);
        });

    });
});
