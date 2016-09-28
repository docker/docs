'use strict';

import {createValidator, required, maxLength, regex} from 'validation';

const createOrganizationValidator = createValidator({
  name: [required, maxLength(30), regex(/^[a-z0-9]+(?:[._-][a-z0-9]+)*$/, 'Lowercase alphanumeric characters and \'._-\' only')]
});
export default createOrganizationValidator;
