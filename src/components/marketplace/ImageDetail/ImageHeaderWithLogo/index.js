import React, { Component, PropTypes } from 'react';
import css from './styles.css';
const { node, string } = PropTypes;

export default class RepositoryTagsListHeader extends Component {
  static propTypes = {
    logo: string,
    rightSide: node,
    title: string.isRequired,
    helpText: node,
  }

  render() {
    const { helpText, logo, rightSide, title } = this.props;
    let img;
    if (logo) {
      img = <img className={css.logo} src={logo} alt="logo" />;
    }
    return (
      <div className={css.header}>
        <div className={css.leftSide}>
          <div>{ img }</div>
          <div className={css.titleSpacing}>
            <div className={css.title}>{title}</div>
            <div className={css.helpText}>{helpText}</div>
          </div>
        </div>
        <div className={css.rightSide}>{rightSide}</div>
      </div>
    );
  }
}
