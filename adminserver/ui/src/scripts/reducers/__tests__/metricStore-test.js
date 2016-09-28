'use strict';

jest.dontMock('../metricStore');

// We can't use ES6 style imports with Jest as they'll be mocked.
// Even if you tell it not to.
// For real.
//
// ♫♬
// Jest has auto mock
// It really does help us out
// Barring ES6
MetricStore = require('../metricStore');

describe('MetricStore', function() {

  describe('url', function() {
    it('returns the correct URL with date appended', function() {
      var url = '/api/v0/admin/metricsData?lastTime=' + new Date(0).toISOString();
      expect(MetricStore.url).toBe(url);
    });
  });

  describe('staleness', function() {
    it('detects stale methods properly', function() {
      var ancient = Date.parse('2011-10-10T14:48:00');
      expect(MetricStore.isStale(ancient)).toBe(true);

      var overFiveMinutesAgo = Date.now() - 5001;
      expect(MetricStore.isStale(overFiveMinutesAgo)).toBe(true);

      var nowAsDate = new Date();
      expect(MetricStore.isStale(nowAsDate)).toBe(false);

      var nowAsTimestamp = Date.now();
      expect(MetricStore.isStale(nowAsTimestamp)).toBe(false);
    });

    it('allows us to set stale durations', function() {
      MetricStore.interval = 100;
      var overFiveMinutesAgo = Date.now() - 101;
      expect(MetricStore.isStale(overFiveMinutesAgo)).toBe(true);
    });
  });
});
