'use strict';

import React, { Component, PropTypes } from 'react';
import FontAwesome from 'components/common/fontAwesome';

import HoverDelete from 'components/common/hoverDelete';
import Pagination from 'components/common/pagination';

import styles from './resource.css';

export default class MemberList extends Component {
  static propTypes = {
    canDelete: PropTypes.bool,
    onDeleteTooltip: PropTypes.func, // Takes in member and returns PropTypes.node
    members: PropTypes.array.isRequired,
    onDelete: PropTypes.func,
    location: PropTypes.object
  }

  render() {
    const {
      canDelete,
      onDeleteTooltip,
      members,
      onDelete
    } = this.props;
    return (
      <Pagination
        location={ location }
        className={ styles.memberList } pageSize={ 10 }>
        {
          members.map((member) => {
            return (
              <div className={ styles.memberItem } key={ member.name }>
                <HoverDelete
                  canDelete={ canDelete }
                  onDelete={ () => onDelete(member) }
                  tooltip={ onDeleteTooltip(member) }
                  >
                  <div>
                      <h3 style={ { wordBreak: 'break-all' } }><FontAwesome icon='fa-user' />{ member.name }</h3>
                  </div>
                </HoverDelete>
              </div>
            );
          })
        }
      </Pagination>
    );
  }
}
