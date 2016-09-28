// TODO Kristie 5/5/16 Get sorting by Created Date once API supports it
// Getting decision on all sorting options from design - these are what the
// API supports
export default [
  {
    value: '',
    label: 'Most Popular',
  },
  {
    value: 'sort=updated_at&order=desc',
    label: 'Recently Updated',
  },
  // {
  //   value: 'sort=last_updated&order=asc',
  //   label: 'Oldest',
  // },
  // {
  //   value: 'sort=popularity&order=asc',
  //   label: 'Least Popular',
  // },
  // {
  //   value: 'sort=source&order=asc',
  //   label: 'Source Ascending',
  // },
  // {
  //   value: 'sort=source',
  //   label: 'Source Descending',
  // },
  // {
  //   value: 'sort=platform&order=asc',
  //   label: 'Platform Ascending',
  // },
  // {
  //   value: 'sort=platform',
  //   label: 'Platform Descending',
  // },
];
