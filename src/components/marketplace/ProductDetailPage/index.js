import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import classnames from 'classnames';
import {
  BackButtonArea,
  CheckIcon,
  FetchingError,
  ImageWithFallback,
  LoadingIndicator,
  Markdown,
  StarIcon,
  Tab,
  Tabs,
} from 'common';
import routes from 'lib/constants/routes';
import { FALLBACK_IMAGE_SRC, FALLBACK_ELEMENT } from 'lib/constants/fallbacks';
import { CS_ENGINE_ID } from 'lib/constants/eusa';
import { formatBucketedNumber } from 'lib/utils/formatNumbers';
import formatCategories from 'lib/utils/format-categories';
import { SECONDARY } from 'lib/constants/variants';
import { SMALL } from 'lib/constants/sizes';
import getLogo from 'lib/utils/get-largest-logo';
const { arrayOf, bool, element, number, object, shape, string } = PropTypes;

// This is a dumb component used to display a Store product, such as an image or
// a bundle. Almost everything about the two content-type pages is the same,
// except for the upper right hand of the page with the versions, subscription,
// and / or purchase information.
// This component expects that section to be passed in as a React element,
// allowing it to stay mostly agnostic
export default class ProductDetailPage extends Component {
  static propTypes = {
    isLoggedInAndWhitelisted: bool,
    imageOrBundle: shape({
      categories: arrayOf(shape({
        name: string,
        label: string,
      })),
      default_version: shape({
        linux: string,
        windows: string,
      }),
      name: string,
      error: string,
      id: string,
      isFetching: bool,
      last_updated: string,
      links: arrayOf(shape({
        label: string,
        url: string,
      })),
      logo_url: object,
      namespace: string,
      platforms: arrayOf(shape({
        name: string,
        label: string,
      })),
      popularity: number,
      publisher: shape({
        id: string,
        name: string,
      }),
      reponame: string,
      screenshots: arrayOf(shape({
        label: string,
        url: string,
      })),
      short_description: string,
      source: string,
    }).isRequired,
    versionsAndSubscriptions: shape({
      isFetching: bool,
      error: bool,
      displayElement: element,
    }),
  }

  state = {
    selectedTab: 0,
  }

  onSelectTab = (e, val, index) => {
    this.setState({ selectedTab: index });
  }

  headerRightSide() {
    const {
      displayElement,
      error,
      isFetching,
    } = this.props.versionsAndSubscriptions;

    const {
      id,
    } = this.props.imageOrBundle;

    if (error) {
      return (
        <div className={css.versionsAndSubscriptionsLoading}>
          Error! Please refresh the page and try again.
        </div>
      );
    }

    if (isFetching) {
      return (
        <div className={css.versionsAndSubscriptionsLoading}>
          <LoadingIndicator />
        </div>
      );
    }

    // If this is the CS Engine bundle, do not show the
    // right side of the header.
    // TODO(mattt) 8/10/2016 - this can go away once
    // CS Engine is no longer treated as a bundle.
    if (id === CS_ENGINE_ID) {
      return (<div></div>);
    }

    return displayElement;
  }

  renderHeader() {
    const {
      categories,
      name,
      id,
      logo_url,
      publisher,
      short_description,
    } = this.props.imageOrBundle;

    const classes = classnames({
      [css.header]: true,
      [css.splitSections]: true,
    });

    const rightSide = this.headerRightSide();

    let verifiedPartner;
    let check;
    // Docker
    if (publisher && publisher.id !== '0') {
      verifiedPartner = 'Docker Verified Partner';
      check = (
        <CheckIcon
          className={css.check}
          variant={SECONDARY}
          size={SMALL}
        />
      );
    }
    const fallbackElement = React.cloneElement(FALLBACK_ELEMENT, {
      className: css.fallbackElement,
    });
    return (
      <div className={classes}>
        <div className={css.rightBorder}>
          <div className={css.inlineWithLogo}>
            <ImageWithFallback
              key={id}
              alt={`${name} logo`}
              className={css.logo}
              fallbackImage={FALLBACK_IMAGE_SRC}
              secondaryFallbackElement={fallbackElement}
              src={getLogo(logo_url)}
            />
            <div className={css.repoDetailColumn}>
              <div className={css.name}>
                {name}
              </div>
              <div className={css.dockerVerifiedPartner}>
                By {publisher.name} {check} {verifiedPartner}
              </div>
            </div>
          </div>
          <div className={css.shortDescription}>{short_description}</div>
          <div className={css.popularityAndRatings}>
            {this.renderPopularity()}
          </div>
          <div className={css.category}>
            <span className={css.categoryTitle}>Categories: </span>
            {formatCategories(categories)}
          </div>
        </div>
        <div className={css.versionsAndSubscriptions}>
          {rightSide}
        </div>
      </div>
    );
  }

