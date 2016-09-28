import React, { Component, PropTypes } from 'react';
import css from './circle.css';
const { number, string } = PropTypes;
export default class CircleSVG extends Component {
  static propTypes = {
    backgroundColor: string,
    displayValue: number.isRequired,
    fillColor: string,
    label: string.isRequired,
    percentageFill: number.isRequired,
    size: number,
    strokeWidth: number,
  }

  static defaultProps = {
    backgroundColor: '#c0c9ce',
    fillColor: '#00cbca',
    strokeWidth: 8,
  }

  render() {
    const {
      backgroundColor,
      displayValue,
      fillColor,
      label,
      size,
      strokeWidth,
      percentageFill,
    } = this.props;
    const center = size / 2;
    const radius = (size / 2) - strokeWidth;
    const circumference = Math.PI * radius * 2;
    let percentFilled = ((100 - percentageFill) / 100) * circumference;
    if (isNaN(percentFilled)) {
      percentFilled = 0;
    }
    const backgroundStyles = {
      stroke: backgroundColor,
      strokeWidth,
    };
    const fillStyles = {
      stroke: fillColor,
      strokeDasharray: circumference,
      strokeDashoffset: percentFilled,
      strokeWidth,
    };
    return (
      <svg
        width={`${size}`}
        height={`${size}`}
        xmlns="http://www.w3.org/2000/svg"
      >
        <circle
          className={css.circle}
          cx={center}
          cy={center}
          r={radius}
          style={backgroundStyles}
        />
        <circle
          className={css.circle}
          cx={center}
          cy={center}
          r={radius}
          style={fillStyles}
        />
        <text y={center}>
          <tspan
            x={center}
            textAnchor="middle"
            className={css.displayValue}
          >
            {displayValue}
          </tspan>
          <tspan
            x={center}
            textAnchor="middle"
            dy="24"
            className={css.label}
          >
            {label}
          </tspan>
        </text>
      </svg>
    );
  }
}
