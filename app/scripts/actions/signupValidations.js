'use strict';

var usernameValidation = function(actionContext, payload) {
  var validationState;
  var username = payload.unsafeUsername;
  if (username.length > 3 && username.length <= 30 && username.match(/[a-z0-9\.\-_]+$/)) {
    validationState = 'SUCCESS';
  } else if (username.length === 0) {
    validationState = 'NOTHING';
  } else {
    validationState = 'ERROR';
  }
  actionContext.dispatch('VALIDATED_SIGNUP_USERNAME', {
      username: username,
      valState: validationState
  });
};

var passValidation = function(actionContext, payload) {
  var validationState;
  var password = payload.unsafePass;
  if (password.length >= 5 && password.match(/[a-zA-Z]/) && password.match(/[0-9]/) &&
      password.match(/[\.-_!@#$%^&\*\(\)\[\]\{\}~:]/)) {
    validationState = 'SUPERSUCCESS';
  } else if (password.length >= 5 && password.match(/[a-zA-Z]/) && password.match(/[0-9]/)) {
    validationState = 'SUCCESS';
  } else if (password.length >= 5) {
    validationState = 'WEAK';
  } else if (password.length === 0) {
    validationState = 'NOTHING';
  } else if (password.length < 5) {
    validationState = 'ERROR';
  }
  actionContext.dispatch('VALIDATED_SIGNUP_PASSWORD', {
      password: password,
      valState: validationState
  });
};
var emailValidation = function(actionContext, payload) {
  var validationState;
  var email = payload.unsafeEmail;
  if (email.length > 3 && email.match(/.@./) && !email.match(/.\ ./)) {
    validationState = 'SUCCESS';
  } else if (email.length === 0) {
    validationState = 'NOTHING';
  } else {
    validationState = 'ERROR';
  }
  actionContext.dispatch('VALIDATED_SIGNUP_EMAIL', {
      email: email,
      valState: validationState
  });
};

module.exports = {
    usernameValidation: usernameValidation,
    passValidation: passValidation,
    emailValidation: emailValidation
};
