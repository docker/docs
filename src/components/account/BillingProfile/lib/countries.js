const countries = require('country-list')();

export const countryOptions = countries.getCodes().map(code => {
  const label = countries.getName(code);
  return { label, value: code };
});

export const countryGetCodeFromName = name => countries.getCode(name);
