const SMALL = 'small';
const SMALL2X = 'small@2x';
const LARGE = 'large';
const LARGE2X = 'large@2x';

// Given a logo_url object, return the highest resolution image available
export default (logo_url) => {
  if (!logo_url) return '';
  if (logo_url[LARGE2X]) return logo_url[LARGE2X];
  if (logo_url[LARGE]) return logo_url[LARGE];
  if (logo_url[SMALL2X]) return logo_url[SMALL2X];
  if (logo_url[SMALL]) return logo_url[SMALL];
  return '';
};
