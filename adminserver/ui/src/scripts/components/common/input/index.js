'use strict';

import React, {
    Component,
    PropTypes
} from 'react';
const {
    string,
    object,
    func,
    bool,
    number
} = PropTypes;
import FontAwesome from 'components/common/fontAwesome';
import Tooltip from 'rc-tooltip';
import styles from './input.css';
import css from 'react-css-modules';
import cn from 'classnames';
import ui from 'redux-ui';

@ui({
    state: {
        visiblePassword: false
    }
})
@css(styles, { allowMultiple: true })
export default class Input extends Component {

    static propTypes = {
        type: string,
        placeholder: string,
        formfield: object,
        ui: object,
        updateUI: func,
        disabled: bool,
        // For numeric inputs, allow a minimum value
        min: number,
        markup: bool,
        isTextarea: bool,
        // Allow a misc onchange if it exists
        onChange: func
    }

    static defaultProps = {
        markup: true
    }

    revealPassword = () => {
        this.props.updateUI({
            visiblePassword: !this.props.ui.visiblePassword
        });
    }

    render() {

        const {
            type,
            placeholder,
            formfield,
            ui: {
                visiblePassword
            },
            disabled,
            min,
            markup,
            isTextarea,
            onChange
        } = this.props;

        const isTouchedError = formfield && formfield.touched && formfield.error;
        const isDirty = formfield && formfield.dirty && !formfield.error;

        return (
            <div styleName={ cn('container', {'touchContainer': isDirty}, {'errorContainer': isTouchedError} ) }>
                { /*  Textarea = same behavior as input but multiple lines */ }
                { isTextarea ?
                    <textarea
                        placeholder={ placeholder }
                        onChange={ onChange }
                        {...formfield }
                        disabled={ disabled } />

                    :
                    <input type={ type === 'password' && visiblePassword ? 'text' : type }
                        placeholder={ placeholder }
                        onChange={ onChange }
                        min={ min }
                        {...formfield }
                        disabled={ disabled } />
                }

                { /*  Displays input states */ }
                <div styleName='icontainer'>
                    { /*  Edited state */ }
                    { markup && isDirty && <div styleName='touched'><FontAwesome icon='fa-check' /></div> }
                    { /*  Error state */ }
                    { markup && isTouchedError &&
                        <div styleName='error'>
                            <Tooltip
                                overlay={ formfield.error }
                                placement='bottom'
                                align={ { overflow: { adjustX: 0, adjustY: 0 } } }
                                trigger={ ['hover'] }>
                                <FontAwesome icon={ 'fa-exclamation' } />
                            </Tooltip>
                        </div>
                    }
                    { /*  Show/hide password toggle */ }
                    { type === 'password' ?
                        <div
                            styleName={ visiblePassword ? 'password hide' : 'password reveal' }
                            onClick={ ::this.revealPassword }>
                            <Tooltip
                                overlay={ <div>{ visiblePassword ? 'Hide password' : 'Show password' }</div> }
                                placement='bottom'
                                align={ { overflow: { adjustX: 0, adjustY: 0 } } }
                                trigger={ ['hover'] }>
                                <FontAwesome icon={ visiblePassword ? 'fa-eye' : 'fa-eye-slash' } />
                            </Tooltip>
                        </div>
                        : undefined
                    }
                </div>
            </div>
        );
    }
}
