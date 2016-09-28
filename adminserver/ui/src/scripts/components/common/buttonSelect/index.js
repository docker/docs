'use strict';

import React, { Component, PropTypes } from 'react';
import FontAwesome from 'components/common/fontAwesome';
const { node, func, object, string } = PropTypes;
import css from 'react-css-modules';
import styles from './buttonSelect.css';
import ui from 'redux-ui';
import classNames from 'classnames';

@ui({
  state: {
    selected: (props) => props.initialChoice
  }
})
@css(styles)
export default class ButtonSelect extends Component {

  static propTypes = {
    onChange: func,
    children: node.isRequired,
    ui: object,
    updateUI: func,
    initialChoice: string.isRequired
  };

  componentWillReceiveProps(next) {
    if (next.initialChoice !== this.props.initialChoice) {
      this.props.updateUI('selected', next.initialChoice);
    }
  }

  markActive = (selected) => () => {
    this.props.updateUI({selected: selected});

    if (this.props.onChange) {
      this.props.onChange(selected);
    }
  };

  render() {
    const {
      ui: {
        selected
      }
    } = this.props;

    return (
      <div styleName='choices'>
        { React.Children.map(this.props.children, (child, i) => {
          return (
              <button
                key={ i }
                type='button'
                className={ selected === child.props.value ? classNames(styles.button, styles.active) : styles.button }
                onClick={ ::this.markActive(child.props.value) }>
                    { (() => {
                        // if the child has children, just use that
                        if (React.Children.count(child.props.children)) {
                            return child;
                        } else {
                            return (
                                <div styleName='buttonContent'>
                                    <div styleName='icon'><FontAwesome icon={ child.props.icon } /></div>
                                    <div styleName='textContainer'>
                                        <div styleName='primaryText'>{ child.props.primaryText }</div>
                                        <div styleName='secondaryText'>{ child.props.secondaryText }</div>
                                    </div>
                                </div>
                            );
                        }
                    })() }
              </button>
          );
        }, this) }
      </div>
    );
  }
}
