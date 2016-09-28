'use strict';

import React, { Component, PropTypes } from 'react';
const { string, node, oneOfType } = PropTypes;
import Tooltip from 'rc-tooltip';
import FontAwesome from 'components/common/fontAwesome';

// base styles for tooltips
import 'rc-tooltip/assets/bootstrap.css';

export default class QTip extends Component {
    static propTypes = {
        tooltip: oneOfType([string, node]),
        placement: string
    }

    static defaultProps = {
        placement: 'right'
    }

    render () {

        const {
            tooltip,
            placement
        } = this.props;

        return (
                <span>
                    <Tooltip
                        overlay={ <div styleName='container'>{ tooltip }</div> }
                        placement={ placement }
                        align={ { overflow: { adjustY: 0 } } }
                        trigger={ ['hover'] }>
                        <FontAwesome icon='fa-question-circle' />
                    </Tooltip>
                </span>
        );
    }
}
