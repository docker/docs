var fs = require('fs');
var assert = require('assert');
/**
 * Writes the stats object out to /stats.json so we can fetch it from
 * the appropriate container.
 *
 * This fails hard (throws an AssertionError) if the filename does
 * not match the regex. This is good because we don't want it shipping
 * to production if we can't reliably find the hash for client.js.
 *
 * The `CLIENT_JS_REGEX` should match `output.filename` in the webpack
 * config below.
 */
module.exports = function() {
  var CLIENT_JS_REGEX = /client..*.js$/;
  this.plugin('done', function(stats) {
    var filename = stats.toJson().assets[0].name;

    console.log(filename);
    console.log(stats.toJson().assets[0]);
    // write stats object
    fs.writeFileSync(
      '/tmp/dux-stats-client-js.json',
      JSON.stringify(stats.toJson()));

    // test `filename`
    assert.strictEqual(true, !!filename.match(CLIENT_JS_REGEX), filename + ' does not match expected client.js regex')

    /**
     * Store the client.js hash in /tmp/.clientjs-hash for use
     * in the server build process
     */
    fs.writeFile('/tmp/.client-js-hash', filename, 'utf8', function(err) {
      if (err) throw err;
    });
  });
}