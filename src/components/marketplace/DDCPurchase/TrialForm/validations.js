export default (values) => {
  const errors = {};
  const REQUIRED = 'Required';
  const {
    accepted,
    company,
    country,
    firstName,
    job,
    lastName,
    phone,
    state,
  } = values;
  if (!firstName || !firstName.trim()) {
    errors.firstName = REQUIRED;
  }
  if (!lastName || !lastName.trim()) {
    errors.lastName = REQUIRED;
  }
  if (!company || !company.trim()) {
    errors.company = REQUIRED;
  }
  if (!job || !job.trim()) {
    errors.job = REQUIRED;
  }
  if (!phone || !phone.trim()) {
    errors.phone = REQUIRED;
  }
  if (!state || !state.trim()) {
    errors.state = REQUIRED;
  }
  if (!country) {
    errors.country = REQUIRED;
  }
  if (accepted !== 'checked') {
    errors.accepted = 'Please accept the terms of use before continuing';
  }
  return errors;
};
