'use strict';

import React, { Component, PropTypes } from 'react';
const { arrayOf, bool, func, object, shape, string } = PropTypes;
import styles from './radioGroup.css';
import Radio from './radio';
import ui from 'redux-ui';
import css from 'react-css-modules';

@ui({
  state: {
    selected: (props) => props.initialChoice
  }
})
@css(styles, {allowMultiple: true})
export default class RadioGroup extends Component {

  static propTypes = {
    choices: arrayOf(shape({
      label: string,
      value: string
    })).isRequired,
    onChange: func,
    ui: object,
    updateUI: func,
    initialChoice: string.isRequired,
    formField: object,
    style: object,
    vertical: bool,
    // Whether to use slim styles making the radiogroup the same size as an
    // input box
    slim: bool,
    id: string
  };

  static defaultProps = {
    vertical: false
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

    if (this.props.formField) {
      this.props.formField.onChange(selected);
    }
  };

  render() {
    const {
      choices,
      ui: {
        selected
        },
      formField,
      vertical,
      slim,
      id
      } = this.props;

    let styleName = 'choices';
    if (slim) {
      styleName += ' slim';
    }
    if (vertical) {
      styleName += ' vertical';
    }

    return (
      <div styleName={ styleName } id={ id }>
        { choices.map((choice, i) => {
          return (
            <span
              className='choice'
              style={ this.props.style }
              key={ i }
              styleName='choice'
              onClick={ ::this.markActive(choice.value) }>
               <Radio active={ selected === choice.value }/>
              { choice.label }
              { formField ? <input type='hidden' {...formField} /> : undefined }
           </span>
          );
        }, this) }
      </div>
    );
  }
}
