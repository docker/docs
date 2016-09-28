'use strict';

export default function bytesToSize(bytes, precision)
{
  const kilobyte = 1000;
  const megabyte = kilobyte * 1000;
  const gigabyte = megabyte * 1000;
  const terabyte = gigabyte * 1000;

  if ((bytes >= 0) && (bytes < kilobyte)) {
    return bytes + ' B';
  } else if ((bytes >= kilobyte) && (bytes < megabyte)) {
    return (bytes / kilobyte).toFixed(precision) + ' KB';
  } else if ((bytes >= megabyte) && (bytes < gigabyte)) {
    return (bytes / megabyte).toFixed(precision) + ' MB';
  } else if ((bytes >= gigabyte) && (bytes < terabyte)) {
    return (bytes / gigabyte).toFixed(precision) + ' GB';
  } else if (bytes >= terabyte) {
    return (bytes / terabyte).toFixed(precision) + ' TB';
  } else {
    return 'Unknown size';
  }
}
