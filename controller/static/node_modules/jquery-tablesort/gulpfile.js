var gulp = require('gulp');
var uglify = require("gulp-uglify");
var rename = require('gulp-rename');


gulp.task('min', function() {
  gulp.src('jquery.tablesort.js')
    .pipe(uglify({preserveComments: 'license'}))
    .pipe(rename('jquery.tablesort.min.js'))
    .pipe(gulp.dest('.'));
});
