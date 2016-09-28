import React, { Component, PropTypes } from 'react';
import css from './styles.css';
const {
  bool,
  func,
  node,
  object,
  shape,
  string,
} = PropTypes;

class EmptySearchResults extends Component {
  static propTypes = {
    isCommunitySearch: bool,
    location: shape({
      pathname: string,
      query: object,
      state: node,
    }),
    rootChangeGlobalSearchValue: func.isRequired,
  }

  static defaultProps = {
    isCommunitySearch: false,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  updateSearch = (query) => () => {
    const { pathname, state } = this.props.location;
    // Update global search bar
    this.props.rootChangeGlobalSearchValue({ value: query.q || '' });
    this.context.router.push({ pathname, query, state });
  }

  renderLink = (description, linkText, onClick) => {
    return (
      <div>
        {description}
        <a className={css.link} onClick={onClick}>
          {linkText}
        </a>
      </div>
    );
  }

  render() {
    const { isCommunitySearch } = this.props;
    const { query } = this.props.location;
    const hasSearchTerm = query && query.q;
    const hasFilter = query &&
      (query.source || query.category || query.platform);
    let description;
    let linkText;
    let onClick;
    if (isCommunitySearch) {
      description = 'There are no results for this search in Docker Hub. ';
      linkText = 'Search Docker Store Images';
      onClick = this.updateSearch({ q: query.q });
    } else if (hasSearchTerm && hasFilter) {
      description = 'Your selected filters may be' +
        ' too narrow for this search term. ';
      linkText = 'Clear filters';
      onClick = this.updateSearch({ q: query.q });
    } else {
      description = 'There are no results for this search. ';
      linkText = 'Try another search';
      onClick = this.updateSearch({});
    }
    return (
      <div className={css.searchHeader}>
        {this.renderLink(description, linkText, onClick)}
      </div>
    );
  }
}

export default EmptySearchResults;
