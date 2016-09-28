package openid

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/jose"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

const defaultSigningKeyExpiresIn = time.Hour * 12

// Service handles API requests relating to signing keys.
type Service struct {
	server.Service

	schemaMgr             schema.Manager
	authorizer            authz.Authorizer
	signingKey            *jose.PrivateKey
	signingKeyCacheMaxAge time.Duration

	// Holds a list of certificates which are implicitly used to verify
	// authentication JWTs from privileged services.
	privilegedRootCerts []*x509.Certificate
}

// NewService returns a new signing keys Service.
func NewService(baseContext context.Context, schemaMgr schema.Manager, privilegedCAFile string, privateKey crypto.PrivateKey, signingKeyCacheMaxAge time.Duration, rootPath string) (*Service, error) {
	// Load root certificates for privileged services.
	privilegedRootCerts, err := loadPrivilegedRootCerts(privilegedCAFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load privileged root certs: %s", err)
	}

	signingKey, err := jose.NewPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to create JWT signing key: %s", err)
	}

	if err := schemaMgr.SaveSigningKey(schema.NewJWK(&signingKey.PublicKey, defaultSigningKeyExpiresIn)); err != nil {
		return nil, fmt.Errorf("unable to save JWT signing key: %s", err)
	}

	// Once each hour, extend the key expiration by another 12 hours.
	go func() {
		for {
			time.Sleep(time.Hour)

			err := schemaMgr.ExtendSigningKeyExpiration(signingKey.ID, defaultSigningKeyExpiresIn)
			for err != nil {
				context.GetLogger(baseContext).Errorf("unable to extend signing key expiration (trying again in 30s): %s", err)

				time.Sleep(time.Second * 30)
				err = schemaMgr.ExtendSigningKeyExpiration(signingKey.ID, defaultSigningKeyExpiresIn)
			}

			context.GetLogger(baseContext).Debugf("successfully extended expiration time for key %s", signingKey.ID)
		}
	}()

	service := &Service{
		Service: server.Service{
			WebService:  new(restful.WebService),
			BaseContext: baseContext,
		},
		schemaMgr:             schemaMgr,
		authorizer:            authz.NewAuthorizer(schemaMgr),
		signingKey:            signingKey,
		signingKeyCacheMaxAge: signingKeyCacheMaxAge,
		privilegedRootCerts:   privilegedRootCerts,
	}

	service.connectRoutes(rootPath)

	return service, nil
}

func loadPrivilegedRootCerts(privilegedCAFile string) ([]*x509.Certificate, error) {
	privilegedRootCAsPEM, err := ioutil.ReadFile(privilegedCAFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read privileged root CA file: %s", err)
	}

	var privilegedRootCerts []*x509.Certificate

	certBlock, privilegedRootCAsPEM := pem.Decode(privilegedRootCAsPEM)
	for certBlock != nil {
		privilegedRootCert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse privileged root cert: %s", err)
		}

		privilegedRootCerts = append(privilegedRootCerts, privilegedRootCert)

		certBlock, privilegedRootCAsPEM = pem.Decode(privilegedRootCAsPEM)
	}

	if len(privilegedRootCerts) == 0 {
		return nil, fmt.Errorf("no valid root certs found in file: %s", privilegedCAFile)
	}

	return privilegedRootCerts, nil
}

// connectRoutes registers all API endpoints on this service with paths
// relative to the given rootPath.
func (s *Service) connectRoutes(rootPath string) {
	s.WebService.Path(rootPath).
		Produces(restful.MIME_JSON).
		Doc("Identity")

	routes := []server.Route{
		s.routeToken(),
		s.routeAuthorize(),
		s.routeLogin(),
		s.routeLogout(),
		s.routeListSigningKeys(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}
