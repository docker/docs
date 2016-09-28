export default (values) => {
  const errors = {};
  const REQUIRED = 'Required';
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
    errors.firstName = REQUIRED;
  }
  if (!lastName || !lastName.trim()) {
    errors.lastName = REQUIRED;
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
  if (!address || !address.trim()) {
    errors.address = REQUIRED;
  }
  if (!city || !city.trim()) {
    errors.city = REQUIRED;
  }
  if (!postCode || !postCode.trim()) {
    errors.postCode = REQUIRED;
  }

  if (!cardNumber) {
    errors.cardNumber = REQUIRED;
  }
  if (!cvv) {
    errors.cvv = REQUIRED;
  }
  if (!expMonth) {
    errors.expMonth = REQUIRED;
  }
  if (!expYear) {
    errors.expYear = REQUIRED;
  }
  const date = new Date();
  if (
    values.expMonth < (date.getMonth() + 1) &&
    values.expYear <= date.getFullYear()
  ) {
    errors.expMonth = REQUIRED;
    errors.expYear = REQUIRED;
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
