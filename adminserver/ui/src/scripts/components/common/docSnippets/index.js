'use strict';

// Snippets of our documentation that we use in tooltips

export let visibilityTooltip = `
Repositories are either public or private.

Public:
- visible to all accounts in the system
- can only be written to by accounts granted explicit write access

Private:
- cannot be discovered by any account unless having explicit read access to it
`;

export let permissionTooltip = `
Type of access granted to members in teams with the corresponding permissions:

|                     | **\`read-only\`** | **\`read-write\`** | **\`admin\`** |
| ------------------- | :-------------: | :--------------: | :---------: |
| view and browse     |        ✓        |         ✓        |      ✓      |
| pull                |        ✓        |         ✓        |      ✓      |
| push                |                 |         ✓        |      ✓      |
| modify/delete tags  |                 |         ✓        |      ✓      |
| edit description    |                 |                  |      ✓      |
| make public/private |                 |                  |      ✓      |
| manage user access  |                 |                  |      ✓      |
`;

export let permissionsByAccessLevel = {
  'read-only': `
- view and browse
- pull
  `,
  'read-write': `
- view and browse
- pull
- push
- modify/delete tags
  `,
  'admin': `
- view and browse
- pull
- push
- modify/delete tags
- edit description
- make public/private
- manage user access
  `
};
