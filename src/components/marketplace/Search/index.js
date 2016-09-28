import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { marketplaceSearch } from 'actions/marketplace';
import { rootChangeGlobalSearchValue } from 'actions/root';
import EmptySearchResults from './EmptySearchResults';
import ImageSearchResult from './ImageSearchResult';
import CommunityImageSearchResult from './CommunityImageSearchResult';
import SearchHeader from './SearchHeader';
import FilterList from './FilterList';
import { getCurrentPage, getCurrentPageSize } from 'lib/utils/pagination';
import css from './styles.css';
import { Card, Pagination, WarningIcon } from 'common';
import { COMMUNITY } from 'lib/constants/searchFilters/sources';
import { PANIC } from 'lib/constants/variants';
import ceil from 'lodash/ceil';
import omit from 'lodash/omit';
import map from 'lodash/map';
import indexOf from 'lodash/indexOf';
import qs from 'qs';
const { array, bool, func, number, object, shape, string } = PropTypes;

const getSelectedFilters = ({ query }) => {
  const queries = qs.parse(query);
  if (!query || !queries) {
    return {};
  }
  const { category = '', platform = '' } = queries;
  return {
    categories: category.split(','),
    platforms: platform.split(','),
  };
};

const mapState = ({ marketplace }, { location }) => {
  const page = getCurrentPage(location);
  const isCommunitySearch = location && location.query
    && location.query.source === COMMUNITY || false;
  const { filters, search } = marketplace;
  return {
    count: search && search.count || 0,
    error: search && search.error || undefined,
    filters: filters || {},
    isCommunitySearch,
    isFetching: search && search.isFetching,
    results: search && search.pages[page] && search.pages[page].results || [],
    selectedFilters: getSelectedFilters(location) || {},
  };
};

const mapDispatch = {
  marketplaceSearch,
  rootChangeGlobalSearchValue,
};

