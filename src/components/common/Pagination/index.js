import React, { PropTypes, Component } from 'react';
import css from './styles.css';
import { range } from 'lodash';
import classnames from 'classnames';
import Page from './Page.js';
const { func, number, oneOfType, string } = PropTypes;

export default class Pagination extends Component {
  static propTypes = {
    className: string,
    currentPage: oneOfType([number, string]).isRequired,
    lastPage: oneOfType([number, string]).isRequired,
    maxVisible: oneOfType([number, string]),
    onChangePage: func.isRequired,
  }

  static defaultProps = {
    maxVisible: 10,
  }

  calculatePagesShown = (currentPage, lastPage, maxVisible) => {
    // Calculate the page numbers that are displayed
    let start = 1;
    let end = lastPage;
    let lastPageShown = start + maxVisible - 1;
    if (currentPage > lastPageShown) {
      // shift window to the right so that the current page shows
      start += currentPage - lastPageShown;
      // recalculate last page shown to account for new start
      lastPageShown = start + maxVisible - 1;
    }
    // adjust the end if it is too far past the max visible
    if (end > lastPageShown && lastPageShown < lastPage) {
      end = lastPageShown;
    }
    // range is exclusive of 2nd param, so add 1 to make range inclusive
    // and return the array of page numbers to display
    return range(start, end + 1);
  }

  mkPage = (pageNumber) => {
    const { currentPage, onChangePage } = this.props;
    return (
      <Page key={pageNumber}
        currentPage={parseInt(currentPage, 10)}
        pageNumber={pageNumber}
        onClick={onChangePage}
      />
    );
  }

  render() {
    const { onChangePage } = this.props;
    // Make sure props are converted to numbers to avoid math errors
    const currentPage = parseInt(this.props.currentPage, 10);
    const lastPage = parseInt(this.props.lastPage, 10);
    const maxVisible = parseInt(this.props.maxVisible, 10);
    let nextPage;
    let previousPage;
    // >> and << components to jump to first or last page
    let jumpToFirstPage;
    let jumpToLastPage;
    const isSinglePage = currentPage === 1 && lastPage === 0;

    if (!isSinglePage && currentPage !== lastPage) {
      nextPage = (
        <Page currentPage={currentPage}
          arrowType="next"
          pageNumber={currentPage + 1}
          onClick={onChangePage}
        />
      );
      jumpToLastPage = (
        <Page currentPage={currentPage}
          arrowType="last"
          pageNumber={lastPage}
          onClick={onChangePage}
        />
      );
    }
    if (currentPage !== 1) {
      previousPage = (
        <Page currentPage={currentPage}
          arrowType="previous"
          pageNumber={currentPage - 1}
          onClick={onChangePage}
        />
      );
      jumpToFirstPage = (
        <Page currentPage={currentPage}
          arrowType="first"
          pageNumber={1}
          onClick={onChangePage}
        />
      );
    }
    const classes = classnames(
      'dpagination',
      { [this.props.className]: !!this.props.className }
    );
    const pages = this.calculatePagesShown(currentPage, lastPage, maxVisible);
    return (
      <div className={classes}>
        <ul className={css.paginationCentered}>
          {jumpToFirstPage}
          {previousPage}
          {pages.map(this.mkPage)}
          {nextPage}
          {jumpToLastPage}
        </ul>
      </div>
    );
  }
}
