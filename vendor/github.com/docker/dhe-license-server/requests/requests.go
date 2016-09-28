package requests

import (
	"errors"
	"time"
)

// Errors

const (
	LicenseNotFound    = "License does not exist"
	InvalidLicense     = "This License is invalid"
	FailedToGenerate   = "Failed to generate License"
	FailedToUpdate     = "Failed to update License"
	InternalServer     = "There was an error with the server. Please retry your request later"
	DuplicateAlias     = "This hub user already has a license with this alias"
	Unauthorized       = "401 Unauthorized"
	InvalidHubUsername = "There is no hub user associated with the provided hub username"
	PackageCloud       = "Failed to execute Package Cloud transaction"
)

var (
	NotFoundError           = errors.New(LicenseNotFound)
	InvalidLicenseError     = errors.New(InvalidLicense)
	FailedToGenerateError   = errors.New(FailedToGenerate)
	FailedToUpdateError     = errors.New(FailedToUpdate)
	InternalServerError     = errors.New(InternalServer)
	DuplicateAliasError     = errors.New(DuplicateAlias)
	UnauthorizedError       = errors.New(Unauthorized)
	InvalidHubUsernameError = errors.New(InvalidHubUsername)
	PackageCloudError       = errors.New(PackageCloud)
)

// Note that only internal requests will pass around this struct, so hopefully
// we will be able to avoid back compatibility issues from changes in the contents
// of a license by updating APIs and their usage simultaneously
type License struct {
	ID         int    `json:"id"`
	KeyID      string `json:"keyId"`
	PrivateKey string `json:"privateKey"`

	Alias      string     `json:"alias"`
	Type       string     `json:"licenseType"`
	Tier       string     `json:"tier"`
	MaxEngines int        `json:"maxEngines"`
	Expiration *time.Time `json:"expiration"`

	HubUUID     string `json:"hubUUID"`
	BillingUUID string `json:"billingUUID"`

	PackageCloudToken   string `json:"packageCloudToken"`
	PackageCloudTokenID int    `json:"packageCloudTokenId"`
}

// License Check

// The CheckLicenseRequest request can also be used to verify ownership
type PrivateCheckLicenseRequest struct {
	KeyID string `json:"keyId"`
}

type PrivateCheckLicenseResponse struct {
	Authorization string `json:"authorization"`
}

type CheckLicenseRequest struct {
	KeyID     string    `json:"keyId"`
	Timestamp time.Time `json:"timestamp"` // In RFC3339 time format
	Token     string    `json:"token"`
}

type CheckLicenseResponse struct {
	Expiration time.Time `json:"expiration"` // In RFC3339 time format
	Token      string    `json:"token"`
	MaxEngines int       `json:"maxEngines"`
	Type       string    `json:"licenseType"`
	Tier       string    `json:"tier"`
}

// Generate License

type GenerateLicenseRequest struct {
	BillingUUID string    `json:"billingUUID"`
	HubUUID     string    `json:"hubUUID"` // You must have at least one of the HubUUID and HubUsername
	HubUsername string    `json:"hubUsername"`
	Alias       string    `json:"alias"`
	MaxEngines  int       `json:"maxEngines"`
	Expiration  time.Time `json:"expiration"` // In RFC3339 time format

	// Only one of the following two needs to be provided (type can be infered from the tier)
	// See util.go for the allowed types and their relationships
	Tier string `json:"tier"`
	Type string `json:"licenseType"`
}

type GenerateLicenseResponse struct {
	KeyID      string `json:"keyId"`
	PrivateKey string `json:"privateKey"`
}

// Update License

type UpdateLicenseRequest struct {
	KeyID      string `json:"keyId"`
	Alias      string `json:"alias"`
	MaxEngines int    `json:"maxEngines"`
	Tier       string `json:"tier"`
}

// Activate License

type ActivateLicenseRequest struct {
	KeyID      string    `json:"keyId"`
	Expiration time.Time `json:"expiration"` // In RFC3339 time format
}

type DeactivateLicenseRequest struct {
	KeyID string `json:"keyId"`
}

// Get Licenses

type GetLicensesForBillingIdResponse struct {
	Licenses []License `json:"licenses"`
}

type GetLicenseForKeyResponse struct {
	License License `json:"license"`
}

type GetLicensesForHubUserResponse struct {
	Licenses []License `json:"licenses"`
}

// Backfill Zendesk licenses

type ZendeskBackfillRequest struct {
	StartLicenseID *int `json:"startLicenseId"`
	EndLicenseID   *int `json:"endLicenseId"`
}
