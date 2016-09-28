'use strict';

import React, { Component, PropTypes } from 'react';
import cn from 'classnames';
import { Link } from 'react-router';
import styles from './breadCrumbs.css';

/**
 * Renders a list of breadcrumbs for any given page.
 *
 * If the `header` prop is passed as `true` this renders in the format of
 * a pageheader.
 *
 * @TODO: Stop using docker-ux' pageHeader element so this can be a nested child
 * of the pageHeader.
 *
 */
export default class BreadCrumbs extends Component {

  static propTypes = {
    header: PropTypes.bool.isRequired,
    items: PropTypes.array.isRequired
  }

  static defaultProps = {
    header: false
  }

  renderItems() {
    const { items } = this.props;
    let children = [];

    items.map((item, index) => {
      if (index !== 0) {
        children.push(<span key={ `separator-${index}` } className={ styles.separator }>&rsaquo;</span>);
      }
      if (Array.isArray(item)) {
        children.push(<Link key={ `item-${index}` } to={ item[1] }>{ item[0] }</Link>);
      } else {
        children.push(<span key={ `item-${index}` }>{ item }</span>);
      }
    });

    return children;
  }

  render() {
    const { header } = this.props;
    const classes = cn({
      [styles.outerContainer]: true,
      [styles.header]: header
    });

    return (
      <div className={ classes }>
        <div className={ styles.wrapper }>
          { this.renderItems() }
        </div>
      </div>
    );
  }
}