@connect(mapState, mapDispatch)
export default class Search extends Component {
  static propTypes = {
    count: number,
    error: string,
    filters: shape({
      categories: object,
      platforms: object,
    }),
    isCommunitySearch: bool,
    isFetching: bool,
    location: shape({
      query: object,
    }),
    marketplaceSearch: func.isRequired,
    results: array,
    rootChangeGlobalSearchValue: func,
    selectedFilters: shape({
      categories: array,
      platforms: array,
    }),
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  onFilterChange = (queryValue) => (selectedFilters) => {
    const { query } = this.props.location;
    // Don't add this filter to the query if none are selected
    const newQuery = omit(query, ['page', queryValue]);
    if (selectedFilters) {
      newQuery[queryValue] = selectedFilters;
    }
    this.updateSearch(newQuery);
  }

  goToPage = (page) => {
    const { query } = this.props.location;
    const newQuery = { ...query, page };
    this.updateSearch(newQuery);
  }

  showCommunityImages = () => {
    const { query } = this.props.location;
    // Only preserve the search term (filters and page don't apply)
    const q = query && query.q || '';
    const newQuery = { q, source: COMMUNITY };
    this.updateSearch(newQuery);
  }

  showDockerVerifiedImages = () => {
    const { query } = this.props.location;
    // Default is to only show docker verified images. only preserve search
    // term b/c filters and pages don't apply
    const q = query && query.q || '';
    this.updateSearch({ q });
  }

  updateSearch = (query) => {
    const { pathname, state } = this.props.location;
    this.context.router.push({ pathname, query, state });
  }

  createFilters = ({ disabled, filters, selectedFilters }) => {
    const createdFilters = map(filters, (label, value) => {
      const isChecked = indexOf(selectedFilters, value) >= 0;
      return { disabled, isChecked, label, value };
    });
    // Sort a list of filters alphabetically by label
    return createdFilters.sort(({ label: l1 }, { label: l2 }) => {
      if (l1 === l2) return 0;
      return l1 > l2 ? 1 : -1;
    });
  }

  renderIsFetching() {
    return <div>Fetching...</div>;
  }

  renderShowCommunityImagesFooter() {
    const showCommunityImages = (
      <a className={css.communityImagesLink} onClick={this.showCommunityImages}>
        View results from Docker Hub
      </a>
    );
    const helpText = [
      'Docker Hub Content is unverified & unscanned content',
      'by the Docker Community',
    ].join(' ');
    return (
      <div className={css.communityImagesBanner}>
        <div className={css.description}>
          Not finding what you're looking for? {showCommunityImages}
        </div>
        <div className={css.helpText}>{helpText}</div>
      </div>
    );
  }

  renderShowDockerVerifiedBanner() {
    const showDockerVerified = (
      <a
        className={css.communityImagesLink}
        onClick={this.showDockerVerifiedImages}
      >
        Back to Docker Verified Images
      </a>
    );
    const helpText = [
      'Docker Hub Content is unverified & unscanned content',
      'by the Docker Community',
    ].join(' ');
    return (
      <Card className={css.showDockerVerifiedBanner} shadow>
        <div>
          <WarningIcon variant={PANIC} className={css.warning} />
        </div>
        <div>
          <div className={css.description}>
            Currently showing results from Docker Hub. {showDockerVerified}
          </div>
          <div className={css.helpText}>{helpText}</div>
        </div>
      </Card>
    );
  }

  renderImage = (image) => {
    const { categories } = this.props.filters;
    if (image.source === COMMUNITY) {
      return <CommunityImageSearchResult key={image.name} image={image} />;
    }
    return (
      <ImageSearchResult
        key={image.id}
        categories={categories}
        image={image}
        location={this.props.location}
      />
    );
  }

  renderFilters() {
    const { filters, isCommunitySearch, selectedFilters } = this.props;
    const categoryFilters = this.createFilters({
      disabled: isCommunitySearch,
      filters: filters.categories,
      selectedFilters: selectedFilters.categories,
    });
    // TODO Kristie 6/18/16 Bring back platform filters when windows is avail
    // const platformFilters = this.createFilters({
    //   disabled: isCommunitySearch,
    //   filters: filters.platforms,
    //   selectedFilters: selectedFilters.platforms,
    // });
    // <FilterList
    //   className={css.filterList}
    //   filters={platformFilters}
    //   onChange={this.onFilterChange('platform')}
    //   title="Available for"
    // />
    return (
      <div>
        <FilterList
          className={css.filterList}
          filters={categoryFilters}
          onChange={this.onFilterChange('category')}
          title="Categories"
        />
      </div>
    );
  }

  renderPagination() {
    const { count, error, isFetching, location } = this.props;
    const page = getCurrentPage(location);
    const pageSize = getCurrentPageSize(location);
    if (isFetching || !count || error) {
      return null;
    }
    const lastPage = ceil(count / pageSize);
    if (lastPage === 1) {
      return null;
    }
    return (
      <Pagination
        className={css.pagination}
        currentPage={page}
        lastPage={lastPage}
        onChangePage={this.goToPage}
      />
    );
  }

  renderSearchResults() {
    const {
      count,
      isCommunitySearch,
      location,
      results,
      rootChangeGlobalSearchValue: changeGlobalSearchValue,
    } = this.props;
    if (count === 0 || !results.length) {
      let footer;
      if (!isCommunitySearch) {
        footer = this.renderShowCommunityImagesFooter();
      }
      return (
        <div>
          <EmptySearchResults
            isCommunitySearch={isCommunitySearch}
            location={location}
            rootChangeGlobalSearchValue={changeGlobalSearchValue}
          />
          {footer}
        </div>
      );
    }
    if (isCommunitySearch) {
      return (
        <div>
          {this.renderShowDockerVerifiedBanner()}
          {results.map(this.renderImage)}
        </div>
      );
    }
    return (
      <div>
        <div className={css.resultsWrapper}>
          {results.map(this.renderImage)}
        </div>
        {this.renderShowCommunityImagesFooter()}
        <hr />
      </div>
    );
  }

  render() {
    const {
      count,
      isFetching,
      location,
      marketplaceSearch: search,
      rootChangeGlobalSearchValue: changeGlobalSearchValue,
    } = this.props;
    return (
      <div className={css.content}>
        <div>
          {this.renderFilters()}
        </div>
        <div>
          <SearchHeader
            count={count}
            isFetching={isFetching}
            location={location}
            marketplaceSearch={search}
            rootChangeGlobalSearchValue={changeGlobalSearchValue}
          />
          {isFetching ? this.renderIsFetching() : this.renderSearchResults()}
          {this.renderPagination()}
        </div>
      </div>
    );
  }
}
