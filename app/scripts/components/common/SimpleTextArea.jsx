'use strict';

import React, {
  Component,
  PropTypes
} from 'react';
import { findDOMNode } from 'react-dom';
import styles from './SimpleTextArea.css';
const { any, func, string, number } = PropTypes;
const debug = require('debug')('SimpleTextArea');

/*
 * A simpler version of DUXInput without errors / alert boxes
 * that is a text area
 */
export default class SimpleTextArea extends Component {
  static propTypes = {
    cols: number,
    name: string,
    onChange: func.isRequired,
    placeholder: string,
    rows: number,
    value: any.isRequired
  }

  static defaultProps = {
    onChange() {
      debug('No onChange function set for Input');
    },
    cols: 2,
    rows: 100,
    value: '',
    placeholder: '',
    name: ''
  }

  focusInput = () => {
    findDOMNode(this.refs.textarea).focus();
  }

  state = {
    value: ''
  }

  _onChange = (e) => {
    // Set local state
    const { value } = e.target;
    this.setState({ value });

    // Pass control to onChange controller. Usually a form.
    this.props.onChange(e);
  }

  componentWillReceiveProps(props) {
    /**
     * State should be tracked globally in a parent component. As
     * such, we need to check the passed in prop against the current
     * (local) state value to avoid overwriting the input value.
     *
     * This prevents a bug which resulted in the cursor being "jumped"
     * to the end of an input after every character.
     */
    if (props.value !== this.state.value) {
      this.setState({
        value: props.value
      });
    }
  }

  render() {
    const { cols,
            rows,
            name,
            placeholder,
            value } = this.props;
    return (
      <div>
        <textarea ref='textarea'
                  className={styles.textarea}
                  name={name}
                  rows={rows}
                  cols={cols}
                  value={this.state.value || value}
                  onChange={this._onChange}
                  placeholder={placeholder} />
      </div>
    );
  }
}
