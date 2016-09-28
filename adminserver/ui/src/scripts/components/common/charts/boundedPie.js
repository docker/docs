'use strict';

/**
 * boundedPie visualises a single piece of bounded data using a pie chart
 *
 */
import React from 'react';
import sector from 'paths-js/sector';
import assign from 'object-assign';

/**
 * This renders a pie chart for a single data point with an upper bound (ie RAM
 * usage; current usage might be 1024mb/4096mb).
 *
 * Example usage:
 *
 *   <BoundedPie
 *     upperBound={4096}
 *     value={1024}
 *     size={300}           // Diameter of pie chart
 *     stroke={3}           // Stroke width
 *     strokeColor={'#aaa'}
 *     fillColor={'#aaa'} />
 */
export default class BoundedPie extends React.Component {

  static propTypes = {
    size: React.PropTypes.number,
    stroke: React.PropTypes.number,
    strokeColor: React.PropTypes.string,
    fillColor: React.PropTypes.string,
    upperBound: React.PropTypes.number,
    value: React.PropTypes.number
  }

  static defaultProps = {
    size: 100, // Diameter of circle
    stroke: 3,
    strokeColor: '#c4cdda',
    fillColor: '#c4cdda',

    upperBound: 0,
    value: 0
  }

  /**
   * Adds defualt properties for a paths.js sector to an object of properties
   * given.
   *
   * In all sectors drawn the outer radius and center is constant.
   *
   * You *must* provide a start and end value to mergeSectorProps to draw
   * a sector.
   */
  mergeSectorProps(props = {}) {
    const radius = this.props.size / 2;
    const stroke = Math.abs(radius - this.props.stroke);

    const defaults = {
      center: [radius, radius],
      r: radius,
      R: Math.abs(stroke)
    };

    return assign({}, defaults, props);
  }

  /**
   * Renders the outer stroke circle comprising two Sectors spanning half of the
   * circumference.
   *
   * Paths.js cannot render a Sector for a full circumference (ie. start: 0,
   * end: Math.PI * 2).  It causes the SVG to be entirely flat.
   *
   * To work around this we render two sectors for the entire circle within
   * a group.
   */
  renderStroke() {
    return (
      <g>
        <path
          d={ sector(this.mergeSectorProps({start: 0, end: Math.PI})).path.print() }
          fill={ this.props.strokeColor }/>
        <path
          d={ sector(this.mergeSectorProps({start: -Math.PI, end: 0})).path.print() }
          fill={ this.props.strokeColor }/>
      </g>
    );
  }

  /**
   * Renders the fill, shading the inner pie chart according to the percentage
   * of this.props.value in this.props.upperBound
   *
   */
  renderFill() {
    const percentage = this.props.value / this.props.upperBound;
    if (isNaN(percentage)) {
      return (<g />);
    }

    // Render a full circle.
    if (percentage >= 1) {
      return (
        <g>
          <path
            d={ sector(this.mergeSectorProps({r: 0, start: 0, end: Math.PI})).path.print() }
            fill={ this.props.fillColor }/>
          <path
            d={ sector(this.mergeSectorProps({r: 0, start: -Math.PI, end: 0})).path.print() }
            fill={ this.props.fillColor }/>
        </g>
      );
    }

    const props = {
      r: 0,
      start: 0,
      end: (Math.PI * 2) * percentage
    };
    return (
      <path
        d={ sector(this.mergeSectorProps(props)).path.print() }
        fill={ this.props.fillColor }/>
    );

  }

  render() {
    return (
      <svg width={ this.props.size } height={ this.props.size }>
        { this.renderStroke() }
        { this.renderFill() }
      </svg>
    );
  }
}
