'use strict';

import React, { PropTypes, Component } from 'react';

/*
 * A simple component that takes one bool prop, show,
 * which will either render or not render children
 */

export default class ToggleChild extends Component {

    static propTypes = {
        children: PropTypes.node,
        show: PropTypes.bool.isRequired
    };

    render() {

        return (
            <span>{
                this.props.show ?
                this.props.children :
                undefined
            }</span>
        );
    }
}

