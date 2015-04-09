package utils

type IScope interface {
	ID() string
	Compare(IScope) bool
}

type IAuthorizer interface {
	// Authorize is expected to set the Authorization on the Context. If
	// Authorization fails, an error should be returned, but additionally,
	// the Authorization on the Context should be set to an instance of
	// NoAuthorization
	Authorize(IContext, ...IScope) error
}

type IAuthorization interface {
	HasScope(IScope) bool
}

// ### START INSECURE AUTHORIZATION TOOLS ###
// THESE ARE FOR DEV PURPOSES ONLY, DO NOT USE IN
// PRODUCTION

// DON'T USE THIS FOR ANYTHING, IT'S VERY INSECURE
type InsecureAuthorizer struct{}

// LIKE I SAID, VERY INSECURE
func (auth *InsecureAuthorizer) Authorize(ctx IContext, scopes ...IScope) error {
	ctx.SetAuthorization(&InsecureAuthorization{})
	return nil
}

// ALSO DON'T USE THIS, IT'S ALSO VERY INSECURE
type InsecureAuthorization struct {
}

// THIS IS JUST INCREDIBLY INSECURE
func (authzn *InsecureAuthorization) HasScope(scope IScope) bool {
	return true
}

// ### END INSECURE AUTHORIZATION TOOLS ###

type NoAuthorization struct{}

func (authzn *NoAuthorization) HasScope(scope IScope) bool {
	return false
}

type SimpleScope string

func (ss SimpleScope) ID() string {
	return string(ss)
}

func (ss SimpleScope) Compare(toCompare IScope) bool {
	return ss.ID() == toCompare.ID()
}

const (
	SSNoAuth SimpleScope = SimpleScope("NoAuth")
	SSCreate             = SimpleScope("Create")
	SSRead               = SimpleScope("Read")
	SSUpdate             = SimpleScope("Update")
	SSDelete             = SimpleScope("Delete")
)