  renderPopularity() {
    const { popularity } = this.props.imageOrBundle;
    if (!popularity) return null;

    const numPulls = formatBucketedNumber(popularity);
    const pullOrPulls = numPulls === '1' ? ' Pull' : ' Pulls';
    return (
      <span className={css.popularity}><b>{numPulls}</b> {pullOrPulls}</span>
    );
  }

  renderDescription() {
    const { full_description } = this.props.imageOrBundle;
    const content = !full_description
      ? 'There is no available description for this repository.'
      : <Markdown className={css.scaleImages} rawMarkdown={full_description} />;
    return (
      <div className={css.description}>
        {content}
      </div>
    );
  }

  renderRatingsAndReviews() {
    return (
      <div className={css.ratingsWrapper}>
        <div>
          <StarIcon className={css.placeholderStar} />
          <StarIcon className={css.placeholderStar} />
          <StarIcon className={css.placeholderStar} />
          <StarIcon className={css.placeholderStar} />
          <StarIcon className={css.placeholderStar} />
        </div>
        <div>Ratings and Reviews</div>
        <div>Coming Soon!</div>
      </div>
    );
  }

  renderTabs() {
    const { selectedTab } = this.state;
    let tabContent;
    if (selectedTab === 0) {
      tabContent = (
        <div>
          {this.renderDescription()}
          {this.renderScreenshots()}
        </div>
      );
    } else {
      tabContent = this.renderRatingsAndReviews();
    }
    return (
      <div className={css.card}>
        <Tabs
          className={css.tabsWrapper}
          selected={selectedTab}
          onSelect={this.onSelectTab}
        >
          <Tab value={0} className={css.tabWidth} isUnderlined>
            Description
          </Tab>
          <Tab value={1} className={css.tabWidth} isUnderlined>
            Reviews
          </Tab>
        </Tabs>
        <div className={css.tabContent}>
          {tabContent}
        </div>
      </div>
    );
  }

  renderScreenshots() {
    const { screenshots } = this.props.imageOrBundle;
    if (!screenshots || !screenshots.length) {
      return null;
    }
    const fallbackElement = (
      <div className={css.fallbackElement}>
        Screenshot temporarily unavailable
      </div>
    );
    return (
      <div className={css.screenshots}>
        {screenshots.map(({ url, label }) => {
          return (
            <ImageWithFallback
              key={url}
              src={url}
              alt={label}
              fallbackImage=""
              className={css.screenshot}
              secondaryFallbackElement={fallbackElement}
            />
          );
        })
        }
      </div>
    );
  }

  render() {
    const {
      error,
      isFetching,
      links,
    } = this.props.imageOrBundle;
    const { isLoggedInAndWhitelisted } = this.props;
    const backButtonProps = {
      pathname: isLoggedInAndWhitelisted ? routes.home() : routes.beta(),
      text: 'Home',
    };
    // TODO Kristie 3/28/16 Proper 404 Page handling
    if (error || isFetching) {
      const content = error ? <FetchingError resource="this product" /> :
        <div>Fetching...</div>;
      return (
        <div className="wrapped">
          <BackButtonArea {...backButtonProps} />
          {content}
        </div>
      );
    }
    const partnerLinks = links && links.map(({ url, label }) => {
      return <a href={url} target="_blank" key={url}>{label}</a>;
    });
    return (
      <div>
        <div className={css.fullWhiteBackground}>
          <div className="wrapped">
            <BackButtonArea {...backButtonProps} />
            {this.renderHeader()}
          </div>
        </div>
        <div className="wrapped">
          <div className={`${css.splitSections} ${css.descriptionSection}`}>
            <div>
              {this.renderTabs()}
            </div>
            <div className={css.partnerLinks}>{partnerLinks}</div>
          </div>
        </div>
      </div>
    );
  }
}
