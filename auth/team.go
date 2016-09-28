package auth

// this is the payload returned from eNZi for a team
type Team struct {
	OrgId             string   `json:"orgId,omitempty"`
	Name              string   `json:"name,omitempty"`
	Id                string   `json:"id,omitempty"`
	FullName          string   `json:"fullName,omitempty"`
	Description       string   `json:"description,omitempty"`
	ManagedMembers    []string `json:"managed_members"`    // XXX Should these be omitempty, or should we have a different serialization for list vs. one?
	DiscoveredMembers []string `json:"discovered_members"` // XXX Should these be omitempty, or should we have a different serialization for list vs. one?
	LdapDN            string   `json:"ldapdn,omitempty"`
	LdapMemberAttr    string   `json:"ldap_member_attr,omitempty"`
}
