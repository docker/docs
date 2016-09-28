'use strict';

import React, {
  Component,
  PropTypes
} from 'react';
import { findDOMNode } from 'react-dom';
const { object, number, string, array, func } = PropTypes;

import has from 'lodash/object/has';
import connectToStores from 'fluxible-addons-react/connectToStores';

import RepoSettingsCollaborators from 'stores/RepoSettingsCollaborators';
import Collaborators from './collaborators/Collaborators.jsx';
import Teams from './collaborators/Teams.jsx';

const debug = require('debug')('RepoCollaborators');

class Wrapper extends Component {
  render() {
    return (
      <div>
        <div className='row'>
          <div className='large-12 columns'>
            {this.props.children}
          </div>
        </div>
      </div>
    );
  }
}

class CollaboratorsWrapper extends Component {

  static propTypes = {
    collaborators: object,
    teams: object,
    allTeams: object,
    namespace: string.isRequired,
    name: string.isRequired
  }

  render() {
    const {
      JWT,
      allTeams,
      newCollaborator,
      collaborators,
      error,
      name,
      namespace,
      teams,
      requests,
      STATUS,
      location,
      history
    } = this.props;

    /**
     * This is how we detect whether the repo is an org's or user's
     * collaborators will only  be populated if the repo belongs to a user
     * allTeams/teams will only be populated if the repo belongs to an org
     */
    if(has(collaborators, 'count')) {
      return (
        <Wrapper>
          <Collaborators {...collaborators}
                         newCollaborator={newCollaborator}
                         error={error}
                         namespace={namespace}
                         name={name}
                         JWT={JWT}
                         requests={requests}
                         STATUS={STATUS}
                         location={location}
                         history={history}/>
        </Wrapper>
      );
    } else {
      return (
        <Wrapper>
          <Teams {...teams}
                 namespace={namespace}
                 name={name}
                 allTeams={allTeams.results}
                 JWT={JWT}
                 STATUS={STATUS}
                 requests={requests}
                 location={location}
                 history={history}/>
        </Wrapper>
      );
    }
  }


}

export default connectToStores(CollaboratorsWrapper,
                               [
                                 RepoSettingsCollaborators
                               ],
                               ({ getStore }, props) => {
                                 return getStore(RepoSettingsCollaborators).getState();
                               });
