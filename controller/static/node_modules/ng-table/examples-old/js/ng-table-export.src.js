angular.module('ngTableExport', [])
.config(['$compileProvider', function($compileProvider) {
    // allow data links
    $compileProvider.aHrefSanitizationWhitelist(/^\s*(https?|ftp|mailto|data):/);
}])
.directive('exportCsv', ['$parse', function ($parse) {
    return {
        restrict: 'A',
        scope: false,
        link: function(scope, element, attrs) {
            var data = '';
            var csv = {
                stringify: function(str) {
                    return '"' +
                        str.replace(/^\s\s*/, '').replace(/\s*\s$/, '') // trim spaces
                            .replace(/"/g,'""') + // replace quotes with double quotes
                        '"';
                },
                generate: function() {
                    data = '';
                    var rows = element.find('tr');
                    angular.forEach(rows, function(row, i) {
                        var tr = angular.element(row),
                            tds = tr.find('th'),
                            rowData = '';
                        if (tr.hasClass('ng-table-filters')) {
                            return;
                        }
                        if (tds.length == 0) {
                            tds = tr.find('td');
                        }
                        if (i != 1) {
                            angular.forEach(tds, function(td, i) {
                                rowData += csv.stringify(angular.element(td).text()) + ';';
                            });
                            rowData = rowData.slice(0, rowData.length - 1); //remove last semicolon
                        }
                        data += rowData + "\n";
                    });
                },
                link: function() {
                    return 'data:text/csv;charset=UTF-8,' + encodeURIComponent(data);
                }
            };
            $parse(attrs.exportCsv).assign(scope.$parent, csv);
        }
    };
}]);