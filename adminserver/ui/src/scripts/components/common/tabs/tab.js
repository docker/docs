'use strict';

import React, {
    Component,
    PropTypes,
    Children,
    cloneElement
} from 'react';

import styles from './tabs.css';

export default class Tab extends Component {

    static propTypes = {
        children: PropTypes.node.isRequired,
        // hook for tests
        id: PropTypes.string
    }

    /**
     * If the given React node is a Link from react-router we need to add our
     * custom 'active' class from css-modules to enable active states.
     */
    maybeAddActiveClass = (node) => {
        const {
            active: activeClassName
        } = styles;

        if (node.type && node.type.displayName === 'Link') {
            return cloneElement(node, { activeClassName });
        }

        return node;
    }

    render() {
        return (
            <li className={ styles.tab } id={ this.props.id }>
                { Children.map(this.props.children, this.maybeAddActiveClass) }
            </li>
        );
    }
}
