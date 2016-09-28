import { createAction } from 'redux-actions';
/* eslint-disable max-len */

// TODO: tests - nathan 06/3/16 Write tests for redirect
export const FINISH_PAGE_TRANSITION = 'FINISH_PAGE_TRANSITION';
export const REDIRECT = 'REDIRECT';
export const ROOT_CHANGE_GLOBAL_SEARCH_VALUE = 'ROOT_CHANGE_GLOBAL_SEARCH_VALUE';
export const START_PAGE_TRANSITION = 'START_PAGE_TRANSITION';

export const redirectTo = createAction(REDIRECT);

export const startPageTransition = createAction(START_PAGE_TRANSITION);
export const finishPageTransition = createAction(FINISH_PAGE_TRANSITION);

export const rootChangeGlobalSearchValue = ({ value }) => {
  return {
    type: ROOT_CHANGE_GLOBAL_SEARCH_VALUE,
    payload: { value },
  };
};

/* eslint-enable */
