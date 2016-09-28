import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
const { func, string, node } = PropTypes;
import css from './styles.css';

/* Server-side-rendering-friendly component that shows a fallback image if
 * the attempted src image does not load.
 */
export default class ImageWithFallback extends Component {
  static propTypes = {
    alt: string,
    className: string,
    fallbackImage: string.isRequired,
    onError: func,
    onLoad: func,
    secondaryFallbackElement: node,
    src: string.isRequired,
  }

  constructor(props) {
    super(props);
    // initialize the image to have the passed in src
    const { src, fallbackImage, secondaryFallbackElement } = props;
    // if src doesn't exist, go directly to trying the fallbackImage
    const shouldTryFallbackImage = !src && !!fallbackImage;
    const shouldTrySecondaryFallbackElement =
      !src && !fallbackImage && !!secondaryFallbackElement;
    this.state = {
      didFallbackFail: shouldTrySecondaryFallbackElement,
      didImageLoad: shouldTrySecondaryFallbackElement,
      imageSrc: shouldTryFallbackImage ? fallbackImage : src,
      isTryingFallbackImage: shouldTryFallbackImage,
    };
  }

  onLoad = () => {
    this.setState({ didImageLoad: true }, () => {
      if (this.props.onLoad) {
        this.props.onLoad();
      }
    });
  }

  setFallback = () => {
    const { fallbackImage, secondaryFallbackElement } = this.props;
    const { isTryingFallbackImage } = this.state;
    // The fallback image has failed
    if (isTryingFallbackImage && secondaryFallbackElement) {
      this.setState({
        didFallbackFail: true,
        imageSrc: '',
        isTryingFallbackImage: false,
      });
    } else {
      // The initial image has failed - try the fallback
      this.setState({
        imageSrc: fallbackImage,
        isTryingFallbackImage: true,
      });
    }
    this.setState({
      imageSrc: fallbackImage,
    });
  }

  render() {
    const {
      alt = '',
      className = '',
      onError,
      secondaryFallbackElement,
    } = this.props;
    const { didFallbackFail, isTryingFallbackImage, didImageLoad } = this.state;
    const showSecondaryFallback = !isTryingFallbackImage && didFallbackFail;
    // In the event of an error loading the image, the onError function on the
    // image will fire, and will bubble up to the parent div, which will call
    // setFallback()
    let shownImage;
    if (showSecondaryFallback) {
      shownImage = secondaryFallbackElement;
    } else {
      const classes = classnames({
        [className]: !!className,
        [css.hideAltText]: true,
        [css.hideFadeInFlash]: !didImageLoad,
        [css.fadeInImage]: didImageLoad,
      });
      // NOTE: Must add the trailing ? on the image urls to trigger the onerror
      // event in FireFox. This has no impact on Chrome
      // eslint-disable-next-line
      // http://stackoverflow.com/questions/26621325/javascript-onerror-event-doesnt-work-in-firefox
      shownImage = (
        <img
          alt={alt}
          className={classes}
          onError={onError}
          onLoad={this.onLoad}
          src={`${this.state.imageSrc}?`}
        />
      );
    }

    return (
      <span onError={this.setFallback}>
        {shownImage}
      </span>
    );
  }
}
