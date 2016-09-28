import React, { Component, PropTypes } from 'react';
import moment from 'moment';
import css from './styles.css';
import { bytesToSize } from 'lib/utils/formatNumbers';
const { number, shape, string } = PropTypes;

export default class TagInfoArea extends Component {
  static propTypes = {
    tag: shape({
      name: string.isRequired,
      full_size: number.isRequired,
      last_updated: string,
    }),
  }

  render() {
    const { tag } = this.props;
    if (!tag) return null;
    const { name, last_updated, full_size } = tag;
    const size = full_size ? bytesToSize(full_size) : null;
    let lastUpdated;
    if (last_updated) {
      lastUpdated = `Last updated ${moment(last_updated).fromNow()}`;
    }
    return (
      <div>
        <div className={css.titleRow}>
          <span className={css.name}>{name}</span>
          <span className={css.tagSize}>{size}</span>
        </div>
        <div className={css.lastUpdated}>{lastUpdated}</div>
      </div>
    );
  }
}
