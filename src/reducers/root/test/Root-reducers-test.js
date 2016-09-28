import rootReducer,
  { DEFAULT_STATE } from 'reducers/root';
import {
  MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS,
  MARKETPLACE_FETCH_MOST_POPULAR,
  MARKETPLACE_FETCH_FEATURED,
} from 'actions/marketplace';
import {
  FINISH_PAGE_TRANSITION,
  ROOT_CHANGE_GLOBAL_SEARCH_VALUE,
  START_PAGE_TRANSITION,
} from 'actions/root';import { expect } from 'chai';

describe('root reducer', () => {
  it('should return the initial state', () => {
    expect(rootReducer(undefined, {})).to.deep.equal(DEFAULT_STATE);
  });

  it('should handle MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS_REQ', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS}_REQ`,
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.autocomplete.isFetching).to.equal(true);
  });

  it('should handle MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS_ACK', () => {
    const namespace = 'libary';
    const reponame = 'ubuntu';
    const name = 'ubuntu';
    const id = '1234';
    const categories = ['infrastructure'];
    const summaries = [{ categories, id, name, namespace, reponame }];
    const action = {
      type: `${MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS}_ACK`,
      payload: {
        count: 1,
        summaries,
      },
    };

    const reducer = rootReducer(undefined, action);
    expect(reducer.autocomplete.isFetching).to.equal(false);
    expect(reducer.autocomplete.suggestions).to.exist;
    expect(reducer.autocomplete.suggestions).to.have.length(1);
    expect(reducer.autocomplete.suggestions[0]).to.deep.equal({
      categories, id, namespace, reponame, name,
    });
  });

  it('should handle MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS_ERR', () => {
    const err = 'This is an error HELLO WORLD';
    const action = {
      type: `${MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS}_ERR`,
      payload: {
        error: err,
      },
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.autocomplete.isFetching).to.equal(false);
    expect(reducer.autocomplete.error).to.equal(err);
  });

  it('should handle ROOT_CHANGE_GLOBAL_SEARCH_VALUE', () => {
    const value = 'query';
    const action = {
      type: ROOT_CHANGE_GLOBAL_SEARCH_VALUE,
      payload: { value },
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.globalSearch).to.equal(value);
  });

  it('should handle START_PAGE_TRANSITION', () => {
    const action = {
      type: START_PAGE_TRANSITION,
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.isPageTransitioning).to.equal(true);
  });

  it('should handle FINISH_PAGE_TRANSITION', () => {
    const action = {
      type: FINISH_PAGE_TRANSITION,
    };
    const STATE = { isPageTransitioning: true };
    const reducer = rootReducer(STATE, action);
    expect(reducer.isPageTransitioning).to.equal(false);
  });

  it('should handle MARKETPLACE_FETCH_MOST_POPULAR_REQ', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_MOST_POPULAR}_REQ`,
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.landingPage.mostPopular.isFetching).to.equal(true);
    expect(reducer.landingPage.mostPopular.images).to.deep.equal([]);
  });

  it('should handle MARKETPLACE_FETCH_MOST_POPULAR_ACK', () => {
    const namespace = 'libary';
    const reponame = 'ubuntu';
    const category = ['infrastructure'];
    // sample image
    const summaries = [{ category, namespace, reponame }];
    const action = {
      type: `${MARKETPLACE_FETCH_MOST_POPULAR}_ACK`,
      payload: {
        summaries,
      },
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.landingPage.mostPopular.isFetching).to.equal(false);
    expect(reducer.landingPage.mostPopular.images).to.exist;
    expect(reducer.landingPage.mostPopular.images[0]).to.exist;
    expect(reducer.landingPage.mostPopular.images[0]).to.deep.equal({
      category, namespace, reponame,
    });
  });


  it('should handle MARKETPLACE_FETCH_MOST_POPULAR_ERR', () => {
    const err = 'This is an error HELLO WORLD';
    const action = {
      type: `${MARKETPLACE_FETCH_MOST_POPULAR}_ERR`,
      payload: {
        error: err,
      },
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.landingPage.mostPopular.isFetching).to.equal(false);
    expect(reducer.landingPage.mostPopular.error).to.exist;
  });

  it('should handle MARKETPLACE_FETCH_FEATURED_REQ', () => {
    const action = {
      type: `${MARKETPLACE_FETCH_FEATURED}_REQ`,
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.landingPage.featured.isFetching).to.equal(true);
    expect(reducer.landingPage.featured.images).to.deep.equal([]);
  });

  it('should handle MARKETPLACE_FETCH_FEATURED_ACK', () => {
    const namespace = 'libary';
    const reponame = 'ubuntu';
    const category = ['infrastructure'];
    // sample image
    const summaries = [{ category, namespace, reponame }];
    const action = {
      type: `${MARKETPLACE_FETCH_FEATURED}_ACK`,
      payload: {
        summaries,
      },
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.landingPage.featured.isFetching).to.equal(false);
    expect(reducer.landingPage.featured.images).to.exist;
    expect(reducer.landingPage.featured.images[0]).to.exist;
    expect(reducer.landingPage.featured.images[0]).to.deep.equal({
      category, namespace, reponame,
    });
  });


  it('should handle MARKETPLACE_FETCH_FEATURED_ERR', () => {
    const err = 'This is an error HELLO WORLD';
    const action = {
      type: `${MARKETPLACE_FETCH_FEATURED}_ERR`,
      payload: {
        error: err,
      },
    };
    const reducer = rootReducer(undefined, action);
    expect(reducer.landingPage.featured.isFetching).to.equal(false);
    expect(reducer.landingPage.featured.error).to.exist;
  });
});
