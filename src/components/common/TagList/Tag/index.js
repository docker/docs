import React, { Component, PropTypes } from 'react';
import './styles.css';

export default class Tag extends Component {
  static propTypes = {
    name: PropTypes.string.isRequired,
    last: PropTypes.bool.isRequired,
  }

  render() {
    const {
      name,
      last,
    } = this.props;

    let tagName = name;

    if (!last) {
      tagName = `${tagName},`;
    }

    return (
      <div className={'dTag'}>
        {tagName}
      </div>
    );
  }
}
