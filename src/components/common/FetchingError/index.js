import React, { PropTypes, Component } from 'react';
import { Card, WarningIcon } from 'common';
import css from './styles.css';
import { DULL } from 'lib/constants/variants';
const { string } = PropTypes;

// TODO Kristie 5/26/16 This should NOT be the final error handling
// Component.
export default class FetchingError extends Component {
  static propTypes = {
    resource: string,
  }

  static defaultProps = {
    resource: 'this data',
  }

  render() {
    const text = [
      'Oops! We are currently unable to fetch',
      `${this.props.resource}.`,
      'Please refresh the page and try again.',
    ].join(' ');
    return (
      <Card shadow>
        <div className={css.content}>
          <WarningIcon className={css.icon} variant={DULL} />
          {text}
        </div>
      </Card>
    );
  }
}
