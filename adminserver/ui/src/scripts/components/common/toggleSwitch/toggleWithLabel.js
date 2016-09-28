'use strict';

import React, { Component, PropTypes } from 'react';
const { object, shape, func, bool } = PropTypes;
import ToggleSwitch from './index.js';
import FontAwesome from 'components/common/fontAwesome';
import InputLabel from 'components/common/inputLabel';
import styles from './toggle.css';
import css from 'react-css-modules';
import ui from 'redux-ui';

@ui({
  state: {
    on: (props) => props.initial,
    touched: false
  }
})
@css(styles)
export default class ToggleWithLabel extends Component {

  static propTypes = {
    labelOptions: object,
    formField: shape({
      onChange: func.isRequired
    }),
    ui: object,
    updateUI: func,
    initial: bool
  }

  handleClick = () => {
    const {
      updateUI,
      formField,
      initial
    } = this.props;

    let on = !this.props.ui.on;
    if (initial && on && !this.props.ui.touched) {
      on = !initial;
    }

    updateUI({
      on
    });
    if (!this.props.ui.touched) {
      updateUI({
        touched: true
      });
    }
    if (formField) {
      formField.onChange(on);
    }
  }

  render() {

    const {
      formField,
      labelOptions,
      ui: {
        on,
        touched
      },
      initial
    } = this.props;

    let checked = on;
    if (initial && !on && !touched) {
      checked = initial;
    }

    const isDirty = formField && formField.dirty && !formField.error;

    return (
      <div styleName='ToggleWithLabel'>
        <InputLabel {...labelOptions}>{ labelOptions.labelText }</InputLabel>
        <ToggleSwitch onClick={ ::this.handleClick } on={ checked }/>
        <input type='checkbox' formfield={ formField }/>
        { isDirty &&
        <div styleName='icontainer'>
          <div styleName='touched'>
            <FontAwesome icon='fa-check'/>
          </div>
        </div>
        }
      </div>
    );
  }
}
