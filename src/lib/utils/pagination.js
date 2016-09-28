import { DEFAULT_SEARCH_PAGE_SIZE } from 'lib/constants/defaults';

export const getCurrentPage = (location) => {
  return location.query && parseInt(location.query.page, 10) || 1;
};

export const getCurrentPageSize = (location) => {
  return location.query && parseInt(location.query.page_size, 10)
    || DEFAULT_SEARCH_PAGE_SIZE;
};
