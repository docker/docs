package utils

import (
	"net/http"
)

type IContext interface {
	// TODO: define a set of standard getters. Using getters
	//       will allow us to easily and transparently cache
	//       fields or load them on demand. Using this interface
	//       will allow people to define their own context struct
	//       that may handle things like caching and lazy loading
	//       differently.
}

type Context struct {
}

func generateContext(r *http.Request) Context {
	return Context{}
}
