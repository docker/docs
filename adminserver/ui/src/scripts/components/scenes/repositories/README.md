# Repositories

This lists all repositories a user has access to.

## Filtering

Filtering refers to the dropdown box in the top left of the scene. An
enterprise user may have hundreds of thousands of namespaces - too many to list
within a dropdown.

Our strategy is:

1. List all namespaces that the current user is a member of
   (`listUserOrganizations`)
2. Populate the typeahead with these namespaces
3. onKeyDown use the autocomplete API to list 10 namespaces from the user's
   input and populate the typeahead
4. When selecting a namespace make a new API request to select namespaces' repos
