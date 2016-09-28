'use strict';

import { Map, Record } from 'immutable';
import { assert } from 'chai';
import { NamespaceRecord, UserRecord } from '../../records.js';
import { normalizedToRecords } from '../records.js';

const normalized = {
    "entities": {
        "org": {
            "fraud": {
                "id": 20,
                "name": "fraud",
                "type": "organization"
            }
        },
        "user": {
            "financegroup": {
                "id": 3,
                "isActive": true,
                "ldapLogin": "CN=FinanceGroup,OU=Groups,DC=ad,DC=dckr,DC=org",
                "name": "financegroup",
                "type": "user"
            }
        }
    },
    "result": [
        {
            "id": "financegroup",
            "schema": "user"
        },
        {
            "id": "fraud",
            "schema": "org"
        }
    ]
};

describe('normalizedToRecords()', () => {

  it('returns a map', () => {
    const result = normalizedToRecords(normalized, {
      'user': NamespaceRecord,
      'org': NamespaceRecord
    });

    assert.isTrue(Map.isMap(result));
    assert.isTrue(Map.isMap(result.get('entities')));
  });

  it('sets items as records', () => {
    const result = normalizedToRecords(normalized, {
      'user': NamespaceRecord,
      'org': NamespaceRecord
    });
    const org = result.getIn(['entities', 'org', 'fraud']);
    const user = result.getIn(['entities', 'user', 'financegroup']);

    assert.isTrue(org.constructor === NamespaceRecord);
    assert.isTrue(user.constructor === NamespaceRecord);
  });

  it('respects the schema map', () => {
    const result = normalizedToRecords(normalized, {
      'user': UserRecord,
      'org': NamespaceRecord
    });
    const org = result.getIn(['entities', 'org', 'fraud']);
    const user = result.getIn(['entities', 'user', 'financegroup']);

    assert.isTrue(org.constructor === NamespaceRecord);
    assert.isTrue(user.constructor === UserRecord);
  });

});
