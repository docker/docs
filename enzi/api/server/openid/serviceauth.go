package openid

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/errors/oautherrors"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/jose"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

func (s *Service) authenticateService(ctx context.Context, r *restful.Request) (service *schema.Service, err error) {
	requiredClientAssertionType := "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	if r.Request.PostFormValue("client_assertion_type") != requiredClientAssertionType {
		return nil, oautherrors.InvalidRequest(fmt.Errorf("form parameter %q MUST be %q", "client_assertion_type", requiredClientAssertionType))
	}

	rawAuthToken := r.Request.PostFormValue("client_assertion")
	if rawAuthToken == "" {
		return nil, oautherrors.InvalidRequest(fmt.Errorf("form parameter 'client_assertion' MUST be a JWT"))
	}

	authToken, err := parseAuthenticationJWT(rawAuthToken)
	if err != nil {
		return nil, oautherrors.InvalidRequest(fmt.Errorf("unable to parse authentication JWT: %s", err))
	}

	service, err = s.validateAuthenticationJWT(ctx, authToken)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return nil, apiErr
		}

		return nil, oautherrors.InvalidClient(fmt.Errorf("unable to validate authentication JWT: %s", err))
	}

	return service, nil
}

type authenticationJWT struct {
	Header authenticationTokenHeader
	Claims authenticationTokenClaims

	SigningInput string
	Signature    string
}

type authenticationTokenHeader struct {
	Type       string   `json:"typ"`
	SigningAlg string   `json:"alg"`
	KeyID      string   `json:"kid"`
	CertChain  []string `json:"x5c"`
}

type authenticationTokenClaims struct {
	Issuer     string   `json:"iss"`
	Subject    string   `json:"sub"`
	Audience   []string `json:"aud"`
	Expiration int64    `json:"exp"`
}

func parseAuthenticationJWT(rawToken string) (*authenticationJWT, error) {
	parts := strings.SplitN(rawToken, ".", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("JWT must have 3 parts separated by '.'")
	}

	rawHeader, err := jose.Base64URLDecode(parts[0])
	if err != nil {
		return nil, fmt.Errorf("unable to base64url-decode JWT header: %s", err)
	}

	var header authenticationTokenHeader
	if err := json.Unmarshal(rawHeader, &header); err != nil {
		return nil, fmt.Errorf("unable to JSON-decode JWT header: %s", err)
	}

	rawClaims, err := jose.Base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("unable to base64url-decode JWT claims: %s", err)
	}

	var claims authenticationTokenClaims
	if err := json.Unmarshal(rawClaims, &claims); err != nil {
		return nil, fmt.Errorf("unable to JSON-decode JWT claims: %s", err)
	}

	return &authenticationJWT{
		Header:       header,
		Claims:       claims,
		SigningInput: strings.Join(parts[:2], "."),
		Signature:    parts[2],
	}, nil
}

func (s *Service) checkAcceptedAudiences(ctx context.Context, acceptedAudiences, jwtAudiences []string) error {
	acceptedAudienceSet := make(map[string]struct{}, len(acceptedAudiences))
	for _, acceptedAudience := range acceptedAudiences {
		acceptedAudienceSet[acceptedAudience] = struct{}{}
	}

	for _, audience := range jwtAudiences {
		if _, ok := acceptedAudienceSet[audience]; ok {
			return nil
		}
	}

	return fmt.Errorf("JWT audience values %q do not match any of those accepted: %q", jwtAudiences, acceptedAudiences)
}

