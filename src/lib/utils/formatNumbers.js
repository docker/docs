import numeral from 'numeral';

// Formatting helpers for popularity counts
export const formatNumber = (n) => {
  if (n < 1000) {
    return numeral(n).format('0a');
  }
  return numeral(n).format('0.0a').toUpperCase();
};

 // Buckets to format the number for popularity:
 // 10M +
 // 5M +
 // 1M +
 // 500k +
 // 100k +
 // 50k +
 // 10k +
export const formatBucketedNumber = (n) => {
  if (typeof n !== 'number' || n < 0) return undefined;
  if (n < 10000) {
    return formatNumber(n);
  } else if (n < 50000) {
    return '10K+';
  } else if (n < 100000) {
    return '50K+';
  } else if (n < 500000) {
    return '100K+';
  } else if (n < 1000000) {
    return '500K+';
  } else if (n < 5000000) {
    return '1M+';
  } else if (n < 10000000) {
    return '5M+';
  }
  return '10M+';
};

/* Tag Size */

export const bytesToSize = (bytes, precision) => {
  const kilobyte = 1000;
  const megabyte = kilobyte * 1000;
  const gigabyte = megabyte * 1000;
  const terabyte = gigabyte * 1000;

  if ((bytes >= 0) && (bytes < kilobyte)) {
    return `${bytes} B`;
  } else if ((bytes >= kilobyte) && (bytes < megabyte)) {
    return `${(bytes / kilobyte).toFixed(precision)} KB`;
  } else if ((bytes >= megabyte) && (bytes < gigabyte)) {
    return `${(bytes / megabyte).toFixed(precision)} MB`;
  } else if ((bytes >= gigabyte) && (bytes < terabyte)) {
    return `${(bytes / gigabyte).toFixed(precision)} GB`;
  } else if (bytes >= terabyte) {
    return `${(bytes / terabyte).toFixed(precision)} TB`;
  }
  return 'Unknown size';
};
