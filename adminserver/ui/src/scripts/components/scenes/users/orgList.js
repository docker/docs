'use strict';

import React, { Component, PropTypes } from 'react';
// Actions
import { connect } from 'react-redux';
import { listUserOrganizations } from 'actions/organizations';
import autoaction from 'autoaction';
import { getUserOrgs } from 'selectors/users';
import { createStructuredSelector } from 'reselect';

const mapState = createStructuredSelector({
  orgs: getUserOrgs
});

/**
 * OrgList is used to show the list of orgs that a user belongs to.
 *
 * We load each user's orgs within a separate API call therefore this is a new
 * component made to handle this.
 *
 */
@autoaction(
  { listUserOrganizations: (props) => ({ name: props.username }) },
  { listUserOrganizations }
)
@connect(mapState)
export default class OrgList extends Component {
  static propTypes = {
    orgs: PropTypes.object
  }

  render() {
    const { orgs } = this.props;

    if (Object.keys(orgs).length === 0) {
      return null;
    }

    return (
      <div>
        {
          Object.keys(orgs).map((orgname, i) => {
            return (
              <span key={ i } style={ { marginRight: '0.5em'} }>
                { orgname }
                {
                  i === Object.keys(orgs).length - 1 ?
                    undefined
                    :
                    <span>,</span>
                }
              </span>
            );
          })
        }
      </div>
    );
  }
}