func (s *Service) validateAuthenticationJWT(ctx context.Context, token *authenticationJWT) (service *schema.Service, err error) {
	// Check that the required header values have been specified.
	if token.Header.Type != "JWT" {
		return nil, fmt.Errorf("JWT header value 'typ' must be 'JWT'")
	}

	if token.Header.SigningAlg == "" {
		return nil, fmt.Errorf("JWT header value 'alg' must be specified")
	}

	// Either KeyID or CertChain must be specified.
	if token.Header.KeyID == "" && len(token.Header.CertChain) == 0 {
		return nil, fmt.Errorf("either JWT header value 'kid' or 'x5c' must be specified")
	}

	// Validate that the token is not yet expired and does not expire too
	// far in the future (we want to enforce that services use short lived
	// tokens).
	expiration := time.Unix(token.Claims.Expiration, 0)
	now := time.Now()

	// Check if expiration is in the past.
	if now.After(expiration) {
		return nil, fmt.Errorf("JWT expired at %d - current time is %d", expiration.Unix(), now.Unix())
	}
	// Check if expiration is too far in the future.
	if expiration.Sub(now) > time.Minute*5 {
		return nil, fmt.Errorf("JWT expires at %d - MUST be less than 5 minutes from now: %s", expiration.Unix(), now.Unix())
	}

	// Ensure that issuer and subject are identical.
	if token.Claims.Issuer != token.Claims.Subject {
		return nil, fmt.Errorf("JWT Issuer claim must be equal to Subject claim - service ID")
	}

	// Ensure that the service that is the issuer/subject of the token
	// exists.
	service, err = s.schemaMgr.GetServiceByID(token.Claims.Issuer)
	if err != nil {
		if err == schema.ErrNoSuchService {
			return nil, fmt.Errorf("no such service with ID: %s", token.Claims.Issuer)
		}

		return nil, errors.Internal(ctx, fmt.Errorf("unable to get service: %s", err))
	}

	// Ensure that the audience is in the list of provider identities for
	// this service.
	if err := s.checkAcceptedAudiences(ctx, service.ProviderIdentities, token.Claims.Audience); err != nil {
		return nil, err
	}

	// Check the signature of the token.
	if err := s.verifyTokenSignature(ctx, service, token); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Service) verifyTokenSignature(ctx context.Context, service *schema.Service, token *authenticationJWT) (err error) {
	var pubKey *jose.PublicKey

	switch {
	case token.Header.KeyID != "":
		// Lookup the signing key for the service using it's registered
		// JWKs URI or use cached key if we've seen it recently.
		pubKey, err = s.getServiceSigningKey(ctx, service, token.Header.KeyID)
	case len(token.Header.CertChain) > 0:
		// Verify the certificate chain in the token. Use the leaf cert
		// public key to verify the token's signature.
		pubKey, err = s.verifyServiceTokenCertChain(ctx, service, token.Header.CertChain)
	default:
		err = fmt.Errorf("JWT header must contain either a `kid` for an advertised signing key or a `x5c` certificate chain")
	}

	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return apiErr
		}

		return fmt.Errorf("unable to get service signing key: %s", err)
	}

	verifier, err := pubKey.Verifier(token.Header.SigningAlg)
	if err != nil {
		return fmt.Errorf("unable to verify token signature: %s", err)
	}

	if err := verifier.Verify(strings.NewReader(token.SigningInput), token.Signature); err != nil {
		return fmt.Errorf("invalid token signature: %s", err)
	}

	return nil
}

