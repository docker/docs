import React, { Component, PropTypes } from 'react';
import moment from 'moment';
import { ClockIcon } from '../Icon';
import { sizes } from 'lib/constants';
import css from './styles.css';

const { TINY: iconSize } = sizes;

const {
  oneOfType,
  instanceOf,
  string,
  number,
} = PropTypes;

const DEFAULT_INTERVAL = 60 * 1000; // one minute

export default class Uptime extends Component {
  static propTypes = {
    interval: number,
    since: oneOfType([string, number, instanceOf(Date)]),
    className: string,
    prefix: string,
  }

  state = {}

  componentWillMount() {
    this.setState({ displayDate: this.displayDate() });
  }

  componentDidMount() {
    this.interval = setInterval(() => {
      this.setState({ displayDate: this.displayDate() });
    }, this.props.interval || DEFAULT_INTERVAL);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  displayDate() {
    if (this.props.since) {
      return moment(new Date(this.props.since)).fromNow();
    }

    return 'never';
  }

  render() {
    const { className = '' } = this.props;
    return (
      <span className={`duptime ${className}`}>
        <ClockIcon size={iconSize} />
        <span className={css.display}>
          { this.props.prefix ? `${this.props.prefix} ` : ''}
          { this.state.displayDate }
        </span>
      </span>
    );
  }
}
