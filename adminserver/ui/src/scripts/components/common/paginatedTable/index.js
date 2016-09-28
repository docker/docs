'use strict';

import React, { Component, PropTypes } from 'react';
import Table from 'components/common/table';
import Pager from 'components/common/pager';
import { getPage } from 'utils';

/*
PaginatedTable

Displays a paginated array in a table.

example use:

   <PaginatedTable

    // how many records to display per page
    perPage={ 10 }

    // table headers is an array of strings
    headers={ ['Username', 'Full Name', 'Teams', ''] }

    // rows should be an array of <tr>
    rows={ (() => {
        return dataArray.map((thing, i) => {
            return (
                <tr key={ i }>
                    <td>{ thing.foo }</td>
                    <td>{ thing.bar }</td>
                    <td>{ thing.wat }</td>
                    <td>{ thing.otherWat }</td>
                </tr>
            );
        });
    })() } />
*/

export default class PaginatedTable extends Component {

    static propTypes = {
        rows: PropTypes.array.isRequired,
        perPage: PropTypes.number,
        headers: PropTypes.array.isRequired
    }

    static defaultProps = {
        perPage: 5
    }

    render () {

        const {
            rows,
            perPage,
            headers
        } = this.props;

        const {
            search
        } = location;

        const page = getPage(search);

        const displayRows = rows.slice(page * perPage, (page + 1) * perPage);

        return (
            <span>
                <Table
                    headers={ headers }>
                    { (() => {
                        return displayRows.map((row) => {
                            return (row);
                        });
                    })() }
                </Table>
                <Pager
                    pageable={ rows }
                    perPage={ perPage }
                />
            </span>
        );
    }
}
