'use strict';

import React, { Component, PropTypes } from 'react';
const { func, node, number, object } = PropTypes;
import cn from 'classnames';
import ui from 'redux-ui';
import { getPage } from 'utils';

import FontAwesome from 'components/common/fontAwesome';

import styles from './pagination.css';

@ui({
  key: 'pagination',
  state: {
    currentPage: 1
  }
})
export default class Pagination extends Component {
  static propTypes = {
    ui: object,
    updateUI: func,

    children: node, // Can be omitted if you use renderSlice
    numItems: number, // Defaults to children.length
    renderSlice: func, // Defaults to array slicing children
    pageSize: number,
    location: object.isRequired
  }

  static defaultProps = {
    pageSize: 10,
    renderSlice: (start, end, children) => children.slice(start, end)
  }

  componentWillMount() {
    const page = getPage(this.props.location.search);
    this.changePageTo(page || this.props.ui.currentPage);
  }

  componentWillReceiveProps(next) {
    this.changePageTo(next.ui.currentPage, next);
  }

  getNumPages(props = undefined) {
    let { children, numItems, pageSize } = props || this.props;
    numItems = (numItems !== undefined) ? numItems : children.length;
    const numPages = Math.ceil(numItems / pageSize);
    return numPages;
  }

  changePageTo(newPage, props = undefined) {
    props = props || this.props;

    const numPages = this.getNumPages(props);
    if (!numPages) {
      // Props haven't fully loaded to calculate numPages yet
      return;
    }

    if (!newPage || newPage < 1) {
      newPage = 1;
    } else if (newPage > numPages) {
      // If current page is somehow greater than last page (for example if deleting items), change page to max page
      newPage = numPages;
    }

    const pathname = this.props.location.pathname;
    const query = this.props.location.query;
    // If query param for page is undefined, we must be on page 1
    // Otherwise it must match the newPage we want to go to
    // If it doesn't match, update our url
    if ((newPage === 1 && query && query.page) ||
        (newPage !== 1 && query && parseInt(query.page) !== newPage)) {
      let newQuery;
      if (newPage === 1) {
        newQuery = { ...query };
        delete newQuery.page;
      } else {
        newQuery = { ...query, page: newPage };
      }
      this.context.router.replace({
          pathname: pathname,
          query: newQuery,
          state: {}
      });
    }

    if (props.ui.currentPage !== newPage) {
      this.props.updateUI('currentPage', newPage);
    }
  }

  renderPageNums(numPages, currentPage) {
    const leading = 3, context = 1, trailing = 3;
    const shouldRenderPageNum = (i) => {
      // Only need to render this page number if it's within
      // - `leading` of start
      // - `trailing` of end
      // - `context` around the current page number.
      return (i <= leading) ||
             (i >= currentPage - context && i <= currentPage + context) ||
             (i > numPages - trailing);
    };

    let pageNums = [];
    let previousElementIsEllipsis = false;
    const onClick = (i) => () => this.changePageTo(i);
    for (let i = 1; i <= numPages; i++) {
      if (!shouldRenderPageNum(i)) {
        // Skip over this `i` if we aren't rendering a page number for it
        if (!previousElementIsEllipsis) {
          // Add a ellipsis if we haven't already
          pageNums.push(<FontAwesome key={ 'ellipsis' + i } icon='fa-ellipsis-h' />);
          previousElementIsEllipsis = true;
        }
        continue;
      }

      // Otherwise render this pageNum
      const className = cn(styles.pageNum, {[styles.active]: currentPage === i });
      pageNums.push(
        <div key={ i } className={ className } onClick={ onClick(i) }>
          { i }
        </div>
      );
      previousElementIsEllipsis = false;
    }

    return pageNums;
  }

  render() {
    const { children, ui: { currentPage }, numItems, renderSlice, pageSize, ...props } = this.props;
    const numPages = this.getNumPages();
    const canPaginate = (numPages > 1);
    let start = 0, end = numItems;
    if (canPaginate) {
      start = (currentPage - 1) * pageSize;
      end = currentPage * pageSize;
    }
    return (
      <div {...props}>

        { renderSlice(start, end, children) }

        { canPaginate &&
          <div className={ styles.footer }>
            <div
              key='previous'
              className={ cn(styles.pageNum, {[styles.disabled]: currentPage === 1 }) }
              onClick={ () => this.changePageTo(currentPage - 1) }
            >
              <FontAwesome icon='fa-chevron-left' />
            </div>
            { this.renderPageNums(numPages, currentPage) }
            <div
              key='next'
              className={ cn(styles.pageNum, {[styles.disabled]: currentPage === numPages }) }
              onClick={ () => this.changePageTo(currentPage + 1) }
            >
              <FontAwesome icon='fa-chevron-right' />
            </div>
          </div>
        }
      </div>
    );
  }
}
