import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import CircleSVG from './CircleSVG';

const { number, string } = PropTypes;

export default class CircleGraph extends Component {
  static propTypes = {
    displayValue: number.isRequired,
    label: string.isRequired,
    percentageFill: number.isRequired,
    size: number,
  }

  static defaultProps = {
    size: 200,
  }

  render() {
    const {
      displayValue,
      label,
      percentageFill,
      size,
    } = this.props;
    return (
      <div className={css.wrapper}>
        <CircleSVG
          percentageFill={percentageFill}
          size={size}
          displayValue={displayValue}
          label={label}
        />
      </div>
    );
  }
}

export default CircleGraph;
