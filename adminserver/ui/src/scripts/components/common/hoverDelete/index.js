'use strict';

import React, { Component } from 'react';
import Tooltip from 'rc-tooltip';
import FontAwesome from 'components/common/fontAwesome';
import styles from './styles.css';

export default class HoverDelete extends Component {
  static propTypes = {
    children: React.PropTypes.any,
    canDelete: React.PropTypes.bool,
    tooltip: React.PropTypes.node,
    onDelete: React.PropTypes.func
  }

  static defaultProps = {
    canDelete: true
  }

  render() {
    let {canDelete, onDelete, tooltip, ...props} = this.props;
    if (!canDelete) {
      return (
        <div { ...props } >
          { this.props.children }
        </div>
      );
    }
    props.className = styles.hoverDelete + ' ' + (props.className || '');
    return (
      <div { ...props } >
        <div className={ styles.deleteButton } onClick={ onDelete }>
          <Tooltip
            mouseEnterDelay={ 0.5 }
            overlay={ <div style={ { maxWidth: 400 } }>{ tooltip }</div> }
            placement='right'
            align={ { overflow: { adjustY: 0 } } }
            trigger={ ['hover'] }>
            <FontAwesome icon='fa-times-circle' />
          </Tooltip>
        </div>
        <div className={ styles.content }>
          { this.props.children }
        </div>
      </div>
    );
  }
}
