import {
  MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS,
  MARKETPLACE_FETCH_MOST_POPULAR,
  MARKETPLACE_FETCH_FEATURED,
} from 'actions/marketplace';
import {
  FINISH_PAGE_TRANSITION,
  ROOT_CHANGE_GLOBAL_SEARCH_VALUE,
  START_PAGE_TRANSITION,
} from 'actions/root';
import cloneDeep from 'lodash/cloneDeep';
import set from 'lodash/set';

export const DEFAULT_STATE = {
  autocomplete: {},
  globalSearch: '',
  isPageTransitioning: false,
  landingPage: {
    mostPopular: {
      isFetching: false,
      images: [],
    },
    featured: {
      isFetching: false,
      images: [],
    },
  },
};

const reduceAutocomplete = (nextState, payload) => {
  const results = payload && payload.summaries || [];
  const suggestions = results.map((image) => {
    // Id is the product id
    const { categories, id, name, namespace, reponame } = image;
    return { categories, id, name, namespace, reponame };
  });
  set(nextState, ['autocomplete', 'isFetching'], false);
  set(nextState, ['autocomplete', 'suggestions'], suggestions);
  return nextState;
};

export default function root(state = DEFAULT_STATE, action) {
  const nextState = cloneDeep(state);
  const { payload, type } = action;

  switch (type) {
    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_MOST_POPULAR
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_MOST_POPULAR}_REQ`:
      nextState.landingPage.mostPopular.isFetching = true;
      return nextState;
    case `${MARKETPLACE_FETCH_MOST_POPULAR}_ACK`:
      nextState.landingPage.mostPopular.images = payload.summaries;
      return nextState;
    case `${MARKETPLACE_FETCH_MOST_POPULAR}_ERR`:
      nextState.landingPage.mostPopular.isFetching = false;
      nextState.landingPage.mostPopular.error =
        'Unable to fetch most popular images';
      return nextState;

    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_FEATURED
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_FEATURED}_REQ`:
      nextState.landingPage.featured.isFetching = true;
      return nextState;
    case `${MARKETPLACE_FETCH_FEATURED}_ACK`:
      nextState.landingPage.featured.images = payload.summaries;
      return nextState;
    case `${MARKETPLACE_FETCH_FEATURED}_ERR`:
      nextState.landingPage.featured.isFetching = false;
      nextState.landingPage.featured.error =
        'Unable to fetch featured images';
      return nextState;

    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS}_REQ`:
      nextState.autocomplete.isFetching = true;
      return nextState;
    case `${MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS}_ACK`:
      return reduceAutocomplete(nextState, payload);
    case `${MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS}_ERR`:
      nextState.autocomplete.isFetching = false;
      nextState.autocomplete.error = payload.error;
      return nextState;

    //--------------------------------------------------------------------------
    // ROOT_CHANGE_GLOBAL_SEARCH_VALUE
    //--------------------------------------------------------------------------
    case `${ROOT_CHANGE_GLOBAL_SEARCH_VALUE}`:
      nextState.globalSearch = payload.value;
      return nextState;

    //--------------------------------------------------------------------------
    // START_PAGE_TRANSITION
    //--------------------------------------------------------------------------
    case `${START_PAGE_TRANSITION}`:
      nextState.isPageTransitioning = true;
      return nextState;

    //--------------------------------------------------------------------------
    // FINISH_PAGE_TRANSITION
    //--------------------------------------------------------------------------
    case `${FINISH_PAGE_TRANSITION}`:
      nextState.isPageTransitioning = false;
      return nextState;
    default:
      return state;
  }
}
