import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { Link } from 'react-router';
import { ArrowBackIcon } from 'common/Icon';
const { func, string } = PropTypes;

export default class BackButtonArea extends Component {
  static propTypes = {
    className: string,
    pathname: string,
    text: string,
    onClick: func,
  }

  render() {
    const { pathname, className, onClick, text } = this.props;
    if (!pathname && !!onClick) {
      return (
        <div className={`${className} ${css.wrapper}`}>
          <a onClick={onClick} className={css.link}>
            <ArrowBackIcon className={css.icon} />
            <span className={css.text}>{text}</span>
          </a>
        </div>
      );
    }
    return (
      <div className={`${className} ${css.wrapper}`}>
        <Link to={pathname} className={css.link}>
          <ArrowBackIcon className={css.icon} />
          <span className={css.text}>{text}</span>
        </Link>
      </div>
    );
  }
}
