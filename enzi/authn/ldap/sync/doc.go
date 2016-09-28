/*
Package sync handles syncing users and groups from an LDAP v3 directory service
to a Docker Trusted Registry user account and team database. The syncing is
done in two phases:

Phase 1: User Sync

First, a set of existing users is created from the data currently in the
database. As we do not expect there to me more than on the order of 100,000
users and an average of 10-20 bytes per username, a set of all usernames should
take up no more than a few megabytes of memory. If we maintain a mapping of
these usernames to LDAP DNs, it still should not take up more than ~100
megabytes dependending on the average DN length. All existing users are loaded
into memory, fetch one page at a time from the database with a page size no
greater than 500 records. Each user's isActive status is also stored in this
mapping.

Using the configured LDAP server and user search filter, a search is performed
for user objects on the directory service. Large datasets will require paging
of search results. Each search result will have a DN (distinguished name) and
a login attribute value. These results are stored in an in-memory set mapping
the login attribute value (username) to the entry's LDAP DN and isActive=true.
At the end of this step, there are two in-memory sets mapping usernames to
user data.

The intersection of these two sets represents those users which already exist
in the system but may need to be updated to be marked as active and ensure the
DN is up to date for team syncing. Those existing users which were not found in
the LDAP search are all marked inactive. Those found via the search that did
not previously exist are all added to the database as active users.

Phase 2: Team Sync

TODO (BrianBland): Document how team syncing works.
*/
package sync
