import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { AngledTitleBox } from 'common';
import ChevronArrows from '../ChevronArrows';
import ImageSearchResult from 'marketplace/Search/ImageSearchResult';
const { array, bool, object, string } = PropTypes;
import classnames from 'classnames';
const noOp = () => {};

export default class FeaturedContentRow extends Component {
  static propTypes = {
    className: string,
    description: string.isRequired,
    headline: string.isRequired,
    images: array.isRequired,
    isFetching: bool,
    location: object.isRequired,
    title: string.isRequired,
  }

  state = {
    firstImageIndex: 0,
  }

  setFirstImageIndex = (i) => () => {
    this.setState({
      firstImageIndex: i,
    });
  }

  renderImage = (image) => {
    const { location } = this.props;
    return (
      <ImageSearchResult
        key={image.id}
        className={css.imageWrapper}
        image={image}
        location={location}
      />
    );
  }

  render() {
    const {
      description,
      headline,
      images = [],
      title,
    } = this.props;
    const { firstImageIndex } = this.state;
    const numImages = images.length;

    let numShowing = 3;
    if (typeof window !== 'undefined' &&
      typeof window.matchMedia !== 'undefined') {
      // Reduce the amount of results showing on a small screen
      const isMedium = !window.matchMedia('(min-width: 960px)').matches;
      const isSmall = !window.matchMedia('(min-width: 820px)').matches;
      if (isSmall) {
        numShowing = 1;
      } else if (isMedium) {
        numShowing = 2;
      }
    }
    const lastImageIndex = firstImageIndex + numShowing - 1;
    const displayImages = images.slice(firstImageIndex, lastImageIndex + 1);
    // The first image showing is the first image available
    const isPreviousDisabled = firstImageIndex === 0;
    // The last image that could possibly be showing is the last image available
    // or there is a not full row
    const isNextDisabled = lastImageIndex >= numImages - 1;
    const onClickNext = isNextDisabled ?
      noOp : this.setFirstImageIndex(lastImageIndex + 1);
    const onClickPrevious = isPreviousDisabled ?
      noOp : this.setFirstImageIndex(firstImageIndex - numShowing);
    const classes = classnames({
      [css.contentRow]: true,
      [this.props.className]: !!this.props.className,
    });
    return (
      <div className={classes}>
        <div className={css.contentDescription}>
          <AngledTitleBox className={css.title} title={title} />
          <div className={css.headline}>{headline}</div>
          <div className={css.description}>{description}</div>
          <ChevronArrows
            isNextDisabled={isNextDisabled}
            isPreviousDisabled={isPreviousDisabled}
            onClickNext={onClickNext}
            onClickPrevious={onClickPrevious}
          />
        </div>
        <div className={css.images}>
          {displayImages.map(this.renderImage)}
        </div>
      </div>
    );
  }
}
