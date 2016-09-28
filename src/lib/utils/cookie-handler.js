
export const readCookie = (name, list) => {
  let cookies;

  if (list) {
    cookies = list;
  } else {
    cookies = typeof document !== 'undefined' ? document.cookie : null;
  }

  if (!cookies) {
    return '';
  }

  const c = cookies.match(`(^|;)\\s*${name}\\s*=\\s*([^;]+)`);
  return c ? c.pop() : '';
};
