'use strict';

import React, {
  PropTypes,
  createClass
} from 'react';
import classnames from 'classnames';
import _ from 'lodash';
var debug = require('debug')('Pagination');

function noop(e) {
  e.preventDefault();
}

var Page = createClass({
  displayName: 'Page',
  propTypes: {
    _onClick: PropTypes.func.isRequired,
    pageNumber: PropTypes.number.isRequired,
    currentPage: PropTypes.number.isRequired
  },
  render() {
    var classes = classnames({
      'current': this.props.currentPage === this.props.pageNumber
    });

    return (
      <li className={classes}>
        <a href='#'
           onClick={this.props._onClick(this.props.pageNumber)}>{this.props.pageNumber}</a>
      </li>
    );
  }
});

function mkPage(pageNumber) {
  return (
    <Page key={pageNumber}
          pageNumber={pageNumber}
          currentPage={this.props.currentPage}
          _onClick={this._onClick}/>
  );
}

export default createClass({
  displayName: 'Pagination',
  propTypes: {
    next: PropTypes.string,
    prev: PropTypes.string,
    currentPage: PropTypes.number.isRequired,
    pageSize: PropTypes.number.isRequired,
    onChangePage: PropTypes.func.isRequired
  },
  _onClick(pageNumber) {
    return (e) => {
      //Check if currentPage is to the right of ellipsis and the last from the beginning or the end
      //based on whether it is the beginning side or end side, update the page ranges
      e.preventDefault();
      this.props.onChangePage(pageNumber);
    };
  },
  render() {
    var paginationComponent;
    // is there a page before this one?
    var previousPageExists = !!this.props.prev;
    // is there a page after this one?
    var nextPageExists = !!this.props.next;
    var currentPage = [this.props.currentPage].map(mkPage, this);
    var prevClasses = classnames({
      'arrow': true,
      'unavailable': !previousPageExists
    });

    var nextClasses = classnames({
      'arrow': true,
      'unavailable': !nextPageExists
    });

    var prevPage = null;
    if (previousPageExists) {
      prevPage = [(
        <li className={prevClasses}
            key='prevPage'
            onClick={this._onClick(this.props.currentPage - 1)}><a href='#'>&laquo;</a></li>
      ), (
        <li key='prevPageNumber'
            onClick={this._onClick(this.props.currentPage - 1)}><a href='#'>{this.props.currentPage - 1}</a></li>
      )];
    }

    var nextPage = null;
    if (nextPageExists) {
      nextPage = [(
        <li key='nextPageNumber'
            onClick={this._onClick(this.props.currentPage + 1)}><a href='#'>{this.props.currentPage + 1}</a></li>
      ), (
        <li className={nextClasses}
            key='nextPage'
            onClick={this._onClick(this.props.currentPage + 1)}><a href='#'>&raquo;</a></li>
      )];
    }

    if (nextPageExists || previousPageExists) {
      paginationComponent = (
        <ul className='pagination'>
          {prevPage}
          {currentPage}
          {nextPage}
        </ul>
      );
    }

    return (
      <div className='pagination-centered'>
        {paginationComponent}
      </div>
    );
  }
});
