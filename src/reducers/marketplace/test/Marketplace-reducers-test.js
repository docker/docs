import marketplaceReducer,
  { DEFAULT_STATE } from 'reducers/marketplace';
import {
  MARKETPLACE_FETCH_BUNDLE_DETAIL,
  MARKETPLACE_FETCH_BUNDLE_SUMMARY,
  MARKETPLACE_FETCH_CATEGORIES,
  MARKETPLACE_FETCH_PLATFORMS,
  MARKETPLACE_FETCH_REPOSITORY_DETAIL,
  MARKETPLACE_FETCH_REPOSITORY_SUMMARY,
  MARKETPLACE_SEARCH,
} from 'actions/marketplace';
import {
  REPOSITORY_FETCH_COMMENTS,
  REPOSITORY_FETCH_IMAGE_DETAIL,
  REPOSITORY_FETCH_IMAGE_TAGS,
} from 'actions/repository';
import {
  NAUTILUS_FETCH_SCAN_DETAIL,
  NAUTILUS_FETCH_TAGS_AND_SCANS,
} from 'actions/nautilus';
import makeRepositoryId from 'lib/utils/repo-image-name';
import { expect } from 'chai';

describe('marketplace reducer', () => {
  const namespace = 'library';
  const reponame = 'ubuntu';
  const id = '1234';
  const repositoryId = makeRepositoryId({ namespace, reponame });
  const tag = 'xenial';
  const page = 1;
  const page_size = 10;
  const status = '404';

  it('should return the initial state', () => {
    expect(marketplaceReducer(undefined, {})).to.deep.equal(DEFAULT_STATE);
  });

  //----------------------------------------------------------------------------
  // MARKETPLACE_FETCH_CATEGORIES
  //----------------------------------------------------------------------------
  it('should handle MARKETPLACE_FETCH_CATEGORIES_ACK', () => {
    const categories = [{ name: 'abc', label: 'ABC' }];
    const action = {
      type: `${MARKETPLACE_FETCH_CATEGORIES}_ACK`,
      payload: categories,
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.filters.categories).to.deep.equal({ abc: 'ABC' });
  });

  //----------------------------------------------------------------------------
  // MARKETPLACE_FETCH_PLATFORMS
  //----------------------------------------------------------------------------
  it('should handle MARKETPLACE_FETCH_PLATFORMS_ACK', () => {
    const platforms = [{ name: 'abc', label: 'ABC' }];
    const action = {
      type: `${MARKETPLACE_FETCH_PLATFORMS}_ACK`,
      payload: platforms,
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.filters.platforms).to.deep.equal({ abc: 'ABC' });
  });

  //----------------------------------------------------------------------------
  // MARKETPLACE_SEARCH
  //----------------------------------------------------------------------------
  it('should handle MARKETPLACE_SEARCH_REQ', () => {
    const action = { type: `${MARKETPLACE_SEARCH}_REQ` };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.search.isFetching).to.equal(true);
  });

  it('should handle MARKETPLACE_SEARCH_ACK', () => {
    const summaries = [{ id }];
    const action = {
      type: `${MARKETPLACE_SEARCH}_ACK`,
      payload: {
        count: 1,
        summaries,
      },
    };

    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.search.isFetching).to.equal(false);
    expect(reducer.search.count).to.equal(1);
    expect(reducer.search.pages).to.exist;
    expect(reducer.search.pages).to.deep.equal({
      1: {
        results: summaries,
      },
    });
  });

  it('should handle MARKETPLACE_SEARCH_ERR', () => {
    const action = {
      type: `${MARKETPLACE_SEARCH}_ERR`,
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.search.isFetching).to.equal(false);
    expect(reducer.search.error).to.equal(status);
  });

  //----------------------------------------------------------------------------
  // MARKETPLACE_FETCH_REPOSITORY_DETAIL
  //----------------------------------------------------------------------------
  it('should handle MARKETPLACE_FETCH_REPOSITORY_DETAIL_REQ', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_REPOSITORY_DETAIL}_REQ`,
      meta: { id },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].isFetching).to.equal(true);
  });

  it('should handle MARKETPLACE_FETCH_REPOSITORY_DETAIL_ACK', () => {
    const payload = {
      id,
      short_description: 'Busybox base image.',
    };
    const action = {
      type: `${MARKETPLACE_FETCH_REPOSITORY_DETAIL}_ACK`,
      meta: { id },
      payload,
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id]).to.deep.equal({
      ...payload, isFetching: false,
    });
  });

  it('should handle MARKETPLACE_FETCH_REPOSITORY_DETAIL_ERR', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_REPOSITORY_DETAIL}_ERR`,
      meta: { id },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.certified[id] = { isFetching: true };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].error)
      .to.equal(status);
  });

  //----------------------------------------------------------------------------
  // MARKETPLACE_FETCH_REPOSITORY_SUMMARY
  //----------------------------------------------------------------------------
  it('should handle MARKETPLACE_FETCH_REPOSITORY_SUMMARY_REQ', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_REPOSITORY_SUMMARY}_REQ`,
      meta: { id },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].isFetching).to.equal(true);
  });

  it('should handle MARKETPLACE_FETCH_REPOSITORY_SUMMARY_ACK', () => {
    const payload = {
      id,
      short_description: 'Busybox base image.',
    };
    const action = {
      type: `${MARKETPLACE_FETCH_REPOSITORY_SUMMARY}_ACK`,
      meta: { id },
      payload,
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id]).to.deep.equal({
      ...payload, isFetching: false,
    });
  });

  it('should handle MARKETPLACE_FETCH_REPOSITORY_SUMMARY_ERR', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_REPOSITORY_SUMMARY}_ERR`,
      meta: { id },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.certified[id] = { isFetching: true };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].error)
      .to.equal(status);
  });

  //----------------------------------------------------------------------------
  // REPOSITORY_FETCH_IMAGE_TAGS
  //----------------------------------------------------------------------------
  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_REQ (certified)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_REQ`,
      meta: { id, isCertified: true, namespace, reponame },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].tags).to.exist;
    expect(reducer.images.certified[id].tags.isFetching).to.equal(true);
  });

  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_REQ (community)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_REQ`,
      meta: { namespace, reponame },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].tags).to.exist;
    expect(reducer.images.community[repositoryId].tags.isFetching)
      .to.equal(true);
  });

  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_ACK (certified)', () => {
    const count = 1;
    const tagObj = { name: tag };
    const payload = { count, results: [tagObj] };
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_ACK`,
      meta: { id, isCertified: true, namespace, reponame },
      payload,
    };

    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].tags).to.exist;
    expect(reducer.images.certified[id].tags.isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].tags.count).to.equal(count);
    expect(reducer.images.certified[id].tags.pages).to.exist;
    expect(reducer.images.certified[id].tags.pages[1]).to.exist;
    expect(reducer.images.certified[id].tags.pages[1].results[0])
      .to.deep.equal(tagObj);
  });

  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_ACK (community)', () => {
    const count = 1;
    const tagObj = { name: tag };
    const payload = { count, results: [tagObj] };
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_ACK`,
      meta: { namespace, reponame },
      payload,
    };

    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].tags).to.exist;
    expect(reducer.images.community[repositoryId].tags.isFetching)
      .to.equal(false);
    expect(reducer.images.community[repositoryId].tags.count).to.equal(count);
    expect(reducer.images.community[repositoryId].tags.pages).to.exist;
    expect(reducer.images.community[repositoryId].tags.pages[1]).to.exist;
    expect(reducer.images.community[repositoryId].tags.pages[1].results[0])
      .to.deep.equal(tagObj);
  });

  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_ERR (certified)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_ERR`,
      meta: { id, isCertified: true, namespace, reponame },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.certified[id] = { tags: { isFetching: true } };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.certified[id].tags).to.exist;
    expect(reducer.images.certified[id].tags.isFetching).to.equal(false);
    expect(reducer.images.certified[id].tags.error).to.equal(status);
  });

  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_ERR (community)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_ERR`,
      meta: { namespace, reponame },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.community[repositoryId] = { tags: { isFetching: true } };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.community[repositoryId].tags).to.exist;
    expect(reducer.images.community[repositoryId].tags.isFetching)
      .to.equal(false);
    expect(reducer.images.community[repositoryId].tags.error).to.equal(status);
  });

  //----------------------------------------------------------------------------
  // MARKETPLACE_FETCH_BUNDLE_DETAIL
  //----------------------------------------------------------------------------
  it('should handle MARKETPLACE_FETCH_BUNDLE_DETAIL_REQ', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_BUNDLE_DETAIL}_REQ`,
      meta: { id },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.bundles[id]).to.exist;
    expect(reducer.bundles[id].isFetching).to.equal(true);
  });

  it('should handle MARKETPLACE_FETCH_BUNDLE_DETAIL_ACK', () => {
    const payload = {
      id,
      short_description: 'DDC bundle!',
    };
    const action = {
      type: `${MARKETPLACE_FETCH_BUNDLE_DETAIL}_ACK`,
      meta: { id },
      payload,
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.bundles[id]).to.exist;
    expect(reducer.bundles[id]).to.deep.equal({
      ...payload, isFetching: false,
    });
  });

  it('should handle MARKETPLACE_FETCH_BUNDLE_DETAIL_ERR', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_BUNDLE_DETAIL}_ERR`,
      meta: { id },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.bundles[id] = { isFetching: true };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.bundles[id]).to.exist;
    expect(reducer.bundles[id].isFetching).to.equal(false);
    expect(reducer.bundles[id].error).to.equal(status);
  });

  //----------------------------------------------------------------------------
  // MARKETPLACE_FETCH_BUNDLE_SUMMARY
  //----------------------------------------------------------------------------
  it('should handle MARKETPLACE_FETCH_BUNDLE_SUMMARY_REQ', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_BUNDLE_SUMMARY}_REQ`,
      meta: { id },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.bundles[id]).to.exist;
    expect(reducer.bundles[id].isFetching).to.equal(true);
  });

  it('should handle MARKETPLACE_FETCH_BUNDLE_SUMMARY_ACK', () => {
    const payload = {
      id,
      short_description: 'DDC bundle!',
    };
    const action = {
      type: `${MARKETPLACE_FETCH_BUNDLE_SUMMARY}_ACK`,
      meta: { id },
      payload,
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.bundles[id]).to.exist;
    expect(reducer.bundles[id]).to.deep.equal({
      ...payload, isFetching: false,
    });
  });

  it('should handle MARKETPLACE_FETCH_BUNDLE_SUMMARY_ERR', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_BUNDLE_SUMMARY}_ERR`,
      meta: { id },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.bundles[id] = { isFetching: true };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.bundles[id]).to.exist;
    expect(reducer.bundles[id].isFetching).to.equal(false);
    expect(reducer.bundles[id].error).to.equal(status);
  });

  //----------------------------------------------------------------------------
  // NAUTILUS_FETCH_SCAN_DETAIL
  //----------------------------------------------------------------------------
  it('should handle NAUTILUS_FETCH_SCAN_DETAIL_REQ', () => {
    const action = {
      type: `${NAUTILUS_FETCH_SCAN_DETAIL}_REQ`,
      meta: { id, namespace, reponame, tag },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].scanDetail).to.exist;
    expect(reducer.images.certified[id].scanDetail[tag]).to.exist;
    expect(reducer.images.certified[id].scanDetail[tag].isFetching)
      .to.equal(true);
  });

  it('should handle NAUTILUS_FETCH_SCAN_DETAIL_ACK', () => {
    const entities = { scan: [], blobs: [] };
    const action = {
      type: `${NAUTILUS_FETCH_SCAN_DETAIL}_ACK`,
      meta: { id, namespace, reponame, tag },
      payload: { entities },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].scanDetail).to.exist;
    expect(reducer.images.certified[id].scanDetail[tag]).to.exist;
    expect(reducer.images.certified[id].scanDetail[tag]).to.deep.equal({
      ...entities, isFetching: false,
    });
  });

  it('should handle NAUTILUS_FETCH_SCAN_DETAIL_ERR', () => {
    const err = 'This is an error HELLO WORLD';
    const action = {
      type: `${NAUTILUS_FETCH_SCAN_DETAIL}_ERR`,
      meta: { id, namespace, reponame, tag },
      payload: {
        error: err,
      },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].scanDetail).to.exist;
    expect(reducer.images.certified[id].scanDetail[tag]).to.exist;
    expect(reducer.images.certified[id].scanDetail[tag].isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].scanDetail[tag].error)
      .to.equal(err);
  });

  //----------------------------------------------------------------------------
  // NAUTILUS_FETCH_TAGS_AND_SCANS
  //----------------------------------------------------------------------------
  it('should handle NAUTILUS_FETCH_TAGS_AND_SCANS_REQ', () => {
    const action = {
      type: `${NAUTILUS_FETCH_TAGS_AND_SCANS}_REQ`,
      meta: { id, namespace, reponame },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.isFetching).to.equal(true);
  });

  it('should handle NAUTILUS_FETCH_TAGS_AND_SCANS_ACK w/ tags & scans', () => {
    const count = 1;
    const tags = { count, results: [{ name: tag }] };
    const scans = [{ namespace, reponame, tag }];
    const action = {
      type: `${NAUTILUS_FETCH_TAGS_AND_SCANS}_ACK`,
      meta: { id, namespace, reponame },
      payload: [tags, scans],
    };

    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].tagsAndScans.count).to.equal(count);
    expect(reducer.images.certified[id].tagsAndScans.pages).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.pages[1]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.pages[1][tag]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.pages[1][tag])
      .to.deep.equal({ namespace, name: tag, reponame, tag });
  });

  it('should handle NAUTILUS_FETCH_TAGS_AND_SCANS_ACK w/ only tags', () => {
    const count = 1;
    const tags = { count, results: [{ name: tag }] };
    const action = {
      type: `${NAUTILUS_FETCH_TAGS_AND_SCANS}_ACK`,
      meta: { id, namespace, reponame },
      payload: [tags],
    };

    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].tagsAndScans.count).to.equal(count);
    expect(reducer.images.certified[id].tagsAndScans.pages).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.pages[1]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.pages[1][tag]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.pages[1][tag])
      .to.deep.equal({ name: tag });
  });

  it('should handle NAUTILUS_FETCH_TAGS_AND_SCANS_ERR', () => {
    const err = 'This is an error HELLO WORLD';
    const action = {
      type: `${NAUTILUS_FETCH_TAGS_AND_SCANS}_ERR`,
      meta: { id, namespace, reponame },
      payload: {
        error: err,
      },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].tagsAndScans).to.exist;
    expect(reducer.images.certified[id].tagsAndScans.isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].tagsAndScans.error)
      .to.equal(err);
  });

  //----------------------------------------------------------------------------
  // REPOSITORY_FETCH_COMMENTS
  //----------------------------------------------------------------------------
  it('should handle REPOSITORY_FETCH_COMMENTS_REQ (community)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_COMMENTS}_REQ`,
      meta: { namespace, page, page_size, reponame },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].comments).to.exist;
    expect(reducer.images.community[repositoryId].comments.isFetching)
      .to.equal(true);
  });

  it('should handle REPOSITORY_FETCH_COMMENTS_REQ (certified)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_COMMENTS}_REQ`,
      meta: { id, isCertified: true, namespace, page, page_size, reponame },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].comments).to.exist;
    expect(reducer.images.certified[id].comments.isFetching)
      .to.equal(true);
  });

  it('should handle REPOSITORY_FETCH_COMMENTS_ACK (community)', () => {
    const comment = {
      id: 240,
      user: 'kristiehoward',
      comment: 'hey hey testing',
      created_on: '2015-09-16T22:04:55.969947Z',
      updated_on: '2015-09-16T22:04:55.993806Z',
    };
    const count = 1;
    const action = {
      type: `${REPOSITORY_FETCH_COMMENTS}_ACK`,
      meta: { namespace, reponame },
      payload: { count, results: [comment] },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.community[repositoryId] = {
      comments: {
        isFetching: true,
      },
    };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].comments).to.exist;
    expect(reducer.images.community[repositoryId].comments.isFetching)
      .to.equal(false);
    expect(reducer.images.community[repositoryId].comments.count)
      .to.equal(count);
    expect(reducer.images.community[repositoryId].comments.pages)
      .to.exist;
    expect(reducer.images.community[repositoryId].comments.pages)
      .to.deep.equal({
        1: {
          results: [comment],
        },
      });
  });

  it('should handle REPOSITORY_FETCH_COMMENTS_ACK (certified)', () => {
    const comment = {
      id: 240,
      user: 'kristiehoward',
      comment: 'hey hey testing',
      created_on: '2015-09-16T22:04:55.969947Z',
      updated_on: '2015-09-16T22:04:55.993806Z',
    };
    const count = 1;
    const action = {
      type: `${REPOSITORY_FETCH_COMMENTS}_ACK`,
      meta: { id, isCertified: true, namespace, page, page_size, reponame },
      payload: { count, results: [comment] },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.certified[id] = { comments: { isFetching: true } };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].comments).to.exist;
    expect(reducer.images.certified[id].comments.isFetching)
      .to.equal(false);
    expect(reducer.images.certified[id].comments.count)
      .to.equal(count);
    expect(reducer.images.certified[id].comments.pages)
      .to.exist;
    expect(reducer.images.certified[id].comments.pages)
      .to.deep.equal({
        1: {
          results: [comment],
        },
      });
  });

  it('should handle REPOSITORY_FETCH_COMMENTS_ERR (community)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_COMMENTS}_ERR`,
      meta: { namespace, reponame },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.community[repositoryId] = {
      isFetching: true,
      comments: {},
    };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].comments).to.exist;
    expect(reducer.images.community[repositoryId].comments.isFetching)
      .to.equal(false);
    expect(reducer.images.community[repositoryId].comments.error)
      .to.equal(status);
  });

  it('should handle REPOSITORY_FETCH_COMMENTS_ERR (certified)', () => {
    const action = {
      type: `${REPOSITORY_FETCH_COMMENTS}_ERR`,
      meta: { id, isCertified: true, namespace, page, page_size, reponame },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.certified[id] = { isFetching: true, comments: {} };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.certified[id]).to.exist;
    expect(reducer.images.certified[id].comments).to.exist;
    expect(reducer.images.certified[id].comments.isFetching).to.equal(false);
    expect(reducer.images.certified[id].comments.error).to.equal(status);
  });

  //----------------------------------------------------------------------------
  // REPOSITORY_FETCH_IMAGE_DETAIL
  //----------------------------------------------------------------------------
  it('should handle REPOSITORY_FETCH_IMAGE_DETAIL_REQ', () => {
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_DETAIL}_REQ`,
      meta: { namespace, reponame },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].isFetching).to.equal(true);
  });

  it('should handle REPOSITORY_FETCH_IMAGE_DETAIL_ACK', () => {
    const short_description = 'Busybox base image.';
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_DETAIL}_ACK`,
      meta: { namespace, reponame },
      payload: {
        namespace,
        name: reponame,
        short_description,
      },
    };
    const reducer = marketplaceReducer(undefined, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].isFetching).to.equal(false);
    expect(reducer.images.community[repositoryId].namespace)
      .to.equal(namespace);
    expect(reducer.images.community[repositoryId].reponame)
      .to.equal(reponame);
    expect(reducer.images.community[repositoryId].short_description)
      .to.equal(short_description);
  });

  it('should handle REPOSITORY_FETCH_IMAGE_DETAIL_ERR', () => {
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_DETAIL}_ERR`,
      meta: { namespace, reponame },
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const REQ_STATE = { ...DEFAULT_STATE };
    REQ_STATE.images.community[repositoryId] = { isFetching: true };
    // ERR must come after REQ, which initializes the entry for this repo
    const reducer = marketplaceReducer(REQ_STATE, action);
    expect(reducer.images.community[repositoryId]).to.exist;
    expect(reducer.images.community[repositoryId].isFetching)
      .to.equal(false);
    expect(reducer.images.community[repositoryId].error)
      .to.equal(status);
  });
});
