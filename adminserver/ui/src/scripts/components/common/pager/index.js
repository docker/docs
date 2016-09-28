'use strict';

import React, { Component, PropTypes } from 'react';
import css from 'react-css-modules';
import styles from './pager.css';
import { getPage } from 'utils';
import FontAwesome from 'components/common/fontAwesome';

const {
  array,
  number,
  object,
  func
  } = PropTypes;

@css(styles)
export default class Pager extends Component {

  static propTypes = {
    pageable: array.isRequired,
    perPage: number.isRequired,
    ui: object,
    updateUI: func,
    onClick: func
  }

  static contextTypes = {
    router: object
  }

  calcPages = () => {
    return Math.ceil(this.props.pageable.length / this.props.perPage);
  }

  choosePage = (choice) => () => {
    this.context.router.push({
      pathname: location.pathname,
      query: {
        page: choice + 1
      }
    });
    if (this.props.onClick) {
      this.props.onClick(choice);
    }
  }

  render() {

    const {
      search
    } = location;

    const page = getPage(search);
    // zero indexed to match page
    const pageCount = this.calcPages() - 1;

    return (
      <ul styleName='pager'>
        <li
          styleName={ page === 0 ? 'disabled' : undefined }
          onClick={ page !== 0 ? ::this.choosePage(page - 1) : undefined }>
          <FontAwesome icon='fa-angle-left'/>
        </li>
        {
          // add one that we removed above
          [...Array(pageCount + 1)].map((n, i) => {
            return (
              <li
                onClick={ ::this.choosePage(i) }
                key={ i }
                styleName={ i === page ? 'active' : undefined }>
                { i + 1 }
              </li>
            );
          })
        }
        <li
          styleName={ page === pageCount ? 'disabled' : undefined }
          onClick={ page !== pageCount ? ::this.choosePage(page + 1) : undefined }>
          <FontAwesome icon='fa-angle-right'/>
        </li>
      </ul>
    );
  }
}
