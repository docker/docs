'use strict';

import React, { Component } from 'react';
import Tooltip from 'rc-tooltip';
import FontAwesome from 'components/common/fontAwesome';

export default class InfoIcon extends Component {
  static propTypes = {
    children: React.PropTypes.node,
    placement: React.PropTypes.string
  }

  static defaultProps = {
    placement: 'right'
  }

  render() {
    const {
      children,
      placement
    } = this.props;
    return (
      <Tooltip
        overlay={ <div style={ { maxWidth: 400 } }>{ children }</div> }
        placement={ placement }
        align={ { overflow: { adjustX: 0, adjustY: 0 } } }
        trigger={ ['hover'] }>
        <FontAwesome icon='fa-question-circle' />
      </Tooltip>
    );
  }
}
