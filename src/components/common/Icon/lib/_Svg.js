import React, { Component, PropTypes } from 'react';
import { values } from 'lodash';
import css from '../styles.css';
import { variants, sizes } from 'lib/constants';

const DEFAULT_SIZE = 24;
const DEFAULT_VIEWBOX = '0 0 24 24';

function iconThunk({ viewBox, width, height }) {
  // eslint-disable-next-line arrow-body-style
  return Shape => {
    return class SvgIcon extends Component {
      static propTypes = {
        size: PropTypes.oneOf(values(sizes)),
        variant: PropTypes.oneOf(values(variants)),
        className: PropTypes.string,
        shapeProps: PropTypes.any,
      }

      render() {
        const { variant, size, shapeProps, className = '' } = this.props;
        const styles = {};
        const classes = 'dicon' +
          ` ${className}` +
          ` ${css[size] || ''}` +
          ` ${css[variant] || ''}`;

        if (width !== DEFAULT_SIZE) {
          styles.width = width;
        }

        if (height !== DEFAULT_SIZE) {
          styles.height = height;
        }

        return (
          <svg
            xmlns="http://www.w3.org/2000/svg"
            preserveAspectRatio="xMidYMid meet"
            className={classes}
            style={styles}
            viewBox={viewBox}
          >
            <Shape {...shapeProps} />
          </svg>
        );
      }
    };
  };
}

export default function mightBeCalledWithSize(ShapeOrSize) {
  if (typeof ShapeOrSize === 'number') {
    // eslint-disable-next-line no-param-reassign
    ShapeOrSize = { viewBox: `0 0 ${ShapeOrSize} ${ShapeOrSize}` };
  }

  if (ShapeOrSize.viewBox || ShapeOrSize.width || ShapeOrSize.height) {
    return iconThunk(ShapeOrSize);
  }

  return iconThunk({
    viewBox: DEFAULT_VIEWBOX,
  })(ShapeOrSize);
}
