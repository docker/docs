export default (values) => {
  const errors = {};
  const {
    firstName,
    lastName,
    company,
  } = values;
  if (!firstName) {
    errors.firstName = 'Required';
  }
  if (!lastName) {
    errors.lastName = 'Required';
  }
  if (!company) {
    errors.company = 'Required';
  }
  return errors;
};
