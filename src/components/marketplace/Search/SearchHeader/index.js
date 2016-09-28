import React, { Component, PropTypes } from 'react';
import { Select } from 'common';
import qs from 'qs';
import css from './styles.css';
import forEach from 'lodash/forEach';
import omit from 'lodash/omit';
import { DEFAULT_SEARCH_PAGE_SIZE } from 'lib/constants/defaults';
import routes from 'lib/constants/routes';
import { sorting } from 'lib/constants/searchFilters';
const {
  bool,
  func,
  node,
  number,
  object,
  shape,
} = PropTypes;

export default class SearchHeader extends Component {
  static propTypes = {
    count: number,
    isFetching: bool,
    location: shape({
      query: object,
      state: node,
    }).isRequired,
    marketplaceSearch: func.isRequired,
    rootChangeGlobalSearchValue: func.isRequired,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  onSortChange = ({ value }) => {
    const { query } = this.props.location;
    const newQuery = omit(query, ['sort', 'order']);
    // value may be 'sort=platform&order=asc'
    const sortParams = qs.parse(value);
    forEach(sortParams, (val, q) => {
      newQuery[q] = val;
    });
    this.updateSearch(newQuery);
  }

  getNumCurrentSearchResults() {
    const { count } = this.props;
    const { query } = this.props.location;
    const page_size = parseInt(
      (query && query.page_size || DEFAULT_SEARCH_PAGE_SIZE),
      10
    );
    const page = parseInt((query && query.page || 1), 10);
    let first = 1;
    let last = page_size;
    if (count < page_size) {
      last = count;
    }
    if (page !== 1) {
      first = 1 + ((page - 1) * page_size);
      last = page * page_size;
    }
    // only one page of results (fewer than page_size results)
    // or showing the last page of results without a full page_size of results
    if (last > count) {
      last = count;
    }

    if (count === 0) {
      first = 0;
      last = 0;
    }
    return <span>{first} - {last} of {count}</span>;
  }

  clearSearch = () => {
    const { query } = this.props.location;
    const newQuery = omit(query, 'q');
    this.updateSearch(newQuery);
  }

  updateSearch = (query) => {
    const { state } = this.props.location;
    // Update global search bar
    this.props.rootChangeGlobalSearchValue({ value: query.q || '' });
    const pathname = routes.search();
    this.context.router.push({ pathname, query, state });
  }

  renderSortSelect = (value) => {
    return (
      <div className={css.select}>
        <Select
          clearable={false}
          options={sorting}
          defaultValue={''}
          value={value}
          onChange={this.onSortChange}
        />
      </div>
    );
  }

  renderCurrentSearch() {
    const { count, isFetching } = this.props;
    if (isFetching) {
      return <div></div>;
    }
    const numResults = this.getNumCurrentSearchResults();
    const results = count === 1 ? 'result' : 'results';
    const { query } = this.props.location;
    const searchQuery = query && query.q || '';
    const currentSearchTerm = <b>{ searchQuery }</b>;
    let clear = (
      <a onClick={this.clearSearch} className={css.clear}>Clear search</a>
    );
    let description = <span>{ results } for { currentSearchTerm }</span>;
    if (searchQuery === '') {
      description = <span>available images</span>;
      clear = '';
    }
    return <div>{ numResults } { description } { clear }</div>;
  }

  render() {
    const { query } = this.props.location;
    // Sorting is one of [sort, sort&order, order]
    const sort = query && query.sort || '';
    const order = query && query.order || '';
    const selectedSortQuery = {};
    if (sort) selectedSortQuery.sort = sort;
    if (order) selectedSortQuery.order = order;
    // map into sorting string (sort=...&order=...)
    const selectedSort = qs.stringify(selectedSortQuery);
    return (
      <div className={css.searchHeader}>
        <div className={css.currentSearch}>
          { this.renderCurrentSearch() }
        </div>
        <div>
          { this.renderSortSelect(selectedSort) }
        </div>
      </div>
    );
  }
}
