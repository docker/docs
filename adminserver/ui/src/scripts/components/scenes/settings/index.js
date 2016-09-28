/**
 * This is the parent container for all settings scenes
 */
'use strict';

import React from 'react';
import { Link } from 'react-router';
import { Tabs, Tab } from 'components/common/tabs';
import Spinner from 'components/common/spinner';
import autoaction from 'autoaction';
import { getSettings } from 'actions/settings';
import consts from 'consts';
import css from 'react-css-modules';
import styles from 'components/scenes/settings/formstyle.css';

@autoaction({
    getSettings: []
}, {
    getSettings
})
@css(styles)
export default class Settings extends React.Component {

    static propTypes = {
        children: React.PropTypes.node
    }

    render() {

        const status = [
            [consts.settings.ALL_SETTINGS]
        ];

        return (
            <Spinner loadingStatus={ status } >
                <div styleName='wrapper'>
                  <Tabs header>
                    <Tab><Link to='/admin/settings/general'>General</Link></Tab>
                    <Tab><Link to='/admin/settings/storage'>Storage</Link></Tab>
                    <Tab><Link to='/admin/settings/gc'>Garbage collection</Link></Tab>
                  </Tabs>

                 { this.props.children }
                </div>
           </Spinner>
    );
  }

}
