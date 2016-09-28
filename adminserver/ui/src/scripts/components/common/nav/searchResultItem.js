'use strict';

import React from 'react';
import styles from './search.css';
import FA from 'components/common/fontAwesome';

const SearchResultItem = ({ data: { data, type }, inputValue }) => {
  let name = (type === 'repo') ? `${data.namespace}/${data.name}` : data.name;
  const parts = name.split(inputValue);

  if (parts.length > 1) {
    const prior = parts[0];
    parts.splice(0, 1);
    const rest = parts.join(inputValue);
    name = (<span>{ prior }<b>{ inputValue }</b>{ rest }</span>);
  } else {
    name = <span>{ name }</span>;
  }

  let icon;
  switch (type) {
    case 'user':
      icon = 'fa-user';
      break;
    case 'org':
      icon = 'fa-users';
      break;
    case 'repo':
      icon = 'fa-book';
  }

  return (
    <div className={ styles.result }>
      <FA icon={ icon } className={ styles.resultIcon } />
      <p>{ name }</p>
    </div>
  );
};

export default SearchResultItem;
