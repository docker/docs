var gulp = require('gulp');
var imagemin = require('gulp-imagemin');
var pngquant = require('imagemin-pngquant');

//Hub2 Images for dev & production (There is a separate task for docker-ux images)
gulp.task('images::dev', function () {
  return gulp.src('app/img/**')
         .pipe(imagemin({
           progressive: true,
           svgoPlugins: [{removeViewBox: false}],
           use: [pngquant({ quality: '65-80', speed: 4 })]
         }))
         .pipe(gulp.dest('app/.build/public/img'));
});

gulp.task('images::prod', function() {
  return gulp.src('app/img/**')
         .pipe(imagemin({
           progressive: true,
           svgoPlugins: [{removeViewBox: false}],
           use: [pngquant({ quality: '65-80', speed: 4 })]
         }))
         .pipe(gulp.dest('.tmp/server/build/img'));
});
