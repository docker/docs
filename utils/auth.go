package utils

type Scope interface {
	ID() interface{}
}

type Authorizer interface {
	Authorize(IContext, ...Scope) error
}
