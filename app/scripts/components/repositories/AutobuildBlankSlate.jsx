'use strict';

import React, { PropTypes } from 'react';
import Card, { Block } from '@dux/element-card';
import FA from 'common/FontAwesome';
import styles from './AutobuildBlankSlate.css';
import { Link } from 'react-router';

var AutobuildBlankSlate = React.createClass({
  displayName: 'AutobuildBlankSlate',
  propTypes: {
    slateItems: React.PropTypes.element
  },
  render() {
    var slateItems = null;
    if (!this.props.slateItems) {
      slateItems = (
        <div>
          <h2>You haven't linked to <FA icon='fa-github'/> GitHub or <FA icon='fa-bitbucket'/> Bitbucket yet.</h2>
          <Link to="/account/authorized-services/" className="button primary">Link Accounts</Link>
        </div>
      );
    } else {
      slateItems = this.props.slateItems;
    }

    return (
      <div className={`row ${styles.row}`}>
        <div className='large-10 large-centered columns'>
          <Card>
            <Block>
              {slateItems}
            </Block>
          </Card>
        </div>
      </div>
    );
  }
});

module.exports = AutobuildBlankSlate;
