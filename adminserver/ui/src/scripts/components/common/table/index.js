'use strict';

import React, {
    Component,
    PropTypes
} from 'react';

export default class Table extends Component {

    static propTypes = {
        headers: PropTypes.array,
        children: PropTypes.node
    }


    render() {
        const {
            headers
        } = this.props;
        return (
            <table cellSpacing='0'>
                <thead>
                    <tr>{
                        headers.map((header, i) => {
                            return (
                                <td key={ i }>{ header }</td>
                            );
                        })
                    }</tr>
                </thead>
                <tbody>{
                    this.props.children
                }</tbody>
            </table>
        );
    }
}
