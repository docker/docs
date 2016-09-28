'use strict';

import React, {
    Component,
    PropTypes
} from 'react';
const {
    string,
    object,
    bool,
    func,
    node,
    oneOfType
} = PropTypes;
import FontAwesome from 'components/common/fontAwesome';
import InputLabel from 'components/common/inputLabel';
import styles from './checkbox.css';
import css from 'react-css-modules';

@css(styles, { allowMultiple: true })
export default class Checkbox extends Component {

    static propTypes = {
        tip: oneOfType([string, node]),
        inputLabel: string,
        formfield: object,
        isChecked: bool,
        onChange: func
    }

    handleClick = () => {
        const value = !this.refs.box.checked;
        this.refs.box.checked = value;

        if (this.props.onChange) {
            this.props.onChange.call(this, value);
        }
        if (this.props.formfield && this.props.formfield.onChange) {
          this.props.formfield.onChange.call(this, value);
        }
    }

    render() {
        const {
            tip,
            inputLabel,
            formfield,
            isChecked,
            onChange
        } = this.props;

        return (
            <div styleName='checkbox'>
                <input
                  ref='box'
                  type='checkbox'
                  checked={ isChecked }
                  onChange={ onChange }
                  {...formfield } />
                <div styleName='check-container' onClick={ ::this.handleClick }>
                    <FontAwesome icon='fa-check' />
                </div>
                { inputLabel ?
                  <span>
                        <InputLabel
                          field={ formfield }
                          tip={ tip }
                          inline={ true }>{ inputLabel }</InputLabel>
                    </span>
                  : undefined }
            </div>
        );
    }
}

