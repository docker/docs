import React, { Component, PropTypes } from 'react';
import {
  Card,
  CheckIcon,
  ImageWithFallback,
  LinuxIcon,
  // StarRating,
  WindowsIcon,
} from 'common';
import { TINY } from 'lib/constants/sizes';
import { FALLBACK_IMAGE_SRC, FALLBACK_ELEMENT } from 'lib/constants/fallbacks';
import { DOWNLOAD_ATTRIBUTES } from 'lib/constants/eusa';
import css from './styles.css';
import { formatBucketedNumber } from 'lib/utils/formatNumbers';
import getLogo from 'lib/utils/get-largest-logo';
import formatCategories from 'lib/utils/format-categories';
import routes from 'lib/constants/routes';
const {
  arrayOf,
  func,
  node,
  number,
  object,
  shape,
  string,
} = PropTypes;

export default class ImageSearchResult extends Component {
  static propTypes = {
    className: string,
    image: shape({
      categories: arrayOf(shape({
        name: string,
        label: string,
      })),
      default_version: shape({
        linux: string,
        windows: string,
      }),
      last_updated: string,
      logo_url: shape({
        small: string,
      }),
      platforms: arrayOf(shape({
        name: string,
        label: string,
      })),
      popularity: number,
      publisher: shape({
        id: string,
        name: string,
      }),
      name: string.isRequired,
      short_description: string,
      source: string,
    }),
    location: shape({
      pathname: string,
      query: object,
      state: node,
    }),
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  onClick = () => {
    const { id } = this.props.image;
    const pathname = routes.imageDetail({ id });
    this.context.router.push({ pathname });
  }

  renderPlatformIcons(platforms) {
    // TODO Kristie 5/16/15 Platform icon for ARM
    let windows;
    let linux;
    if (!platforms) {
      return null;
    }
    platforms.forEach(({ name }) => {
      if (name === 'windows') {
        windows = <WindowsIcon className={css.icon} />;
      }
      if (name === 'linux') {
        linux = <LinuxIcon className={css.icon} />;
      }
    });
    return (
      <div className={css.platformIcons}>
        {linux}
        {windows}
      </div>
    );
  }

  renderPopularity(popularity) {
    if (!popularity) {
      return null;
    }
    const numPulls = formatBucketedNumber(popularity);
    const pullOrPulls = numPulls === '1' ? ' Pull' : ' Pulls';
    return (
      <div className={css.popularity}>
        {`${numPulls} ${pullOrPulls}`}
      </div>
    );
  }

  renderRatingsAndPopularity({ popularity }) {
    const popularityArea = this.renderPopularity(popularity);
    // <StarRating rating={rating} disabled />
    return (
      <div className={css.ratingsAndPopularity}>
        {popularityArea}
      </div>
    );
  }

  renderCardFooter({ platforms, download_attribute }) {
    const priceType = download_attribute === DOWNLOAD_ATTRIBUTES.POST_EUSA ?
      'Partner Licensed' : 'FREE';
    return (
      <div className={css.cardFooter}>
        <div>{priceType}</div>
        {this.renderPlatformIcons(platforms)}
      </div>
    );
  }

  renderCardHeader({
    name,
    logo_url,
    popularity,
    publisher,
  }) {
    const fallbackElement = React.cloneElement(FALLBACK_ELEMENT, {
      className: css.fallbackElement,
    });
    const logo = (
      <div>
        <ImageWithFallback
          alt={`${name} logo`}
          className={css.logo}
          fallbackImage={FALLBACK_IMAGE_SRC}
          secondaryFallbackElement={fallbackElement}
          src={getLogo(logo_url)}
        />
      </div>
    );
    const check = <CheckIcon size={TINY} className={css.check} />;
    return (
      <div>
        <div className={css.cardHeader}>
          {logo}
          <div className={css.name}>{name}</div>
          <div className={css.publisherName}>
            {publisher && publisher.name && check}
            {publisher && publisher.name}
          </div>
        </div>
        {this.renderRatingsAndPopularity({ popularity })}
      </div>
    );
  }

  renderCardCenter({ categories, short_description }) {
    const showCategories = formatCategories(categories);
    let truncatedShortDescription = short_description;
    const maxChars = 90;
    if (short_description && short_description.length > maxChars) {
      truncatedShortDescription =
        `${short_description.substring(0, maxChars)}...`;
    }
    return (
      <div className={`${css.cardCenter} ${css.noResizing}`}>
        <div className={`${css.shortDescription} ${css.noResizing}`}>
          {truncatedShortDescription}
        </div>
        <div className={`${css.category} ${css.noResizing}`}>
          {showCategories}
        </div>
      </div>
    );
  }

  render() {
    const { className, image } = this.props;
    const classes = className ? `${css.cardWrapper} ${className}`
      : `${css.cardWrapper}`;
    return (
      <div onClick={this.onClick} className={classes}>
        <Card shadow hover>
          {this.renderCardHeader(image)}
          <hr />
          {this.renderCardCenter(image)}
          <hr />
          {this.renderCardFooter(image)}
        </Card>
      </div>
    );
  }
}
