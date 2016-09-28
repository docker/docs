'use strict';

import React, { PropTypes, Component } from 'react';

/*
* A simple component to make onMouseOver and onMouseOut toggles less verbose
*
* expects one prop, hover, a function that is called on both mouse over and mouse out
*
*/

export default class Hover extends Component {
    static propTypes = {
        children: PropTypes.node,
        hover: PropTypes.func.isRequired
    };

    render() {

        return (
            <span
                onMouseOver={ this.props.hover }
                onMouseOut={ this.props.hover }>{
                this.props.children
            }</span>
        );
    }
}
