// jest.unmock('Marketplace/lib/index.js');
jest.autoMockOff();
const { DEFAULT_SEARCH_PAGE_SIZE } = require('lib/constants/defaults');
const routes = require('lib/constants/routes');
const {
  formatNumber,
  formatBucketedNumber,
} = require('lib/utils/formatNumbers');
const mkRepoTitle = require('lib/utils/repo-image-name');

describe('marketplace lib consts and helper functions', () => {
  it('should have a DEFAULT_SEARCH_PAGE_SIZE of 12', () => {
    expect(DEFAULT_SEARCH_PAGE_SIZE).toEqual(12);
  });

  it('should appropriately generate routes', () => {
    const { search, detail } = routes;
    const namespace = 'library';
    const reponame = 'redis';
    const searchRoute = '/search';
    const detailRoute = `/images/${namespace}/${reponame}`;
    expect(search).toBeDefined();
    expect(detail).toBeDefined();
    expect(search()).toEqual(searchRoute);
    expect(detail({ namespace, reponame })).toEqual(detailRoute);
  });

  it('should format numbers with formatNumber', () => {
    expect(formatNumber(100)).toEqual('100');
    expect(formatNumber(1003)).toEqual('1.0K');
  });

  it('should format numbers with formatBucketedNumber', () => {
    expect(formatBucketedNumber(2000)).toEqual('2.0K');
    expect(formatBucketedNumber(49999)).toEqual('10K+');
    expect(formatBucketedNumber(99999)).toEqual('50K+');
    expect(formatBucketedNumber(499999)).toEqual('100K+');
    expect(formatBucketedNumber(999999)).toEqual('500K+');
    expect(formatBucketedNumber(4999999)).toEqual('1M+');
    expect(formatBucketedNumber(9999999)).toEqual('5M+');
    expect(formatBucketedNumber(10000001)).toEqual('10M+');
    expect(formatBucketedNumber(undefined)).toEqual(undefined);
    expect(formatBucketedNumber('abc')).toEqual(undefined);
  });

  it('should format reponames correctly with mkRepoTitle', () => {
    let namespace = 'library';
    let reponame = 'redis';
    expect(mkRepoTitle({ namespace, reponame })).toEqual(reponame);
    namespace = 'dockerWoo';
    expect(mkRepoTitle({ namespace, reponame }))
      .toEqual(`${namespace}/${reponame}`);
    reponame = 'lal';
    expect(mkRepoTitle({ reponame }))
      .toEqual('');
  });
});
