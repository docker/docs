import inRange from 'lodash/inRange';
import {
  STEP0,
  STEP1,
  STEP2,
  STEP3,
  STEP4,
} from 'lib/constants/publisherSteps';

const stepList = [STEP0, STEP1, STEP2, STEP3, STEP4];
export const getProductStep = (productStatus) => {
  const productStep = stepList.indexOf(productStatus) > -1 ?
    stepList.indexOf(productStatus) : 0;
  return productStep;
};

export const getCurrentStep = (step, productStep) => {
  // NOTE: step should be an integer between 0-4
  let currentStep = productStep;
  if (Number.isInteger(step) && inRange(step, 0, productStep + 1)) {
    currentStep = step;
  }
  return currentStep;
};
