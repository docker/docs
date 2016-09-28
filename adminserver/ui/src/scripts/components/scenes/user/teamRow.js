'use strict';

import React, { Component, PropTypes } from 'react';
const { string, object } = PropTypes;
import autoaction from 'autoaction';
import { getTeam } from 'actions/teams';
import { teamDetails } from 'selectors/teams';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';


const mapState = createStructuredSelector({
  team: teamDetails
});

@connect(mapState)
@autoaction({
  getTeam: (props) => {
    return [props.org, `id:${props.id}`];
  }
}, {
  getTeam
})
export default class TeamRow extends Component {

  static propTypes = {
    org: string,
    team: object,
    params: object
  };


  render () {

    const {
      team,
      org
      } = this.props;

    return (
      <tr>
        <td>
          { org }
        </td>
        <td>
          { team.get('name') }
        </td>
      </tr>
    );
  }
}
