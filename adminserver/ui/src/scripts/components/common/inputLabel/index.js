'use strict';

import React, { Component, PropTypes } from 'react';
const { node, string, bool, oneOfType } = PropTypes;
import styles from './inputLabel.css';
import QTip from 'components/common/qtip';
import classNames from 'classnames';

export default class InputLabel extends Component {

    static propTypes = {
        children: node.isRequired,
        field: string,
        isOptional: bool,
        hint: string,
        tip: oneOfType([string, node]),
        inline: bool
    }

    render() {

        const {
            children,
            isOptional,
            hint,
            tip,
            inline,
            field
        } = this.props;

        return (
            <div className={ inline ? classNames(styles.inputLabel, styles.inline) : styles.inputLabel }>
                <label htmlFor={ field }>{ children }</label>
                { tip ? <QTip tooltip={ tip } /> : undefined }
                { isOptional ? <span className={ styles.optional }>(optional)</span> : undefined }
                { /* \u2013 is an endash: "–" */ }
                { hint ? <span className='hint'>{ '\u2013 ' + hint }</span> : undefined }
            </div>
        );

    }

}
