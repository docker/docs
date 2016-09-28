'use strict';

import React, { Component, PropTypes } from 'react';
import Tooltip from 'rc-tooltip';
import cn from 'classnames';
import styles from './label.css';

export default class TagLabel extends Component {
  static propTypes = {
    // TODO: Add an icon prop
    children: PropTypes.node.isRequired,
    tooltip: PropTypes.node,
    variant: PropTypes.oneOf([ // Note: remember to change TagLabelList also
        'accessLevel', // The access level of the current team
        'members', // List of members in teams card
        'repoAccessLevel', // The access level on the current repo
        'repository', // List of repos in orgs card
        'repositoryTag', // List of tags in repo details
        'role', // Whether you are a member of a team/org or owner of an org
        'selectedMembers' // List of selected members in create teams form
    ])
  }

  render() {
    const {
      children,
      tooltip,
      variant
    } = this.props;
    const classes = cn([styles.label, styles[variant]]);
    const content = <span className={ classes }>{ children }</span>;
    if (tooltip) {
      return (
        <Tooltip
          overlay={ <div style={ { maxWidth: 400 } }>{ tooltip }</div> }
          placement='right'
          align={ { overflow: { adjustY: 0 } } }
          trigger={ ['hover'] }>
          { content }
        </Tooltip>
      );
    }
    return content;
  }
}
