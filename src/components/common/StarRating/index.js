import React, { Component, PropTypes } from 'react';
import { StarIcon } from '../Icon';
import css from './styles.css';
import map from 'lodash/map';
import range from 'lodash/range';
const { bool, func, number, string } = PropTypes;
import { SMALL } from 'lib/constants/sizes';
import classnames from 'classnames';

export default class StarRating extends Component {
  static propTypes = {
    className: string,
    disabled: bool,
    filledStarClassName: string,
    maxStars: number,
    onStarClick: func,
    rating: number,
    starClassName: string,
  }

  static defaultProps = {
    disabled: false,
    maxStars: 5,
  }

  onClick = (starNum) => {
    if (!this.props.disabled && this.props.onStarClick) {
      this.props.onStarClick(starNum);
    }
  }

  makeStar = (starNum) => {
    const { disabled, rating, starClassName, filledStarClassName } = this.props;
    const isFilled = starNum <= rating;
    const classes = classnames('dStar', {
      [css.star]: !isFilled,
      [css.filled]: isFilled,
      [css.clickable]: !disabled,
      [starClassName]: !!starClassName && !isFilled,
      [filledStarClassName]: !!filledStarClassName && isFilled,
    });
    return (
      <StarIcon
        key={starNum}
        className={classes}
        onClick={this.onClick}
        size={SMALL}
      />
    );
  }

  render() {
    const { className, maxStars } = this.props;
    return (
      <div className={className}>
        { map(range(1, maxStars + 1), this.makeStar) }
      </div>
    );
  }
}
