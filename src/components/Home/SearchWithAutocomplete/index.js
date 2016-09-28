import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import css from './styles.css';
import { AutocompleteSearchBar } from 'common';
import {
  marketplaceFetchAutocompleteSuggestions,
  marketplaceSearch,
} from 'actions/marketplace';
import { rootChangeGlobalSearchValue } from 'actions/root';
import formatCategories from 'lib/utils/format-categories';
import routes from 'lib/constants/routes';
import debounce from 'lodash/debounce';
import { DEFAULT_DEBOUNCE_TIME } from 'lib/constants/defaults';
const { array, bool, func, object, shape } = PropTypes;

const mapStateToProps = ({ root }) => {
  const { autocomplete } = root;
  return { autocomplete };
};

const dispatcher = {
  marketplaceFetchAutocompleteSuggestions,
  marketplaceSearch,
  rootChangeGlobalSearchValue,
};

/*
 * Search to be used in the LandingPage with the same autocomplete functionality
 * as the global search bar
 */
@connect(mapStateToProps, dispatcher)
export default class SearchWithAutocomplete extends Component {
  static propTypes = {
    autocomplete: shape({
      isFetching: bool,
      suggestions: array,
    }),
    location: object.isRequired,
    marketplaceFetchAutocompleteSuggestions: func.isRequired,
    marketplaceSearch: func.isRequired,
    rootChangeGlobalSearchValue: func.isRequired,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  constructor(props) {
    super(props);
    // Make sure there is one debounce function per component instance
    // http://stackoverflow.com/a/28046731/5965502
    this.debouncedFetchAutocompleteSuggestions = debounce(
      this.debouncedFetchAutocompleteSuggestions,
      DEFAULT_DEBOUNCE_TIME
    );
  }

  state = {
    searchQuery: '',
  }

  componentWillUnmount() {
    this.debouncedFetchAutocompleteSuggestions.cancel();
  }

  onSearchQueryChange = (e, value) => {
    // Remove this synthetic event from the pool so that we can still access the
    // event asyncronously (for debouncing)
    // https://facebook.github.io/react/docs/events.html#event-pooling
    e.persist();
    // change the value showing in the search bar
    this.setState({ searchQuery: value });
    // fetch new suggestions (debounced)
    this.debouncedFetchAutocompleteSuggestions(e, value);
  }

  onSelectAutosuggestItem = (value, item) => {
    const { id } = item;
    // Jump to the product detail page for this result
    const detail = routes.imageDetail({ id });
    this.context.router.push(detail);
  }

  debouncedFetchAutocompleteSuggestions = (e, value) => {
    this.props.marketplaceFetchAutocompleteSuggestions({ q: value });
  }

  // Search Bar form has been submitted
  search = (q) => {
    // change the value in the global search bar
    // this.props.rootChangeGlobalSearchValue({ value });
    // fire search action and transition to search results page
    this.props.marketplaceSearch({ q });
    const pathname = routes.search();
    const { state } = this.props.location;
    // search from global search bar will have a query (q) (no page num or size)
    const query = { q };
    this.context.router.push({ pathname, query, state });
  }

  renderAutocompleteItem = (item, isHighlighted) => {
    const { id, name, categories } = item;
    const catNames = formatCategories(categories);
    const catText = catNames ? ` in ${catNames}` : '';
    const itemClass = isHighlighted ? css.highlightedResult : css.result;
    return (
      <div className={itemClass} key={id} id={id}>
        <span className={css.resultName}>{name}</span>
        <span className={css.resultCategories}>{catText}</span>
      </div>
    );
  };

  render() {
    const { autocomplete } = this.props;
    const { searchQuery } = this.state;
    const { suggestions = [] } = autocomplete;
    const menuTitle = (
      <div className={css.menuTitle}>
        Suggested Results
      </div>
    );
    const getItemValue = (item) => item.id;
    const classNames = {
      icon: css.icon,
      input: css.input,
      wrapper: css.wrapper,
    };
    return (
      <AutocompleteSearchBar
        classNames={classNames}
        getItemValue={getItemValue}
        id="landingPage-search"
        items={suggestions}
        menuTitle={menuTitle}
        onChange={this.onSearchQueryChange}
        onSelect={this.onSelectAutosuggestItem}
        onSubmit={this.search}
        ref="autocomplete"
        renderItem={this.renderAutocompleteItem}
        value={searchQuery}
      />
    );
  }
}
