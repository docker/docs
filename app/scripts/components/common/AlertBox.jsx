'use strict';

import React, { Component, PropTypes } from 'react';
const { oneOf, func } = PropTypes;
import { intents } from 'dux/dux/utils/props';
import classnames from 'classnames';
import _ from 'lodash';

export default class AlertBox extends Component {
  static propTypes = {
    intent: oneOf(intents),
    onClick: func
  }

  static defaultProps = {
    onClick() {}
  }

  render() {
    const { intent } = this.props;

    const classes = classnames({
      'alert-box': true,
      [intent]: _.includes(intents, intent)
    });

    return (
      <div className={classes}
           onClick={this.props.onClick}>
        {this.props.children}
      </div>
    );
  }
}
