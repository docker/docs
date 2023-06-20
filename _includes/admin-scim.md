This section is for administrators who want to enable System for Cross-domain Identity Management (SCIM) 2.0 for their business. It is available for Docker Business customers.

SCIM provides automated user provisioning and de-provisioning for your Docker organization or company through your identity provider (IdP).  Once you enable SCIM in Docker Hub and your IdP, any user assigned to the Docker application in the IdP is automatically provisioned in Docker Hub and added to the organization or company.

Similarly, if a user gets unassigned from the Docker application in the IdP, the user is removed from the organization or company in Docker Hub. SCIM also synchronizes changes made to a user's attributes in the IdP, for instance the user’s first name and last name.

The following provisioning features are supported:
 - Creating new users
 - Push user profile updates
 - Remove users
 - Deactivate users
 - Re-activate users
 - Group mapping

The table below lists the supported attributes. Note that your attribute mappings must match for SSO to prevent duplicating your members.

| Attribute    | Description
|:---------------------------------------------------------------|:-------------------------------------------------------------------------------------------|
| username             | Unique identifier of the user (email)                                   |
| givenName                            | User’s first name |
| familyName |User’s surname                                              |