func (s *Service) verifyServiceTokenCertChain(ctx context.Context, service *schema.Service, certChain []string) (*jose.PublicKey, error) {
	if len(certChain) == 0 {
		return nil, fmt.Errorf("empty x509 certificate chain")
	}

	// Load root certificates from the service config.
	roots := x509.NewCertPool()
	if service.CABundle != "" && !roots.AppendCertsFromPEM([]byte(service.CABundle)) {
		return nil, fmt.Errorf("unable to load certs from service CA Bundle")
	}

	if service.Privileged {
		for _, privilegedRootCert := range s.privilegedRootCerts {
			roots.AddCert(privilegedRootCert)
		}
	}

	// Ensure the first element is encoded correctly.
	// Note: these certs use the standard b64-encoding, not PEM format.
	leafCertDer, err := base64.StdEncoding.DecodeString(certChain[0])
	if err != nil {
		return nil, fmt.Errorf("unable to decode leaf certificate: %s", err)
	}

	// And that it is a valid x509 certificate.
	leafCert, err := x509.ParseCertificate(leafCertDer)
	if err != nil {
		return nil, fmt.Errorf("unable to parse leaf certificate: %s", err)
	}

	// The rest of the certs in the chain (if any) are intermediates.
	intermediates := x509.NewCertPool()
	for i := 1; i < len(certChain); i++ {
		intermediateCertDer, err := base64.StdEncoding.DecodeString(certChain[i])
		if err != nil {
			return nil, fmt.Errorf("unable to decode intermediate certificate: %s", err)
		}

		intermediateCert, err := x509.ParseCertificate(intermediateCertDer)
		if err != nil {
			return nil, fmt.Errorf("unable to parse intermediate certificate: %s", err)
		}

		intermediates.AddCert(intermediateCert)
	}

	verifyOpts := x509.VerifyOptions{
		Intermediates: intermediates,
		Roots:         roots,
		// The service should be able to reuse certs that is has other
		// uses for and not just x509.KeyUsageDigitalSignature which
		// this use case would strictly fall under. Most commonly, the
		// service will be reusing a key/cert they have been issued for
		// client or server authentication.
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	// TODO: this call returns certificate chains which we ignore for now,
	// but we should check them for revocations if we have the ability
	// later. For now, it is the service's responsibility to update their
	// registered CA bundle if that ever changes.
	if _, err = leafCert.Verify(verifyOpts); err != nil {
		return nil, fmt.Errorf("unable to verify certificate chain: %s", err)
	}

	// Get the public key from the leaf certificate.
	leafKey, err := jose.NewPublicKey(leafCert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to get public key from leaf certificate: %s", err)
	}

	return leafKey, nil
}

func (s *Service) getServiceSigningKey(ctx context.Context, service *schema.Service, keyID string) (*jose.PublicKey, error) {
	serviceKey, err := s.schemaMgr.GetServiceKey(service.ID, keyID)
	if err != nil && err != schema.ErrNoSuchServiceKey {
		return nil, errors.Internal(ctx, fmt.Errorf("unable to get service key: %s", err))
	}

	if err == nil && time.Now().Before(serviceKey.JWK.Expiration) {
		pubKey, err := serviceKey.JWK.PublicKey()
		if err != nil {
			return pubKey, errors.Internal(ctx, fmt.Errorf("unable to convert service key to usable public key: %s", err))
		}

		return pubKey, nil
	}

	// Need to fetch the latest JWK set for this service. Try each
	// registered JWKsURI until one succeeds or they all fail.
	var (
		success bool
		keySet  jose.JWKSet
		maxAge  time.Duration
		lastErr error
	)
	for _, jwksURI := range service.JWKsURIs {
		keySet, maxAge, err = s.fetchServiceKeys(ctx, jwksURI, service.CABundle)
		if err != nil {
			lastErr = fmt.Errorf("unable to fetch service keys from %s: %s", jwksURI, err)
			continue
		}

		success = true
		break
	}

	if !success {
		return nil, lastErr
	}

	expiration := time.Now().Add(maxAge)

	var selectedKey *jose.PublicKey

	newKeys := make([]schema.JWK, len(keySet.Keys))
	for i, key := range keySet.Keys {
		newKeys[i] = schema.JWK{
			ID:          key.ID,
			KeyType:     key.KeyType,
			Modulus:     key.Modulus,
			Exponent:    key.Exponent,
			Curve:       key.Curve,
			XCoordinate: key.XCoordinate,
			YCoordinate: key.YCoordinate,
			Expiration:  expiration,
		}

		if key.ID == keyID {
			selectedKey = key
		}
	}

	if err := s.schemaMgr.SaveServiceKeys(service.ID, newKeys...); err != nil {
		return nil, errors.Internal(ctx, fmt.Errorf("unable to save new service keys: %s", err))
	}

	if selectedKey == nil {
		return nil, fmt.Errorf("unable to find key ID %q at %s", keyID, service.JWKsURIs)
	}

	return selectedKey, nil
}

func (s *Service) fetchServiceKeys(ctx context.Context, jwksURI, caBundle string) (jwkSet jose.JWKSet, maxAge time.Duration, err error) {
	var tlsConfig *tls.Config
	if caBundle != "" {
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM([]byte(caBundle)) {
			return jwkSet, 0, fmt.Errorf("unable to load service CA bundle")
		}
		tlsConfig = &tls.Config{RootCAs: caPool}
	}

	client := server.NewHTTPClient(tlsConfig)

	resp, err := client.Get(jwksURI)
	if err != nil {
		return jwkSet, 0, fmt.Errorf("HTTP error: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return jwkSet, 0, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
		return jwkSet, 0, fmt.Errorf("unable to decode JWK set: %s", err)
	}

	// Use Cache-Control 'max-age' value if specified.
	cacheControlOpts := strings.Split(resp.Header.Get("Cache-Control"), ", ")
	for _, opt := range cacheControlOpts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.ToLower(parts[0]) == "max-age" {
			if val, err := strconv.ParseUint(parts[1], 10, 32); err == nil {
				maxAge = time.Second * time.Duration(val)
				break // Found a valid max-age value.
			}
		}
	}

	// We do not want to cache keys longer than 24 hours, but do want to
	// cache them for at least one hour.
	if maxAge > time.Hour*24 {
		maxAge = time.Hour * 24
	} else if maxAge < time.Hour {
		maxAge = time.Hour
	}

	return jwkSet, maxAge, nil
}
