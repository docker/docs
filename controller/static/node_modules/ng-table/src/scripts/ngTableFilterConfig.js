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
        .provider('ngTableFilterConfig', ngTableFilterConfigProvider);

    ngTableFilterConfigProvider.$inject = [];

    function ngTableFilterConfigProvider(){
        var config;
        var defaultConfig = {
            defaultBaseUrl: 'ng-table/filters/',
            defaultExt: '.html',
            aliasUrls: {}
        };

        this.$get = ngTableFilterConfig;
        this.resetConfigs = resetConfigs;
        this.setConfig = setConfig;

        init();

        /////////

        function init(){
            resetConfigs();
        }

        function resetConfigs(){
            config = defaultConfig;
        }

        function setConfig(customConfig){
            var mergeConfig = angular.extend({}, config, customConfig);
            mergeConfig.aliasUrls = angular.extend({}, config.aliasUrls, customConfig.aliasUrls);
            config = mergeConfig;
        }

        /////////

        ngTableFilterConfig.$inject = [];

        function ngTableFilterConfig(){

            var publicConfig;

            var service = {
                config: publicConfig,
                getTemplateUrl: getTemplateUrl,
                getUrlForAlias: getUrlForAlias
            };
            Object.defineProperty(service, "config", {
                get: function(){
                    return publicConfig = publicConfig || angular.copy(config);
                },
                enumerable: true
            });

            return service;

            /////////

            function getTemplateUrl(filterValue, filterKey){
                if (angular.isObject(filterValue)){
                    filterValue = filterValue.id;
                }
                if (filterValue.indexOf('/') !== -1){
                    return filterValue;
                }

                return service.getUrlForAlias(filterValue, filterKey);
            }

            function getUrlForAlias(aliasName/*, filterKey*/){
                return config.aliasUrls[aliasName] || config.defaultBaseUrl + aliasName + config.defaultExt;
            }
        }
    }
})();
