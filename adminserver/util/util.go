package util

import (
	"net/http"

	"github.com/docker/dhe-deploy/garant/authn"
	gorillacontext "github.com/gorilla/context"
)

// GetAuthenticatedUser gets the "user" value from the request using
// gorilla/context. If there is no set value, an anonymous user object will
// be returned.
func GetAuthenticatedUser(request *http.Request) *authn.User {
	val := gorillacontext.Get(request, "user")
	user, _ := val.(*authn.User)

	return user
}
