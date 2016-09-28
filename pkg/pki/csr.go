package pki

type CertificateSigningRequest struct {
	CertificateRequest string `json:"certificate_request,omitempty"`
	Profile            string `json:"profile"`
}
