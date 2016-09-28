var gulp = require('gulp');
require('gulp-grunt')(gulp); // add all the gruntfile tasks to gulp

var Dgeni = require('dgeni');
var dgeni_config = require('./docs/config');
var webserver = require("gulp-webserver");
var browserify = require('gulp-browserify');
var sequence = require('gulp-sequence');
var wait = require('gulp-wait');
var ghPages = require('gulp-gh-pages');

//serverStream
var serverStream;

gulp.task('build:doc-assets', function () {
    gulp.src(['./docs/app/**/*']).pipe(gulp.dest('./dist/docs/app'))
});

gulp.task('compile:doc-src', function () {
    gulp.src('dist/docs/app/js/app.js')
        .pipe(browserify({insertGlobals: true}))
        .pipe(gulp.dest('dist/docs/app/build'));
});

//Include all needed vendor CSS require by vendor libs
gulp.task('vendor-css', function () {
    return gulp.src([
        './app/vendor/angular-material/angular-material.css',
        './Content/font-awesome/css/*.css'
    ]).pipe(gulp.dest('./docs/release/app/vendor/css'));
});

gulp.task('run:dgeni', function () {
    var dgeni = new Dgeni([dgeni_config]);
    return dgeni.generate().then(function (docs) {
        console.log(docs.length, 'docs generated');
    })
});

gulp.task('run:server', function () {
    serverStream = gulp.src('./dist/docs')
        .pipe(webserver({
                port: 3030,
                livereload:true,
                open:true
        }));
});

gulp.task('watch:files', function () {
    gulp.watch(['docs/**/*.*'], ['build:doc-assets', 'run:dgeni', 'compile:doc-src'])
});

gulp.task('kill-server', function () {
    if(serverStream) serverStream.emit('kill');
});

gulp.task('deploy:docs', function(){
    return gulp.src('./dist/docs/**/*').pipe(ghPages(
        {
            remoteUrl: 'https://github.com/bigmeech/ng-table.git',
            branch:'gh-pages',
            origin: 'origin',
            push:true,
            message:"Update" + new Date()
        }
    ))
});

gulp.task('dev', function (callback) {
    sequence('grunt-default', 'build:doc-assets', 'run:dgeni', 'compile:doc-src', 'kill-server', 'run:server', 'watch:files')(callback);
});

gulp.task('default', function (callback) {
    sequence('grunt-default', 'build:doc-assets', 'run:dgeni', 'compile:doc-src')(callback);
});
