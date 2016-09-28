'use strict';

import React, { Component, PropTypes } from 'react';
import MG from 'metrics-graphics';
import cn from 'classnames';
import ReactDOM from 'react-dom';

import styles from './line.css';

export default class LineChart extends Component {
  static propTypes = {
    data: PropTypes.array.isRequired,

    height: PropTypes.number,
    width: PropTypes.number,

    xAxis: PropTypes.bool,
    yAxis: PropTypes.bool,
    xAccessor: PropTypes.string.isRequired,
    yAccessor: React.PropTypes.string.isRequired,

    title: PropTypes.string,
    smallText: PropTypes.bool,

    area: PropTypes.bool,
    className: PropTypes.oneOfType([PropTypes.string, PropTypes.array])
  }

  static defaultProps = {
    area: true,
    xAxis: true,
    yAxis: true,
    smallText: false
  }

  componentDidMount() {
    this.drawGraph(this.props);
  }

  // Re-renders the graph if the last data item has changed
  componentWillReceiveProps(next) {
    // Small performance optimisation on components that often re-render:
    // re-rendering always passes down new props, and because we have no state
    // in render() we have to manually detect whether to redraw our graph.
    //
    // this compares the last entries of nextProps and currentProps to see if
    // they're the same.  If they are, don't re-render the graph.
    //
    // Once we move to using raw SVG graphs within render() using
    // this.props.data directly this will not be necessary.
    //
    // TODO: Fix this shit man.  This is naaaaasty
    const newLast = next.data[next.data.length - 1];
    const last = this.props.data[this.props.data.length - 1];
    if (JSON.stringify(newLast) !== JSON.stringify(last)) {
      this.drawGraph(next);
    }
  }

  drawGraph(props) {
    const el = ReactDOM.findDOMNode(this.refs.graph);
    const legendEl = ReactDOM.findDOMNode(this.refs.legend);
    let settings = {
      title: props.title,
      data: props.data,
      full_height: true,
      full_width: true,
      x_axis: props.xAxis,
      x_accessor: props.xAccessor,
      y_accessor: props.yAccessor,
      target: el,
      small_text: props.smallText,
      interpolate: 'linear',
      area: props.area,
      legend: props.legend,
      legend_target: legendEl,
      animate_on_load: false,
      transition_on_update: false
    };

    if (this.props.height) {
      settings.full_height = false;
      settings.height = this.props.height;
    }

    if (this.props.width) {
      settings.full_width = false;
      settings.width = this.props.width;
    }

    this.graph = MG.data_graphic(settings);
  }

  render() {
    const legendClasses = cn({
      'legend': true,
      [styles.legend]: true,
      [styles.legendSmall]: this.props.smallText,
      [styles.legendWithoutAxis]: !this.props.xAxis
    });

    return (
      <div className={ styles.container }>
        <div className={ this.props.className } ref='graph' />
        <span ref='legend' className={ legendClasses }></span>
      </div>
    );
  }
}
