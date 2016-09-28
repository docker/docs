const countries = require('country-list')();

export const countryOptions = countries.getCodes()
.map(code => {
  const label = countries.getName(code);
  return { label, value: code };
})
.filter(c => c.value !== 'US');
countryOptions.unshift({ label: countries.getName('US'), value: 'US' });

export const countryGetCodeFromName = name => countries.getCode(name);
