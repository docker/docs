'use strict';

import React from 'react';
// components
import InputLabel from 'components/common/inputLabel';
import RadioGroup from 'components/common/radioGroup';

const RepoVisibility = ({ formField }) => (
  <div>
    <InputLabel>Visibility</InputLabel>
    <RadioGroup
      slim
      initialChoice={ formField.initialValue }
      choices={ [
        { label: 'Public', value: 'public' },
        { label: 'Private', value: 'private' }
      ] }
      formField={ formField } />
  </div>
);

export default RepoVisibility;
