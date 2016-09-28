export default (values) => {
  const errors = {};
  const {
    accepted,
    address,
    cardNumber,
    city,
    country,
    cvv,
    expMonth,
    expYear,
    firstName,
    lastName,
    postCode,
    state,
  } = values;
  if (!firstName || !firstName.trim()) {
    errors.firstName = 'Required';
  }
  if (!lastName || !lastName.trim()) {
    errors.lastName = 'Required';
  }
  if (!state || !state.trim()) {
    errors.state = 'Required';
  }
  if (!country) {
    errors.country = 'Required';
  }
  if (accepted !== 'checked') {
    errors.accepted = 'Please accept the terms of use before continuing';
  }
  if (!address || !address.trim()) {
    errors.address = 'Required';
  }
  if (!city || !city.trim()) {
    errors.city = 'Required';
  }
  if (!postCode || !postCode.trim()) {
    errors.postCode = 'Required';
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
