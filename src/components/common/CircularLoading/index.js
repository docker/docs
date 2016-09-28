import React, { PropTypes } from 'react';
import { values } from 'lodash/object';
import { sizes } from 'lib/constants';
import css from './styles.css';

const CircularLoading = (props) => {
  const loadClasses = {
    tiny: css.tinyLoader,
    small: css.smallLoader,
    regular: css.regularLoader,
    large: css.largeLoader,
    xlarge: css.xlargeLoader,
  };
  const loaderClass = loadClasses[props.size];

  return (
    <div className={`${loaderClass} ${props.className}`}>
      { props.children }
    </div>
  );
};

CircularLoading.propTypes = {
  children: PropTypes.node,
  className: PropTypes.string,
  size: PropTypes.oneOf(values(sizes)),
};

CircularLoading.defaultProps = {
  size: 'regular',
};

export default CircularLoading;
