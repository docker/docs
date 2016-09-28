// Categories is an array of category objects like { name, label }
export default (categories) => {
  if (!categories) {
    return '';
  }
  // Map snakecase strings to category names to show
  return categories.map(({ label }) => label).join(', ');
};
