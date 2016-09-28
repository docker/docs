import React, { Component, PropTypes } from 'react';
import css from './styles.css';
const { string } = PropTypes;

export default class AngledTitleBox extends Component {
  static propTypes = {
    className: string,
    title: string.isRequired,
  }
  render() {
    return (
      <div className={this.props.className || ''}>
        <span className={css.box}>
          { this.props.title }
        </span>
      </div>
    );
  }
}
