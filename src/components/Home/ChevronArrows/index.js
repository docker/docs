import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { ChevronIcon } from 'common';
import classnames from 'classnames';
import { LARGE } from 'lib/constants/sizes';
const { bool, func } = PropTypes;

export default class ChevronArrows extends Component {
  static propTypes = {
    isNextDisabled: bool,
    isPreviousDisabled: bool,
    onClickNext: func.isRequired,
    onClickPrevious: func.isRequired,
  }

  _onClickNext = () => {
    if (!this.props.isNextDisabled) {
      this.props.onClickNext();
    }
  }

  _onClickPrevious = () => {
    if (!this.props.isPreviousDisabled) {
      this.props.onClickPrevious();
    }
  }

  render() {
    const { isNextDisabled, isPreviousDisabled } = this.props;
    const nextClasses = classnames({
      [css.icon]: true,
      [css.disabled]: isNextDisabled,
    });
    const previousClasses = classnames({
      [css.icon]: true,
      [css.disabled]: isPreviousDisabled,
    });
    return (
      <div>
        <div className={css.previous} onClick={this._onClickPrevious}>
          <ChevronIcon className={previousClasses} size={LARGE} />
        </div>
        <div className={css.next} onClick={this._onClickNext}>
          <ChevronIcon className={nextClasses} size={LARGE} />
        </div>
      </div>
    );
  }
}
