'use strict';

import React from 'react';
import { Link } from 'react-router';
import css from 'react-css-modules';
import styles from './teamList.css';

/**
 * Renders the list of teams within the left panel of the organization page.
 * Note that this automatically renders the 'Owners' team; you do not need to
 * include this in the teams array.
 *
 * @param teams        array   Array of TeamRecords
 * @param orgName      string  Organization name for the teams
 * @param params       object  Object of router state so we cna detect whether
 *                             a team is selected
 */
const TeamList = ({ teams = [], orgName = '' }) => {

  return (
    <div styleName='list'>
      { teams.map(t => (
        <Link to={ `/orgs/${orgName}/teams/${t.name}` }
          key={ t.name }
          activeClassName={ styles.active }>
            { t.name }
        </Link>
      )) }
    </div>
  );
};

export default css(TeamList, styles);
