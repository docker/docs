import { REDIRECT } from 'actions/root';

export default () => next => action => {
  if (action.type === REDIRECT) {
    window.location.href = action.payload;
    return Promise.resolve();
  }

  return next(action);
};
