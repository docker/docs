'use strict';

import React, { Component, PropTypes } from 'react';
const { string, bool } = PropTypes;
import { Link } from 'react-router';
import Tooltip from 'rc-tooltip';
import FA from 'components/common/fontAwesome';
import styles from './leftNav.css';
import css from 'react-css-modules';

@css(styles)
export default class LeftNavItem extends Component {
  static propTypes = {
    page: string.isRequired,
    pageName: string.isRequired,
    icon: string.isRequired,
    navExpanded: bool.isRequired
  }

  render() {
    const {
      page,
      pageName,
      icon,
      navExpanded
    } = this.props;
    return (
      <div>
        { navExpanded ?
          <Link to={ page } activeClassName={ styles.active }>
            <div styleName='iconContainer'>
              <FA icon={ icon } />
            </div>
            { pageName }
          </Link>
          :
          <Tooltip
            overlay={ <div style={ { maxWidth: 400 } }>{ pageName }</div> }
            placement='right'
            align={ { overflow: { adjustX: 0, adjustY: 0 } } }
            trigger={ ['hover'] }>
            <Link to={ page } activeClassName={ styles.active }>
              <div styleName='iconContainer'>
                <FA icon={ icon } />
              </div>
              { navExpanded ? pageName : '' }
            </Link>
          </Tooltip>
        }
      </div>
    );
  }
}
