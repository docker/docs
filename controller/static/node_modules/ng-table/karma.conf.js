// Karma configuration file
// See http://karma-runner.github.io/0.10/config/configuration-file.html
module.exports = function (config) {
    config.set({
        basePath: '',

        frameworks: ['jasmine'],

        // list of files / patterns to load in the browser
        files: [
            // libraries
            'bower_components/lodash/lodash.js',
            'bower_components/angular/angular.js',
            'bower_components/angular-mocks/angular-mocks.js',

            // directive
            './dist/ng-table.js',

            // tests
            'test/*.js'
            //'test/tableParamsSpec.js'
            //'test/tableControllerSpec.js'
        ],

        // generate js files from html templates
        preprocessors: {
            '*.js': 'coverage'
        },

        reporters: ['progress', 'coverage'],

        autoWatch: true,
        browsers: ['Chrome'],
        coverageReporter: {
            type: 'lcov',
            dir: 'out/coverage'
        }
    });
};
