export default (values) => {
  const errors = {};
  const {
    cardNumber,
    cvv,
    expMonth,
    expYear,
    firstName,
    lastName,
  } = values;
  if (!firstName) {
    errors.firstName = 'Required';
  }
  if (!lastName) {
    errors.lastName = 'Required';
  }
  if (!cardNumber) {
    errors.cardNumber = 'Required';
  }
  if (!cvv) {
    errors.cvv = 'Required';
  }
  if (!expMonth) {
    errors.expMonth = 'Required';
  }
  if (!expYear) {
    errors.expYear = 'Required';
  }
  const date = new Date();
  if (
    values.expMonth < (date.getMonth() + 1) &&
    values.expYear <= date.getFullYear()
  ) {
    errors.expMonth = 'Required';
    errors.expYear = 'Required';
  }
  /*
  TODO: feature - nathan 05/19/16
  Luhn card number check.
  cvv check.
  expiration check.
  billing address validations check
  email check?
  account check? (already has a plan associated)
  */
  return errors;
};
