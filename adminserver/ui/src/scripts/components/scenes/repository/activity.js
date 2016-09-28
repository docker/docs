'use strict';

import React, { PropTypes, Component } from 'react';
import css from 'react-css-modules';
import styles from './repository.css';
import SVG from 'components/common/svg';

@css(styles)
export default class RepositoryActivityTab extends Component {
  render() {
    const colors = ['#8f9ea8', '#ef4a53', '#00cbca'];
    return (
      <div styleName='activity'>
        {
          Array.apply(null, Array(25)).map((_, i) => {
            const choice = Math.floor(Math.random() * 3);
            return (
              <ActivityRow key={ i } color={ colors[choice] }/>
            );
          })
        }
      </div>
    );
  }
}

@css(styles)
class ActivityRow extends Component {
  static propTypes = {
    color: PropTypes.string,
    styles: PropTypes.object
  }

  render() {
    return (
      <div styleName='activityRow'>
        <div styleName='ago'>5 minutes ago</div>
        <div styleName='circleDisplay'>
          <hr />
          <SVG
            className={ this.props.styles.svg }
            viewBox='0 0 12 12'>
            <circle
              styleName='circle'
              r='3'
              cx='6px'
              cy='6px'
              style={ { stroke: this.props.color } }></circle>
          </SVG>
        </div>
        <div styleName='action'>
          <SVG
            className={ this.props.styles.svg }
            viewBox='0 0 6 15'>
            <polygon points='7,0 0,7 7,15'></polygon>
          </SVG>
          <strong>admin (you)</strong> changed the repository key and root key
        </div>
      </div>
    );
  }
}
