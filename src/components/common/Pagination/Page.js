import React, { PropTypes, Component } from 'react';
import css from './styles.css';
import classnames from 'classnames';
import { ChevronIcon, DoubleChevronIcon } from 'common';
const { func, number, oneOf } = PropTypes;
import { SMALL } from 'lib/constants/sizes';
const FIRST = 'first';
const LAST = 'last';
const NEXT = 'next';
const PREVIOUS = 'previous';
const types = [FIRST, LAST, NEXT, PREVIOUS];

export default class Page extends Component {
  static propTypes = {
    arrowType: oneOf(types),
    currentPage: number.isRequired,
    onClick: func.isRequired,
    pageNumber: number.isRequired,
  }

  onClick = () => {
    const { pageNumber } = this.props;
    this.props.onClick(pageNumber);
  }

  render() {
    const { arrowType, currentPage, pageNumber } = this.props;
    const classes = classnames({
      dpage: !arrowType,
      [css.page]: true,
      [css.currentPage]: currentPage === pageNumber,
      [css.chevron]: !!arrowType,
      [css.firstPage]: arrowType === FIRST,
      [css.nextPage]: arrowType === NEXT,
      [css.previousPage]: arrowType === PREVIOUS,
    });
    let icon;
    if (arrowType) {
      if (arrowType === PREVIOUS || arrowType === NEXT) {
        icon = <ChevronIcon size={SMALL} />;
      } else {
        icon = <DoubleChevronIcon size={SMALL} />;
      }
    }
    const display = !!arrowType ? icon : pageNumber;
    return (
      <li className={classes} onClick={this.onClick}>
        {display}
      </li>
    );
  }
}
